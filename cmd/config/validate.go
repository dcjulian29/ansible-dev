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

package config

import (
	"fmt"
	"path/filepath"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// warnAboutDirectory reports problems with a directory a path setting points
// at, naming the setting as key. It warns rather than fails: the value is still
// saved, because a path may legitimately be configured on a machine before the
// directory is created, or written into a dotfile shared across machines where
// it does not exist everywhere.
//
// Without this the mistake is silent — a typo or a relative path persists
// happily and only surfaces much later as a confusing failure from
// "role compare" or "runbook compare".
//
// A "~"-rooted path is checked in its expanded form and counts as absolute: it
// is the portable way to write a home-relative directory in a dotfile shared
// between machines, so it must not be reported as a relative-path mistake.
func warnAboutDirectory(key, path string) {
	expanded := filesystem.ExpandHome(path)

	if !filepath.IsAbs(expanded) {
		fmt.Println(textformat.Warn(fmt.Sprintf(
			"%s '%s' is not an absolute path; it will be resolved against the "+
				"working directory of whichever command reads it", key, path)))
	}

	if !filesystem.DirectoryExist(expanded) {
		fmt.Println(textformat.Warn(fmt.Sprintf(
			"%s '%s' does not exist or is not a directory", key, path)))
	}
}
