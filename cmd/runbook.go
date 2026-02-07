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
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	runbookCmd = &cobra.Command{
		Use:   "runbook",
		Args:  cobra.ExactArgs(1),
		Short: "Run ansible runbook in the Ansible development vagrant environment",
		Long:  "Run ansible runbook in the Ansible development vagrant environment",
		Run: func(cmd *cobra.Command, args []string) {
			playFromFlags.Name = args[0]
			playFromFlags.AskBecomePass, _ = cmd.Flags().GetBool("becomepass")
			playFromFlags.AskVaultPass, _ = cmd.Flags().GetBool("vaultpass")
			playFromFlags.FlushCache, _ = cmd.Flags().GetBool("flushcache")
			playFromFlags.Step, _ = cmd.Flags().GetBool("step")
			playFromFlags.Verbose, _ = cmd.Flags().GetBool("verbose")

			execute_runbook(playFromFlags)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			ensureAnsibleDirectory()
			ensureVagrantfile()
		},
	}
)

func init() {
	rootCmd.AddCommand(runbookCmd)

	runbookCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	runbookCmd.Flags().Bool("ask-vault-password", false, "ask for vault password")
	runbookCmd.Flags().Bool("ask-become-password", false, "ask for privilege escalation password")
	runbookCmd.Flags().Bool("flush-cache", false, "clear the fact cache for every host in inventory")
	runbookCmd.Flags().BoolP("step", "s", false, "one-step-at-a-time: confirm each task before running")
}

func execute_runbook(play Play) {
	var param []string

	if play.FlushCache {
		param = append(param, "--flush-cache")
	}

	if play.AskVaultPass {
		param = append(param, "--ask-vault-password")
	}

	if play.Verbose {
		param = append(param, "-v")
	}

	if play.Step {
		param = append(param, "--step")
	}

	param = append(param, "playbooks/runbook.yml")

	executeExternalProgram("ansible-playbook", param...)
}
