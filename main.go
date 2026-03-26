/*
Copyright © 2026 Julian Easterling julian@julianscorner.com

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

// ansible-dev is a CLI tool that streamlines Ansible role and playbook
// development inside Vagrant-based environments. It wraps common
// ansible-galaxy, ansible-playbook, and Vagrant workflows behind a
// single, ergonomic command-line interface built with Cobra.
//
// Usage:
//
//	ansible-dev <command> [flags]
//
// Run "ansible-dev --help" for the full list of available commands.
package main

import "github.com/dcjulian29/ansible-dev/cmd"

func main() {
	cmd.Execute()
}
