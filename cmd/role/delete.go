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
	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/spf13/cobra"
)

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <role>",
		Short: "Remove Ansible role files from the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			role := args[0]
			folder, _ := ansible.RoleFolder(role)

			if ansible.RoleFolderExists(role) {
				if err := filesystem.RemoveDirectory(folder); err != nil {
					return err
				}
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	return cmd
}
