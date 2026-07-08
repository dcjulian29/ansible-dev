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
	"path/filepath"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// newCmd creates the Cobra command for "ansible-dev role new", which
// scaffolds a new Ansible role (typically backed by "ansible-galaxy role init").
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
// After scaffolding, the ansible-galaxy "tests" folder is removed and the
// embedded role template (LICENSE, README, lint configuration, GitHub
// workflows, meta/main.yml, ...) is overlaid with !!ROLE_NAME!! / !!ROLE_DESC!!
// substituted. When --publish is set, the role is additionally copied to the
// directory named by the ANSIBLE_ROLES environment variable, committed to a new
// git repository, pushed to a freshly-created public GitHub repository, and
// recorded in requirements.yml.
//
// Flags:
//   - --force, -f:       force overwrite of an existing role directory.
//     When set, the current role folder is deleted before
//     scaffolding (default false).
//   - --verbose, -v:     forward the verbose flag to [ansible.NewRole] so
//     that ansible-galaxy prints additional debug messages during
//     initialization (default false).
//   - --description, -d: description text substituted for !!ROLE_DESC!! in
//     the template and used for the published repository (default empty).
//   - --publish, -p:     create and push a public GitHub repository for the
//     role via git and gh, then add it to requirements.yml (default false).
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

			if err := ansible.NewRole(role, verbose); err != nil {
				return err
			}

			// ansible-galaxy init creates a "tests" folder that is not wanted
			// in the published role, so remove it before overlaying templates.
			tests := filepath.Join(folder, "tests")
			if filesystem.DirectoryExist(tests) {
				if err := filesystem.RemoveDirectory(tests); err != nil {
					return err
				}
			}

			description, _ := cmd.Flags().GetString("description")

			if err := ansible.ApplyRoleTemplate(folder, role, description); err != nil {
				return err
			}

			publish, _ := cmd.Flags().GetBool("publish")
			if !publish {
				return nil
			}

			if err := ansible.PublishRole(folder, role, description); err != nil {
				return err
			}

			base := ansible.BaseRoleName(role)

			requirements, _ := ansible.ReadRequirements()
			requirements.Roles = append(requirements.Roles, ansible.Role{
				Name:   role,
				Source: fmt.Sprintf("https://github.com/dcjulian29/ansible-role-%s.git", base),
			})

			if err := ansible.SaveRequirements(requirements); err != nil {
				return err
			}

			msg := fmt.Sprintf("role '%s' published and added to requirements.yml", role)
			fmt.Println(textformat.Info(msg))

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().BoolP("force", "f", false, "force overwriting an existing role")
	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	cmd.Flags().StringP("description", "d", "", "description of the role (fills the template)")
	cmd.Flags().BoolP("publish", "p", false, "create and push a public GitHub repository for the role")

	return cmd
}
