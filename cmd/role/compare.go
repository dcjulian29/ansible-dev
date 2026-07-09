/*
Copyright © 2026 Julian Easterling

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
	"fmt"
	"os"
	"strings"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/settings"
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
//   - The upstream repository directory is the configured roles_path
//     (see "ansible-dev config roles-path"). An error is returned if it is
//     not set.
//
// For each subdirectory in the local roles path, the command looks for a
// matching directory under the configured roles path (falling back to a name
// with the "dcjulian29." prefix stripped). If a match is found, it performs a
// file-by-file hash comparison, excluding the configured role_ignore
// substrings (nothing is excluded when that list is empty).
//
// When differences are detected the command opens the diff tool configured for
// the current operating system (see "ansible-dev config diff-program"),
// substituting the role_filter into its argument template. Unless --no-diff is
// given, it is an error when no diff program is configured for this OS.
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

			repoPath, err := settings.RolesPath()
			if err != nil {
				return err
			}

			repoFolder := strings.ReplaceAll(repoPath, "\\", sep)

			ignored, err := settings.RoleIgnore()
			if err != nil {
				return err
			}

			var launch func(left, right string) error

			if !nodiff {
				diff, err := settings.Diff()
				if err != nil {
					return err
				}

				launch = func(left, right string) error {
					program, args := diff.Command(diff.RoleFilter, left, right)

					return execute.ExternalProgram(program, args...)
				}
			}

			folder, err := ansible.RootRoleFolder()
			if err != nil {
				return err
			}

			workingFolder := strings.ReplaceAll(pwd+sep+folder, "/./", sep)
			workingFolder = strings.ReplaceAll(workingFolder, "\\./", sep)

			entries, err := os.ReadDir(workingFolder)
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				return fmt.Errorf("no files found in '%s'", workingFolder)
			}

			home := ansible.HomeFolder()

			for _, e := range entries {
				workingEntry := workingFolder + sep + e.Name()
				workingEntry = strings.ReplaceAll(workingEntry, "/./", sep)
				workingEntry = strings.ReplaceAll(workingEntry, "\\./", sep)
				repoEntry := repoFolder + sep + e.Name()

				if !filesystem.DirectoryExist(repoEntry) {
					repoEntry = strings.Replace(repoEntry, "dcjulian29.", "", 1)

					if !filesystem.DirectoryExist(repoEntry) {
						continue
					}
				}

				if _, err := ansible.ComparePair(
					workingEntry, repoEntry, ignored, checksum, nodiff, launch, home,
				); err != nil {
					return err
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
