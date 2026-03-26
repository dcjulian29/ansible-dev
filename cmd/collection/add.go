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

// addCmd creates the Cobra command for "ansible-dev collection add", which
// appends a new Ansible Galaxy collection to the requirements.yml file.
//
// Usage:
//
//	ansible-dev collection add <collection> [flags]
//
// The positional argument <collection> is the fully qualified collection
// name (e.g. "community.general"). If no --source flag is provided, the
// collection name is used as the source, implying it will be fetched from
// Ansible Galaxy. The collection type is always set to "galaxy".
//
// Flags:
//   - --source, -s: override the source URL or Galaxy reference for the
//     collection. Defaults to the collection name when omitted.
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

			if len(name) > 0 && len(source) == 0 {
				source = name // Ansible Galaxy Collection
			}

			requirements, _ := ansible.ReadRequirements()

			collection := ansible.Collection{
				Name:    name,
				Source:  source,
				Type:    "galaxy",
				Version: version,
			}

			requirements.Collections = append(requirements.Collections, collection)

			if err := ansible.SaveRequirements(requirements); err != nil {
				return err
			}

			msg := fmt.Sprintf("collection '%s' added to requirements.yml but must be restored before use", name)
			fmt.Println(color.Info(msg))

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().StringP("source", "s", "", "source of the collection")
	cmd.Flags().StringP("version", "v", "", "version of the collection")

	return cmd
}
