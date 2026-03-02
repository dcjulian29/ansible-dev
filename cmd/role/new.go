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
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/spf13/cobra"
)

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <role>",
		Short: "Create a new Ansible role in the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}

			role := args[0]

			force, _ := cmd.Flags().GetBool("force")
			folder, _ := ansible.RoleFolder(role)

			if ansible.RoleFolderExists(role) {
				if !force {
					return fmt.Errorf("role '%s' exists. Use '--force' to replace", role)
				}

				filesystem.RemoveDirectory(folder)
			}

			verbose, _ := cmd.Flags().GetBool("verbose")

			return ansible.NewRole(role, verbose)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().BoolP("force", "f", false, "force overwriting an existing role")
	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")

	return cmd
}
