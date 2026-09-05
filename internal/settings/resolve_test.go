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
	"strings"
	"testing"
)

func TestUnsetError_WithoutLegacyVariablePointsAtSetter(t *testing.T) {
	t.Setenv(LegacyRolesEnv, "")

	err := unsetError("roles_path", "roles-path", LegacyRolesEnv)

	if !strings.Contains(err.Error(), "config roles-path <dir>") {
		t.Errorf("error = %q, want it to name the roles-path setter", err)
	}

	if strings.Contains(err.Error(), "import-env") {
		t.Errorf("error = %q, should not mention import-env when the variable is unset", err)
	}
}

func TestUnsetError_WithLegacyVariablePointsAtImport(t *testing.T) {
	t.Setenv(LegacyRolesEnv, "/path/to/ansible/roles")

	err := unsetError("roles_path", "roles-path", LegacyRolesEnv)

	if !strings.Contains(err.Error(), "config import-env") {
		t.Errorf("error = %q, want it to point at import-env", err)
	}

	if !strings.Contains(err.Error(), LegacyRolesEnv) {
		t.Errorf("error = %q, want it to name %s", err, LegacyRolesEnv)
	}
}
