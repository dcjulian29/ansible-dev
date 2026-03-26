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

// Package inventory implements the "ansible-dev inventory" (aliased as
// "inv") command, which displays the Ansible inventory for the current
// development environment by delegating to the ansible-inventory CLI tool.
package inventory

import (
	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for
// "ansible-dev inventory". The command is also aliased as "inv".
//
// The command delegates to "ansible-inventory" and supports two display
// modes controlled by the --variables flag:
//
// Graph mode (default, --variables not set):
//
//	Runs "ansible-inventory --graph" to display the inventory as an
//	indented host/group tree.
//
// List mode (--variables set):
//
//	Runs "ansible-inventory --list" to output the full inventory with
//	host variables included. The output format defaults to JSON but can
//	be changed with --toml or --yaml (mutually exclusive).
//
// Flags:
//   - --variables:   include host variables and switch to list mode
//     (default false).
//   - --toml:        output in TOML format instead of JSON; requires
//     --variables (mutually exclusive with --yaml).
//   - --yaml, -y:    output in YAML format instead of JSON; requires
//     --variables (mutually exclusive with --toml).
//
// A PreRunE hook validates the environment with two checks:
//  1. [ansible.EnsureAnsibleDirectory] confirms ansible.cfg is present.
//  2. [vagrant.EnsureVagrantfile] confirms a Vagrantfile is present.
//
// An error is returned if either pre-flight check fails or if
// ansible-inventory exits with a non-zero status.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inventory",
		Aliases: []string{"inv"},
		Short:   "Show inventory information for the Ansible development vagrant environment",
		RunE: func(cmd *cobra.Command, _ []string) error {
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
		PreRunE: func(_ *cobra.Command, _ []string) error {
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
