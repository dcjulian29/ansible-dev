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

// Role describes a single Ansible role dependency declared in the
// requirements.yml file.
//
// Fields:
//   - Name:    the role name as it appears in the Galaxy namespace or local path.
//   - Source:  an optional URL or Galaxy reference where the role is hosted.
//   - Version: an optional version constraint string (e.g. "v1.2.0").
type Role struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"src"`
	Version string `yaml:"version"`
}
