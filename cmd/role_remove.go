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

	"github.com/spf13/cobra"
)

var roleRemoveCmd = &cobra.Command{
	Use:   "remove <role>",
	Short: "Remove Ansible role from requirements.yml",
	Long:  "Remove Ansible role from requirements.yml",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		var changed []Role

		role := args[0]
		requirements, _ := readRequirementsFile()

		for i, r := range requirements.Roles {
			if r.Name == role {
				changed = append(requirements.Roles[:i], requirements.Roles[i+1:]...)
			}
		}

		if len(changed) > 0 {
			requirements.Roles = changed
			writeRequirementsFile(requirements)
			fmt.Println(Info("role '%s' removed from requirements.yml", role))
			return
		}

		fmt.Println(Fatal("role '%s' not present.", role))

		if r, _ := cmd.Flags().GetBool("purge"); r {
			remove_role(role)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
	},
}

func init() {
	roleCmd.AddCommand(roleRemoveCmd)

	roleRemoveCmd.Flags().Bool("purge", false, "remove role files too")
}
