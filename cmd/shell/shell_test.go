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

package shell

import (
	"strings"
	"testing"
)

// run is not named execute: that would collide with the go-toolbox execute
// package shell.go imports.
func run(args ...string) (string, error) {
	cmd := NewCommand()

	var out strings.Builder

	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)

	err := cmd.Execute()

	return out.String(), err
}

// A bare "shell" and a "shell --" with nothing after it both print the help
// without reaching PreRunE.
func TestNewCommand_BarePrintsHelp(t *testing.T) {
	for _, args := range [][]string{nil, {"--"}} {
		out, err := run(args...)
		if err != nil || !strings.Contains(out, "Usage:") {
			t.Errorf("%q: error = %v, output = %q; want the help and no error", args, err, out)
		}
	}
}

// With a command after "--" the Args check passes and the command carries on
// to PreRunE. The test runs in cmd/shell, which has no ansible.cfg, so that is
// where it stops, before ansible could be called.
func TestNewCommand_CommandPassesThrough(t *testing.T) {
	_, err := run("--", "cat", "/etc/hostname")
	if err == nil || err.Error() != "not an Ansible development directory" {
		t.Errorf("error = %v, want the PreRunE error", err)
	}
}
