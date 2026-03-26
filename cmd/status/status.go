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

// Package status implements the "ansible-dev status" command, which
// displays the current state of the Vagrant virtual machines in the
// Ansible development environment.
package status

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for "ansible-dev status",
// which outputs the current state of the Vagrant-managed virtual machines
// by delegating directly to "vagrant status".
//
// The command has no flags or positional arguments. Its output mirrors
// what "vagrant status" produces — typically a table listing each VM's
// name, provider, and current state (e.g. running, poweroff, not created).
//
// A PreRunE hook performs two checks before execution:
//  1. [ansible.EnsureAnsibleDirectory] — verifies the current directory
//     is a valid Ansible project. Returns a simplified "not an Ansible
//     development directory" message on failure.
//  2. [vagrant.EnsureVagrantfile] — confirms a Vagrantfile is present
//     so that "vagrant status" has a valid configuration to report on.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Output status of the Ansible development vagrant environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			return execute.ExternalProgram("vagrant", "status")
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
