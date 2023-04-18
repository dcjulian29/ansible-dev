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
	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore role(s) from file(s), URL(s) or Ansible Galaxy",
		Long:  `Restore role(s) from file(s), URL(s) or Ansible Galaxy.`,
		Run: func(cmd *cobra.Command, args []string) {
			restore_roles()
		},
	}
)

func init() {
	roleCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "tell Ansible to print more debug messages")
}

func restore_roles() {
	ensureAnsibleDirectory()
	ensureRequirementsFile()

	param := []string{"install", "-r", "requirements.yml"}

	if verbose {
		param = []string{"install", "-v", "-r", "requirements.yml"}
	}

	executeExternalProgram("ansible-galaxy", param...)
	ensureWorkingDirectoryAndExit()
}
