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

package settings

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// exampleFile is the shipped example configuration, relative to this package.
const exampleFile = "../../ansible-dev.yml.example"

// loadExample decodes the example configuration with unknown fields rejected,
// so a key that no longer exists on Config fails the test rather than being
// silently ignored.
func loadExample(t *testing.T) Config {
	t.Helper()

	data, err := os.ReadFile(exampleFile)
	if err != nil {
		t.Fatalf("reading %s: %v", exampleFile, err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		t.Fatalf("%s does not match Config: %v", exampleFile, err)
	}

	return cfg
}

func TestExampleFile_DecodesIntoConfig(t *testing.T) {
	cfg := loadExample(t)

	if cfg.RolesPath == "" {
		t.Error("roles_path is empty; the example should show a sample path")
	}

	if cfg.RunbooksPath == "" {
		t.Error("runbooks_path is empty; the example should show a sample path")
	}

	if len(cfg.Diff) == 0 {
		t.Error("diff has no operating system entries")
	}

	for name, diff := range cfg.Diff {
		if diff.Program == "" {
			t.Errorf("diff[%s] has no program", name)
		}
	}
}

// TestExampleFile_IgnoreEntriesHaveNoEscapes guards against reintroducing
// regex-style escapes in the ignore lists. Entries are matched with a plain
// strings.Contains against the full path, so an entry such as "\.git" only
// matches on Windows, where the backslash is the path separator, and silently
// compares .git directories everywhere else.
func TestExampleFile_IgnoreEntriesHaveNoEscapes(t *testing.T) {
	cfg := loadExample(t)

	lists := map[string][]string{
		"role_ignore":    cfg.RoleIgnore,
		"runbook_ignore": cfg.RunbookIgnore,
	}

	for name, list := range lists {
		if len(list) == 0 {
			t.Errorf("%s is empty; the example should show sample entries", name)
		}

		for _, entry := range list {
			if strings.ContainsAny(entry, `\/`) {
				t.Errorf("%s entry %q contains a path separator; entries are "+
					"matched as plain substrings and must be portable", name, entry)
			}
		}
	}
}
