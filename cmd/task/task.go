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

// Package task implements the "ansible-dev tasks" command, which lists
// all tasks that would be executed for a given Ansible role, optionally
// filtered by tags.
package task

import (
	"errors"
	"strings"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for "ansible-dev tasks",
// which lists every task that would be executed for the specified Ansible
// role. The command requires exactly one positional argument.
//
// Usage:
//
//	ansible-dev tasks <role> [flags]
//
// The positional argument <role> is the name of the role to inspect (e.g.
// "dcjulian29.docker"). The command proceeds in two steps:
//
//  1. Generate a temporary playbook — calls [ansible.GenerateRolePlay]
//     to create a minimal playbook at .tmp/play.yml that includes only
//     the specified role.
//
//  2. List tasks — invokes "ansible-playbook --list-tasks .tmp/play.yml"
//     to parse the generated playbook and print all tasks it would
//     execute. When the --tags flag is provided, the "--tags" argument
//     is prepended so that only tasks matching the specified tags are
//     shown.
//
// Flags:
//   - --tags: a comma-separated list of tag values used to filter the
//     output to only plays and tasks tagged with those values
//     (default empty). May also be specified multiple times.
//
// This command complements "ansible-dev tags", which lists available tag
// names. Use "tags" to discover tags, then "tasks --tags <tag>" to
// preview which tasks a filtered run would execute.
//
// A PreRunE hook performs two checks before execution:
//  1. [ansible.EnsureAnsibleDirectory] — verifies the current directory
//     is a valid Ansible project. Returns a simplified "not an Ansible
//     development directory" message on failure.
//  2. [vagrant.EnsureVagrantfile] — confirms a Vagrantfile is present.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks <role>",
		Short: "List all tasks that would be executed in the role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			err := ansible.GenerateRolePlay(args[0])
			if err != nil {
				return err
			}

			var param []string

			tags, _ := cmd.Flags().GetStringSlice("tags")

			if len(tags) > 0 {
				param = append(param, "--tags")
				param = append(param, strings.Join(tags, ","))
			}

			param = append(param, "--list-tasks", ".tmp/play.yml")

			err = execute.ExternalProgram("ansible-playbook", param...)
			if err != nil {
				return err
			}

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if err := ansible.EnsureAnsibleDirectory(); err != nil {
				return errors.New("not an Ansible development directory")
			}

			return vagrant.EnsureVagrantfile()
		},
	}

	cmd.Flags().StringSlice("tags", []string{}, "only plays and task tagged with these values")

	return cmd
}
