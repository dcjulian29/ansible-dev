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
	"path/filepath"

	"gopkg.in/ini.v1"
)

// CollectionsFolder returns the base directory under which Ansible collections
// are installed. It reads the "collections_path" key from the [defaults]
// section of the local ansible.cfg and appends "ansible_collections", yielding
// the standard install layout in which a collection lives at
// <collections_path>/ansible_collections/<namespace>/<name>.
//
// An error is returned if ansible.cfg cannot be loaded, the [defaults] section
// is missing, or the "collections_path" key is not defined.
func CollectionsFolder() (string, error) {
	cfg, err := ini.Load("ansible.cfg")
	if err != nil {
		return "", err
	}

	section, err := cfg.GetSection("defaults")
	if err != nil {
		return "", err
	}

	path, err := section.GetKey("collections_path")
	if err != nil {
		return "", err
	}

	return filepath.Join(path.String(), "ansible_collections"), nil
}
