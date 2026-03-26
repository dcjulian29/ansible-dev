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
	"errors"
	"strings"

	"gopkg.in/ini.v1"
)

// GetInventory reads the "hosts.ini" file in the current directory, parses
// the [vagrant] INI section, and returns a slice of Inventory entries—one
// per host defined in that section. An error is returned if the file cannot
// be loaded or the [vagrant] section is missing.
func GetInventory() ([]Inventory, error) {
	inv, err := ini.Load("hosts.ini")
	if err != nil {
		return []Inventory{}, err
	}

	section, err := inv.GetSection("vagrant")
	if err != nil {
		return []Inventory{}, errors.New("can't find the 'vagrant' section in the hosts.ini file")
	}

	inventory := []Inventory{}

	for _, vm := range section.KeyStrings() {
		i := Inventory{
			Name:    strings.Split(vm, " ")[0],
			Address: section.Key(vm).String(),
		}

		inventory = append(inventory, i)
	}

	return inventory, nil
}
