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

// requirementsProject makes an empty directory the working directory and, when
// content is non-empty, seeds requirements.yml in it. Both the reader and the
// writer resolve that name relative to the working directory.
func requirementsProject(t *testing.T, content string) {
	t.Helper()

	t.Chdir(t.TempDir())

	if content != "" {
		if err := os.WriteFile("requirements.yml", []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRequirementsFileExist(t *testing.T) {
	requirementsProject(t, "")

	if RequirementsFileExist() {
		t.Error("RequirementsFileExist() = true, want false in an empty directory")
	}

	if err := os.WriteFile("requirements.yml", []byte("---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !RequirementsFileExist() {
		t.Error("RequirementsFileExist() = false, want true once the file is present")
	}
}

func TestEnsureRequirementsFile(t *testing.T) {
	requirementsProject(t, "")

	err := EnsureRequirementsFile()
	if err == nil {
		t.Fatal("EnsureRequirementsFile() = nil, want an error when the file is missing")
	}

	if !strings.Contains(err.Error(), "requirements.yml") {
		t.Errorf("error = %q, want it to name requirements.yml", err)
	}

	requirementsProject(t, "---\n")

	if err := EnsureRequirementsFile(); err != nil {
		t.Errorf("EnsureRequirementsFile() = %v, want nil once the file exists", err)
	}
}

func TestReadRequirements_ParsesRolesAndCollections(t *testing.T) {
	requirementsProject(t, `---
collections:
  - name: dcjulian29.workstation
    type: git
    source: https://github.com/dcjulian29/ansible-runbook-workstation.git
roles:
  - name: dcjulian29.nginx
    src: https://github.com/dcjulian29/ansible-role-nginx.git
`)

	got, err := ReadRequirements()
	if err != nil {
		t.Fatalf("ReadRequirements() error = %v, want nil", err)
	}

	if len(got.Roles) != 1 || got.Roles[0].Name != "dcjulian29.nginx" {
		t.Errorf("Roles = %+v, want one entry named dcjulian29.nginx", got.Roles)
	}

	if len(got.Collections) != 1 || got.Collections[0].Name != "dcjulian29.workstation" {
		t.Errorf("Collections = %+v, want one entry named dcjulian29.workstation", got.Collections)
	}
}

func TestReadRequirements_ErrorsWhenMissing(t *testing.T) {
	requirementsProject(t, "")

	if _, err := ReadRequirements(); err == nil {
		t.Error("ReadRequirements() error = nil, want an error when the file is absent")
	}
}

func TestReadRequirements_ErrorsOnMalformedYAML(t *testing.T) {
	requirementsProject(t, "roles: [unterminated\n")

	if _, err := ReadRequirements(); err == nil {
		t.Error("ReadRequirements() error = nil, want a parse error")
	}
}

// An empty document is valid YAML and must read as an empty set rather than
// failing, since that is what a freshly initialized project contains.
func TestReadRequirements_EmptyDocumentIsNotAnError(t *testing.T) {
	requirementsProject(t, "---\n")

	got, err := ReadRequirements()
	if err != nil {
		t.Fatalf("ReadRequirements() error = %v, want nil", err)
	}

	if len(got.Roles) != 0 || len(got.Collections) != 0 {
		t.Errorf("Requirements = %+v, want both lists empty", got)
	}
}

func TestSaveRequirements_RoundTrips(t *testing.T) {
	requirementsProject(t, "")

	want := Requirements{
		Roles: []Role{{
			Name:   "acme.nginx",
			Source: "https://github.com/acme/ansible-role-nginx.git",
		}},
	}

	if err := SaveRequirements(want); err != nil {
		t.Fatalf("SaveRequirements() error = %v, want nil", err)
	}

	got, err := ReadRequirements()
	if err != nil {
		t.Fatalf("ReadRequirements() error = %v, want nil", err)
	}

	if len(got.Roles) != 1 {
		t.Fatalf("Roles = %+v, want one entry", got.Roles)
	}

	if got.Roles[0].Name != want.Roles[0].Name || got.Roles[0].Source != want.Roles[0].Source {
		t.Errorf("Roles[0] = %+v, want %+v", got.Roles[0], want.Roles[0])
	}
}

// TestSaveRequirements_TruncatesExistingContent pins the O_TRUNC behavior:
// writing a shorter document over a longer one must not leave a tail behind,
// which would produce a file that no longer parses.
func TestSaveRequirements_TruncatesExistingContent(t *testing.T) {
	requirementsProject(t, "")

	long := Requirements{
		Roles: []Role{
			{Name: "acme.one", Source: "https://github.com/acme/ansible-role-one.git"},
			{Name: "acme.two", Source: "https://github.com/acme/ansible-role-two.git"},
			{Name: "acme.three", Source: "https://github.com/acme/ansible-role-three.git"},
		},
	}

	if err := SaveRequirements(long); err != nil {
		t.Fatal(err)
	}

	short := Requirements{
		Roles: []Role{{Name: "acme.one", Source: "https://github.com/acme/ansible-role-one.git"}},
	}

	if err := SaveRequirements(short); err != nil {
		t.Fatal(err)
	}

	got, err := ReadRequirements()
	if err != nil {
		t.Fatalf("ReadRequirements() error = %v, want nil after overwriting", err)
	}

	if len(got.Roles) != 1 {
		t.Errorf("Roles = %+v, want exactly one entry after the shorter save", got.Roles)
	}
}
