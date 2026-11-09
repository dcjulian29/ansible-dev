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

package tag

import (
	"strings"
	"testing"
)

// Neither case gets past the Args check, so PreRunE and RunE never run and
// ansible-playbook is not called.
func TestNewCommand_ArgumentCheck(t *testing.T) {
	execute := func(args ...string) (string, error) {
		cmd := NewCommand()

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
