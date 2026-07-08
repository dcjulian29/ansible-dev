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
	"os"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/ansible-dev/internal/templates"
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// NewRunbook renders the embedded runbook scaffolding into a new directory
// named by the ANSIBLE_RUNBOOKS environment variable joined with name,
// substituting !!RUNBOOK_NAME!! and !!RUNBOOK_DESC!!. It returns the absolute
// path of the created directory.
//
// Unlike a role, a runbook has no ansible-galaxy skeleton: the embedded
// template is the entire scaffold. An error is returned if ANSIBLE_RUNBOOKS is
// unset or the destination already exists.
func NewRunbook(name, description string) (string, error) {
	runbooks := os.Getenv("ANSIBLE_RUNBOOKS")
	if len(runbooks) == 0 {
		return "", fmt.Errorf("the 'ANSIBLE_RUNBOOKS' environment variable is not defined")
	}

	dest := filepath.Join(strings.ReplaceAll(runbooks, "\\", string(os.PathSeparator)), name)

	if filesystem.DirectoryExist(dest) {
		return "", fmt.Errorf("runbook already exists at '%s'", dest)
	}

	src, err := templates.Runbook()
	if err != nil {
		return "", err
	}

	err = ApplyTemplate(src, dest, map[string]string{
		"!!RUNBOOK_NAME!!": name,
		"!!RUNBOOK_DESC!!": description,
	})
	if err != nil {
		return "", err
	}

	return dest, nil
}
