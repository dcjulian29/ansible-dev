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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"up"},
		Short:   "Starts and potentially provision the Ansible development vagrant environment",
		RunE: func(cmd *cobra.Command, args []string) error {
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
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
