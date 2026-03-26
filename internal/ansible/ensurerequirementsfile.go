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
	"fmt"
)

// EnsureRequirementsFile verifies that a requirements.yml file exists in the
// current working directory. It returns a non-nil error if the file is
// missing, making it suitable as a pre-flight check before operations that
// depend on requirements.yml.
func EnsureRequirementsFile() error {
	exist := RequirementsFileExist()
	if !exist {
		return fmt.Errorf("requirements.yml file is not present")
	}

	return nil
}
