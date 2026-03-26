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

// Package upgrade implements the "ansible-dev upgrade" command, which
// updates and prunes the Vagrant boxes used by the Ansible development
// environment.
package upgrade

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for "ansible-dev upgrade",
// which updates the Vagrant boxes referenced by the development environment
// to their latest available versions by delegating to [vagrant.Upgrade].
//
// The command has no flags or positional arguments. Under the hood,
// [vagrant.Upgrade] typically runs "vagrant box update" followed by
// "vagrant box prune" to download newer box images and remove outdated
// ones.
//
// Note: this command updates the base box images only — it does not
// recreate or re-provision running VMs. After upgrading boxes, use
// "ansible-dev destroy" followed by "ansible-dev start" to rebuild VMs
// from the updated images.
//
// A PreRunE hook performs two checks before execution:
//  1. [ansible.EnsureAnsibleDirectory] — verifies the current directory
//     is a valid Ansible project. Returns a simplified "not an Ansible
//     development directory" message on failure.
//  2. [vagrant.EnsureVagrantfile] — confirms a Vagrantfile is present
//     so that Vagrant knows which boxes to update.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Updates the boxes that the Ansible development vagrant environment uses",
		RunE: func(_ *cobra.Command, _ []string) error {
			return vagrant.Upgrade()
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
