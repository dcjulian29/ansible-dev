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
	"strings"

	"github.com/spf13/cobra"
)

var (
	shellCommand string

	shellCmd = &cobra.Command{
		Use:   "shell [flags] -- <command>",
		Short: "Execute shell command in the Ansible development vagrant environment",
		Long:  "Execute shell command in the Ansible development vagrant environment",
		Args: func(cmd *cobra.Command, args []string) error {
			shellCommand = strings.Join(args, " ")

			return nil
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			if len(shellCommand) == 0 {
				cmd.Help()
				return
			}

			limit := "ansibledev"

			if r, _ := cmd.Flags().GetBool("test"); r {
				limit = "vagrant"
			}

			param := []string{
				"-i", "hosts.ini",
				"-l", limit,
				"-m", "shell",
				"-a", shellCommand,
				"all",
			}

			executeExternalProgram("ansible", param...)
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
	rootCmd.AddCommand(shellCmd)

	shellCmd.Flags().BoolP("development", "d", true, "only execute on the development VMs")
	shellCmd.Flags().BoolP("test", "t", false, "only execute on the test VMs")

	shellCmd.MarkFlagsMutuallyExclusive("development", "test")
}
