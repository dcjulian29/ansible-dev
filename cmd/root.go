/*
Copyright © 2026 Julian Easterling julian@julianscorner.com

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

// Package cmd is the root command package for the ansible-dev CLI. It
// wires together every subcommand, configures the Cobra root command,
// and provides the [Execute] entry point invoked by main.
//
// Subcommands are organized into individual packages under cmd/ and
// registered during init. They cover the full lifecycle of an Ansible
// development environment:
//
//   - collection: manage Ansible collections in requirements.yml.
//   - destroy:    tear down the Vagrant environment.
//   - initialize: scaffold a new Ansible project.
//   - inventory:  display the host inventory.
//   - ping:       verify host reachability.
//   - play:       provision roles against Vagrant hosts.
//   - reset:      reset the development environment.
//   - restore:    install dependencies from requirements.yml.
//   - role:       manage Ansible roles (add, compare, delete, list, new, remove).
//   - runbook:    execute the project's runbook playbook.
//   - shell:      run ad-hoc shell commands on all hosts.
//   - start/up:   boot and optionally provision VMs.
//   - status:     show Vagrant VM state.
//   - stop/down:  gracefully halt VMs.
//   - tag:        list tags defined in a role.
//   - task:       list tasks that would execute for a role.
//   - upgrade:    update and prune Vagrant boxes.
package cmd

import (
	"fmt"
	"os"

	"github.com/dcjulian29/ansible-dev/cmd/collection"
	"github.com/dcjulian29/ansible-dev/cmd/destroy"
	"github.com/dcjulian29/ansible-dev/cmd/initialize"
	"github.com/dcjulian29/ansible-dev/cmd/inventory"
	"github.com/dcjulian29/ansible-dev/cmd/ping"
	"github.com/dcjulian29/ansible-dev/cmd/play"
	"github.com/dcjulian29/ansible-dev/cmd/reset"
	"github.com/dcjulian29/ansible-dev/cmd/restore"
	"github.com/dcjulian29/ansible-dev/cmd/role"
	"github.com/dcjulian29/ansible-dev/cmd/runbook"
	"github.com/dcjulian29/ansible-dev/cmd/shell"
	"github.com/dcjulian29/ansible-dev/cmd/start"
	"github.com/dcjulian29/ansible-dev/cmd/status"
	"github.com/dcjulian29/ansible-dev/cmd/stop"
	"github.com/dcjulian29/ansible-dev/cmd/tag"
	"github.com/dcjulian29/ansible-dev/cmd/task"
	"github.com/dcjulian29/ansible-dev/cmd/upgrade"
	"github.com/dcjulian29/go-toolbox/color"
	"github.com/spf13/cobra"
	"go.szostok.io/version/extension"
)

// rootCmd is the top-level Cobra command for the ansible-dev CLI. When
// invoked without a subcommand it prints the help text. Both
// SilenceErrors and SilenceUsage are enabled so that error formatting
// is handled exclusively by [Execute].
var rootCmd = &cobra.Command{
	Use:   "ansible-dev",
	Short: "ansible-dev enables development of Ansible playbooks, roles, and runbooks.",
	Long: `ansible-dev integrates with Vagrant to enable developers to define, develop, and test Ansible
playbooks, roles, and runbooks.

By utilizing Ansible, developers/operators can automate the deployment of software applications
across multiple hosting providers, reducing the time and effort required to manage complex
infrastructure environments.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		return nil
	},
}

// Execute is the main entry point for the CLI, called from main.main.
// It adds the auto-generated "version" subcommand (provided by
// go.szostok.io/version) with an upgrade notice for the
// "dcjulian29/ansible-dev" GitHub repository, then delegates to
// [cobra.Command.Execute].
//
// If execution returns an error, the error message is printed to stderr
// using [color.Fatal] formatting and the process exits with code 1.
func Execute() {
	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("dcjulian29", "ansible-dev"),
		),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "\n"+color.Fatal(err.Error()))
		os.Exit(1)
	}
}

// init registers all subcommands on the root command during package
// initialization. Each subcommand is provided by a dedicated package
// under cmd/ and exposes a NewCommand factory function that returns a
// configured [cobra.Command].
func init() {
	rootCmd.AddCommand(collection.NewCommand())
	rootCmd.AddCommand(destroy.NewCommand())
	rootCmd.AddCommand(initialize.NewCommand())
	rootCmd.AddCommand(inventory.NewCommand())
	rootCmd.AddCommand(ping.NewCommand())
	rootCmd.AddCommand(play.NewCommand())
	rootCmd.AddCommand(reset.NewCommand())
	rootCmd.AddCommand(restore.NewCommand())
	rootCmd.AddCommand(role.NewCommand())
	rootCmd.AddCommand(runbook.NewCommand())
	rootCmd.AddCommand(shell.NewCommand())
	rootCmd.AddCommand(start.NewCommand())
	rootCmd.AddCommand(status.NewCommand())
	rootCmd.AddCommand(stop.NewCommand())
	rootCmd.AddCommand(tag.NewCommand())
	rootCmd.AddCommand(task.NewCommand())
	rootCmd.AddCommand(upgrade.NewCommand())
}
