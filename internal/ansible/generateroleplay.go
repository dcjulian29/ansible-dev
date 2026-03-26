/*
Copyright © 2026 Julian Easterling julian@julianscorner.com

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
	"errors"
	"fmt"
	"os"

	"github.com/dcjulian29/go-toolbox/filesystem"
)

// GenerateRolePlay creates a temporary Ansible playbook at ".tmp/play.yml"
// that applies a single role identified by roleName. The generated YAML
// targets the "all" host group. An error is returned if the ".tmp" directory
// cannot be created or the file cannot be written.
func GenerateRolePlay(roleName string) error {
	if !filesystem.DirectoryExists(".tmp") {
		err := os.Mkdir(".tmp", 0755)
		if err != nil {
			return err
		}
	}

	if filesystem.FileExists(".tmp/play.yml") {
		err := os.Remove(".tmp/play.yml")
		if err != nil {
			return err
		}
	}

	file, err := os.Create(".tmp/play.yml")
	if err != nil {
		return errors.New("can't create temporary play file")
	}

	defer file.Close() //nolint

	content := `---
- name: Test Ansible Role
  hosts: all
  any_errors_fatal: true
  become: true

  roles:
`

	content = fmt.Sprintf("%s%s", content, fmt.Sprintf("    - %s\n", roleName))

	if _, err = file.WriteString(content); err != nil {
		return err
	}

	return nil
}
