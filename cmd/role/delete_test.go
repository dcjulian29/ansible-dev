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
	"path/filepath"
	"strings"
	"testing"

	"github.com/dcjulian29/ansible-dev/internal/settings"
)

// Neither case gets past the Args check, so PreRunE and RunE never run and
// nothing is deleted. Before the check moved into Args, "delete a b" removed a
// and silently ignored b.
func TestDeleteCmd_ArgumentCheck(t *testing.T) {
	execute := func(args ...string) (string, error) {
		cmd := deleteCmd()

		var out strings.Builder

		cmd.SetOut(&out)
		cmd.SetErr(&out)
		cmd.SetArgs(args)

		err := cmd.Execute()

		return out.String(), err
	}

	out, err := execute()
	if err != nil || !strings.Contains(out, "Usage:") {
		t.Errorf("bare: error = %v, output = %q; want the help and no error", err, out)
	}

	want := "expects a role name, received 2 arguments"
	if _, err := execute("a", "b"); err == nil || err.Error() != want {
		t.Errorf("two names: error = %v, want %q", err, want)
	}
}

// "role new x" creates roles/<namespace>.x, and "role delete x" used to look
// for roles/x, find nothing and exit 0. A role installed under a bare name
// must still be found as typed, and a missing one must be an error.
func TestDeleteCmd_ResolvesTheFolder(t *testing.T) {
	cfg, err := settings.Load()
	if err != nil {
		t.Fatal(err)
	}

	// Restored afterwards: the configuration is a process-wide cache, and the
	// namespace tests depend on it holding nothing.
	original := cfg
	t.Cleanup(func() { settings.Save(&original) }) //nolint:errcheck

	cfg.Namespace = "acme"
	if err := settings.Save(&cfg); err != nil {
		t.Fatal(err)
	}

	t.Chdir(t.TempDir())

	if err := os.WriteFile("ansible.cfg", []byte("[defaults]\nroles_path = ./roles\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, folder := range []string{"acme.x", "legacy"} {
		if err := os.MkdirAll(filepath.Join("roles", folder), 0o750); err != nil {
			t.Fatal(err)
		}
	}

	run := func(role string) error {
		cmd := deleteCmd()
		cmd.SetArgs([]string{role})

		return cmd.Execute()
	}

	for role, folder := range map[string]string{"x": "acme.x", "legacy": "legacy"} {
		if err := run(role); err != nil {
			t.Errorf("delete %s: %v", role, err)
		}

		if _, err := os.Stat(filepath.Join("roles", folder)); !os.IsNotExist(err) {
			t.Errorf("delete %s: roles/%s still exists (stat error = %v)", role, folder, err)
		}
	}

	if err := run("missing"); err == nil || !strings.Contains(err.Error(), "not present") {
		t.Errorf("delete missing: error = %v, want a not-present error", err)
	}
}
