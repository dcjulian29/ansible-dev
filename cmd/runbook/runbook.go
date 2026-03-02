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
package runbook

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runbook",
		Short: "Run ansible runbook in the Ansible development vagrant environment",
		RunE: func(cmd *cobra.Command, args []string) error {
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
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
