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

	"github.com/dcjulian29/go-toolbox/color"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/network"
)

// Up starts a named Vagrant VM and waits for it to become reachable on
// the network at the specified address.
//
// It performs the following steps:
//
//  1. Runs "vagrant up <name>" to boot (or resume) the VM.
//  2. Polls addr with ICMP pings until the host responds or the retry
//     limit is reached.
//
// The retry loop attempts up to 20 pings. A progress dot is printed for
// each failed attempt. If the VM responds within the limit, a green
// "[Found]" status is printed; otherwise a red "[NotFound]" status is
// printed and an error is returned.
//
// Parameters:
//   - name: the Vagrant machine name as defined in the Vagrantfile and
//     the [vagrant] section of hosts.ini.
//   - addr: the expected IP address of the VM (e.g. "192.168.57.42").
//
// An error is returned if "vagrant up" fails or the VM does not respond
// to pings within the retry limit.
func Up(name, addr string) error {
	fmt.Printf(color.Yellow("\nBringing '%s' online...\n\n"), name)

	if err := execute.ExternalProgram("vagrant", "up", name); err != nil {
		return err
	}

	fmt.Printf(color.Yellow("\nSearching for '%s' at %s..."), name, addr)

	found := false
	count := 0

	for !found {
		found = network.Ping(addr)

		if found {
			fmt.Println(color.Green(" [Found]"))
		} else {
			if count < 20 {
				fmt.Print(".")
				count++
			} else {
				fmt.Println(color.Red(" [NotFound]"))
				return fmt.Errorf("can't find the '%s' VM at %s", name, addr)
			}
		}
	}

	return nil
}
