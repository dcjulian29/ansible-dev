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
	"os"
	"strings"
	"testing"
)

// developmentRequirements has the shape of the development repository's real
// requirements.yml: a git collection with "type:" and "source:" (the keys the
// old struct round trip dropped), Galaxy entries with nothing but a name, and
// a comment.
const developmentRequirements = `---
# Installed by "ansible-dev restore".
collections:
  - name: ansible.posix
  - name: dcjulian29.updateos
    type: git
    source: https://github.com/dcjulian29/ansible-runbook-updateos.git
roles:
  - name: dcjulian29.base
    src: https://github.com/dcjulian29/ansible-role-base.git
  - name: geerlingguy.java
    src: geerlingguy.java
`

func openRequirements(t *testing.T) *RequirementsDocument {
	t.Helper()

	doc, err := OpenRequirements()
	if err != nil {
		t.Fatalf("OpenRequirements() error = %v, want nil", err)
	}

	return doc
}

func save(t *testing.T, doc *RequirementsDocument) string {
	t.Helper()

	if err := doc.Save(); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	return read(t, requirementsFile)
}

// Opening and saving without a change must give back the same bytes; any
// difference here would show up as a diff in every edit.
func TestRequirementsDocument_SaveWithoutChangesIsByteIdentical(t *testing.T) {
	requirementsProject(t, developmentRequirements)

	if got := save(t, openRequirements(t)); got != developmentRequirements {
		t.Errorf("requirements.yml changed on a plain round trip:\n%s", got)
	}
}

// Adding a role appends exactly its lines and leaves everything else alone.
func TestRequirementsDocument_AddRoleChangesOnlyItsEntry(t *testing.T) {
	requirementsProject(t, developmentRequirements)

	doc := openRequirements(t)

	err := doc.AddRole(Role{Name: "dcjulian29.nginx", Source: "https://github.com/dcjulian29/ansible-role-nginx.git"})
	if err != nil {
		t.Fatalf("AddRole() error = %v, want nil", err)
	}

	want := developmentRequirements +
		"  - name: dcjulian29.nginx\n" +
		"    src: https://github.com/dcjulian29/ansible-role-nginx.git\n"

	if got := save(t, doc); got != want {
		t.Errorf("requirements.yml =\n%s\nwant\n%s", got, want)
	}

	if _, err := os.Stat(requirementsFile + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("the temporary file was left behind (stat error = %v)", err)
	}
}

// Empty fields are left out, and a numeric-looking version stays a string.
func TestRequirementsDocument_AddRoleOmitsEmptyFieldsAndQuotesVersions(t *testing.T) {
	requirementsProject(t, "---\nroles:\n")

	doc := openRequirements(t)

	if err := doc.AddRole(Role{Name: "geerlingguy.pip", Source: "geerlingguy.pip"}); err != nil {
		t.Fatal(err)
	}

	if err := doc.AddRole(Role{Name: "acme.x", Source: "acme.x", Version: "1.0"}); err != nil {
		t.Fatal(err)
	}

	want := "---\nroles:\n" +
		"  - name: geerlingguy.pip\n    src: geerlingguy.pip\n" +
		"  - name: acme.x\n    src: acme.x\n    version: \"1.0\"\n"

	if got := save(t, doc); got != want {
		t.Errorf("requirements.yml =\n%s\nwant\n%s", got, want)
	}
}

// A collection is written with "source:", the key ansible-galaxy reads for
// collections, in the same order the existing git entries use.
func TestRequirementsDocument_AddCollectionWritesSource(t *testing.T) {
	requirementsProject(t, developmentRequirements)

	doc := openRequirements(t)

	err := doc.AddCollection(Collection{
		Name:   "dcjulian29.reboot_system",
		Type:   "git",
		Source: "https://github.com/dcjulian29/ansible-runbook-reboot_system.git",
	})
	if err != nil {
		t.Fatalf("AddCollection() error = %v, want nil", err)
	}

	added := "    source: https://github.com/dcjulian29/ansible-runbook-updateos.git\n" +
		"  - name: dcjulian29.reboot_system\n" +
		"    type: git\n" +
		"    source: https://github.com/dcjulian29/ansible-runbook-reboot_system.git\n" +
		"roles:\n"

	if got := save(t, doc); !strings.Contains(got, added) {
		t.Errorf("collection not appended to the end of collections:\n%s", got)
	}
}

func TestRequirementsDocument_AddRejectsADuplicate(t *testing.T) {
	requirementsProject(t, developmentRequirements)

	doc := openRequirements(t)

	if err := doc.AddRole(Role{Name: "dcjulian29.base", Source: "x"}); err == nil {
		t.Error("AddRole() error = nil, want an error for a role already listed")
	}

	if err := doc.AddCollection(Collection{Name: "ansible.posix"}); err == nil {
		t.Error("AddCollection() error = nil, want an error for a collection already listed")
	}
}

// Removing the only role must still be saved. The old code skipped the save
// when the list came out empty, so the entry stayed and the command exited 0.
func TestRequirementsDocument_RemoveTheOnlyRole(t *testing.T) {
	requirementsProject(t, "---\nroles:\n  - name: acme.x\n    src: acme.x\n")

	doc := openRequirements(t)

	found, err := doc.RemoveRole("acme.x")
	if err != nil || !found {
		t.Fatalf("RemoveRole() = %v, %v; want true, nil", found, err)
	}

	if got := save(t, doc); strings.Contains(got, "acme.x") {
		t.Errorf("the role is still in requirements.yml:\n%s", got)
	}
}

func TestRequirementsDocument_RemoveLeavesTheRestAlone(t *testing.T) {
	requirementsProject(t, developmentRequirements)

	doc := openRequirements(t)

	found, err := doc.RemoveCollection("dcjulian29.updateos")
	if err != nil || !found {
		t.Fatalf("RemoveCollection() = %v, %v; want true, nil", found, err)
	}

	want := strings.Replace(developmentRequirements,
		"  - name: dcjulian29.updateos\n"+
			"    type: git\n"+
			"    source: https://github.com/dcjulian29/ansible-runbook-updateos.git\n", "", 1)

	if got := save(t, doc); got != want {
		t.Errorf("requirements.yml =\n%s\nwant\n%s", got, want)
	}
}

func TestRequirementsDocument_RemoveReportsAMissingName(t *testing.T) {
	requirementsProject(t, "---\n")

	found, err := openRequirements(t).RemoveRole("acme.x")
	if err != nil || found {
		t.Errorf("RemoveRole() = %v, %v; want false, nil", found, err)
	}
}

// A project whose requirements.yml is just "---" gets the key created.
func TestRequirementsDocument_AddToAnEmptyDocument(t *testing.T) {
	requirementsProject(t, "---\n")

	doc := openRequirements(t)

	if err := doc.AddRole(Role{Name: "acme.x", Source: "acme.x"}); err != nil {
		t.Fatal(err)
	}

	want := "---\nroles:\n  - name: acme.x\n    src: acme.x\n"

	if got := save(t, doc); got != want {
		t.Errorf("requirements.yml =\n%s\nwant\n%s", got, want)
	}
}

// Removing the last role saves "roles: []", and a document can be written
// "{}". Adding to either must give block YAML, not the inline
// "roles: [{name: acme.x, src: acme.x}]" the parsed flow style produced.
func TestRequirementsDocument_AddToAFlowStyleList(t *testing.T) {
	want := "---\nroles:\n  - name: acme.x\n    src: acme.x\n"

	for name, content := range map[string]string{
		"empty list":    "---\nroles: []\n",
		"empty mapping": "--- {}\n",
	} {
		t.Run(name, func(t *testing.T) {
			requirementsProject(t, content)

			doc := openRequirements(t)

			if err := doc.AddRole(Role{Name: "acme.x", Source: "acme.x"}); err != nil {
				t.Fatal(err)
			}

			if got := save(t, doc); got != want {
				t.Errorf("requirements.yml =\n%s\nwant\n%s", got, want)
			}
		})
	}
}

// A file holding only comments used to lose them on the first add: yaml.v3
// keeps them as a foot comment after "---" and drops them without it, and
// the empty document was rebuilt from scratch.
func TestRequirementsDocument_AddKeepsCommentsOfACommentOnlyFile(t *testing.T) {
	added := "roles:\n  - name: acme.x\n    src: acme.x\n"

	for name, tc := range map[string]struct{ content, want string }{
		"with ---":    {"---\n# pinned for prod\n", "---\n# pinned for prod\n\n" + added},
		"without ---": {"# pinned for prod\n", "---\n# pinned for prod\n\n" + added},
		"two blocks":  {"---\n# top\n\n# second\n", "---\n# top\n\n# second\n\n" + added},
		"--- alone":   {"---\n", "---\n" + added},
	} {
		t.Run(name, func(t *testing.T) {
			requirementsProject(t, tc.content)

			doc := openRequirements(t)

			if err := doc.AddRole(Role{Name: "acme.x", Source: "acme.x"}); err != nil {
				t.Fatal(err)
			}

			if got := save(t, doc); got != tc.want {
				t.Errorf("requirements.yml =\n%s\nwant\n%s", got, tc.want)
			}
		})
	}
}

// Each of these must stop the command: opening an empty document instead
// would let the next Save replace the file with a single entry.
func TestOpenRequirements_Errors(t *testing.T) {
	for name, content := range map[string]string{
		"missing":   "",
		"malformed": "roles: [unterminated\n",
		"bare list": "- name: acme.x\n",
		"not lists": "roles: acme.x\n",
	} {
		t.Run(name, func(t *testing.T) {
			requirementsProject(t, content)

			doc, err := OpenRequirements()
			if err == nil && name == "not lists" {
				// The file opens; the error comes when the roles key is used.
				err = doc.AddRole(Role{Name: "acme.y"})
			}

			if err == nil {
				t.Error("error = nil, want one")
			}
		})
	}
}
