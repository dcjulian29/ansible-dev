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
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// addCmd creates the Cobra command for "ansible-dev collection add", which
// appends a new Ansible Galaxy collection to the requirements.yml file.
//
// Usage:
//
//	ansible-dev collection add <collection> [flags]
//
// The positional argument <collection> is the fully qualified collection
// name (e.g. "community.general"). Without --source the entry is just the
// name, which ansible-galaxy installs from Galaxy. With --source it is
// written as a git collection (type: git, source: <url>), like the
// dcjulian29 runbook entries.
//
// The entry is appended through [ansible.OpenRequirements], which leaves the
// rest of requirements.yml untouched. A collection already listed under the
// same name is an error, as is a requirements.yml that is missing or does not
// parse.
//
// Flags:
//   - --source, -s: the git repository URL to install the collection from.
//     Omit it for a Galaxy collection.
//   - --version, -v: pin the collection to a specific version string.
//     Defaults to the empty string (latest).
//
// The command prints a reminder that the collection must be restored
// (installed) before it can be used. If no argument is supplied, the help
// text is displayed instead.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <collection>",
		Short: "Add Ansible collection to requirements.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			name := args[0]

			source, _ := cmd.Flags().GetString("source")
			version, _ := cmd.Flags().GetString("version")

			requirements, err := ansible.OpenRequirements()
			if err != nil {
				return err
			}

			collection := ansible.Collection{Name: name, Version: version}

			if source != "" {
				collection.Type = "git"
				collection.Source = source
			}

			// A collection already listed is refused rather than added twice.
			if err := requirements.AddCollection(collection); err != nil {
				return err
			}

			if err := requirements.Save(); err != nil {
				return err
			}

			msg := fmt.Sprintf("collection '%s' added to requirements.yml but must be restored before use", name)
			fmt.Println(textformat.Info(msg))

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().StringP("source", "s", "", "git repository URL of the collection (omit for Galaxy)")
	cmd.Flags().StringP("version", "v", "", "version of the collection")

	return cmd
}
