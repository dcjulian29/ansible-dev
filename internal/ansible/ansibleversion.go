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
	"regexp"

	"github.com/dcjulian29/go-toolbox/execute"
)

// coreVersion matches the first line of "ansible-galaxy --version", for
// example "ansible-galaxy [core 2.21.5]". Only the numeric X.Y.Z is kept, so a
// pre-release such as "2.22.0rc1" records as "2.22.0".
var coreVersion = regexp.MustCompile(`\[core (\d+\.\d+\.\d+)`)

// AnsibleVersion returns the ansible-core version that ansible-galaxy reports,
// such as "2.21.5".
//
// It is what a new role records as min_ansible_version: the role is written
// and tested on this version, so it is the oldest one it can honestly claim.
// ansible-galaxy is asked rather than ansible because it is the program
// "role new" already depends on to scaffold the role, so if it can scaffold,
// it can answer.
func AnsibleVersion() (string, error) {
	output, err := execute.ExternalProgramCapture("ansible-galaxy", "--version")
	if err != nil {
		return "", fmt.Errorf("could not run ansible-galaxy to read the Ansible version: %w", err)
	}

	return parseAnsibleVersion(output)
}

// versionFormat is what a version given by hand must look like: "2.21" or
// "2.21.5". Anything else would land in meta/main.yml as written and only
// surface later as a lint or Galaxy failure.
var versionFormat = regexp.MustCompile(`^\d+\.\d+(\.\d+)?$`)

// ValidateAnsibleVersion rejects a hand-supplied min_ansible_version that is
// not a plain "X.Y" or "X.Y.Z" number, such as "v2.21" or "2.21.x".
func ValidateAnsibleVersion(version string) error {
	if !versionFormat.MatchString(version) {
		return fmt.Errorf("ansible version '%s' must look like 2.21 or 2.21.5", version)
	}

	return nil
}

// parseAnsibleVersion extracts the ansible-core version from the output of
// "ansible-galaxy --version". It is separate from [AnsibleVersion] so the
// parsing can be tested without Ansible installed.
func parseAnsibleVersion(output string) (string, error) {
	match := coreVersion.FindStringSubmatch(output)
	if match == nil {
		return "", fmt.Errorf("could not find the ansible-core version in %q", output)
	}

	return match[1], nil
}
