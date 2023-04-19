/*
Copyright © 2023 Julian Easterling <julian@julianscorner.com>

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
	"strings"

	"github.com/spf13/cobra"
)

var (
	tasksCmd = &cobra.Command{
		Use:   "tasks <role>",
		Short: "List all tasks that would be executed in the role",
		Long:  "List all tasks that would be executed in the role",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}

			generate_play(args[0])

			var param []string

			tags, _ := cmd.Flags().GetStringSlice("tags")

			if len(tags) > 0 {
				param = append(param, fmt.Sprintf("--tags %s", strings.Join(tags, ",")))
			}

			param = append(param, "--list-tasks", ".tmp/play.yml")

			executeExternalProgram("ansible-playbook", param...)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			ensureAnsibleDirectory()
			ensureVagrantfile()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			ensureWorkingDirectoryAndExit()
		},
	}
)

func init() {
	rootCmd.AddCommand(tasksCmd)

	tasksCmd.Flags().StringSlice("tags", []string{}, "only plays and task tagged with these values")
}
