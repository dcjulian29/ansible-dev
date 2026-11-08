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
	"reflect"
	"testing"
)

func TestParsePlatforms_NoValuesMeansTheDefaults(t *testing.T) {
	got, err := ParsePlatforms(nil)
	if err != nil {
		t.Fatalf("ParsePlatforms(nil) error = %v, want nil", err)
	}

	if !reflect.DeepEqual(got, DefaultPlatforms()) {
		t.Errorf("ParsePlatforms(nil) = %v, want %v", got, DefaultPlatforms())
	}
}

func TestParsePlatforms_ParsesNamesAndVersions(t *testing.T) {
	got, err := ParsePlatforms([]string{"debian:trixie, bookworm", "EL:10", "Amazon Linux:2"})
	if err != nil {
		t.Fatalf("ParsePlatforms() error = %v, want nil", err)
	}

	// Names come back in the schema's spelling and versions are trimmed.
	want := []Platform{
		{Name: "Debian", Versions: []string{"trixie", "bookworm"}},
		{Name: "EL", Versions: []string{"10"}},
		{Name: "Amazon Linux", Versions: []string{"2"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParsePlatforms() = %v, want %v", got, want)
	}
}

func TestParsePlatforms_Rejects(t *testing.T) {
	for name, values := range map[string][]string{
		"no colon":        {"Debian"},
		"unknown name":    {"Debain:trixie"},
		"empty name":      {":trixie"},
		"empty version":   {"Debian:"},
		"trailing comma":  {"Debian:trixie,"},
		"platform twice":  {"EL:9", "el:10"},
		"one bad of many": {"Debian:trixie", "Rocky"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParsePlatforms(values); err == nil {
				t.Errorf("ParsePlatforms(%q) error = nil, want an error", values)
			}
		})
	}
}

// The rendered block has to line up under "platforms:" in meta/main.yml, and
// every version must be quoted so an EL version stays a string.
func TestRenderPlatforms(t *testing.T) {
	got := renderPlatforms(DefaultPlatforms())

	want := "    - name: Debian\n" +
		"      versions:\n" +
		"        - \"trixie\"\n" +
		"    - name: EL\n" +
		"      versions:\n" +
		"        - \"10\""

	if got != want {
		t.Errorf("renderPlatforms() =\n%s\nwant\n%s", got, want)
	}
}
