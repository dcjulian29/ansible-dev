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

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/network"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// Up starts a named Vagrant VM, discovers its IP address via
// "vagrant ssh-config", updates hosts.ini with the discovered address,
// and waits for the VM to become reachable from the host.
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
// An error is returned if "vagrant up" fails, the IP cannot be
// discovered via ssh-config, hosts.ini cannot be updated, or the VM
// does not respond to pings within the retry limit.
func Up(name string) error {
	fmt.Printf(textformat.Yellow("\nBringing '%s' online...\n\n"), name)

	if err := execute.ExternalProgram("vagrant", "up", name); err != nil {
		return err
	}

	addr, err := getSSHHost(name)
	if err != nil {
		return err
	}

	if err := updateHostsIni(name, addr); err != nil {
		return fmt.Errorf("failed to update hosts.ini for %s: %w", name, err)
	}

	fmt.Printf(textformat.Yellow("\nSearching for '%s' at %s..."), name, addr)

	found := false
	count := 0

	for !found {
		found = network.Ping(addr)

		if found {
			fmt.Println(textformat.Green(" [Found]"))
		} else {
			if count < 20 {
				fmt.Print(".")
				count++
			} else {
				fmt.Println(textformat.Red(" [NotFound]"))
				return fmt.Errorf("can't find the '%s' VM at %s", name, addr)
			}
		}
	}

	return nil
}

func getSSHHost(name string) (string, error) {
	out, err := execute.ExternalProgramCapture("vagrant", "ssh-config", name)
	if err != nil {
		return "", fmt.Errorf("vagrant ssh-config %s: %w", name, err)
	}

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "HostName ") {
			return strings.TrimPrefix(line, "HostName "), nil
		}
	}

	return "", fmt.Errorf("HostName not found in ssh-config output for '%s'", name)
}

func updateHostsIni(name, addr string) error {
	return UpdateInventoryAddress("hosts.ini", name, addr)
}
