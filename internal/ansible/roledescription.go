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

import "strings"

// RoleDescriptionPrefix is prepended to a role description so the text supplied
// with --description reads as a verb phrase completing the sentence: "install
// and configure nginx" becomes "Ansible role to install and configure nginx".
const RoleDescriptionPrefix = "Ansible role to "

// RoleDescription composes the description written into meta/main.yml, the
// generated README, and the published repository. Composing it once here is
// what keeps those three in agreement; each used to carry its own wording, so
// one role was described three slightly different ways.
//
// When prefix is false the description is used exactly as given, for a
// description that is already a complete sentence and would otherwise read as
// "Ansible role to Manages the nginx configuration".
//
// An empty description yields an empty string either way — a bare prefix
// describes nothing.
func RoleDescription(description string, prefix bool) string {
	description = strings.TrimSpace(description)

	if description == "" || !prefix {
		return description
	}

	return RoleDescriptionPrefix + description
}
