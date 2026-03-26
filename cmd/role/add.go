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
	"fmt"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/color"
	"github.com/spf13/cobra"
)

// addCmd creates the Cobra command for "ansible-dev role add", which
// appends a new Ansible role to the requirements.yml file.
//
// Usage:
//
//	ansible-dev role add <role> [flags]
//
// The positional argument <role> is the role name (e.g.
// "geerlingguy.docker"). If no --source flag is provided, the role name
// is used as the source, implying it will be fetched from Ansible Galaxy.
//
// Flags:
//   - --source, -s: override the source URL or Galaxy reference for the
//     role. Defaults to the role name when omitted.
//   - --version, -v: pin the role to a specific version string.
//     Defaults to the empty string (latest).
//
// The command prints a reminder that the role must be restored (installed
// via "ansible-dev restore") before it can be used. If no argument is
// supplied, the help text is displayed instead.
//
// Note: this command only modifies requirements.yml — it does not
// download or install the role files. Use "ansible-dev restore" to
// install declared dependencies. This mirrors the behavior of the
// sibling "ansible-dev collection add" command.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <role>",
		Short: "Add Ansible role to requirements.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			name := args[0]

			source, _ := cmd.Flags().GetString("source")
			version, _ := cmd.Flags().GetString("version")

			if len(name) > 0 && len(source) == 0 {
				source = name // Ansible Galaxy Role
			}

			requirements, _ := ansible.ReadRequirements()

			role := ansible.Role{
				Name:    name,
				Source:  source,
				Version: version,
			}

			requirements.Roles = append(requirements.Roles, role)

			if err := ansible.SaveRequirements(requirements); err != nil {
				return err
			}

			msg := fmt.Sprintf("role '%s' added to requirements.yml but must be restored before use", name)
			fmt.Println(color.Info(msg))

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().StringP("source", "s", "", "source of the role")
	cmd.Flags().StringP("version", "v", "", "version of the role")

	return cmd
}
