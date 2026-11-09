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

package ansible

import (
	"github.com/dcjulian29/ansible-dev/internal/settings"
)

// InstalledRoleName returns the name of role's folder under roles_path.
//
// The name as typed wins when that folder exists, so a role installed under a
// bare name (one never created by "role new") is still found. Otherwise a
// bare name is qualified with the configured namespace, since "role new"
// always creates "<namespace>.<name>": without this, "role new nginx" followed
// by "role delete nginx" looked for roles/nginx and removed nothing.
//
// When neither folder exists, or no namespace is configured, the name is
// returned as typed, so the caller's "not present" error names what the user
// wrote.
func InstalledRoleName(role string) string {
	if RoleFolderExists(role) || RoleNamespace(role) != "" {
		return role
	}

	namespace, err := settings.Namespace()
	if err != nil {
		return role
	}

	if qualified := QualifiedRoleName(namespace, role); RoleFolderExists(qualified) {
		return qualified
	}

	return role
}
