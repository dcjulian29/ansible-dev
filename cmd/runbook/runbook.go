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

// Package runbook implements the "ansible-dev runbook" command, which
// executes the project's fixed runbook playbook against the Vagrant
// development environment.
package runbook

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for "ansible-dev runbook",
// which runs the project's runbook playbook against the Vagrant-managed
// development hosts via [ansible.ExecuteRunbook].
//
// The command constructs an [ansible.Play] with the name "Runbook" and
// populates its execution options from the provided flags. The vault
// password prompt is explicitly disabled (AskVaultPass = false).
//
// Flags:
//   - --verbose, -v:          tell Ansible to print more debug messages
//     (default false). Maps to [ansible.Play.Verbose].
//   - --ask-become-password:  prompt for the privilege escalation (sudo)
//     password at runtime (default false). Maps to
//     [ansible.Play.AskBecomePass].
//   - --flush-cache:          clear the fact cache for every host in the
//     inventory before execution (default false). Maps to
//     [ansible.Play.FlushCache].
//   - --step, -s:             run in one-step-at-a-time mode, confirming
//     each task before it executes (default false). Maps to
//     [ansible.Play.Step].
//
// A PreRunE hook performs two checks before execution:
//  1. [ansible.EnsureAnsibleDirectory] — verifies the current directory
//     is a valid Ansible project. Returns a simplified "not an Ansible
//     development directory" message on failure.
//  2. [vagrant.EnsureVagrantfile] — confirms a Vagrantfile is present,
//     since the runbook targets Vagrant-managed hosts.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runbook",
		Short: "Run ansible runbook in the Ansible development vagrant environment",
		RunE: func(cmd *cobra.Command, _ []string) error {
			play := ansible.Play{
				Name: "Runbook",
			}

			play.AskBecomePass, _ = cmd.Flags().GetBool("becomepass")
			play.AskVaultPass = false
			play.FlushCache, _ = cmd.Flags().GetBool("flushcache")
			play.Step, _ = cmd.Flags().GetBool("step")
			play.Verbose, _ = cmd.Flags().GetBool("verbose")

			if err := ansible.ExecuteRunbook(play); err != nil {
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

	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	cmd.Flags().Bool("ask-become-password", false, "ask for privilege escalation password")
	cmd.Flags().Bool("flush-cache", false, "clear the fact cache for every host in inventory")
	cmd.Flags().BoolP("step", "s", false, "one-step-at-a-time: confirm each task before running")

	return cmd
}
