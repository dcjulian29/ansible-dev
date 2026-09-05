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

package role

import (
	"fmt"
	"path/filepath"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// newCmd creates the Cobra command for "ansible-dev role new", which
// scaffolds a new Ansible role (typically backed by "ansible-galaxy role init").
//
// Usage:
//
//	ansible-dev role new <role> [flags]
//
// The positional argument <role> is the name of the role to create, either bare
// ("nginx") or qualified with a namespace ("dcjulian29.nginx"). A qualified name
// supplies the namespace; see --namespace below. If a role directory already
// exists at the resolved path (checked via [ansible.RoleFolderExists]):
//   - Without --force: the command returns an error prompting the user
//     to pass --force to replace the existing role.
//   - With --force: the existing directory is removed via
//     [filesystem.RemoveDirectory] before the new role is scaffolded.
//
// If no argument is supplied, the help text is displayed instead.
//
// After scaffolding, the ansible-galaxy "tests" folder is removed and the
// embedded role template (LICENSE, README, lint configuration, GitHub
// workflows, meta/main.yml, ...) is overlaid with !!ROLE_NAME!! / !!ROLE_DESC!!
// substituted. When --publish is set, the role is additionally copied to the
// configured roles_path directory, committed to a new git repository, pushed to
// a freshly-created public GitHub repository, and recorded in requirements.yml.
//
// Flags:
//   - --namespace, -n:   the Galaxy/GitHub namespace the role belongs to. It
//     fills meta/main.yml, the generated README, the published repository
//     owner, and the requirements.yml entry. Defaults to the namespace in
//     <role> when one is given, otherwise the configured namespace (see
//     "ansible-dev config namespace"); the command fails when none is
//     available.
//   - --force, -f:       force overwrite of an existing role directory.
//     When set, the current role folder is deleted before
//     scaffolding (default false).
//   - --verbose, -v:     forward the verbose flag to [ansible.NewRole] so
//     that ansible-galaxy prints additional debug messages during
//     initialization (default false).
//   - --description, -d: what the role does, written as a verb phrase. It is
//     composed by [ansible.RoleDescription] into the text substituted for
//     !!ROLE_DESC!! in the template and used for the published repository
//     (default empty).
//   - --no-prefix:       use the description verbatim rather than prefixing it
//     with "Ansible role to" (default false).
//   - --publish, -p:     create and push a public GitHub repository for the
//     role via git and gh, then add it to requirements.yml (default false).
//
// The composed description is printed before the template is applied, so the
// added phrase is visible at the point it is chosen rather than only in the
// generated files.
// resolveNamespace determines the namespace a new role belongs to, in
// decreasing order of how explicitly it was chosen: the --namespace flag, then
// a namespace written into the role argument itself ("acme.nginx"), then the
// configured default.
//
// Reading it from the role argument is what keeps the existing habit of typing
// a fully-qualified name working without also having to pass a flag, and lets a
// one-off role for another organization be created without touching the
// configuration. When none of the three supplies a value the command stops:
// there is no default worth guessing, and the namespace is difficult to change
// once it has been published.
func resolveNamespace(cmd *cobra.Command, role string) (string, error) {
	namespace, _ := cmd.Flags().GetString("namespace")

	if namespace == "" {
		namespace = ansible.RoleNamespace(role)
	}

	if namespace == "" {
		configured, err := settings.Namespace()
		if err != nil {
			return "", err
		}

		namespace = configured
	}

	if err := ansible.ValidateNamespace(namespace); err != nil {
		return "", err
	}

	return namespace, nil
}

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <role>",
		Short: "Create a new Ansible role in the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			role := args[0]

			// Resolved before anything is created: the namespace is baked into
			// meta/main.yml, the README, the published repository and the
			// requirements entry, so discovering it is missing after
			// ansible-galaxy has already scaffolded the role would leave a
			// half-made role behind.
			namespace, err := resolveNamespace(cmd, role)
			if err != nil {
				return err
			}

			qualified := ansible.QualifiedRoleName(namespace, role)

			fmt.Println(textformat.Info(fmt.Sprintf("creating role '%s'", qualified)))

			force, _ := cmd.Flags().GetBool("force")
			folder, _ := ansible.RoleFolder(role)

			if ansible.RoleFolderExists(role) {
				if !force {
					return fmt.Errorf("role '%s' exists. Use '--force' to replace", role)
				}

				if err := filesystem.RemoveDirectory(folder); err != nil {
					return err
				}
			}

			verbose, _ := cmd.Flags().GetBool("verbose")

			if err := ansible.NewRole(role, verbose); err != nil {
				return err
			}

			// ansible-galaxy init creates a "tests" folder that is not wanted
			// in the published role, so remove it before overlaying templates.
			tests := filepath.Join(folder, "tests")
			if filesystem.DirectoryExist(tests) {
				if err := filesystem.RemoveDirectory(tests); err != nil {
					return err
				}
			}

			raw, _ := cmd.Flags().GetString("description")
			noprefix, _ := cmd.Flags().GetBool("no-prefix")
			description := ansible.RoleDescription(raw, !noprefix)

			// The composed description is what lands in meta/main.yml, the
			// README, and the published repository, so show it rather than
			// leaving the caller to discover the added phrase afterwards.
			if description != "" {
				fmt.Println(textformat.Info(fmt.Sprintf("described as '%s'", description)))
			}

			if err := ansible.ApplyRoleTemplate(folder, namespace, role, description); err != nil {
				return err
			}

			publish, _ := cmd.Flags().GetBool("publish")
			if !publish {
				return nil
			}

			if err := ansible.PublishRole(folder, namespace, role, description); err != nil {
				return err
			}

			base := ansible.BaseRoleName(role)

			// Recorded under the qualified name even when the role argument was
			// bare, so requirements.yml matches what the role calls itself in
			// meta/main.yml.
			requirements, _ := ansible.ReadRequirements()
			requirements.Roles = append(requirements.Roles, ansible.Role{
				Name:   qualified,
				Source: fmt.Sprintf("https://github.com/%s/ansible-role-%s.git", namespace, base),
			})

			if err := ansible.SaveRequirements(requirements); err != nil {
				return err
			}

			msg := fmt.Sprintf("role '%s' published and added to requirements.yml", qualified)
			fmt.Println(textformat.Info(msg))

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().BoolP("force", "f", false, "force overwriting an existing role")
	cmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	cmd.Flags().StringP("namespace", "n", "",
		"Galaxy/GitHub namespace for the role (default: the namespace in <role>, else the configured one)")
	cmd.Flags().StringP("description", "d", "",
		"what the role does, as a verb phrase completing \"Ansible role to ...\"")
	cmd.Flags().Bool("no-prefix", false,
		"use the description verbatim instead of prefixing it with \"Ansible role to\"")
	cmd.Flags().BoolP("publish", "p", false, "create and push a public GitHub repository for the role")

	return cmd
}
