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

// removeCmd creates the Cobra command for "ansible-dev collection remove",
// which removes a named Ansible collection from the requirements.yml file.
//
// Usage:
//
//	ansible-dev collection remove <collection>
//
// The positional argument <collection> is matched against the Name field of
// each [ansible.Collection] entry in requirements.yml. If a match is found,
// the entry is removed and the updated requirements are written back to
// disk via [ansible.SaveRequirements]. A confirmation message is printed to
// stdout on success.
//
// If no argument is supplied, the help text is displayed. If the named
// collection is not present in requirements.yml, an error is returned.
//
// Note: this command only modifies requirements.yml — it does not
// uninstall the collection files from the collections_path. Use
// "ansible-dev collection purge" to remove installed artifacts.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
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
				if err := ansible.SaveRequirements(requirements); err != nil {
					return err
				}

				msg := fmt.Sprintf("collection '%s' removed from requirements.yml", collection)
				fmt.Println(color.Info(msg))
				return nil
			}

			return fmt.Errorf("collection '%s' not present", collection)
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	return cmd
}
