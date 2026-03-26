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

	"github.com/dcjulian29/go-toolbox/color"
)

// ApplyRoles iterates over the supplied role names, generates a temporary
// single-role playbook for each, and executes it via ansible-playbook. The
// optional tags slice limits task execution and verbose enables increased
// output. Execution stops on the first error encountered.
func ApplyRoles(roles []string, tags []string, verbose bool) error {
	play := Play{
		Tags:       tags,
		FlushCache: true,
		Verbose:    verbose,
	}

	if len(roles) > 0 {
		for _, role := range roles {
			fmt.Println(color.Info(fmt.Sprintf("\nApplying the '%s' role...", role)))

			err := GenerateRolePlay(role)
			if err != nil {
				return err
			}

			play.Name = role

			err = ExecutePlay(play)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
