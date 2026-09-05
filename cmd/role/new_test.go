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
	"runtime"
	"strings"
	"testing"

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

func TestNewCmd_RegistersItsFlags(t *testing.T) {
	cmd := newCmd()

	for _, name := range []string{"namespace", "description", "no-prefix", "force", "verbose", "publish"} {
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
