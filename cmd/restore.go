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

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore Ansible collections and roles from files, URLs or Ansible Galaxy",
	Long:  "Restore Ansible collections and roles from files, URLs or Ansible Galaxy",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		force, _ := cmd.Flags().GetBool("force")

		param := []string{"install"}

		if verbose {
			param = append(param, "-v")
		}

		if force {
			param = append(param, "--force")
		}

		param = append(param, "-r", "requirements.yml")

		executeExternalProgram("ansible-galaxy", param...)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
		ensureRequirementsFile()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		ensureWorkingDirectoryAndExit()
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().Bool("force", false, "force overwriting existing roles or collections")
	restoreCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
}
