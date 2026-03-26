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

// Package tag implements the "ansible-dev tags" command, which lists all
// available Ansible tags defined within a given role.
package tag

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for "ansible-dev tags",
// which enumerates every tag defined in the specified Ansible role.
//
// Usage:
//
//	ansible-dev tags <role>
//
// The positional argument <role> is the name of the role to inspect (e.g.
// "dcjulian29.docker"). The command proceeds in two steps:
//
//  1. Generate a temporary playbook — calls [ansible.GenerateRolePlay]
//     to create a minimal playbook at .tmp/play.yml that includes only
//     the specified role.
//
//  2. List tags — invokes "ansible-playbook --list-tags .tmp/play.yml"
//     to parse the generated playbook and print all tags it contains.
//     If the playbook execution fails, a descriptive error is returned.
//
// If no argument is supplied, the help text is displayed instead.
//
// This command is useful for discovering available tags before running
// "ansible-dev play" or "ansible-dev start" with the --tag flag.
//
// A PreRunE hook performs two checks before execution:
//  1. [ansible.EnsureAnsibleDirectory] — verifies the current directory
//     is a valid Ansible project. Returns a simplified "not an Ansible
//     development directory" message on failure.
//  2. [vagrant.EnsureVagrantfile] — confirms a Vagrantfile is present.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags <role>",
		Short: "List all available tags in the role",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			err := ansible.GenerateRolePlay(args[0])
			if err != nil {
				return err
			}

			err = execute.ExternalProgram("ansible-playbook", "--list-tags", ".tmp/play.yml")
			if err != nil {
				return errors.New("can't execute playbook to list tags")
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

	return cmd
}
