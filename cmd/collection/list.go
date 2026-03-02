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
package collection

import (
	"os"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show list of Ansible collections in the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			required, _ := cmd.Flags().GetBool("requirements")
			verbose, _ := cmd.Flags().GetBool("verbose")

			if required {
				requirements, err := ansible.ReadRequirements()
				if err != nil {
					return err
				}

				table := tablewriter.NewTable(os.Stdout, tablewriter.WithTrimSpace(tw.Off))

				for _, c := range requirements.Collections {
					row := []string{c.Name, c.Source, c.Type, c.Version}
					if err := table.Append(row); err != nil {
						return nil
					}
				}

				return table.Render()
			} else {
				param := []string{"collection", "list"}

				if verbose {
					param = []string{"collection", "list", "-v"}
				}

				err := execute.ExternalProgram("ansible-galaxy", param...)
				if err != nil {
					return err
				}
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	cmd.Flags().BoolP("requirements", "r", false, "show only collections from requirements.yml")

	return cmd
}
