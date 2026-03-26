/*
Copyright © 2026 Julian Easterling <julian@julianscorner.com>

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
	"gopkg.in/ini.v1"
)

// RootRoleFolder returns the base directory path where Ansible roles are
// stored, as configured by the "roles_path" key in the [defaults] section
// of the local ansible.cfg file. An error is returned if ansible.cfg cannot
// be loaded, the [defaults] section is missing, or the "roles_path" key
// is not defined.
func RootRoleFolder() (string, error) {
	cfg, err := ini.Load("ansible.cfg")
	if err != nil {
		return "", err
	}

	section, err := cfg.GetSection("defaults")
	if err != nil {
		return "", err
	}

	path, err := section.GetKey("roles_path")
	if err != nil {
		return "", err
	}

	return path.String(), nil
}
