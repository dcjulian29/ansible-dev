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

// Play describes the parameters for a single ansible-playbook invocation.
//
// Fields:
//   - Name:          human-readable identifier for the play (typically a role or playbook name).
//   - Tags:          optional Ansible tags used to limit which tasks are executed.
//   - AskVaultPass:  when true, the --ask-vault-password flag is passed to ansible-playbook.
//   - AskBecomePass: when true, the --ask-become-pass flag is passed to ansible-playbook.
//   - FlushCache:    when true, the --flush-cache flag is passed to clear the fact cache.
//   - Step:          when true, the --step flag is passed so each task must be confirmed.
//   - Verbose:       when true, the -v flag is passed for increased output verbosity.
type Play struct {
	Name          string
	Tags          []string
	AskVaultPass  bool
	AskBecomePass bool
	FlushCache    bool
	Step          bool
	Verbose       bool
}
