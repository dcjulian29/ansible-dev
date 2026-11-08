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
	"strconv"
	"strings"
)

// Platform is one entry under galaxy_info.platforms in meta/main.yml: a
// platform name and the versions of it the role supports.
type Platform struct {
	Name     string
	Versions []string
}

// platformNames are the platform names the ansible-lint meta schema accepts
// (src/ansiblelint/schemas/meta.json, the "const" of each *PlatformModel). A
// name outside this list fails schema[meta]. Copied from ansible-lint main on
// 2026-11-07; the list rarely changes, and checking against it here turns a
// typo into an error before the role exists rather than a lint failure after.
//
// Versions are deliberately not checked: each platform has its own list in
// the schema, and ansible-lint already reports a bad one.
var platformNames = []string{
	"AIX", "AlmaLinux", "Alpine", "Amazon", "Amazon Linux", "ArchLinux",
	"Astra Linux", "ClearLinux", "Cumulus", "Debian", "DellOS", "Devuan",
	"DragonFlyBSD", "EL", "Fedora", "FreeBSD", "GenericBSD", "GenericLinux",
	"GenericUNIX", "Gentoo", "HardenedBSD", "IOS", "Junos", "Kali", "MacOSX",
	"Mageia", "NXOS", "NetBSD", "OpenBSD", "OpenWrt", "OracleLinux", "PAN-OS",
	"Rocky", "SLES", "SmartOS", "Solaris", "Synology", "TMOS", "Ubuntu",
	"Void Linux", "Windows", "aos", "eos", "macOS", "opensuse", "os10",
	"vCenter", "vSphere",
}

// DefaultPlatforms returns the platforms a new role declares when none are
// given: Debian 13 (trixie) and EL 10, the pair roles are developed and
// tested on.
func DefaultPlatforms() []Platform {
	return []Platform{
		{Name: "Debian", Versions: []string{"trixie"}},
		{Name: "EL", Versions: []string{"10"}},
	}
}

// ParsePlatforms turns --platform values of the form NAME:VERSION[,VERSION...]
// (for example "Debian:trixie,bookworm") into platforms, in the order given.
// No values means [DefaultPlatforms].
//
// A platform named twice is an error rather than merged: its versions belong
// together in one value, and silently combining them would hide a typo such
// as "EL:9" given where "EL:10" was meant.
func ParsePlatforms(values []string) ([]Platform, error) {
	if len(values) == 0 {
		return DefaultPlatforms(), nil
	}

	platforms := make([]Platform, 0, len(values))
	seen := make(map[string]bool)

	for _, value := range values {
		platform, err := parsePlatform(value)
		if err != nil {
			return nil, err
		}

		if seen[platform.Name] {
			return nil, fmt.Errorf(
				"platform '%s' is given more than once; list its versions together, as %s:a,b",
				platform.Name, platform.Name)
		}

		seen[platform.Name] = true
		platforms = append(platforms, platform)
	}

	return platforms, nil
}

// parsePlatform parses one NAME:VERSION[,VERSION...] value.
func parsePlatform(value string) (Platform, error) {
	name, list, found := strings.Cut(value, ":")
	if !found {
		return Platform{}, fmt.Errorf(
			"platform '%s' must be NAME:VERSION[,VERSION...], such as Debian:trixie", value)
	}

	canonical, err := platformName(strings.TrimSpace(name))
	if err != nil {
		return Platform{}, err
	}

	var versions []string

	for _, version := range strings.Split(list, ",") {
		version = strings.TrimSpace(version)
		if version == "" {
			return Platform{}, fmt.Errorf("platform '%s' has an empty version", value)
		}

		versions = append(versions, version)
	}

	return Platform{Name: canonical, Versions: versions}, nil
}

// platformName matches name against [platformNames] ignoring case and returns
// the schema's spelling, so "debian" or "el" is written as "Debian" or "EL".
func platformName(name string) (string, error) {
	for _, known := range platformNames {
		if strings.EqualFold(name, known) {
			return known, nil
		}
	}

	return "", fmt.Errorf("'%s' is not a Galaxy platform name; ansible-lint accepts: %s",
		name, strings.Join(platformNames, ", "))
}

// renderPlatforms writes platforms as the YAML list that follows "platforms:"
// in meta/main.yml, indented to sit under galaxy_info, without a trailing
// newline (the template supplies it).
//
// Every version is quoted. The schema only accepts strings, and an unquoted
// EL version such as 10 is read as a number and fails schema[meta]; quoting
// all of them means no version can trip that, whatever the platform.
func renderPlatforms(platforms []Platform) string {
	var b strings.Builder

	for i, platform := range platforms {
		if i > 0 {
			b.WriteString("\n")
		}

		fmt.Fprintf(&b, "    - name: %s\n      versions:", platform.Name)

		for _, version := range platform.Versions {
			fmt.Fprintf(&b, "\n        - %s", strconv.Quote(version))
		}
	}

	return b.String()
}
