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

// Package role implements the "ansible-dev role" command group, which
// provides subcommands for managing Ansible roles in the development
// environment. Operations include adding, comparing, creating, deleting,
// listing, and removing roles.
package role

import (
	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for the "role" command
// group. The command is also aliased as "roles" for convenience.
//
// When invoked without a subcommand it prints the help text. A PreRunE
// hook calls [ansible.EnsureAnsibleDirectory] to verify that the current
// working directory contains an ansible.cfg file before any subcommand
// executes.
//
// The following subcommands are registered:
//   - add:     add an existing role to requirements.yml.
//   - compare: compare local role files against their upstream source.
//   - delete:  delete a role's directory from the roles path.
//   - list:    list roles declared in requirements.yml or installed on disk.
//   - new:     scaffold a new role via ansible-galaxy init.
//   - remove:  remove a role entry from requirements.yml.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "role",
		Aliases: []string{"roles"},
		Short:   "Provide management of ansible roles in the development environment",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.AddCommand(addCmd())
	cmd.AddCommand(compareCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(newCmd())
	cmd.AddCommand(removeCmd())

	return cmd
}
