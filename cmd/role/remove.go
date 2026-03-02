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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().Bool("purge", false, "remove role files too")

	return cmd
}
