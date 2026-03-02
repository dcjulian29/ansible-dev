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
package role

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/color"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/spf13/cobra"
)

func compareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare Installed Ansible roles with the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			checksum, _ := cmd.Flags().GetBool("checksum")
			nodiff, _ := cmd.Flags().GetBool("no-diff")
			sep := string(os.PathSeparator)
			pwd, _ := os.Getwd()

			var userFolder string

			if runtime.GOOS == "windows" {
				userFolder = strings.ReplaceAll(os.Getenv("USERPROFILE"), "\\", sep)
			} else {
				userFolder = strings.ReplaceAll(os.Getenv("HOME"), "\\", sep)
			}

			repoFolder := strings.ReplaceAll(os.Getenv("ANSIBLE_ROLES"), "\\", sep)

			folder, err := ansible.RootRoleFolder()
			if err != nil {
				return err
			}

			workingFolder := strings.ReplaceAll(pwd+sep+folder, "/./", sep)
			workingFolder = strings.ReplaceAll(workingFolder, "\\./", sep)

			if len(repoFolder) == 0 {
				return errors.New("the Ansible roles directory is not defined in environment")
			}

			entries, err := os.ReadDir(workingFolder)
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				return fmt.Errorf("no files found in '%s'", workingFolder)
			}

			for _, e := range entries {
				workingEntry := workingFolder + sep + e.Name()
				workingEntry = strings.ReplaceAll(workingEntry, "/./", sep)
				workingEntry = strings.ReplaceAll(workingEntry, "\\./", sep)
				repoEntry := strings.Replace(workingFolder, workingFolder, repoFolder+sep+e.Name(), 1)

				if !filesystem.DirectoryExists(repoEntry) {
					repoEntry = strings.Replace(repoEntry, "dcjulian29.", "", 1)

					if !filesystem.DirectoryExists(repoEntry) {
						repoEntry = ""
					}
				}

				if len(repoEntry) > 0 {
					source := strings.Replace(workingEntry, userFolder, "~", 1)
					dest := strings.Replace(repoEntry, userFolder, "~", 1)
					ignored := []string{"\\.git", "\\.github", ".galaxy_install_info", ".ansible"}

					fmt.Println("'" + source + "' --> '" + dest + "'")

					_, workingFile, err := filesystem.ScanDirectory(workingEntry, ignored)
					if err != nil {
						return err
					}

					_, repoFile, err := filesystem.ScanDirectory(repoEntry, ignored)
					if err != nil {
						return err
					}

					compare := false

					if len(workingFile) != len(repoFile) {
						compare = true
					}

					for _, f := range workingFile {
						f2 := strings.Replace(f, workingEntry, repoEntry, 1)

						var h1, h2 string

						if filesystem.FileExists(f) {
							h1, err = filesystem.FileHash(f)
							if err != nil {
								return err
							}
						}

						if filesystem.FileExists(f2) {
							h2, err = filesystem.FileHash(f2)
							if err != nil {
								return err
							}
						}

						if h1 != h2 {
							compare = true
						}

						if checksum {
							file := strings.Replace(f, workingEntry+sep, "", 1)
							if h1 == h2 {
								fmt.Println(color.Green(fmt.Sprintf("%s: %s == %s", file, h1, h2)))
							} else {
								fmt.Println(color.Red(fmt.Sprintf("%s: %s != %s\n", file, h1, h2)))
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

						return execute.ExternalProgram(program, params...)
					}
				}
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().Bool("checksum", false, "show only file checksums")
	cmd.Flags().Bool("no-diff", false, "do not open diff tool to compare")

	return cmd
}
