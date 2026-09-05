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
	"fmt"
	"strings"
)

// ValidateNamespace rejects a namespace that cannot be composed into a role
// name or a repository name.
//
// Only the characters that actually break composition are rejected: a "." would
// produce a three-part role name that Ansible cannot resolve, and a "/" or a
// space would corrupt the "<namespace>/ansible-role-<name>" repository. Galaxy
// itself prefers lowercase letters and underscores, but that is a publishing
// convention rather than something this tool needs to enforce, and rejecting a
// name GitHub would have accepted is worse than passing it through.
func ValidateNamespace(namespace string) error {
	if strings.TrimSpace(namespace) == "" {
		return fmt.Errorf("namespace cannot be empty")
	}

	for _, r := range namespace {
		valid := r == '_' || r == '-' ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')

		if !valid {
			return fmt.Errorf(
				"namespace '%s' contains '%c'; only letters, digits, '_' and '-' are allowed",
				namespace, r)
		}
	}

	return nil
}

// RoleNamespace returns the namespace portion of a fully-qualified role name,
// so "dcjulian29.nginx" yields "dcjulian29". It is the other half of
// [BaseRoleName] and returns an empty string when the name carries no
// namespace, which is how a caller distinguishes "the user named a namespace"
// from "fall back to the configured one".
func RoleNamespace(role string) string {
	if i := strings.IndexByte(role, '.'); i >= 0 {
		return role[:i]
	}

	return ""
}

// QualifiedRoleName joins a namespace and a role into the Galaxy form
// "<namespace>.<name>".
//
// Any namespace already on the role is replaced rather than stacked, so the
// result is the same whether the caller wrote "nginx" or "dcjulian29.nginx" —
// which is what lets the role argument stay in whichever form the user prefers
// without producing "dcjulian29.dcjulian29.nginx".
func QualifiedRoleName(namespace, role string) string {
	return namespace + "." + BaseRoleName(role)
}
