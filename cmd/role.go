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
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Roles struct {
	Roles []Role `yaml:"roles"`
}

type Role struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"src"`
	Version string `yaml:"version"`
}

var (
	verbose bool

	roleCmd = &cobra.Command{
		Use:   "role",
		Short: "Manage ansible role",
		Long: `Manage ansible role requirements by allowing additional roles to be
added, created, deleted, and restored.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	rootCmd.AddCommand(roleCmd)
}

func ensureRequirementsFile() {
	if !fileExists("requirements.yml") {
		fmt.Println("ERROR: Requirements file is not present!")
		ensureWorkingDirectoryAndExit()
	}
}

func readRequirementsFile() []Role {
	ensureRequirementsFile()

	file, err := os.Open("requirements.yml")
	if err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	var roles Roles
	if err := yaml.Unmarshal(data, &roles); err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	return roles.Roles
}

func writeRequirementsFile(roles []Role) {
	file, err := os.OpenFile("requirements.yml", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	defer file.Close()

	var roleChild Roles
	roleChild.Roles = roles

	data, err := yaml.Marshal(roleChild)

	if err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	if _, err := file.Write(data); err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}
}
