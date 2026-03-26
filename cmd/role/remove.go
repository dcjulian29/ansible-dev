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

// removeCmd creates the Cobra command for "ansible-dev role remove",
// which removes an Ansible role entry from requirements.yml.
//
// Usage:
//
//	ansible-dev role remove <role> [flags]
//
// The positional argument <role> is the role name to remove (e.g.
// "dcjulian29.docker"). The command reads the current requirements via
// [ansible.ReadRequirements], searches for a role whose Name field
// matches the argument exactly, and splices it out of the list. If no
// matching role is found, an error is returned.
//
// After a successful removal the updated manifest is persisted via
// [ansible.SaveRequirements] and an informational message is printed.
//
// If no argument is supplied, the help text is displayed instead.
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			var changed []ansible.Role

			role := args[0]
			requirements, _ := ansible.ReadRequirements()
			found := false

			for i, r := range requirements.Roles {
				if r.Name == role {
					changed = append(requirements.Roles[:i], requirements.Roles[i+1:]...)
					found = true
				}
			}

			if !found {
				return fmt.Errorf("role '%s' not present", role)
			}

			if len(changed) > 0 {
				requirements.Roles = changed
				if err := ansible.SaveRequirements(requirements); err != nil {
					return err
				}

				msg := fmt.Sprintf("role '%s' removed from requirements.yml", role)
				fmt.Println(color.Info(msg))
			}

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
