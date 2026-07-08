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

package runbook

import (
	"fmt"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// newCmd creates the Cobra command for "ansible-dev runbook new", which
// scaffolds a new standalone Ansible runbook repository from the embedded
// runbook template.
//
// Usage:
//
//	ansible-dev runbook new <runbook> [flags]
//
// The runbook is rendered into the directory named by the ANSIBLE_RUNBOOKS
// environment variable joined with <runbook>, with !!RUNBOOK_NAME!! and
// !!RUNBOOK_DESC!! substituted. When --publish is set, the directory is
// committed to a new git repository and pushed to a freshly-created public
// GitHub repository named "ansible-runbook-<runbook>".
//
// Flags:
//   - --description, -d: description text substituted for !!RUNBOOK_DESC!! in
//     the template and used for the published repository (default empty).
//   - --publish, -p:     create and push a public GitHub repository for the
//     runbook via git and gh (default false).
//
// If no argument is supplied, the help text is displayed instead.
func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <runbook>",
		Short: "Scaffold a new Ansible runbook repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			name := args[0]
			description, _ := cmd.Flags().GetString("description")

			dir, err := ansible.NewRunbook(name, description)
			if err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("runbook '%s' created at '%s'", name, dir)))

			publish, _ := cmd.Flags().GetBool("publish")
			if !publish {
				return nil
			}

			if err := ansible.PublishRunbook(dir, name, description); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("runbook '%s' published", name)))

			return nil
		},
	}

	cmd.Flags().StringP("description", "d", "", "description of the runbook (fills the template)")
	cmd.Flags().BoolP("publish", "p", false, "create and push a public GitHub repository for the runbook")

	return cmd
}
