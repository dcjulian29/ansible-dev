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
package shell

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/vagrant"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/spf13/cobra"
)

var shellCommand string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell [flags] -- <command>",
		Short: "Execute shell command in the Ansible development vagrant environment",
		Args: func(cmd *cobra.Command, args []string) error {
			shellCommand = strings.Join(args, " ")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			param := []string{
				"-i", "hosts.ini",
				"-m", "shell",
				"-a", fmt.Sprintf("\"%s\"", shellCommand),
				"all",
			}

			fmt.Println(param)
			return execute.ExternalProgram("ansible", param...)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ansible.EnsureAnsibleDirectory(); err != nil {
				return errors.New("not an Ansible development directory")
			}

			return vagrant.EnsureVagrantfile()
		},
	}

	return cmd
}
