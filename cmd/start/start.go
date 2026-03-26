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

// Package start implements the "ansible-dev start" (aliased as "up")
// command, which boots the Vagrant development environment and
// optionally provisions it with Ansible roles.
package start

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/spf13/cobra"
)

var (
	roles   []string
	tags    []string
	verbose bool
)

// NewCommand creates and returns the Cobra command for "ansible-dev start",
// which boots every host in the project's inventory and optionally
// provisions them with the specified Ansible roles. The command is also
// aliased as "up" for convenience.
//
// The command proceeds in two phases:
//
//  1. Start VMs — reads the inventory via [ansible.GetInventory] and
//     calls [vagrant.Up] for each host, passing the host name and IP
//     address. Each VM is started sequentially; the command aborts on
//     the first failure.
//
//  2. Provision (optional) — calls [ansible.ApplyRoles] with the
//     collected roles, tags, and verbose flag. When no roles are
//     specified (empty --role), ApplyRoles controls whether
//     provisioning is skipped or a default set is applied.
//
// Flags:
//   - --role, -r:    one or more Ansible roles to provision the VMs
//     with after startup. Accepts repeated flags or comma-separated
//     values (default empty).
//   - --tag, -t:     one or more Ansible tags to filter tasks within the
//     provisioned roles. Accepts repeated flags or comma-separated
//     values (default empty).
//   - --verbose, -v: tell Ansible to print more debug messages during
//     provisioning (default false).
//
// A PreRunE hook performs two checks before execution:
//  1. [ansible.EnsureAnsibleDirectory] — verifies the current directory
//     is a valid Ansible project. Returns a simplified "not an Ansible
//     development directory" message on failure.
//  2. [vagrant.EnsureVagrantfile] — confirms a Vagrantfile is present.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"up"},
		Short:   "Starts and potentially provision the Ansible development vagrant environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			inventory, err := ansible.GetInventory()
			if err != nil {
				return err
			}

			for _, host := range inventory {
				if err := vagrant.Up(host.Name, host.Address); err != nil {
					return err
				}
			}

			err = ansible.ApplyRoles(roles, tags, verbose)
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

	cmd.Flags().StringSliceVarP(&roles, "role", "r", []string{}, "provision the VMs with the specified role(s)")
	cmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "apply the role(s) with the specified tag(s)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "tell Ansible to print more debug messages")

	return cmd
}
