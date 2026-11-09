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
	"strings"
	"testing"
)

// Neither case gets past the Args check, so PreRunE and RunE never run and
// requirements.yml is not touched.
func TestAddCmd_MissingRoleName(t *testing.T) {
	execute := func(args ...string) (string, error) {
		cmd := addCmd()

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

	if _, err := execute("-s", "https://example.com/role.git"); err == nil || err.Error() != "expects a role name" {
		t.Errorf("flags only: error = %v, want %q", err, "expects a role name")
	}
}
