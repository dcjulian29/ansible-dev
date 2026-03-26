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
	"os"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// ExecuteRunbook runs the fixed runbook playbook located at
// "playbooks/runbook.yml" using ansible-playbook. Unlike [ExecutePlay],
// which operates on the generated ".tmp/play.yml" file, ExecuteRunbook
// targets the project's canonical runbook directly and does not support
// tag-based filtering.
//
// Command-line flags are derived from the supplied [Play] struct:
//
//   - [Play.FlushCache]:    passes --flush-cache to clear the fact cache.
//   - [Play.AskVaultPass]:  passes --ask-vault-password for encrypted vaults.
//   - [Play.Verbose]:       passes -v for increased output verbosity.
//   - [Play.Step]:          passes --step to confirm each task interactively.
//
// The [Play.Tags] and [Play.Name] fields are ignored by this function.
//
// If an "ansible.log" file exists in the current directory it is removed
// before execution to ensure a clean log for the run. An error is returned
// if the log cannot be removed or if ansible-playbook exits with a
// non-zero status.
func ExecuteRunbook(play Play) error {
	var param []string

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

	param = append(param, "playbooks/runbook.yml")

	if filesystem.FileExists("ansible.log") {
		if err := os.Remove("ansible.log"); err != nil {
			return err
		}
	}

	return execute.ExternalProgram("ansible-playbook", param...)
}
