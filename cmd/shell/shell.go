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

// Package shell implements the "ansible-dev shell" command, which
// executes an ad-hoc shell command on all hosts in the Vagrant
// development environment using the Ansible "shell" module.
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

// NewCommand creates and returns the Cobra command for "ansible-dev shell",
// which runs an ad-hoc shell command on every host defined in the
// project's hosts.ini inventory file.
//
// Usage:
//
//	ansible-dev shell [flags] -- <command>
//
// Everything after the "--" separator is treated as the shell command to
// execute. Multiple words are joined into a single string and passed to
// the Ansible "shell" module via its -a argument. For example:
//
//	ansible-dev shell -- uptime
//	ansible-dev shell -- cat /etc/hostname
//
// Under the hood, the command invokes:
//
//	ansible -i hosts.ini -m shell -a "<command>" all
//
// targeting the "all" host pattern, which runs the command on every host
// in the inventory. If no arguments are supplied, the help text is
// displayed instead.
//
// A PreRunE hook performs two checks before execution:
//  1. [ansible.EnsureAnsibleDirectory] — verifies the current directory
//     is a valid Ansible project. Returns a simplified "not an Ansible
//     development directory" message on failure.
//  2. [vagrant.EnsureVagrantfile] — confirms a Vagrantfile is present,
//     since the command targets Vagrant-managed hosts.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell [flags] -- <command>",
		Short: "Execute shell command in the Ansible development vagrant environment",
		Args: func(_ *cobra.Command, args []string) error {
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
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if err := ansible.EnsureAnsibleDirectory(); err != nil {
				return errors.New("not an Ansible development directory")
			}

			return vagrant.EnsureVagrantfile()
		},
	}

	return cmd
}
