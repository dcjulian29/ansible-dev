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

package config

import (
	"fmt"

	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/spf13/cobra"
)

// showCmd creates "ansible-dev config show", which prints the current
// configuration file rendered as YAML.
func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Print the current configuration as YAML",
		RunE: func(_ *cobra.Command, _ []string) error {
			out, err := settings.Show()
			if err != nil {
				return err
			}

			fmt.Print(out)

			return nil
		},
	}
}
