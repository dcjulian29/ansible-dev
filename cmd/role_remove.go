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
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var (
	searchRoleName string
	removeRole     bool

	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a role from requirements.yml",
		Long:  `Remove a role from requirements.yml file.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				if len(searchRoleName) > 0 {
					fmt.Println("ERROR: Can't supply role name as argument and flag!")
					ensureWorkingDirectoryAndExit()
				}

				searchRoleName = args[0]
			}

			if len(searchRoleName) == 0 {
				if removeRole {
					fmt.Println("ERROR: Can't remove role files without specifying the name!")
					ensureWorkingDirectoryAndExit()
				}

				cmd.Help()
				ensureWorkingDirectoryAndExit()
			}

			remove_requirement_role()

			if removeRole {
				remove_role()
			}

			ensureWorkingDirectoryAndExit()
		},
	}
)

func init() {
	roleCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVar(&searchRoleName, "name", "", "name of role")
	removeCmd.Flags().BoolVar(&removeRole, "purge", false, "remove role files too")
}

func remove_requirement_role() {
	ensureAnsibleDirectory()

	changed := false
	roles := readRequirementsFile()
	var newRoles []Role

	for i, r := range roles {
		if r.Name == searchRoleName {
			changed = true
			newRoles = append(roles[:i], roles[i+1:]...)
		}
	}

	if changed {
		writeRequirementsFile(newRoles)
		fmt.Println("Role", searchRoleName, "removed.")
		return
	}

	fmt.Println("WARN: Role", searchRoleName, "not present.")
}

func remove_role() {
	ensureAnsibleDirectory()

	cfg, err := ini.Load("ansible.cfg")
	if err != nil {
		fmt.Println(err)
		return
	}

	section, err := cfg.GetSection("defaults")
	if err != nil {
		fmt.Println(err)
		return
	}

	rolePath, err := section.GetKey("roles_path")
	if err != nil {
		fmt.Println(err)
		return
	}

	roleFolder := filepath.Join(rolePath.String(), searchRoleName)

	files, err := filepath.Glob(filepath.Join(roleFolder, "*"))
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(files) == 0 {
		fmt.Println("WARN: Role", searchRoleName, "files not present.")
		return
	}

	removeDir(roleFolder)

	fmt.Println("Role", searchRoleName, "files purged.")
}
