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

package ansible

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// HomeFolder returns the current user's home directory (USERPROFILE on Windows,
// HOME elsewhere) with backslashes normalized to the platform separator. It is
// used to abbreviate printed paths to "~".
func HomeFolder() string {
	sep := string(os.PathSeparator)

	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(os.Getenv("USERPROFILE"), "\\", sep)
	}

	return strings.ReplaceAll(os.Getenv("HOME"), "\\", sep)
}

// ComparePair compares the files under primaryDir against their counterparts
// under secondaryDir by content hash. This is the shared engine behind both
// "role compare" and "runbook compare"; the callers differ only in how they
// pair an installed directory with its canonical source.
//
// Any path containing one of the ignore substrings is skipped on both sides.
// When checksum is true, a per-file hash line is printed — green when the two
// copies match, red when they differ. When a difference is detected and noDiff
// is false, a graphical diff tool is launched (WinMerge on Windows using the
// diffFilter file filter, Meld elsewhere). When homeFolder is non-empty it is
// abbreviated to "~" in the printed header.
//
// It returns true when any difference (file count or content) was found.
func ComparePair(
	primaryDir, secondaryDir string,
	ignore []string,
	checksum, noDiff bool,
	diffFilter, homeFolder string,
) (bool, error) {
	sep := string(os.PathSeparator)

	header := func(p string) string {
		if len(homeFolder) > 0 {
			return strings.Replace(p, homeFolder, "~", 1)
		}

		return p
	}

	fmt.Printf("'%s' --> '%s'\n", header(primaryDir), header(secondaryDir))

	_, primaryFiles, err := filesystem.ScanDirectory(primaryDir, ignore)
	if err != nil {
		return false, err
	}

	_, secondaryFiles, err := filesystem.ScanDirectory(secondaryDir, ignore)
	if err != nil {
		return false, err
	}

	differ := len(primaryFiles) != len(secondaryFiles)

	for _, f := range primaryFiles {
		other := strings.Replace(f, primaryDir, secondaryDir, 1)

		var h1, h2 string

		if filesystem.FileExist(f) {
			if h1, err = filesystem.FileHash(f); err != nil {
				return false, err
			}
		}

		if filesystem.FileExist(other) {
			if h2, err = filesystem.FileHash(other); err != nil {
				return false, err
			}
		}

		if h1 != h2 {
			differ = true
		}

		if checksum {
			name := strings.Replace(f, primaryDir+sep, "", 1)
			if h1 == h2 {
				fmt.Println(textformat.Green(fmt.Sprintf("%s: %s == %s", name, h1, h2)))
			} else {
				fmt.Println(textformat.Red(fmt.Sprintf("%s: %s != %s", name, h1, h2)))
			}
		}
	}

	if differ && !noDiff {
		if err := launchDiff(primaryDir, secondaryDir, diffFilter); err != nil {
			return differ, err
		}
	}

	return differ, nil
}

// launchDiff opens a graphical diff between the canonical source (secondaryDir,
// shown on the left) and the installed copy (primaryDir, shown on the right).
// On Windows it uses WinMerge with the named file filter; elsewhere it uses
// Meld.
func launchDiff(primaryDir, secondaryDir, filter string) error {
	var (
		program string
		params  []string
	)

	if runtime.GOOS == "windows" {
		program = "C:\\Program Files\\WinMerge\\winmergeu.exe"
		params = []string{"/r", "/m", "Full", "/u", "/f", filter, secondaryDir, primaryDir}
	} else {
		program = "/usr/bin/meld"
		params = []string{secondaryDir, primaryDir}
	}

	return execute.ExternalProgram(program, params...)
}
