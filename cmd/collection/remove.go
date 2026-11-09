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

package collection

import (
	"fmt"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/cmdargs"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// removeCmd creates the Cobra command for "ansible-dev collection remove",
// which removes a named Ansible collection from the requirements.yml file.
//
// Usage:
//
//	ansible-dev collection remove <collection>
//
// The positional argument <collection> is matched against the name of each
// collections entry in requirements.yml, opened via
// [ansible.OpenRequirements]. Every match is removed and the file is saved,
// leaving the rest of it untouched, even when no collections are left. A
// confirmation message is printed to stdout on success. A requirements.yml
// that is missing or does not parse is an error.
//
// The argument is checked by [cmdargs.Required]: typed bare, the command
// prints its help and exits 0; given more than one collection name, it fails
// without changing anything. If the named collection is not present in
// requirements.yml, an error is returned.
//
// Note: this command only modifies requirements.yml — it does not
// uninstall the collection files from the collections_path.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <collection>",
		Short: "Remove Ansible collection from requirements.yml",
		// Checked in Args rather than RunE so a bare "collection remove" prints
		// its help even outside a development folder, and "collection remove a
		// b" is refused instead of quietly ignoring b.
		Args: cmdargs.Required("a collection name"),
		RunE: func(_ *cobra.Command, args []string) error {
			collection := args[0]

			requirements, err := ansible.OpenRequirements()
			if err != nil {
				return err
			}

			found, err := requirements.RemoveCollection(collection)
			if err != nil {
				return err
			}

			if !found {
				return fmt.Errorf("collection '%s' not present", collection)
			}

			if err := requirements.Save(); err != nil {
				return err
			}

			msg := fmt.Sprintf("collection '%s' removed from requirements.yml", collection)
			fmt.Println(textformat.Info(msg))

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	return cmd
}
