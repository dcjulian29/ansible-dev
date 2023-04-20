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

var collectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show list of Ansible collections in the development environment",
	Long:  "Show list of Ansible collections in the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		required, _ := cmd.Flags().GetBool("requirements")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if required {
			ensureRequirementsFile()
			list_required_collections()
		} else {
			list_all_collections(verbose)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		ensureWorkingDirectoryAndExit()
	},
}

func init() {
	collectionCmd.AddCommand(collectionListCmd)

	collectionListCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	collectionListCmd.Flags().BoolP("requirements", "r", false, "show only collections from requirements.yml")
}

func list_all_collections(verbose bool) {
	param := []string{"collection", "list"}

	if verbose {
		param = []string{"collection", "list", "-v"}
	}

	executeExternalProgram("ansible-galaxy", param...)
}

func list_required_collections() {
	requirements, _ := readRequirementsFile()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Source", "Type", "Version"})
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

	for _, c := range requirements.Collections {
		row := []string{c.Name, c.Source, c.Type, c.Version}
		table.Append(row)
	}

	table.Render()
}
