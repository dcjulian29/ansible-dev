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
	"runtime"
	"strings"
	"testing"

	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/spf13/cobra"
)

// TestMain points the home directory at an empty sandbox, so the configuration
// singleton finds no ansible-dev.yml. That is deliberate: it makes the "no
// namespace configured" branch reachable here, which it is not in the settings
// package, where a configuration has to exist for the other tests. It also
// keeps these tests from reading the developer's real configuration, which
// would make the outcome depend on the machine.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ansible-dev-role")
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

// command builds the real "role new" command so the tests exercise the flags as
// registered, not a stand-in.
func command(t *testing.T, args ...string) *cobra.Command {
	t.Helper()

	cmd := newCmd()

	if err := cmd.ParseFlags(args); err != nil {
		t.Fatalf("parsing %v: %v", args, err)
	}

	return cmd
}

func TestResolveNamespace_FlagWins(t *testing.T) {
	// The flag beats a namespace written into the role argument, so a
	// deliberate override is never silently ignored.
	got, err := resolveNamespace(command(t, "--namespace", "acme"), "dcjulian29.nginx")
	if err != nil {
		t.Fatalf("resolveNamespace() error = %v, want nil", err)
	}

	if got != "acme" {
		t.Errorf("resolveNamespace() = %q, want %q", got, "acme")
	}
}

func TestResolveNamespace_ShorthandFlagIsAccepted(t *testing.T) {
	got, err := resolveNamespace(command(t, "-n", "acme"), "nginx")
	if err != nil {
		t.Fatalf("resolveNamespace() error = %v, want nil", err)
	}

	if got != "acme" {
		t.Errorf("resolveNamespace() = %q, want %q", got, "acme")
	}
}

func TestResolveNamespace_FallsBackToTheRoleArgument(t *testing.T) {
	got, err := resolveNamespace(command(t), "acme.nginx")
	if err != nil {
		t.Fatalf("resolveNamespace() error = %v, want nil", err)
	}

	if got != "acme" {
		t.Errorf("resolveNamespace() = %q, want %q", got, "acme")
	}
}

// With no flag, no namespace on the role, and no configuration, the command
// must stop rather than invent one — the role would otherwise be published
// under a namespace nobody chose.
func TestResolveNamespace_ErrorsWhenNoneIsAvailable(t *testing.T) {
	_, err := resolveNamespace(command(t), "nginx")
	if err == nil {
		t.Fatal("resolveNamespace() error = nil, want an error when no namespace is available")
	}

	if !strings.Contains(err.Error(), "config namespace") {
		t.Errorf("error = %q, want it to name the setter", err)
	}

	if !strings.Contains(err.Error(), "--namespace") {
		t.Errorf("error = %q, want it to mention the flag", err)
	}
}

func TestResolveNamespace_RejectsAnInvalidFlagValue(t *testing.T) {
	if _, err := resolveNamespace(command(t, "--namespace", "acme/corp"), "nginx"); err == nil {
		t.Error("resolveNamespace() error = nil, want an error for a namespace containing '/'")
	}
}

// A namespace taken from the role argument is validated too, not just one given
// with the flag.
func TestResolveNamespace_RejectsAnInvalidRoleNamespace(t *testing.T) {
	if _, err := resolveNamespace(command(t), " .nginx"); err == nil {
		t.Error("resolveNamespace() error = nil, want an error for a whitespace namespace")
	}
}

// The cmdargs package tests the rule itself; this pins that "role new" uses it.
// Neither case gets past the Args check, so PreRunE and RunE never run and
// nothing is created.
func TestNewCmd_MissingRoleName(t *testing.T) {
	execute := func(args ...string) (string, error) {
		cmd := newCmd()

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

	if _, err := execute("-d", "install nginx"); err == nil || err.Error() != "expects a role name" {
		t.Errorf("flags only: error = %v, want %q", err, "expects a role name")
	}
}

// An existing published clone used to be found only by PublishRole, after
// --force had already deleted the development copy and scaffolded a new one.
// It must stop the command first, with the existing role untouched.
func TestNewCmd_PublishStopsBeforeForceWhenTheCloneExists(t *testing.T) {
	published := t.TempDir()

	if err := os.Mkdir(filepath.Join(published, "x"), 0o750); err != nil {
		t.Fatal(err)
	}

	cfg, err := settings.Load()
	if err != nil {
		t.Fatal(err)
	}

	// Restored afterwards: the configuration is a process-wide cache, and the
	// namespace tests above depend on it holding nothing.
	original := cfg
	t.Cleanup(func() { settings.Save(&original) }) //nolint:errcheck

	cfg.RolesPath = published
	if err := settings.Save(&cfg); err != nil {
		t.Fatal(err)
	}

	t.Chdir(t.TempDir())

	if err := os.WriteFile("ansible.cfg", []byte("[defaults]\nroles_path = ./roles\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	keep := filepath.Join("roles", "acme.x", "keep.txt")

	if err := os.MkdirAll(filepath.Dir(keep), 0o750); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(keep, []byte("existing role"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := newCmd()
	cmd.SetArgs([]string{"acme.x", "--force", "--publish", "--version", "2.21"})

	err = cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "published role already exists") {
		t.Errorf("role new error = %v, want the existing published role reported", err)
	}

	if _, err := os.Stat(keep); err != nil {
		t.Errorf("--force removed the development role before the check: %v", err)
	}
}

// An existing role without --force used to be refused only after
// "ansible-galaxy --version" had run (a container start on Windows). With PATH
// empty that call would fail, so getting "exists" proves the refusal now comes
// first; "creating role" is no longer printed before it either.
func TestNewCmd_ExistingRoleIsRefusedBeforeAnsibleRuns(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := os.WriteFile("ansible.cfg", []byte("[defaults]\nroles_path = ./roles\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join("roles", "acme.x"), 0o750); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", "")

	cmd := newCmd()
	cmd.SetArgs([]string{"acme.x"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "exists. Use '--force'") {
		t.Errorf("role new error = %v, want the existing role refused before ansible-galaxy runs", err)
	}
}

func TestNewCmd_RegistersItsFlags(t *testing.T) {
	cmd := newCmd()

	for _, name := range []string{"namespace", "description", "no-prefix", "force", "verbose", "publish", "version", "platform"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("flag --%s is not registered", name)
		}
	}

	if f := cmd.Flags().ShorthandLookup("n"); f == nil || f.Name != "namespace" {
		t.Error("-n is not bound to --namespace")
	}

	if f := cmd.Flags().ShorthandLookup("d"); f == nil || f.Name != "description" {
		t.Error("-d is not bound to --description")
	}
}
