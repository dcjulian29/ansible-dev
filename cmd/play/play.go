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
package play

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/spf13/cobra"
)

var playFromFlags ansible.Play

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "play <role>",
		Args:  cobra.ExactArgs(1),
		Short: "Provision ansible role to Ansible development vagrant environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			playFromFlags.Name = args[0]
			playFromFlags.Tags, _ = cmd.Flags().GetStringSlice("tags")
			playFromFlags.AskBecomePass, _ = cmd.Flags().GetBool("becomepass")
			playFromFlags.AskVaultPass, _ = cmd.Flags().GetBool("vaultpass")
			playFromFlags.FlushCache, _ = cmd.Flags().GetBool("flushcache")
			playFromFlags.Step, _ = cmd.Flags().GetBool("step")
			playFromFlags.Verbose, _ = cmd.Flags().GetBool("verbose")

			err := ansible.GenerateRolePlay(playFromFlags.Name)
			if err != nil {
				return err
			}
			err = ansible.ExecutePlay(playFromFlags)
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

	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	cmd.Flags().Bool("ask-vault-password", false, "ask for vault password")
	cmd.Flags().Bool("ask-become-password", false, "ask for privilege escalation password")
	cmd.Flags().Bool("flush-cache", false, "clear the fact cache for every host in inventory")
	cmd.Flags().BoolP("step", "s", false, "one-step-at-a-time: confirm each task before running")
	cmd.Flags().StringSlice("tags", []string{}, "only plays and task tagged with these values")

	return cmd
}
