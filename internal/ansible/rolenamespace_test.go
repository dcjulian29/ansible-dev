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
	"strings"
	"testing"
)

func TestRoleNamespace(t *testing.T) {
	tests := []struct {
		role string
		want string
	}{
		{role: "dcjulian29.nginx", want: "dcjulian29"},
		{role: "acme.nginx", want: "acme"},
		{role: "nginx", want: ""},
		{role: "", want: ""},
	}

	for _, tt := range tests {
		if got := RoleNamespace(tt.role); got != tt.want {
			t.Errorf("RoleNamespace(%q) = %q, want %q", tt.role, got, tt.want)
		}
	}
}

// RoleNamespace and BaseRoleName split a qualified name between them; a role
// resolved from one should never keep the other's half.
func TestRoleNamespace_ComplementsBaseRoleName(t *testing.T) {
	const role = "acme.nginx"

	if ns, base := RoleNamespace(role), BaseRoleName(role); ns+"."+base != role {
		t.Errorf("RoleNamespace(%q)+BaseRoleName(%q) = %q, want %q", role, role, ns+"."+base, role)
	}
}

// TestQualifiedRoleName_IsIdempotent pins the behavior that makes the role
// argument accept either form: qualifying an already-qualified name must not
// stack namespaces.
func TestQualifiedRoleName_IsIdempotent(t *testing.T) {
	tests := []struct {
		namespace string
		role      string
		want      string
	}{
		{namespace: "dcjulian29", role: "nginx", want: "dcjulian29.nginx"},
		{namespace: "dcjulian29", role: "dcjulian29.nginx", want: "dcjulian29.nginx"},
		{namespace: "acme", role: "dcjulian29.nginx", want: "acme.nginx"},
	}

	for _, tt := range tests {
		got := QualifiedRoleName(tt.namespace, tt.role)
		if got != tt.want {
			t.Errorf("QualifiedRoleName(%q, %q) = %q, want %q",
				tt.namespace, tt.role, got, tt.want)
		}

		if again := QualifiedRoleName(tt.namespace, got); again != tt.want {
			t.Errorf("QualifiedRoleName applied twice = %q, want %q", again, tt.want)
		}
	}
}

func TestValidateNamespace_Accepts(t *testing.T) {
	for _, ns := range []string{"dcjulian29", "acme", "acme_corp", "acme-corp", "a1"} {
		if err := ValidateNamespace(ns); err != nil {
			t.Errorf("ValidateNamespace(%q) = %v, want nil", ns, err)
		}
	}
}

func TestValidateNamespace_Rejects(t *testing.T) {
	// A "." would produce a three-part role name and a "/" or space would
	// corrupt the repository name, which is the whole reason to validate.
	for _, ns := range []string{"", "   ", "acme.corp", "acme/corp", "acme corp"} {
		err := ValidateNamespace(ns)
		if err == nil {
			t.Errorf("ValidateNamespace(%q) = nil, want an error", ns)
			continue
		}

		if !strings.Contains(err.Error(), "namespace") {
			t.Errorf("ValidateNamespace(%q) error = %q, want it to name the setting", ns, err)
		}
	}
}
