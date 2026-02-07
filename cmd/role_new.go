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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var roleNewCmd = &cobra.Command{
	Use:   "new <role>",
	Short: "Create a new Ansible role in the development environment",
	Long:  "Create a new Ansible role in the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		role := args[0]
		force, _ := cmd.Flags().GetBool("force")
		folder, _ := roleFolder(role)

		if roleFolderExists(role) {
			if !force {
				fmt.Println(Fatal("Role exists! Use -force to replace."))
				return
			}

			removeDir(folder)
		}

		param := []string{
			"init",
			role,
			"--init-path",
			strings.ReplaceAll(rootRoleFolder(), "\\", "/"),
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		if verbose {
			param = append(param, "--verbose")
		}

		executeExternalProgram("ansible-galaxy", param...)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		ensureWorkingDirectoryAndExit()
	},
}

func init() {
	roleCmd.AddCommand(roleNewCmd)

	roleNewCmd.Flags().BoolP("force", "f", false, "force overwriting an existing role")
	roleNewCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
}
