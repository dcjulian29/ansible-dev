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

// Package ping implements the "ansible-dev ping" command, which verifies
// that all hosts in the development inventory are reachable and responsive
// to Ansible connections.
package ping

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for "ansible-dev ping".
//
// The command runs "ansible -i hosts.ini -m ping all", which executes the
// Ansible ping module against every host in the local inventory. This is a
// connectivity and authentication check — it confirms that Ansible can SSH
// into each Vagrant VM and receive a "pong" response, rather than
// performing an ICMP ping.
//
// A PreRunE hook validates the environment with two checks:
//  1. [ansible.EnsureAnsibleDirectory] confirms ansible.cfg is present.
//     A custom error message ("not an Ansible development directory") is
//     returned on failure.
//  2. [vagrant.EnsureVagrantfile] confirms a Vagrantfile is present.
//
// An error is returned if either pre-flight check fails or if the ansible
// command exits with a non-zero status (indicating one or more hosts are
// unreachable).
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping the Ansible development vagrant environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := execute.ExternalProgram("ansible", "-i", "hosts.ini", "-m", "ping", "all")
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

	return cmd
}
