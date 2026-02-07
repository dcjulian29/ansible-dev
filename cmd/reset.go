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
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Reset the Ansible development vagrant environment",
		Long:  "Reset the Ansible development vagrant environment",
		Run: func(cmd *cobra.Command, args []string) {
			if r, _ := cmd.Flags().GetBool("no-recreate"); !r {
				vagrant_destroy()
			}

			vagrant_up(cmd)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			ensureAnsibleDirectory()
			ensureVagrantfile()
		},
	}
)

func init() {
	rootCmd.AddCommand(resetCmd)

	resetCmd.Flags().Bool("base", true, "provision the VMs with the base role minimal tag")
	resetCmd.Flags().String("role", "", "provision the VMs with the specified role")
	resetCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	resetCmd.Flags().Bool("no-recreate", false, "do not recreate the VMs")

	resetCmd.MarkFlagsMutuallyExclusive("development", "test")
}
