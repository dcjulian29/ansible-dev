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
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// deleteCmd creates the Cobra command for "ansible-dev role delete",
// which removes an installed Ansible role's files from the local
// development environment.
//
// Usage:
//
//	ansible-dev role delete <role>
//
// The positional argument <role> is the role name (e.g.
// "dcjulian29.docker"). A bare name is resolved by
// [ansible.InstalledRoleName]: the folder as typed if it exists, otherwise
// "<configured namespace>.<name>", the form "role new" creates. The folder is
// then removed entirely using [filesystem.RemoveDirectory]. If no folder
// exists, the command fails with "role '<name>' folder not present".
//
// The argument is checked by [cmdargs.Required]: typed bare, the command
// prints its help and exits 0; given more than one role name, it fails
// without deleting anything.
//
// Note: this command only deletes the on-disk role files. It does not
// modify requirements.yml. To remove a role from the manifest, use
// "ansible-dev role remove" instead. The two commands can be used
// together to fully clean up a role.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <role>",
		Short: "Remove Ansible role files from the development environment",
		// Checked in Args rather than RunE so a bare "role delete" prints its
		// help even outside a development folder, and "role delete a b" is
		// refused before anything is deleted instead of quietly ignoring b.
		Args: cmdargs.Required("a role name"),
		RunE: func(_ *cobra.Command, args []string) error {
			// "role new nginx" creates roles/<namespace>.nginx, so a bare name
			// is resolved to that folder when no roles/nginx exists.
			role := ansible.InstalledRoleName(args[0])

			// An ansible.cfg without a readable roles_path must stop the
			// command: RoleFolderExists hits the same error and reports false,
			// so ignoring it here made a delete that removed nothing exit 0
			// as if it had succeeded.
			folder, err := ansible.RoleFolder(role)
			if err != nil {
				return err
			}

			// A missing folder used to exit 0 silently, so a mistyped or
			// unqualified name looked like a successful delete.
			if !ansible.RoleFolderExists(role) {
				return fmt.Errorf("role '%s' folder not present", role)
			}

			if err := filesystem.RemoveDirectory(folder); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("role '%s' files were deleted", role)))

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	return cmd
}
