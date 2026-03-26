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
	"fmt"
	"os"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// ExecutePlay runs the ansible-playbook command using the temporary playbook
// located at ".tmp/play.yml". Command-line flags are derived from the fields
// of the supplied Play struct. If an "ansible.log" file exists it is removed
// before execution. An error is returned if the log cannot be removed or if
// ansible-playbook exits with a non-zero status.
func ExecutePlay(play Play) error {
	var param []string

	if len(play.Tags) > 0 {
		param = append(param, "--tags")
		param = append(param, strings.Join(play.Tags, ","))
	}

	if play.FlushCache {
		param = append(param, "--flush-cache")
	}

	if play.AskVaultPass {
		param = append(param, "--ask-vault-password")
	}

	if play.Verbose {
		param = append(param, "-v")
	}

	if play.Step {
		param = append(param, "--step")
	}

	param = append(param, ".tmp/play.yml")

	if filesystem.FileExists("ansible.log") {
		err := os.Remove("ansible.log")
		if err != nil {
			return fmt.Errorf("can't remove ansible.log: %v", err)
		}
	}

	if filesystem.FileExists(".tmp/play.yml") {
		err := execute.ExternalProgram("ansible-playbook", param...)
		if err != nil {
			return fmt.Errorf("can't execute playbook: %v", err)
		}
	}

	return nil
}
