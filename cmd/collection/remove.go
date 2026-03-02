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
package collection

import (
	"fmt"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/color"
	"github.com/spf13/cobra"
)

func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <collection>",
		Short: "Remove Ansible collection from requirements.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			var changed []ansible.Collection

			collection := args[0]
			requirements, _ := ansible.ReadRequirements()

			for i, r := range requirements.Collections {
				if r.Name == collection {
					changed = append(requirements.Collections[:i], requirements.Collections[i+1:]...)
				}
			}

			if len(changed) > 0 {
				requirements.Collections = changed
				ansible.SaveRequirements(requirements)
				msg := fmt.Sprintf("collection '%s' removed from requirements.yml", collection)
				fmt.Println(color.Info(msg))
				return nil
			}

			return fmt.Errorf("collection '%s' not present.", collection)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	return cmd
}
