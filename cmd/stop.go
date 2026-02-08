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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var (
	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stops the Ansible development vagrant environment",
		Long:  "Stops the Ansible development vagrant environment",
		Run: func(cmd *cobra.Command, args []string) {
			vagrant_down()
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			ensureAnsibleDirectory()
			ensureVagrantfile()
		},
	}
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

func vagrant_down() {
	inv, err := ini.Load("hosts.ini")
	cobra.CheckErr(err)

	section, err := inv.GetSection("vagrant")
	cobra.CheckErr(err)

	for _, vm := range section.KeyStrings() {
		name := strings.Split(vm, " ")[0]

		fmt.Printf(Yellow("\nStopping '%s'...\n\n"), name)

		executeExternalProgram("vagrant", "halt", name)
	}

}
