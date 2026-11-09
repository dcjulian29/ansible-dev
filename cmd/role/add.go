/*
Copyright © 2026 Julian Easterling

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
	"github.com/dcjulian29/ansible-dev/internal/cmdargs"
	"github.com/dcjulian29/go-toolbox/textformat"
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
// The role is appended through [ansible.OpenRequirements], which leaves the
// rest of requirements.yml untouched. A role already listed under the same
// name is an error, as is a requirements.yml that is missing or does not
// parse.
//
// The command prints a reminder that the role must be restored (installed
// via "ansible-dev restore") before it can be used. The argument is checked by
// [cmdargs.Required]: typed bare, the command prints its help and exits 0;
// given flags but no role name, it fails with "expects a role name".
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
		// Checked in Args rather than RunE so a bare "role add" prints its help
		// even outside a development folder, where PreRunE would fail first.
		Args: cmdargs.Required("a role name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			source, _ := cmd.Flags().GetString("source")
			version, _ := cmd.Flags().GetString("version")

			if len(name) > 0 && len(source) == 0 {
				source = name // Ansible Galaxy Role
			}

			requirements, err := ansible.OpenRequirements()
			if err != nil {
				return err
			}

			role := ansible.Role{
				Name:    name,
				Source:  source,
				Version: version,
			}

			// A role already listed is refused rather than added a second time.
			if err := requirements.AddRole(role); err != nil {
				return err
			}

			if err := requirements.Save(); err != nil {
				return err
			}

			msg := fmt.Sprintf("role '%s' added to requirements.yml but must be restored before use", name)
			fmt.Println(textformat.Info(msg))

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
