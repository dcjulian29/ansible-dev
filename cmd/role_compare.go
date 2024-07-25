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
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var roleCompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare Installed Ansible roles with the development environment",
	Long:  "Compare Installed Ansible roles with the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		checksum, _ := cmd.Flags().GetBool("checksum")
		nodiff, _ := cmd.Flags().GetBool("no-diff")

		pwd, _ := os.Getwd()

		userFolder := strings.Replace(os.Getenv("USERPROFILE"), "\\", "\\", -1)
		repoFolder := strings.Replace(os.Getenv("ANSIBLE_ROLES"), "\\", "\\", -1)
		workingFolder := strings.Replace(pwd+"\\"+rootRoleFolder(), "\\./", "\\", -1)

		if len(repoFolder) == 0 {
			fmt.Println(fmt.Errorf("the Ansible development role directory is not defined"))
			return
		}

		entries, err := os.ReadDir(workingFolder)

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, e := range entries {
			workingEntry := workingFolder + "\\" + e.Name()

			if err != nil {
				fmt.Println("error in accessing working folder:", err)
			}

			repoEntry := strings.Replace(workingEntry, workingFolder, repoFolder, 1)

			if !dirExists(repoEntry) {
				repoEntry = strings.Replace(repoEntry, "dcjulian29.", "", 1)

				if !dirExists(repoEntry) {
					repoEntry = ""
				}
			}

			if len(repoEntry) > 0 {
				source := strings.Replace(workingEntry, userFolder, "~", 1)
				dest := strings.Replace(repoEntry, userFolder, "~", 1)

				fmt.Println("'" + source + "' --> '" + dest + "'")

				ignored := []string{"\\.git\\", "\\.github\\"}

				_, workingFile := scanDirectory(workingEntry, ignored)
				_, repoFile := scanDirectory(repoEntry, ignored)
				compare := false

				if len(workingFile) != len(repoFile) {
					compare = true
				}

				for _, f := range workingFile {
					if strings.Contains(f, ".galaxy_install_info") {
						continue
					}

					f2 := strings.Replace(f, workingEntry, repoEntry, 1)
					h1 := fileHash(f)
					h2 := fileHash(f2)

					if h1 != h2 {
						compare = true
					}

					if checksum {
						if h1 == h2 {
							fmt.Printf(Green("%s: %s == %s\n"), strings.Replace(f, workingEntry+"\\", "", 1), h1, h2)
						} else {
							fmt.Printf(Red("%s: %s != %s\n"), strings.Replace(f, workingEntry+"\\", "", 1), h1, h2)
						}
					}
				}

				if compare && !nodiff {
					program := "C:\\Program Files\\WinMerge\\winmergeu.exe"
					params := []string{"/r", "/m", "Full", "/u", "/f", "AnsibleRoles", repoEntry, workingEntry}

					executeExternalProgram(program, params...)
				}
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		ensureWorkingDirectoryAndExit()
	},
}

func init() {
	roleCmd.AddCommand(roleCompareCmd)

	roleCompareCmd.Flags().Bool("checksum", false, "show only file checksums")
	roleCompareCmd.Flags().Bool("no-diff", false, "do not open diff tool to compare")
}
