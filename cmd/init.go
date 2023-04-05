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

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an development environment for Ansible development",
	Long: `Initialize an development environment for Ansible development by creating the folder
structure and generating the needed files to quickly set up a virtual environment
ready for development. Vagrant can be used to manage the environment and connect
to troubleshoot and/or validate.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing development environment...")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("force", "f", false, "Overwrite an existing development environment")
}
