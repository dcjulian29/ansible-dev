/*
Copyright © 2026 Julian Easterling <julian@julianscorner.com>

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

// Collection describes a single Ansible Galaxy collection dependency
// declared in the requirements.yml file.
//
// Fields:
//   - Name:    the fully qualified collection name (e.g. "community.general").
//   - Source:  an optional URL or Galaxy server where the collection is hosted.
//   - Type:    the source type (e.g. "galaxy", "git", "url").
//   - Version: an optional version constraint string (e.g. ">=2.0.0").
type Collection struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"src"`
	Type    string `yaml:"type"`
	Version string `yaml:"version"`
}
