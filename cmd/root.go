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
