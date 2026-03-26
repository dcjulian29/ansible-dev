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
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// ReadRequirements reads and parses the requirements.yml file in the current
// working directory into a [Requirements] struct. It first calls
// [EnsureRequirementsFile] to confirm the file exists, then deserializes the
// YAML content.
//
// An error is returned if the file is missing, cannot be opened, cannot be
// read, or contains invalid YAML that does not conform to the
// [Requirements] schema.
func ReadRequirements() (Requirements, error) {
	err := EnsureRequirementsFile()
	if err != nil {
		return Requirements{}, err
	}

	var requirements Requirements

	file, err := os.Open("requirements.yml")
	if err != nil {
		return Requirements{}, err
	}

	defer file.Close() //nolint:errcheck

	data, err := io.ReadAll(file)
	if err != nil {
		return Requirements{}, err
	}

	err = yaml.Unmarshal(data, &requirements)
	if err != nil {
		return Requirements{}, err
	}

	return requirements, nil
}
