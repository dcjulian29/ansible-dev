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
)

// Upgrade updates the Vagrant box for the current environment and then
// removes stale versions that are no longer in use. It performs two
// sequential operations:
//
//  1. Runs "vagrant box update" to download the latest version of the
//     box defined in the Vagrantfile.
//  2. Runs "vagrant box prune --force --keep-active-boxes" to delete
//     older box versions while retaining any version still referenced
//     by an active VM.
//
// An error is returned if either command exits with a non-zero status.
// Execution stops at the first failure, so if the update step fails the
// prune step is skipped.
func Upgrade() error {
	update := []string{"box", "update"}
	prune := []string{"box", "prune", "--force", "--keep-active-boxes"}

	err := execute.ExternalProgram("vagrant", update...)
	if err != nil {
		return err
	}

	err = execute.ExternalProgram("vagrant", prune...)
	if err != nil {
		return err
	}

	return nil
}
