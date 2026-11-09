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

package cmdargs

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// run executes a command shaped like "role new": Required("a role name") as
// its Args check, a PreRunE that always fails (standing in for
// EnsureAnsibleDirectory outside a development folder), and a RunE that
// records whether it ran. It returns what was printed, whether RunE ran, and
// the error.
func run(t *testing.T, args ...string) (string, bool, error) {
	t.Helper()

	ran := false

	cmd := &cobra.Command{
		Use:           "new <role>",
		Args:          Required("a role name"),
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("not an ansible development folder")
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			ran = true
			return nil
		},
	}

	cmd.Flags().StringP("description", "d", "", "")

	var out bytes.Buffer

	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)

	err := cmd.Execute()

	return out.String(), ran, err
}

func TestRequired_BareCommandPrintsHelp(t *testing.T) {
	out, ran, err := run(t)

	// No error means exit 0, and the failing PreRunE was never reached.
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if !strings.Contains(out, "Usage:") {
		t.Errorf("output = %q, want the help", out)
	}

	if ran {
		t.Error("RunE ran, want it skipped when the help is shown")
	}
}

func TestRequired_FlagsWithoutArgumentNamesWhatIsMissing(t *testing.T) {
	_, ran, err := run(t, "-d", "something")

	if err == nil || err.Error() != "expects a role name" {
		t.Fatalf("Execute() error = %v, want %q", err, "expects a role name")
	}

	if ran {
		t.Error("RunE ran, want it skipped")
	}
}

func TestRequired_TooManyArgumentsSaysHowMany(t *testing.T) {
	_, _, err := run(t, "one", "two")

	if err == nil || err.Error() != "expects a role name, received 2 arguments" {
		t.Fatalf("Execute() error = %v, want the count in the message", err)
	}
}

// With exactly one argument Required passes, so the command carries on to
// PreRunE; its failure here shows the check let the command through.
func TestRequired_OneArgumentPassesThrough(t *testing.T) {
	_, _, err := run(t, "nginx")

	if err == nil || err.Error() != "not an ansible development folder" {
		t.Fatalf("Execute() error = %v, want the PreRunE error", err)
	}
}

// TestAtLeastOne builds a command shaped like "shell": AtLeastOne as its Args
// check and a RunE that records the words it received. The failing PreRunE of
// the Required tests is left out; Args runs before it either way.
func TestAtLeastOne(t *testing.T) {
	execute := func(args ...string) (string, []string, error) {
		var got []string

		cmd := &cobra.Command{
			Use:           "shell -- <command>",
			Args:          AtLeastOne("a shell command after --"),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(_ *cobra.Command, args []string) error {
				got = args
				return nil
			},
		}

		cmd.Flags().Bool("verbose", false, "")

		var out bytes.Buffer

		cmd.SetOut(&out)
		cmd.SetErr(&out)
		cmd.SetArgs(args)

		err := cmd.Execute()

		return out.String(), got, err
	}

	for _, bare := range [][]string{nil, {"--"}} {
		out, got, err := execute(bare...)
		if err != nil || !strings.Contains(out, "Usage:") || got != nil {
			t.Errorf("%q: error = %v, ran with %q; want the help and no error", bare, err, got)
		}
	}

	if _, _, err := execute("--verbose"); err == nil || err.Error() != "expects a shell command after --" {
		t.Errorf("flags only: error = %v, want the missing command named", err)
	}

	_, got, err := execute("--", "cat", "/etc/hostname")
	if err != nil || strings.Join(got, " ") != "cat /etc/hostname" {
		t.Errorf("words: error = %v, args = %q; want every word passed through", err, got)
	}
}
