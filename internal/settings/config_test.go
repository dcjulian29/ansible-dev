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

// sandboxHome is the fake home directory every test in this package runs
// against. TestMain points the platform's home variable at it before any test
// runs, so nothing here can reach the developer's real
// ~/.config/ansible-dev.yml.
var sandboxHome string

// TestMain redirects the home directory for the whole package.
//
// os.UserHomeDir, which the configuration package calls to locate the file, is
// consulted on each call rather than at startup, so setting the variable here
// is enough to relocate the configuration. It is done in TestMain rather than
// per-test because [Load] caches on first use: whichever test reached it first
// would otherwise fix the cache to the real file.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ansible-dev-settings")
	if err != nil {
		panic(err)
	}

	sandboxHome = dir

	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", dir) //nolint:errcheck
	} else {
		os.Setenv("HOME", dir) //nolint:errcheck
	}

	code := m.Run()

	os.RemoveAll(dir) //nolint:errcheck
	os.Exit(code)
}

// TestConfigFile_RoundTrips exercises Path, Load, Save and Show as one
// sequence. They share the package-level singleton, whose cache is populated
// once per process, so splitting them into separate tests would make each
// depend on which ran first.
func TestConfigFile_RoundTrips(t *testing.T) {
	path, err := Path()
	if err != nil {
		t.Fatalf("Path() error = %v, want nil", err)
	}

	// Guard before anything is written: if the redirect above ever stops
	// working, fail here rather than overwriting a real configuration.
	if !strings.HasPrefix(path, sandboxHome) {
		t.Fatalf("Path() = %q, which is outside the sandbox %q; refusing to write", path, sandboxHome)
	}

	if want := filepath.Join(sandboxHome, ".config", "ansible-dev.yml"); path != want {
		t.Errorf("Path() = %q, want %q", path, want)
	}

	// A missing file is not an error: a first run has no configuration yet.
	empty, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil when the file does not exist", err)
	}

	if empty.Namespace != "" || empty.RolesPath != "" || len(empty.Diff) != 0 {
		t.Errorf("Load() = %+v, want a zero configuration", empty)
	}

	cfg := Config{
		Namespace:     "acme",
		RolesPath:     "~/code/ansible/roles",
		RunbooksPath:  "~/code/ansible/runbooks",
		RoleIgnore:    []string{".git"},
		RunbookIgnore: []string{"MANIFEST.json"},
		Diff: map[string]DiffTool{
			runtime.GOOS: {Program: "diff", AdditionalArgs: []string{"{left}", "{right}"}},
		},
	}

	if err := Save(&cfg); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("configuration file was not written: %v", err)
	}

	// Save refreshes the cache, so the next Load sees the new values rather
	// than the empty result cached above.
	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if got.Namespace != cfg.Namespace || got.RolesPath != cfg.RolesPath {
		t.Errorf("Load() = %+v, want the saved values", got)
	}

	if len(got.RoleIgnore) != 1 || got.RoleIgnore[0] != ".git" {
		t.Errorf("RoleIgnore = %v, want [.git]", got.RoleIgnore)
	}

	if got.CurrentDiff().Program != "diff" {
		t.Errorf("CurrentDiff().Program = %q, want %q", got.CurrentDiff().Program, "diff")
	}

	// The resolver reads the same cache. Its unset branch cannot be reached in
	// this process once a value is cached, so that half is covered where no
	// configuration exists at all — see the cmd/role tests.
	namespace, err := Namespace()
	if err != nil {
		t.Fatalf("Namespace() error = %v, want nil once one is configured", err)
	}

	if namespace != "acme" {
		t.Errorf("Namespace() = %q, want %q", namespace, "acme")
	}

	shown, err := Show()
	if err != nil {
		t.Fatalf("Show() error = %v, want nil", err)
	}

	// Show renders the configuration with the YAML keys the file uses, which is
	// what makes "config show" output something a user can paste back.
	for _, want := range []string{"namespace: acme", "roles_path:", "runbooks_path:"} {
		if !strings.Contains(shown, want) {
			t.Errorf("Show() missing %q:\n%s", want, shown)
		}
	}

	// The persisted file must carry the same keys, not Go field names.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "namespace: acme") {
		t.Errorf("saved file does not use the yaml keys:\n%s", data)
	}
}

func TestSave_RejectsANilConfiguration(t *testing.T) {
	if err := Save(nil); err == nil {
		t.Error("Save(nil) = nil, want an error")
	}
}
