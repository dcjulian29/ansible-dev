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

package config

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/spf13/cobra"
)

// TestMain redirects the home directory so these tests read and write a
// throwaway configuration instead of the developer's own. The commands here
// call settings.Save, so without this the suite would rewrite a real file.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ansible-dev-config")
	if err != nil {
		panic(err)
	}

	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", dir) //nolint:errcheck
	} else {
		os.Setenv("HOME", dir) //nolint:errcheck
	}

	code := m.Run()

	os.RemoveAll(dir) //nolint:errcheck
	os.Exit(code)
}

// run executes a command with the given arguments, capturing anything it prints
// so the suite's output stays readable.
func run(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()

	original := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	done := make(chan string, 1)

	go func() {
		var sb strings.Builder

		buf := make([]byte, 4096)

		for {
			n, readErr := r.Read(buf)
			if n > 0 {
				sb.Write(buf[:n])
			}

			if readErr != nil {
				break
			}
		}

		done <- sb.String()
	}()

	cmd.SetArgs(args)
	cmd.SetOut(w)
	cmd.SetErr(w)

	runErr := cmd.Execute()

	w.Close() //nolint:errcheck

	os.Stdout = original

	out := <-done

	r.Close() //nolint:errcheck

	return out, runErr
}

func TestNewCommand_RegistersEverySubcommand(t *testing.T) {
	cmd := NewCommand()

	registered := map[string]bool{}
	for _, c := range cmd.Commands() {
		registered[c.Name()] = true
	}

	want := []string{
		"show", "path", "import-env",
		"namespace", "roles-path", "runbooks-path",
		"role-ignore", "runbook-ignore",
		"diff-program", "diff-role-filter", "diff-runbook-filter", "diff-args",
	}

	for _, name := range want {
		if !registered[name] {
			t.Errorf("subcommand %q is not registered", name)
		}
	}
}

func TestIgnoreCmd_RegistersItsOperations(t *testing.T) {
	cmd := ignoreCmd("role-ignore", "",
		func(c *settings.Config) []string { return c.RoleIgnore },
		func(c *settings.Config, v []string) { c.RoleIgnore = v },
	)

	registered := map[string]bool{}
	for _, c := range cmd.Commands() {
		registered[c.Name()] = true
	}

	for _, name := range []string{"list", "add", "remove", "clear"} {
		if !registered[name] {
			t.Errorf("ignore subcommand %q is not registered", name)
		}
	}
}

// roleIgnore builds a fresh command group over the role ignore list. A new one
// is built per step because Cobra remembers the args of a previous Execute.
func roleIgnore() *cobra.Command {
	return ignoreCmd("role-ignore", "",
		func(c *settings.Config) []string { return c.RoleIgnore },
		func(c *settings.Config, v []string) { c.RoleIgnore = v },
	)
}

func currentRoleIgnore(t *testing.T) []string {
	t.Helper()

	cfg, err := settings.Load()
	if err != nil {
		t.Fatal(err)
	}

	return cfg.RoleIgnore
}

// TestIgnoreList_AddRemoveClear walks the whole lifecycle in one test. The
// commands share the configuration singleton, so a sequence in a single test is
// deterministic where separate tests would depend on execution order.
func TestIgnoreList_AddRemoveClear(t *testing.T) {
	if _, err := run(t, roleIgnore(), "add", ".git"); err != nil {
		t.Fatalf("add: %v", err)
	}

	if _, err := run(t, roleIgnore(), "add", ".github"); err != nil {
		t.Fatalf("add: %v", err)
	}

	if got := currentRoleIgnore(t); len(got) != 2 {
		t.Fatalf("ignore list = %v, want two entries", got)
	}

	// A duplicate is reported and changes nothing, rather than growing the list
	// with a second copy that would never match differently.
	out, err := run(t, roleIgnore(), "add", ".git")
	if err != nil {
		t.Fatalf("duplicate add: %v", err)
	}

	if !strings.Contains(out, "already in the ignore list") {
		t.Errorf("duplicate add output = %q, want a warning", out)
	}

	if got := currentRoleIgnore(t); len(got) != 2 {
		t.Errorf("ignore list = %v, want the duplicate to be rejected", got)
	}

	// Removing something absent is likewise reported, not silently accepted.
	out, err = run(t, roleIgnore(), "remove", ".vscode")
	if err != nil {
		t.Fatalf("remove missing: %v", err)
	}

	if !strings.Contains(out, "not in the ignore list") {
		t.Errorf("remove-missing output = %q, want a warning", out)
	}

	if got := currentRoleIgnore(t); len(got) != 2 {
		t.Errorf("ignore list = %v, want it unchanged", got)
	}

	if _, err := run(t, roleIgnore(), "remove", ".git"); err != nil {
		t.Fatalf("remove: %v", err)
	}

	got := currentRoleIgnore(t)
	if len(got) != 1 || got[0] != ".github" {
		t.Errorf("ignore list = %v, want [.github]", got)
	}

	out, err = run(t, roleIgnore(), "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if !strings.Contains(out, ".github") || strings.Contains(out, ".git\n.git") {
		t.Errorf("list output = %q, want just the remaining entry", out)
	}

	if _, err := run(t, roleIgnore(), "clear"); err != nil {
		t.Fatalf("clear: %v", err)
	}

	if got := currentRoleIgnore(t); len(got) != 0 {
		t.Errorf("ignore list = %v, want it empty after clear", got)
	}
}

func TestNamespaceCmd_RejectsAnInvalidNamespace(t *testing.T) {
	if _, err := run(t, namespaceCmd(), "acme.corp"); err == nil {
		t.Error("namespace command error = nil, want an error for a value containing '.'")
	}
}

// A rejected value must not reach the file, or a later read would fail on
// something the setter refused.
func TestNamespaceCmd_SetsAndShows(t *testing.T) {
	if _, err := run(t, namespaceCmd(), "acme"); err != nil {
		t.Fatalf("setting the namespace: %v", err)
	}

	if _, err := run(t, namespaceCmd(), "not/valid"); err == nil {
		t.Fatal("expected the invalid value to be rejected")
	}

	out, err := run(t, namespaceCmd())
	if err != nil {
		t.Fatalf("showing the namespace: %v", err)
	}

	if !strings.Contains(out, "acme") || strings.Contains(out, "not/valid") {
		t.Errorf("namespace = %q, want the rejected value not to have been stored", out)
	}
}
