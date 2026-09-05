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
	"path/filepath"
	"strings"
	"testing"
)

// ansibleProject creates a project directory containing the given ansible.cfg
// and makes it the working directory for the test. The readers in this file all
// resolve "ansible.cfg" relative to the working directory, so every one of them
// needs this.
func ansibleProject(t *testing.T, cfg string) string {
	t.Helper()

	dir := t.TempDir()

	if cfg != "" {
		if err := os.WriteFile(filepath.Join(dir, "ansible.cfg"), []byte(cfg), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Chdir(dir)

	return dir
}

const defaultsCfg = "[defaults]\nroles_path = ./roles\ncollections_path = ./collections\n"

func TestEnsureAnsibleDirectory_AcceptsAProjectDirectory(t *testing.T) {
	ansibleProject(t, defaultsCfg)

	if err := EnsureAnsibleDirectory(); err != nil {
		t.Errorf("EnsureAnsibleDirectory() = %v, want nil", err)
	}
}

func TestEnsureAnsibleDirectory_RejectsOtherDirectories(t *testing.T) {
	ansibleProject(t, "")

	err := EnsureAnsibleDirectory()
	if err == nil {
		t.Fatal("EnsureAnsibleDirectory() = nil, want an error without ansible.cfg")
	}

	if !strings.Contains(err.Error(), "ansible development folder") {
		t.Errorf("error = %q, want it to say this is not an ansible folder", err)
	}
}

func TestRootRoleFolder_ReadsTheDefaultsSection(t *testing.T) {
	ansibleProject(t, defaultsCfg)

	got, err := RootRoleFolder()
	if err != nil {
		t.Fatalf("RootRoleFolder() error = %v, want nil", err)
	}

	if got != "./roles" {
		t.Errorf("RootRoleFolder() = %q, want %q", got, "./roles")
	}
}

func TestRootRoleFolder_ErrorsWhenTheKeyIsMissing(t *testing.T) {
	tests := []struct {
		name string
		cfg  string
	}{
		{name: "no defaults section", cfg: "[galaxy]\nserver = https://galaxy.ansible.com\n"},
		{name: "no roles_path key", cfg: "[defaults]\ncollections_path = ./collections\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ansibleProject(t, tt.cfg)

			if _, err := RootRoleFolder(); err == nil {
				t.Error("RootRoleFolder() error = nil, want an error")
			}
		})
	}
}

func TestRootRoleFolder_ErrorsWithoutAnsibleCfg(t *testing.T) {
	ansibleProject(t, "")

	if _, err := RootRoleFolder(); err == nil {
		t.Error("RootRoleFolder() error = nil, want an error when ansible.cfg is absent")
	}
}

func TestRoleFolder_JoinsTheRoleOntoTheRootFolder(t *testing.T) {
	ansibleProject(t, defaultsCfg)

	got, err := RoleFolder("dcjulian29.nginx")
	if err != nil {
		t.Fatalf("RoleFolder() error = %v, want nil", err)
	}

	// The role is appended verbatim: the directory keeps whatever form the
	// caller used, namespace included.
	want := filepath.Join("roles", "dcjulian29.nginx")
	if got != want {
		t.Errorf("RoleFolder() = %q, want %q", got, want)
	}
}

func TestRoleFolderExists(t *testing.T) {
	dir := ansibleProject(t, defaultsCfg)

	if err := os.MkdirAll(filepath.Join(dir, "roles", "nginx"), 0o755); err != nil {
		t.Fatal(err)
	}

	if !RoleFolderExists("nginx") {
		t.Error("RoleFolderExists(\"nginx\") = false, want true")
	}

	if RoleFolderExists("absent") {
		t.Error("RoleFolderExists(\"absent\") = true, want false")
	}
}

// RoleFolderExists swallows the resolution error, so a directory with no
// ansible.cfg must report false rather than panicking on the empty path.
func TestRoleFolderExists_FalseWithoutAnsibleCfg(t *testing.T) {
	ansibleProject(t, "")

	if RoleFolderExists("nginx") {
		t.Error("RoleFolderExists() = true, want false when ansible.cfg is absent")
	}
}

func TestCollectionsFolder_AppendsAnsibleCollections(t *testing.T) {
	ansibleProject(t, defaultsCfg)

	got, err := CollectionsFolder()
	if err != nil {
		t.Fatalf("CollectionsFolder() error = %v, want nil", err)
	}

	// The install layout is <collections_path>/ansible_collections/<ns>/<name>,
	// so the reader owns the fixed middle segment.
	want := filepath.Join("collections", "ansible_collections")
	if got != want {
		t.Errorf("CollectionsFolder() = %q, want %q", got, want)
	}
}

func TestCollectionsFolder_ErrorsWhenTheKeyIsMissing(t *testing.T) {
	ansibleProject(t, "[defaults]\nroles_path = ./roles\n")

	if _, err := CollectionsFolder(); err == nil {
		t.Error("CollectionsFolder() error = nil, want an error without collections_path")
	}
}
