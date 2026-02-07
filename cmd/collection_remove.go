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

	"github.com/spf13/cobra"
)

var collectionRemoveCmd = &cobra.Command{
	Use:   "remove <collection>",
	Short: "Remove Ansible collection from requirements.yml",
	Long:  "Remove Ansible collection from requirements.yml",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		var changed []Collection

		collection := args[0]
		requirements, _ := readRequirementsFile()

		for i, r := range requirements.Collections {
			if r.Name == collection {
				changed = append(requirements.Collections[:i], requirements.Collections[i+1:]...)
			}
		}

		if len(changed) > 0 {
			requirements.Collections = changed
			writeRequirementsFile(requirements)
			fmt.Println(Info("collection '%s' removed from requirements.yml", collection))
			return
		}

		fmt.Println(Warn("collection '%s' not present.", collection))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
	},
}

func init() {
	collectionCmd.AddCommand(collectionRemoveCmd)
}
