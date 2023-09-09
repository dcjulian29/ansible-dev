/*
Copyright © 2023 Julian Easterling <julian@julianscorner.com>

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
	pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Ping the Ansible development vagrant environment",
		Long:  "Ping the Ansible development vagrant environment",
		Run: func(cmd *cobra.Command, args []string) {
			if len(shellCommand) == 0 {
				cmd.Help()
				return
			}

			limit := "ansibledev"

			if r, _ := cmd.Flags().GetBool("provision"); r {
				limit = "provisiontest"
			}

			if r, _ := cmd.Flags().GetBool("test"); r {
				limit = "vagrant"
			}

			executeExternalProgram("ansible", "-i", "hosts.ini", limit, "-m", "ping")
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			ensureAnsibleDirectory()
			ensureVagrantfile()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			ensureWorkingDirectoryAndExit()
		},
	}
)

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().BoolP("development", "d", true, "only ping the development VMs")
	pingCmd.Flags().BoolP("provision", "p", false, "only ping the provision VM")
	pingCmd.Flags().BoolP("test", "t", false, "only ping the test VMs")

	pingCmd.MarkFlagsMutuallyExclusive("development", "provision", "test")
}
