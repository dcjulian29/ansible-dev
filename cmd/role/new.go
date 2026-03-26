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

// newCmd creates the Cobra command for "ansible-dev role new", which
// scaffolds a new Ansible role in the development environment using
// [ansible.NewRole] (typically backed by "ansible-galaxy role init").
//
// Usage:
//
//	ansible-dev role new <role> [flags]
//
// The positional argument <role> is the name of the role to create (e.g.
// "dcjulian29.nginx"). If a role directory already exists at the resolved
// path (checked via [ansible.RoleFolderExists]):
//   - Without --force: the command returns an error prompting the user
//     to pass --force to replace the existing role.
//   - With --force: the existing directory is removed via
//     [filesystem.RemoveDirectory] before the new role is scaffolded.
//
// If no argument is supplied, the help text is displayed instead.
//
// Flags:
//   - --force, -f:   force overwrite of an existing role directory.
//     When set, the current role folder is deleted before
//     scaffolding (default false).
//   - --verbose, -v: forward the verbose flag to [ansible.NewRole] so
//     that ansible-galaxy prints additional debug messages during
//     initialization (default false).
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <role>",
		Short: "Create a new Ansible role in the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			role := args[0]

			force, _ := cmd.Flags().GetBool("force")
			folder, _ := ansible.RoleFolder(role)

			if ansible.RoleFolderExists(role) {
				if !force {
					return fmt.Errorf("role '%s' exists. Use '--force' to replace", role)
				}

				if err := filesystem.RemoveDirectory(folder); err != nil {
					return err
				}
			}

			verbose, _ := cmd.Flags().GetBool("verbose")

			return ansible.NewRole(role, verbose)
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().BoolP("force", "f", false, "force overwriting an existing role")
	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")

	return cmd
}
