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

// removeCmd creates the Cobra command for "ansible-dev role remove",
// which removes an Ansible role entry from requirements.yml.
//
// Usage:
//
//	ansible-dev role remove <role> [flags]
//
// The positional argument <role> is the role name to remove (e.g.
// "dcjulian29.docker"). The command opens requirements.yml via
// [ansible.OpenRequirements] and removes every roles entry whose name
// matches the argument exactly, leaving the rest of the file untouched. If no
// matching role is found, an error is returned; so is a requirements.yml that
// is missing or does not parse.
//
// After a successful removal the file is saved, even when no roles are left,
// and an informational message is printed.
//
// The argument is checked by [cmdargs.Required]: typed bare, the command
// prints its help and exits 0; given flags but no role name, it fails with
// "expects a role name".
//
// Flags:
//   - --purge: additionally remove the role's installed files from disk
//     via [ansible.RemoveRole]. Without this flag only the
//     requirements.yml entry is removed (default false).
//
// Note: this command modifies the manifest (requirements.yml) only. To
// delete installed role files without touching the manifest, use
// "ansible-dev role delete". Passing --purge combines both operations
// in a single invocation.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <role>",
		Short: "Remove Ansible role from requirements.yml",
		// Checked in Args rather than RunE so a bare "role remove" prints its
		// help even outside a development folder, where PreRunE would fail
		// first, and "role remove --purge" names the missing role.
		Args: cmdargs.Required("a role name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			role := args[0]

			requirements, err := ansible.OpenRequirements()
			if err != nil {
				return err
			}

			found, err := requirements.RemoveRole(role)
			if err != nil {
				return err
			}

			if !found {
				return fmt.Errorf("role '%s' not present", role)
			}

			if err := requirements.Save(); err != nil {
				return err
			}

			msg := fmt.Sprintf("role '%s' removed from requirements.yml", role)
			fmt.Println(textformat.Info(msg))

			if r, _ := cmd.Flags().GetBool("purge"); r {
				return ansible.RemoveRole(role)
			}

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().Bool("purge", false, "remove role files too")

	return cmd
}
