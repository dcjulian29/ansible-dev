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

import (
	"os"

	"gopkg.in/yaml.v3"
)

// SaveRequirements serializes the given [Requirements] struct to YAML and
// writes it to requirements.yml in the current working directory. If the
// file already exists its contents are truncated and overwritten; otherwise
// a new file is created with mode 0644.
//
// An error is returned if the file cannot be opened, the struct cannot be
// marshalled to YAML, or the write fails.
func SaveRequirements(requirements Requirements) error {
	file, err := os.OpenFile("requirements.yml", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close() //nolint:errcheck

	data, err := yaml.Marshal(requirements)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
