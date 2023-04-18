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
	"fmt"

	"github.com/spf13/cobra"
)

var (
	roleName    string
	roleSource  string
	roleVersion string

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a role to requirements.yml",
		Long: `Add a role to requirements.yml file.

The role will need to still need to be restored before usage.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				if len(roleName) > 0 || len(roleSource) > 0 {
					fmt.Println("ERROR: Can't supply role name as argument and flag!")
					ensureWorkingDirectoryAndExit()
				}

				roleName = args[0]
				roleSource = args[0]
			}

			if len(roleName) > 0 && len(roleSource) == 0 {
				roleSource = roleName
			}

			if len(roleName) == 0 {
				cmd.Help()
				ensureWorkingDirectoryAndExit()
			}

			add_requirement_role()
		},
	}
)

func init() {
	roleCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&roleName, "name", "", "name of role")
	addCmd.Flags().StringVarP(&roleSource, "source", "s", "", "source of the role")
	addCmd.Flags().StringVarP(&roleVersion, "version", "v", "", "version of the role")
}

func add_requirement_role() {
	ensureAnsibleDirectory()

	var roles []Role

	if fileExists("requirements.yml") {
		roles = readRequirementsFile()
	}

	role := Role{
		Name:    roleName,
		Source:  roleSource,
		Version: roleVersion,
	}

	roles = append(roles, role)

	writeRequirementsFile(roles)

	fmt.Println("Role", roleName, "added.")

	ensureWorkingDirectoryAndExit()
}
