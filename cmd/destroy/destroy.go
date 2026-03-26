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

// Package destroy implements the "ansible-dev destroy" command, which tears
// down the Vagrant-based virtual machine environment and removes all local
// development artifacts.
package destroy

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for "ansible-dev destroy".
//
// The command delegates to [vagrant.Destroy], which force-destroys all
// managed VMs and removes the ansible.log, .vagrant, and .tmp artifacts
// from the current directory.
//
// A PreRunE hook performs two validations before execution:
//  1. Calls [ansible.EnsureAnsibleDirectory] to confirm the current
//     directory contains an ansible.cfg file.
//  2. Calls [vagrant.EnsureVagrantfile] to confirm a Vagrantfile is
//     present.
//
// An error is returned if either pre-flight check fails or if the
// destroy operation itself encounters a failure.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the Ansible development vagrant environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			return vagrant.Destroy()
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
