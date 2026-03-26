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
	"path/filepath"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
)

// NewRole scaffolds a new Ansible role using "ansible-galaxy init". The role
// is created inside the configured roles directory. When verbose is true the
// --verbose flag is forwarded to ansible-galaxy. An error is returned if the
// role folder path cannot be resolved or ansible-galaxy exits with non-zero
// status.
func NewRole(role string, verbose bool) error {
	path, err := RoleFolder(role)
	if err != nil {
		return err
	}

	param := []string{
		"init",
		role,
		"--init-path",
		strings.ReplaceAll(filepath.Dir(path), "\\", "/"),
	}

	if verbose {
		param = append(param, "--verbose")
	}

	if err := execute.ExternalProgram("ansible-galaxy", param...); err != nil {
		return err
	}

	return nil
}
