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
	"github.com/dcjulian29/ansible-dev/internal/cmdargs"
	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

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

// newCmd creates the Cobra command for "ansible-dev role new", which
// scaffolds a new Ansible role (typically backed by "ansible-galaxy role init").
//
// Usage:
//
//	ansible-dev role new <role> [flags]
//
// The positional argument <role> is the name of the role to create, either bare
// ("nginx") or qualified with a namespace ("dcjulian29.nginx"). A qualified name
// supplies the namespace; see --namespace below. Either way the role is created
// in roles_path as "<namespace>.<name>". If a role directory already exists at
// that path (checked via [ansible.RoleFolderExists]):
//   - Without --force: the command returns an error prompting the user
//     to pass --force to replace the existing role.
//   - With --force: the existing directory is removed via
//     [filesystem.RemoveDirectory] before the new role is scaffolded.
//
// Exactly one argument is required, checked by [cmdargs.Required]: typed bare,
// the command prints its help and exits 0; given flags but no role name, it
// fails with "expects a role name". Both happen even outside a development
// folder.
//
// After scaffolding, the ansible-galaxy "tests" folder is removed, the
// placeholder main.yml in "handlers" and "vars" is deleted (the folders stay),
// and the
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
//   - --version:         the min_ansible_version written into meta/main.yml,
//     as "2.21" or "2.21.5". Defaults to the ansible-core version that
//     ansible-galaxy reports, i.e. the one the role is developed on.
//   - --platform:        a platform the role supports, as
//     NAME:VERSION[,VERSION...] (for example "Debian:trixie,bookworm");
//     repeat it for each platform. Names are checked against the Galaxy list
//     and versions are always quoted. Defaults to Debian trixie and EL 10.
//
// All flags are checked before an existing role is removed by --force. With
// --publish, so is requirements.yml: a missing or unparseable file, or a role
// already listed in it, stops the command before anything is created or
// pushed. The entry is saved only after publishing succeeds.
//
// The composed description is printed before the template is applied, so the
// added phrase is visible at the point it is chosen rather than only in the
// generated files.
func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <role>",
		Short: "Create a new Ansible role in the development environment",
		// A bare "role new" prints the help; with flags but no name it fails
		// with "expects a role name". Checked in Args so both happen even
		// outside a development folder, before PreRunE can fail first.
		Args: cmdargs.Required("a role name"),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// Every flag is checked here, before --force can remove an existing
			// role, so a bad value never costs the role it was meant to replace.

			// min_ansible_version in meta/main.yml: --version when given,
			// otherwise the installed version, since that is what the role is
			// developed and tested on. An Ansible that cannot report its
			// version stops the command here too.
			ansibleVersion, _ := cmd.Flags().GetString("version")

			if ansibleVersion != "" {
				if err := ansible.ValidateAnsibleVersion(ansibleVersion); err != nil {
					return err
				}
			} else {
				ansibleVersion, err = ansible.AnsibleVersion()
				if err != nil {
					return err
				}
			}

			// galaxy_info.platforms in meta/main.yml: each --platform value, or
			// Debian trixie + EL 10 when none is given.
			values, _ := cmd.Flags().GetStringArray("platform")

			platforms, err := ansible.ParsePlatforms(values)
			if err != nil {
				return err
			}

			qualified := ansible.QualifiedRoleName(namespace, role)

			// With --publish, the requirements.yml entry is added now, in
			// memory, and saved only once publishing succeeds. Doing it here,
			// before --force or ansible-galaxy touch anything and before the
			// repository is pushed, means a missing or broken file, or a role
			// already listed, stops the command while there is nothing to undo.
			publish, _ := cmd.Flags().GetBool("publish")

			var requirements *ansible.RequirementsDocument

			if publish {
				// An existing published clone (or an unset roles_path) would
				// otherwise only stop PublishRole after --force had already
				// deleted and re-scaffolded the development copy.
				if _, err := ansible.PublishDestination(role); err != nil {
					return err
				}

				requirements, err = ansible.OpenRequirements()
				if err != nil {
					return err
				}

				// Recorded under the qualified name even when the role argument
				// was bare, so requirements.yml matches what the role calls
				// itself in meta/main.yml.
				err = requirements.AddRole(ansible.Role{
					Name: qualified,
					Source: fmt.Sprintf("https://github.com/%s/ansible-role-%s.git",
						namespace, ansible.BaseRoleName(role)),
				})
				if err != nil {
					return err
				}
			}

			fmt.Println(textformat.Info(fmt.Sprintf("creating role '%s'", qualified)))

			// The development copy is always named "<namespace>.<name>", like
			// every other role under roles_path, whichever form was typed. The
			// existence check, --force and ansible-galaxy must all use this one
			// name: checking one folder and removing another would let --force
			// delete the wrong role. (The published clone keeps the bare name;
			// see [ansible.PublishRole].)
			force, _ := cmd.Flags().GetBool("force")

			folder, err := ansible.RoleFolder(qualified)
			if err != nil {
				return err
			}

			if ansible.RoleFolderExists(qualified) {
				if !force {
					return fmt.Errorf("role '%s' exists. Use '--force' to replace", qualified)
				}

				if err := filesystem.RemoveDirectory(folder); err != nil {
					return err
				}
			}

			verbose, _ := cmd.Flags().GetBool("verbose")

			// ansible-galaxy init accepts a dotted role name and creates the
			// folder under exactly that name, so no rename is needed afterwards.
			if err := ansible.NewRole(qualified, verbose); err != nil {
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

			// handlers/ and vars/ stay, since a new role often grows into
			// them, but their generated main.yml holds nothing except an
			// MIT-0 SPDX header (which contradicts the Apache-2.0 license and
			// fails yaml[comments]) and a placeholder comment. Remove the
			// files and leave the folders empty; git does not track an empty
			// folder, so nothing unused is published. defaults/ and tasks/
			// keep theirs: the template overlay below replaces each main.yml
			// with a clean one.
			for _, unused := range []string{"handlers", "vars"} {
				if err := filesystem.RemoveFile(filepath.Join(folder, unused, "main.yml")); err != nil {
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

			if err := ansible.ApplyRoleTemplate(folder, namespace, role, description, ansibleVersion, platforms); err != nil {
				return err
			}

			if !publish {
				return nil
			}

			if err := ansible.PublishRole(folder, namespace, role, description); err != nil {
				return err
			}

			// The entry was added before anything was created; save it now
			// that the repository it points at exists.
			if err := requirements.Save(); err != nil {
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
	cmd.Flags().String("version", "",
		"min_ansible_version for meta/main.yml, as 2.21 or 2.21.5 (default: the installed ansible-core version)")
	// StringArray, not StringSlice: a slice flag splits on commas, which would
	// cut "Debian:trixie,bookworm" into two invalid values.
	cmd.Flags().StringArray("platform", nil,
		"supported platform for meta/main.yml as NAME:VERSION[,VERSION...]; repeatable (default: Debian:trixie and EL:10)")

	return cmd
}
