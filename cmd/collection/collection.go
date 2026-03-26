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

// Package collection implements the "ansible-dev collection" command group,
// which provides subcommands for managing Ansible Galaxy collections in the
// development environment. Available subcommands include add, list, purge,
// and remove.
package collection

import (
	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for the "collection"
// command group. The command is also aliased as "collections" for
// convenience.
//
// When invoked without a subcommand it prints the help text. A PreRunE
// hook calls [ansible.EnsureAnsibleDirectory] to verify that the current
// working directory contains an ansible.cfg file before any subcommand
// executes.
//
// The following subcommands are registered:
//   - add:    add a new collection to requirements.yml.
//   - list:   list all collections declared in requirements.yml.
//   - purge:  remove all installed collection artifacts.
//   - remove: remove a collection from requirements.yml.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection",
		Aliases: []string{"collections"},
		Short:   "Provide management of ansible collections in the development environment",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.AddCommand(addCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(purgeCmd())
	cmd.AddCommand(removeCmd())

	return cmd
}
