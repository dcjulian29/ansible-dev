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
	"runtime"
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
		sep := string(os.PathSeparator)

		pwd, _ := os.Getwd()

		var userFolder string

		if runtime.GOOS == "windows" {
			userFolder = strings.Replace(os.Getenv("USERPROFILE"), "\\", sep, -1)
		} else {
			userFolder = strings.Replace(os.Getenv("HOME"), "\\", sep, -1)
		}

		repoFolder := strings.Replace(os.Getenv("ANSIBLE_ROLES"), "\\", sep, -1)
		workingFolder := strings.Replace(pwd+sep+rootRoleFolder(), "/./", sep, -1)
		workingFolder = strings.Replace(workingFolder, "\\./", sep, -1)

		if len(repoFolder) == 0 {
			fmt.Println(fmt.Errorf("the Ansible development role directory is not defined"))
			return
		}

		entries, err := os.ReadDir(workingFolder)

		if err != nil {
			fmt.Println(err)
			return
		}

		if len(entries) == 0 {
			fmt.Println(fmt.Errorf("No files found in '" + workingFolder + "'"))
		}

		for _, e := range entries {
			workingEntry := workingFolder + sep + e.Name()
			workingEntry = strings.Replace(workingEntry, "/./", sep, -1)
			workingEntry = strings.Replace(workingEntry, "\\./", sep, -1)

			if err != nil {
				fmt.Println("error in accessing working folder:", err)
			}

			repoEntry := strings.Replace(workingFolder, workingFolder, repoFolder+sep+e.Name(), 1)

			if !dirExists(repoEntry) {
				repoEntry = strings.Replace(repoEntry, "dcjulian29.", "", 1)

				if !dirExists(repoEntry) {
					repoEntry = ""
				}
			}

			if len(repoEntry) > 0 {
				source := strings.Replace(workingEntry, userFolder, "~", 1)
				dest := strings.Replace(repoEntry, userFolder, "~", 1)
				ignored := []string{"\\.git", "\\.github", ".galaxy_install_info"}

				fmt.Println("'" + source + "' --> '" + dest + "'")

				_, workingFile := scanDirectory(workingEntry, ignored)
				_, repoFile := scanDirectory(repoEntry, ignored)
				compare := false

				if len(workingFile) != len(repoFile) {
					compare = true
				}

				for _, f := range workingFile {
					f2 := strings.Replace(f, workingEntry, repoEntry, 1)

					var h1, h2 string

					if fileExists(f) {
						h1 = fileHash(f)
					} else {
						if checksum {
							fmt.Printf(Red("%s: ==> %s\n"), strings.Replace(f, workingEntry+sep, "", 1))
						}
					}

					if fileExists(f2) {
						h2 = fileHash(f2)
					} else {
						if checksum {
							fmt.Printf(Red("%s: <== %s\n"), strings.Replace(f2, workingEntry+sep, "", 1))
						}
					}

					if h1 != h2 {
						compare = true
					}

					if checksum {
						if h1 == h2 {
							fmt.Printf(Green("%s: %s == %s\n"), strings.Replace(f, workingEntry+sep, "", 1), h1, h2)
						} else {
							fmt.Printf(Red("%s: %s != %s\n"), strings.Replace(f, workingEntry+sep, "", 1), h1, h2)
						}
					}
				}

				if compare && !nodiff {
					var program string
					var params []string

					if runtime.GOOS == "windows" {
						program = "C:\\Program Files\\WinMerge\\winmergeu.exe"
						params = []string{"/r", "/m", "Full", "/u", "/f", "AnsibleRoles", repoEntry, workingEntry}
					} else {
						program = "/usr/bin/meld"
						params = []string{repoEntry, workingEntry}
					}

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
