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
