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

var roleAddCmd = &cobra.Command{
	Use:   "add <role>",
	Short: "Add Ansible role to requirements.yml",
	Long:  "Add Ansible role to requirements.yml",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		name := args[0]

		source, _ := cmd.Flags().GetString("source")
		version, _ := cmd.Flags().GetString("version")

		if len(name) > 0 && len(source) == 0 {
			source = name // Ansible Galaxy Role
		}

		var requirements Requirements

		if fileExists("requirements.yml") {
			requirements, _ = readRequirementsFile()
		}

		role := Role{
			Name:    name,
			Source:  source,
			Version: version,
		}

		requirements.Roles = append(requirements.Roles, role)

		writeRequirementsFile(requirements)

		fmt.Println(Info("role '%s' added to requirements.yml but must be restored before use.", name))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
	},
}

func init() {
	roleCmd.AddCommand(roleAddCmd)

	roleAddCmd.Flags().StringP("source", "s", "", "source of the role")
	roleAddCmd.Flags().StringP("version", "v", "", "version of the role")
}
