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
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var roleCmd = &cobra.Command{
	Use:     "role",
	Aliases: []string{"roles"},
	Short:   "Provide management of ansible roles in the development environment",
	Long:    "Provide management of ansible roles in the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(roleCmd)
}

func rootRoleFolder() string {
	ensureAnsibleDirectory()

	cfg, err := ini.Load("ansible.cfg")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	section, err := cfg.GetSection("defaults")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	rolePath, err := section.GetKey("roles_path")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return rolePath.String()
}

func roleFolder(role string) (string, error) {
	folder := filepath.Join(rootRoleFolder(), role)

	return folder, nil
}

func roleFolderExists(role string) bool {
	folder, err := roleFolder(role)

	if err != nil {
		return false
	}

	return dirExists(folder)
}

func remove_role(role string) {
	if !roleFolderExists(role) {
		fmt.Println(Warn("WARN: Role '%s' folder not present.", role))
		return
	}

	folder, err := roleFolder(role)
	if err != nil {
		fmt.Println(err)
		return
	}

	files, err := filepath.Glob(filepath.Join(folder, "*"))
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(files) == 0 {
		fmt.Println(Warn("WARN: Role '%s' files not present.", role))
		return
	}

	removeDir(folder)

	fmt.Println(Info("Role '%s' files were deleted.", role))
}
