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

// Package restore implements the "ansible-dev restore" command, which
// installs all Ansible Galaxy collections and roles declared in the
// project's requirements.yml file.
package restore

import (
	"errors"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the Cobra command for
// "ansible-dev restore".
//
// The command runs "ansible-galaxy install -r requirements.yml" to
// download and install every collection and role listed in the project's
// requirements file. This is the counterpart to the collection and role
// "add" subcommands, which only modify the manifest without installing
// artifacts.
//
// Flags:
//   - --force:      force overwriting of already-installed roles or
//     collections (passes --force to ansible-galaxy; default false).
//   - --verbose, -v: enable verbose ansible-galaxy output (passes -v;
//     default false).
//
// A PreRunE hook validates the environment with two checks:
//  1. [ansible.EnsureAnsibleDirectory] confirms ansible.cfg is present.
//     A custom error message ("not an Ansible development directory") is
//     returned on failure.
//  2. [ansible.EnsureRequirementsFile] confirms requirements.yml exists.
//
// An error is returned if either pre-flight check fails or if
// ansible-galaxy exits with a non-zero status.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore Ansible collections and roles from the requirements.yml file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			force, _ := cmd.Flags().GetBool("force")

			param := []string{"install"}

			if verbose {
				param = append(param, "-v")
			}

			if force {
				param = append(param, "--force")
			}

			param = append(param, "-r", "requirements.yml")

			if err := execute.ExternalProgram("ansible-galaxy", param...); err != nil {
				return err
			}

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if err := ansible.EnsureAnsibleDirectory(); err != nil {
				return errors.New("not an Ansible development directory")
			}

			return ansible.EnsureRequirementsFile()
		},
	}

	cmd.Flags().Bool("force", false, "force overwriting existing roles or collections")
	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")

	return cmd
}
