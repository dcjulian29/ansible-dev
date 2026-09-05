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

func TestRemoveRole_DeletesTheRoleTree(t *testing.T) {
	dir := ansibleProject(t, defaultsCfg)

	role := filepath.Join(dir, "roles", "nginx")
	writeFile(t, role, "tasks/main.yml", "---\n")

	captureStdout(t, func() {
		if err := RemoveRole("nginx"); err != nil {
			t.Errorf("RemoveRole() error = %v, want nil", err)
		}
	})

	if _, err := os.Stat(role); !os.IsNotExist(err) {
		t.Errorf("role directory still present: %v", err)
	}
}

func TestRemoveRole_ErrorsWhenTheRoleIsAbsent(t *testing.T) {
	ansibleProject(t, defaultsCfg)

	err := RemoveRole("absent")
	if err == nil {
		t.Fatal("RemoveRole() error = nil, want an error")
	}

	if !strings.Contains(err.Error(), "not present") {
		t.Errorf("error = %q, want it to say the folder is not present", err)
	}
}

// An empty role directory is reported rather than removed: it usually means the
// role was already cleaned up, and deleting silently would hide that.
func TestRemoveRole_ErrorsWhenTheRoleFolderIsEmpty(t *testing.T) {
	dir := ansibleProject(t, defaultsCfg)

	if err := os.MkdirAll(filepath.Join(dir, "roles", "hollow"), 0o755); err != nil {
		t.Fatal(err)
	}

	err := RemoveRole("hollow")
	if err == nil {
		t.Fatal("RemoveRole() error = nil, want an error for an empty role folder")
	}

	if !strings.Contains(err.Error(), "files not present") {
		t.Errorf("error = %q, want it to mention the missing files", err)
	}
}
