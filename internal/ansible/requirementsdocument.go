/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ansible

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// requirementsFile is resolved against the working directory. Every command
// that edits it runs EnsureAnsibleDirectory in PreRunE, which requires
// ansible.cfg in the working directory, so that directory is the development
// root.
const requirementsFile = "requirements.yml"

// RequirementsDocument is requirements.yml held as a yaml.Node tree rather
// than decoded into [Requirements]. Editing the tree changes only the entry
// being added or removed: keys the structs do not know (a collection's
// "source:"), the order of entries, and comments all survive, where a
// decode-and-remarshal round trip dropped or rewrote them.
type RequirementsDocument struct {
	doc *yaml.Node
}

// OpenRequirements reads requirements.yml for editing. A missing file, a read
// failure, or YAML that does not parse is an error: carrying on with an empty
// document would make the next Save replace the whole file with the one entry
// being added. An empty document ("---" alone) is valid and opens as an empty
// mapping.
func OpenRequirements() (*RequirementsDocument, error) {
	if err := EnsureRequirementsFile(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(requirementsFile)
	if err != nil {
		return nil, err
	}

	var doc yaml.Node

	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("reading %s: %w", requirementsFile, err)
	}

	// An empty file decodes to a zero node and "---" alone to a document
	// holding a null scalar; both become a document with an empty mapping so
	// entries can be added to it.
	if doc.Kind == 0 || (len(doc.Content) == 1 && doc.Content[0].Tag == "!!null") {
		doc = yaml.Node{
			Kind:    yaml.DocumentNode,
			Content: []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}},
		}
	}

	// The legacy roles-only format (a bare list) and anything else that is
	// not a mapping is refused rather than rewritten.
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s must be a mapping with collections and roles keys", requirementsFile)
	}

	return &RequirementsDocument{doc: &doc}, nil
}

// Save writes the document back to requirements.yml with a leading "---" and
// a 2-space indent, matching the file's own style and the repository's
// yamllint rules (yaml.v3 defaults to a 4-space indent and no "---").
//
// The bytes go to a temporary file that is then renamed over
// requirements.yml, so a failed write leaves the original intact instead of
// truncated.
func (r *RequirementsDocument) Save() error {
	var buf bytes.Buffer

	buf.WriteString("---\n")

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	if err := enc.Encode(r.doc); err != nil {
		return err
	}

	if err := enc.Close(); err != nil {
		return err
	}

	tmp := requirementsFile + ".tmp"

	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return err
	}

	if err := os.Rename(tmp, requirementsFile); err != nil {
		os.Remove(tmp) //nolint:errcheck

		return err
	}

	return nil
}

// section returns the list under key ("roles" or "collections"). A key with
// nothing after it ("roles:") reads as null and is treated as an empty list.
// When create is set, a missing or null key becomes an empty list so an entry
// can be appended; otherwise nil is returned for it. A key holding anything
// other than a list is an error rather than something to overwrite.
func (r *RequirementsDocument) section(key string, create bool) (*yaml.Node, error) {
	root := r.doc.Content[0]

	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value != key {
			continue
		}

		value := root.Content[i+1]

		switch {
		case value.Kind == yaml.SequenceNode:
			return value, nil
		case value.Tag == "!!null":
			if create {
				*value = yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

				return value, nil
			}

			return nil, nil
		default:
			return nil, fmt.Errorf("%s in %s is not a list", key, requirementsFile)
		}
	}

	if !create {
		return nil, nil
	}

	list := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		list)

	return list, nil
}

// AddRole appends role to the roles list, written as name, src and version
// with empty fields left out. A role already listed under the same name is an
// error, rather than a second entry ansible-galaxy would install twice.
func (r *RequirementsDocument) AddRole(role Role) error {
	return r.add("roles", "role", entry("name", role.Name, "src", role.Source, "version", role.Version))
}

// AddCollection appends collection to the collections list, written as name,
// type, source and version with empty fields left out. Collections take
// "source:", not the "src:" roles use. A collection already listed under the
// same name is an error.
func (r *RequirementsDocument) AddCollection(collection Collection) error {
	return r.add("collections", "collection", entry(
		"name", collection.Name,
		"type", collection.Type,
		"source", collection.Source,
		"version", collection.Version))
}

// RemoveRole removes every roles entry named name and reports whether there
// was one. Removing the last role leaves an empty list, which is still saved.
func (r *RequirementsDocument) RemoveRole(name string) (bool, error) {
	return r.remove("roles", name)
}

// RemoveCollection removes every collections entry named name and reports
// whether there was one.
func (r *RequirementsDocument) RemoveCollection(name string) (bool, error) {
	return r.remove("collections", name)
}

func (r *RequirementsDocument) add(key, kind string, item *yaml.Node) error {
	list, err := r.section(key, true)
	if err != nil {
		return err
	}

	name := entryName(item)

	for _, existing := range list.Content {
		if entryName(existing) == name {
			return fmt.Errorf("%s '%s' is already in %s", kind, name, requirementsFile)
		}
	}

	// An empty list is written "[]" (Save does that itself once the last entry
	// is removed), and the encoder keeps a parsed node's flow style, so an
	// entry added to it would come out inline as "roles: [{name: ...}]".
	// Clearing the flag writes the list, and a root mapping read as "{}", as
	// block YAML like the rest of the file.
	list.Style &^= yaml.FlowStyle
	r.doc.Content[0].Style &^= yaml.FlowStyle

	list.Content = append(list.Content, item)

	return nil
}

func (r *RequirementsDocument) remove(key, name string) (bool, error) {
	list, err := r.section(key, false)
	if err != nil || list == nil {
		return false, err
	}

	kept := list.Content[:0]

	for _, item := range list.Content {
		if entryName(item) != name {
			kept = append(kept, item)
		}
	}

	found := len(kept) < len(list.Content)
	list.Content = kept

	return found, nil
}

// entry builds one list item from key/value pairs, in the order given and
// leaving out empty values, so a role without a version gets no
// `version: ""` line. Values are tagged as strings, which makes the encoder
// quote one such as "1.0" that would otherwise read back as a number.
func entry(pairs ...string) *yaml.Node {
	item := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}

	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] == "" {
			continue
		}

		item.Content = append(item.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: pairs[i]},
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: pairs[i+1]})
	}

	return item
}

// entryName returns the name of a list item, or "" for an item without one
// (a bare "- src: ..." entry, which ansible-galaxy also accepts).
func entryName(item *yaml.Node) string {
	if item.Kind != yaml.MappingNode {
		return ""
	}

	for i := 0; i+1 < len(item.Content); i += 2 {
		if item.Content[i].Value == "name" {
			return item.Content[i+1].Value
		}
	}

	return ""
}
