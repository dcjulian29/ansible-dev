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

package vagrant

import (
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// Destroy tears down the Vagrant environment and removes all local
// artifacts produced during a development session. It performs the
// following steps in order:
//
//  1. Runs "vagrant destroy --force" to stop and delete all managed VMs.
//  2. Removes the "ansible.log" file created by ansible-playbook runs.
//  3. Removes the ".vagrant" directory that stores Vagrant machine state.
//  4. Removes the ".tmp" directory used for generated playbook files
//     (see [ansible.GenerateRolePlay] and [ansible.GeneratePlaybookPlay]).
//
// Execution stops and the first encountered error is returned if any step
// fails. A nil return indicates the environment was fully cleaned up.
func Destroy() error {
	param := []string{"destroy", "--force"}

	if err := execute.ExternalProgram("vagrant", param...); err != nil {
		return err
	}

	if err := filesystem.RemoveFile("ansible.log"); err != nil {
		return err
	}

	if err := filesystem.RemoveDirectory(".vagrant"); err != nil {
		return err
	}

	if err := filesystem.RemoveDirectory(".tmp"); err != nil {
		return err
	}

	return nil
}
