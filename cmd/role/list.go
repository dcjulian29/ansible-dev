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
package role

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
		Short: "Show list of Ansible roles in the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			required, _ := cmd.Flags().GetBool("requirements")
			verbose, _ := cmd.Flags().GetBool("verbose")

			if required {
				if err := ansible.EnsureRequirementsFile(); err != nil {
					return err
				}

				requirements, _ := ansible.ReadRequirements()
				table := tablewriter.NewTable(os.Stdout, tablewriter.WithTrimSpace(tw.Off))

				for _, r := range requirements.Roles {
					row := []string{r.Name, r.Source, r.Version}
					if err := table.Append(row); err != nil {
						return err
					}
				}

				if err := table.Render(); err != nil {
					return err
				}
			} else {
				param := []string{"role", "list"}

				if verbose {
					param = []string{"role", "list", "-v"}
				}

				return execute.ExternalProgram("ansible-galaxy", param...)
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	cmd.Flags().BoolP("requirements", "r", false, "show only roles from requirements.yml")

	return cmd
}
