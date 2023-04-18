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
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	required bool

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Show the list of roles",
		Long:  `Show the name and version of each role installed in the roles_path.`,
		Run: func(cmd *cobra.Command, args []string) {
			list_roles()
		},
	}
)

func init() {
	roleCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "tell Ansible to print more debug messages")
	listCmd.Flags().BoolVarP(&required, "required", "r", false, "show only roles from requirements.yml")
}

func list_roles() {
	ensureAnsibleDirectory()

	if required {
		ensureRequirementsFile()
		list_required_roles()
	} else {
		list_all_roles()
	}

	ensureWorkingDirectoryAndExit()
}

func list_all_roles() {
	param := "list"

	if verbose {
		param = "list -v"
	}

	executeExternalProgram("ansible-galaxy", param)
}

func list_required_roles() {
	roles := readRequirementsFile()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Source", "Version"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, r := range roles {
		row := []string{r.Name, r.Source, r.Version}
		table.Append(row)
	}

	table.Render()
}
