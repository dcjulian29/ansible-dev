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
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// homeEnv is the variable os.UserHomeDir consults, which differs by platform.
func homeEnv() string {
	if runtime.GOOS == "windows" {
		return "USERPROFILE"
	}

	return "HOME"
}

// TestResolvePath_ExpandsTilde is the regression test for a "~"-rooted setting
// being used verbatim: every existence check against it failed, so "role
// compare" skipped every role and exited 0 without printing anything.
func TestResolvePath_ExpandsTilde(t *testing.T) {
	home := t.TempDir()
	t.Setenv(homeEnv(), home)

	want := filepath.Join(home, "code", "roles")
	if err := os.MkdirAll(want, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := resolvePath("~/code/roles", "roles_path", "roles-path", LegacyRolesEnv)
	if err != nil {
		t.Fatalf("resolvePath() error = %v, want nil", err)
	}

	if got != want {
		t.Errorf("resolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePath_AcceptsExistingAbsolutePath(t *testing.T) {
	dir := t.TempDir()

	got, err := resolvePath(dir, "roles_path", "roles-path", LegacyRolesEnv)
	if err != nil {
		t.Fatalf("resolvePath() error = %v, want nil", err)
	}

	if got != dir {
		t.Errorf("resolvePath() = %q, want %q", got, dir)
	}
}

func TestResolvePath_MissingDirectoryIsAnError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	_, err := resolvePath(missing, "roles_path", "roles-path", LegacyRolesEnv)
	if err == nil {
		t.Fatal("resolvePath() error = nil, want an error for a missing directory")
	}

	if !strings.Contains(err.Error(), "roles_path") ||
		!strings.Contains(err.Error(), missing) {
		t.Errorf("error = %q, want it to name the setting and the offending path", err)
	}
}

// A path that had to be expanded should report both forms: the value as written
// in the file is what the reader has to go and edit, and the expansion is what
// explains why it was not found.
func TestResolvePath_MissingTildePathReportsBothForms(t *testing.T) {
	home := t.TempDir()
	t.Setenv(homeEnv(), home)

	_, err := resolvePath("~/code/nope", "roles_path", "roles-path", LegacyRolesEnv)
	if err == nil {
		t.Fatal("resolvePath() error = nil, want an error for a missing directory")
	}

	if !strings.Contains(err.Error(), "~/code/nope") {
		t.Errorf("error = %q, want it to quote the configured value", err)
	}

	if !strings.Contains(err.Error(), filepath.Join(home, "code", "nope")) {
		t.Errorf("error = %q, want it to show the expanded path", err)
	}
}

func TestResolvePath_EmptyValueReportsUnset(t *testing.T) {
	t.Setenv(LegacyRolesEnv, "")

	_, err := resolvePath("", "roles_path", "roles-path", LegacyRolesEnv)
	if err == nil {
		t.Fatal("resolvePath() error = nil, want an error for an unset value")
	}

	if !strings.Contains(err.Error(), "config roles-path <dir>") {
		t.Errorf("error = %q, want the unset message naming the setter", err)
	}
}

func TestUnsetError_WithoutLegacyVariablePointsAtSetter(t *testing.T) {
	t.Setenv(LegacyRolesEnv, "")

	err := unsetError("roles_path", "roles-path", LegacyRolesEnv)

	if !strings.Contains(err.Error(), "config roles-path <dir>") {
		t.Errorf("error = %q, want it to name the roles-path setter", err)
	}

	if strings.Contains(err.Error(), "import-env") {
		t.Errorf("error = %q, should not mention import-env when the variable is unset", err)
	}
}

func TestUnsetError_WithLegacyVariablePointsAtImport(t *testing.T) {
	t.Setenv(LegacyRolesEnv, "/path/to/ansible/roles")

	err := unsetError("roles_path", "roles-path", LegacyRolesEnv)

	if !strings.Contains(err.Error(), "config import-env") {
		t.Errorf("error = %q, want it to point at import-env", err)
	}

	if !strings.Contains(err.Error(), LegacyRolesEnv) {
		t.Errorf("error = %q, want it to name %s", err, LegacyRolesEnv)
	}
}
