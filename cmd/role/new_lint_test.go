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

package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// requireLintTools skips the test when ansible-galaxy, yamllint or ansible-lint
// is missing, except in CI on Linux, where the workflow installs them and a
// missing one means the step broke rather than that the machine lacks Ansible.
// ansible-lint does not run on Windows, so the Windows CI job always skips.
func requireLintTools(t *testing.T) {
	t.Helper()

	for _, tool := range []string{"ansible-galaxy", "yamllint", "ansible-lint"} {
		if _, err := exec.LookPath(tool); err != nil {
			if os.Getenv("CI") != "" && runtime.GOOS != "windows" {
				t.Fatalf("%s is not on PATH; the CI workflow should install it", tool)
			}

			t.Skipf("%s is not on PATH", tool)
		}
	}
}

// lint runs a linter in the role folder and fails the test with its output.
func lint(t *testing.T, dir, name string, args ...string) {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("%s failed on the scaffold: %v\n%s", name, err, out)
	}
}

// A freshly scaffolded role must pass its own lint configuration. The
// templates once shipped "skip_list: [ ]", an unquoted EL version and a
// .yamllint ansible-lint rejects, and ansible-galaxy init left MIT-0 SPDX
// headers that fail yaml[comments]; each was only found by hand after a real
// role had been created. This runs the whole "role new" pipeline, galaxy
// leftovers included, so any of them coming back fails here.
func TestNewCmd_ScaffoldPassesLint(t *testing.T) {
	requireLintTools(t)

	dir := t.TempDir()
	t.Chdir(dir)

	// The relative roles_path matters on Windows: the docker wrappers mount
	// the current folder, so ansible-galaxy cannot see an absolute host path.
	cfg := "[defaults]\nroles_path = ./roles\n"
	if err := os.WriteFile("ansible.cfg", []byte(cfg), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir("roles", 0o750); err != nil {
		t.Fatal(err)
	}

	// The namespace comes from the role argument, since TestMain leaves no
	// configuration to supply one.
	cmd := newCmd()
	cmd.SetArgs([]string{"dcjulian29.linttest", "-d", "exercise the scaffold"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("role new: %v", err)
	}

	role := filepath.Join(dir, "roles", "dcjulian29.linttest")

	lint(t, role, "yamllint", ".")
	// --strict because the .yamllint incompatibility was only a warning, and
	// ansible-lint exits 0 on warnings.
	lint(t, role, "ansible-lint", "--profile", "production", "--strict")
}
