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

import "testing"

func TestParseAnsibleVersion(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"first line only", "ansible-galaxy [core 2.21.5]", "2.21.5"},
		{"pre-release keeps the numeric part", "ansible-galaxy [core 2.22.0rc1]", "2.22.0"},
		{
			// The real command prints the config file, module paths and more
			// after the version line; only the first [core ...] counts.
			"full output",
			"ansible-galaxy [core 2.21.5]\n  config file = /etc/ansible/ansible.cfg\n" +
				"  python version = 3.14.0\n  jinja version = 3.1.6",
			"2.21.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAnsibleVersion(tt.output)
			if err != nil {
				t.Fatalf("parseAnsibleVersion() error = %v, want nil", err)
			}

			if got != tt.want {
				t.Errorf("parseAnsibleVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Output without a version must be an error, never an empty string: an empty
// min_ansible_version would be written into meta/main.yml and fail its lint.
func TestParseAnsibleVersion_ErrorsWithoutAVersion(t *testing.T) {
	for _, output := range []string{"", "command not found", "ansible-galaxy [core]"} {
		if _, err := parseAnsibleVersion(output); err == nil {
			t.Errorf("parseAnsibleVersion(%q) error = nil, want an error", output)
		}
	}
}

func TestValidateAnsibleVersion(t *testing.T) {
	for _, ok := range []string{"2.21", "2.21.5", "2.18.4"} {
		if err := ValidateAnsibleVersion(ok); err != nil {
			t.Errorf("ValidateAnsibleVersion(%q) error = %v, want nil", ok, err)
		}
	}

	for _, bad := range []string{"", "2", "v2.21", "2.21.x", "2.21.5rc1", " 2.21"} {
		if err := ValidateAnsibleVersion(bad); err == nil {
			t.Errorf("ValidateAnsibleVersion(%q) error = nil, want an error", bad)
		}
	}
}
