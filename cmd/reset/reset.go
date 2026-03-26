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

// Package reset implements the "ansible-dev reset" command, which
// destroys and recreates the entire Vagrant development environment from
// scratch, optionally provisioning one or more Ansible roles afterward.
package reset

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

// NewCommand creates and returns the Cobra command for "ansible-dev reset".
//
// The command performs a full environment reset in three phases:
//
//  1. Destroy – calls [vagrant.Destroy] to force-destroy all VMs and
//     remove local artifacts (ansible.log, .vagrant, .tmp).
//  2. Recreate – reads the inventory via [ansible.GetInventory] and calls
//     [vagrant.Up] for each host, booting the VM and waiting for network
//     reachability.
//  3. Provision (optional) – if one or more --role flags were provided,
//     calls [ansible.ApplyRoles] to generate and execute a temporary
//     playbook for each role.
//
// Flags:
//   - --role:       one or more Ansible roles to apply after the VMs are
//     brought online. May be specified multiple times or as a
//     comma-separated list (default: none).
//   - --tag:        one or more Ansible tags to filter task execution
//     within the applied roles. May be specified multiple times or as a
//     comma-separated list (default: none).
//   - --verbose, -v: enable verbose ansible-playbook output
//     (default false).
//
// A PreRunE hook validates the environment with two checks:
//  1. [ansible.EnsureAnsibleDirectory] confirms ansible.cfg is present.
//     A custom error message ("not an Ansible development directory") is
//     returned on failure.
//  2. [vagrant.EnsureVagrantfile] confirms a Vagrantfile is present.
//
// Execution is fail-fast: an error at any phase stops the command and
// returns immediately.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset the Ansible development vagrant environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := vagrant.Destroy(); err != nil {
				return err
			}

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

	cmd.Flags().StringSliceVar(&roles, "role", []string{}, "provision the VMs with the specified role(s)")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "apply the role with the specified tag(s)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "tell Ansible to print more debug messages")

	return cmd
}
