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
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var roleListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show list of Ansible roles in the development environment",
	Long:  "Show list of Ansible roles in the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		required, _ := cmd.Flags().GetBool("requirements")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if required {
			ensureRequirementsFile()
			list_required_roles()
		} else {
			list_all_roles(verbose)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
	},
}

func init() {
	roleCmd.AddCommand(roleListCmd)

	roleListCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	roleListCmd.Flags().BoolP("requirements", "r", false, "show only roles from requirements.yml")
}

func list_all_roles(verbose bool) {
	param := []string{"role", "list"}

	if verbose {
		param = []string{"role", "list", "-v"}
	}

	executeExternalProgram("ansible-galaxy", param...)
}

func list_required_roles() {
	requirements, _ := readRequirementsFile()
	table := tablewriter.NewTable(os.Stdout, tablewriter.WithTrimSpace(tw.Off))

	for _, r := range requirements.Roles {
		row := []string{r.Name, r.Source, r.Version}
		table.Append(row)
	}

	table.Render()
}
