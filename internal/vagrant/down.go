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
	"fmt"
	"strings"

	"github.com/dcjulian29/go-toolbox/color"
	"github.com/dcjulian29/go-toolbox/execute"
	"gopkg.in/ini.v1"
)

// Down gracefully stops every Vagrant VM listed in the [vagrant] section of
// the local hosts.ini inventory file by running "vagrant halt <name>" for
// each entry.
//
// Host names are extracted from the INI keys by taking the first
// whitespace-delimited token (matching the convention used by
// [ansible.GetInventory]). A yellow status line is printed to stdout before
// each VM is halted.
//
// An error is returned if hosts.ini cannot be loaded, the [vagrant] section
// is missing, or any individual "vagrant halt" invocation fails. Execution
// stops at the first failure, leaving remaining VMs in their current state.
func Down() error {
	inv, err := ini.Load("hosts.ini")
	if err != nil {
		return err
	}

	section, err := inv.GetSection("vagrant")
	if err != nil {
		return err
	}

	for _, vm := range section.KeyStrings() {
		name := strings.Split(vm, " ")[0]

		fmt.Printf(color.Yellow("\nStopping '%s'...\n\n"), name)

		err := execute.ExternalProgram("vagrant", "halt", name)
		if err != nil {
			return err
		}
	}

	return nil
}
