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
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// GalaxyInfo holds the subset of an Ansible collection's galaxy.yml that is
// needed to locate a collection: its namespace and name, which together
// determine the install path <collections_path>/ansible_collections/
// <namespace>/<name>.
type GalaxyInfo struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
}

// ReadGalaxyInfo reads and parses the galaxy.yml file in dir into a
// [GalaxyInfo]. An error is returned if the file cannot be read or does not
// contain valid YAML.
func ReadGalaxyInfo(dir string) (GalaxyInfo, error) {
	var info GalaxyInfo

	data, err := os.ReadFile(filepath.Join(dir, "galaxy.yml"))
	if err != nil {
		return info, err
	}

	if err := yaml.Unmarshal(data, &info); err != nil {
		return info, err
	}

	return info, nil
}
