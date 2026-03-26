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

// Package stop implements the "ansible-dev stop" (aliased as "down")
// command, which gracefully shuts down all Vagrant virtual machines in
// the Ansible development environment.
package stop

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for "ansible-dev stop",
// which gracefully halts every Vagrant VM in the development environment
// by delegating to [vagrant.Down]. The command is also aliased as "down"
// for convenience.
//
// The command has no flags or positional arguments. It complements the
// "ansible-dev start" / "up" command: "start" boots and optionally
// provisions the VMs, while "stop" / "down" shuts them down without
// destroying their state. To fully tear down the environment (including
// removing VM disk images), use "ansible-dev destroy" instead.
//
// A PreRunE hook performs two checks before execution:
//  1. [ansible.EnsureAnsibleDirectory] — verifies the current directory
//     is a valid Ansible project. Returns a simplified "not an Ansible
//     development directory" message on failure.
//  2. [vagrant.EnsureVagrantfile] — confirms a Vagrantfile is present
//     so that Vagrant knows which VMs to halt.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Aliases: []string{"down"},
		Short:   "Stops the Ansible development vagrant environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			return vagrant.Down()
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
