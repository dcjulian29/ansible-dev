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
	"bufio"
	"fmt"
	"os"
	"strings"
)

// UpdateInventoryAddress rewrites the ansible_host value for the named
// machine in the given INI inventory file. It is an unexported helper
// used by updateHostsIni.
func UpdateInventoryAddress(filename, name, addr string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close() //nolint:errcheck

	var lines []string

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == name {
			updated := []string{fields[0]}
			for _, field := range fields[1:] {
				if strings.HasPrefix(field, "ansible_host=") {
					updated = append(updated, fmt.Sprintf("ansible_host=%s", addr))
				} else {
					updated = append(updated, field)
				}
			}

			line = strings.Join(updated, " ")
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}
