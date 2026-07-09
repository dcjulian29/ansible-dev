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

package settings

import (
	"strings"
	"testing"
)

func TestCommand_SubstitutesFilterLeftRight(t *testing.T) {
	d := DiffTool{
		Program:        "winmergeu.exe",
		AdditionalArgs: []string{"/r", "/f", "{filter}", "{left}", "{right}"},
	}

	program, args := d.Command("AnsibleRoles", "SRC", "DST")

	if program != "winmergeu.exe" {
		t.Errorf("program = %q, want %q", program, "winmergeu.exe")
	}

	got := strings.Join(args, " ")
	want := "/r /f AnsibleRoles SRC DST"

	if got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
}

func TestCommand_EmptyFilterDropsFlagAndPlaceholder(t *testing.T) {
	d := DiffTool{
		Program:        "winmergeu.exe",
		AdditionalArgs: []string{"/r", "/f", "{filter}", "{left}", "{right}"},
	}

	_, args := d.Command("", "SRC", "DST")

	got := strings.Join(args, " ")
	want := "/r SRC DST" // both "/f" and "{filter}" dropped

	if got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
}

func TestCommand_EmptyFilterDropsEmbeddedToken(t *testing.T) {
	d := DiffTool{
		AdditionalArgs: []string{"--filter={filter}", "{left}", "{right}"},
	}

	_, args := d.Command("", "SRC", "DST")

	got := strings.Join(args, " ")
	want := "SRC DST" // "--filter={filter}" dropped whole, no preceding flag removed

	if got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
}

func TestCommand_NoPlaceholdersPassThrough(t *testing.T) {
	d := DiffTool{AdditionalArgs: []string{"{left}", "{right}"}}

	_, args := d.Command("", "SRC", "DST")

	if strings.Join(args, " ") != "SRC DST" {
		t.Errorf("args = %q, want %q", args, "SRC DST")
	}
}

func TestSetters_PopulateCurrentOSEntry(t *testing.T) {
	var cfg Config

	cfg.SetDiffProgram("prog")
	cfg.SetRoleDiffFilter("roles")
	cfg.SetRunbookDiffFilter("runbooks")
	cfg.SetDiffAdditionalArgs([]string{"a", "b"})

	got := cfg.CurrentDiff() // reads the same runtime.GOOS key the setters wrote

	if got.Program != "prog" || got.RoleFilter != "roles" || got.RunbookFilter != "runbooks" {
		t.Errorf("CurrentDiff() = %+v, want program/roles/runbooks set", got)
	}

	if strings.Join(got.AdditionalArgs, ",") != "a,b" {
		t.Errorf("AdditionalArgs = %v, want [a b]", got.AdditionalArgs)
	}
}
