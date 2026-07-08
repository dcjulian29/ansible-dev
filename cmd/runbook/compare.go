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

package runbook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// compareCmd creates the Cobra command for "ansible-dev runbook compare", the
// runbook analogue of "role compare". Where role compare walks the local roles
// directory, runbook compare is driven from the source side: it walks each
// runbook repository under the configured runbooks_path (see "ansible-dev
// config runbooks-path") and compares it against the collection installed in
// the current project.
//
// For each runbook directory it reads galaxy.yml to obtain the namespace and
// name (so no namespace is hard-coded), locates the installed collection at
// <collections_path>/ansible_collections/<namespace>/<name>, and delegates the
// file-by-file hash comparison to [ansible.ComparePair].
//
// The ignore set is the configured runbook_ignore list (nothing is excluded
// when it is empty). To keep the checksum comparison and the visual diff in
// agreement, configure it to match the diff tool's filter — for example the
// source repo's SCM/per-repo files (galaxy.yml, README.md, .devcontainer) and
// the installed copy's runtime artifacts (MANIFEST.json, FILES.json). Runbooks
// present under runbooks_path but not installed are reported and skipped.
//
// Flags:
//   - --checksum: print per-file hash comparisons.
//   - --no-diff:  do not launch the graphical diff tool on differences.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the current
// directory is a valid Ansible project.
func compareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare installed runbook collections with their source repositories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			checksum, _ := cmd.Flags().GetBool("checksum")
			nodiff, _ := cmd.Flags().GetBool("no-diff")
			sep := string(os.PathSeparator)
			pwd, _ := os.Getwd()

			runbooksPath, err := settings.RunbooksPath()
			if err != nil {
				return err
			}

			runbooksFolder := strings.ReplaceAll(runbooksPath, "\\", sep)

			ignored, err := settings.RunbookIgnore()
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
					program, args := diff.Command(diff.RunbookFilter, left, right)

					return execute.ExternalProgram(program, args...)
				}
			}

			collections, err := ansible.CollectionsFolder()
			if err != nil {
				return err
			}

			collectionsFolder := strings.ReplaceAll(pwd+sep+collections, "/./", sep)
			collectionsFolder = strings.ReplaceAll(collectionsFolder, "\\./", sep)

			entries, err := os.ReadDir(runbooksFolder)
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				return fmt.Errorf("no files found in '%s'", runbooksFolder)
			}

			home := ansible.HomeFolder()

			for _, e := range entries {
				if !e.IsDir() {
					continue
				}

				sourceEntry := runbooksFolder + sep + e.Name()

				info, err := ansible.ReadGalaxyInfo(sourceEntry)
				if err != nil {
					fmt.Println(textformat.Yellow(fmt.Sprintf("skipping '%s': %s", e.Name(), err)))
					continue
				}

				installedEntry := filepath.Join(collectionsFolder, info.Namespace, info.Name)

				if !filesystem.DirectoryExist(installedEntry) {
					fmt.Println(textformat.Yellow(fmt.Sprintf(
						"runbook '%s.%s' is not installed at '%s'", info.Namespace, info.Name, installedEntry)))

					continue
				}

				if _, err := ansible.ComparePair(
					installedEntry, sourceEntry, ignored, checksum, nodiff, launch, home,
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
