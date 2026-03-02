/*
Copyright © 2026 Julian Easterling <julian@julianscorner.com>

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
package inventory

import (
	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inventory",
		Aliases: []string{"inv"},
		Short:   "Show inventory information for the Ansible development vagrant environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			param := []string{}
			variables, _ := cmd.Flags().GetBool("variables")

			if variables {
				if r, _ := cmd.Flags().GetBool("toml"); r {
					param = append(param, "--toml")
				}

				if r, _ := cmd.Flags().GetBool("yaml"); r {
					param = append(param, "--yaml")
				}

				param = append(param, "--list")
			} else {
				param = append(param, "--graph")
			}

			err := execute.ExternalProgram("ansible-inventory", param...)
			if err != nil {
				return err
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ansible.EnsureAnsibleDirectory(); err != nil {
				return err
			}

			return vagrant.EnsureVagrantfile()
		},
	}

	cmd.Flags().Bool("variables", false, "include host variables")
	cmd.Flags().Bool("toml", false, "Use TOML format instead of default JSON")
	cmd.Flags().BoolP("yaml", "y", false, "Use YAML format instead of default JSON")

	cmd.MarkFlagsMutuallyExclusive("toml", "yaml")

	return cmd
}
