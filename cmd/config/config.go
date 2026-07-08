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

// Package config implements the "ansible-dev config" command group, which
// views and edits the ansible-dev configuration file (~/.config/ansible-dev.yml).
package config

import (
	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/spf13/cobra"
)

// NewCommand returns the Cobra command for the "config" group. Its subcommands
// display the configuration (show), set the required repository paths
// (roles-path, runbooks-path), manage the compare ignore lists (role-ignore,
// runbook-ignore), and configure the diff tool for the current operating system
// (diff-program, diff-role-filter, diff-runbook-filter, diff-args).
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and manage the ansible-dev configuration file",
	}

	cmd.AddCommand(showCmd())
	cmd.AddCommand(rolesPathCmd())
	cmd.AddCommand(runbooksPathCmd())

	cmd.AddCommand(ignoreCmd(
		"role-ignore",
		"Manage the 'role compare' ignore list",
		func(c *settings.Config) []string { return c.RoleIgnore },
		func(c *settings.Config, v []string) { c.RoleIgnore = v },
	))
	cmd.AddCommand(ignoreCmd(
		"runbook-ignore",
		"Manage the 'runbook compare' ignore list",
		func(c *settings.Config) []string { return c.RunbookIgnore },
		func(c *settings.Config, v []string) { c.RunbookIgnore = v },
	))

	cmd.AddCommand(diffProgramCmd())
	cmd.AddCommand(diffRoleFilterCmd())
	cmd.AddCommand(diffRunbookFilterCmd())
	cmd.AddCommand(diffArgsCmd())

	return cmd
}
