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

// compareCmd creates the Cobra command for "ansible-dev role compare",
// which compares each installed role in the development environment against
// its upstream source in the canonical Ansible roles repository.
//
// The command determines paths from two sources:
//   - The local roles directory is resolved via [ansible.RootRoleFolder]
//     (the "roles_path" setting in ansible.cfg) relative to the current
//     working directory.
//   - The upstream repository directory is read from the ANSIBLE_ROLES
//     environment variable. An error is returned if this variable is not
//     set.
//
// For each subdirectory in the local roles path, the command looks for a
// matching directory in ANSIBLE_ROLES (falling back to a name with the
// "dcjulian29." prefix stripped). If a match is found, it performs a
// file-by-file hash comparison, excluding .git, .github,
// .galaxy_install_info, and .ansible paths.
//
// When differences are detected the command opens a graphical diff tool:
//   - Windows: WinMerge (C:\Program Files\WinMerge\winmergeu.exe) with
//     recursive comparison and the "AnsibleRoles" file filter.
//   - Linux/macOS: Meld (/usr/bin/meld).
//
// Flags:
//   - --checksum:  print per-file hash comparisons to stdout. Matching
//     files are shown in green; differing files are shown in red
//     (default false).
//   - --no-diff:   skip launching the graphical diff tool even when
//     differences are detected (default false). Useful for CI or when
//     only checksum output is desired.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
func compareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare Installed Ansible roles with the development environment",
		RunE: func(cmd *cobra.Command, _ []string) error {
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

						if err := execute.ExternalProgram(program, params...); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	cmd.Flags().Bool("checksum", false, "show only file checksums")
	cmd.Flags().Bool("no-diff", false, "do not open diff tool to compare")

	return cmd
}
