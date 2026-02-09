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
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	force bool

	initCmd = &cobra.Command{
		Use:     "init",
		Aliases: []string{"initialize"},
		Short:   "Initialize an Ansible development vagrant environment",
		Long:    "Initialize an Ansible development vagrant environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Yellow("Initializing development environment..."))

			create_folder()
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite an existing development environment")
}

func create_folder() {
	if workingDirectory != folderPath {
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			cobra.CheckErr(os.MkdirAll(folderPath, 0755))
		}
	}

	if fileExists("ansible.cfg") && !force {
		fmt.Println(Fatal("folder for ansible development already exist and force was not provided"))
		os.Exit(1)
	}

	ansible_cfg()
	ansible_lint()

	fmt.Println("  ...  collections/")
	cobra.CheckErr(ensureDir("collections"))

	group_vars()
	host_vars()
	inventory_file()
	runbook()

	fmt.Println("  ...  roles/")
	cobra.CheckErr(ensureDir("roles"))

	yaml_ignore()
	yaml_lint()
	vagrant_file()

}

func ansible_cfg() {
	fmt.Println("  ...  ansible.cfg")

	file, err := os.Create("ansible.cfg")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`[defaults]
any_errors_fatal            = true
collections_path            = ./collections
duplicate_dict_key          = error
interpreter_python          = auto_silent
inventory                   = ./hosts.ini
log_path                    = ./ansible.log
roles_path                  = ./roles
callback_result_format      = yaml
verbosity                   = 1
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func ansible_lint() {
	fmt.Println("  ...  .ansible-lint")
	file, err := os.Create(".ansible-lint")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`---
enable_list:
  - args
  - empty-string-compare
  - no-log-password
  - no-same-owner
  - name[prefix]
  - yaml
exclude_paths:
  - collections/
  - roles/
kinds:
  - playbook: "playbooks/*.yml"
profile: production
skip_list:
  - experimental
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func group_vars() {
	fmt.Println("  ...  group_vars/")
	cobra.CheckErr(ensureDir("group_vars/"))

	fmt.Println("  ...    vagrant.yml")

	file, err := os.Create("group_vars/vagrant.yml")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte("---\nname: value")

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func host_vars() {
	fmt.Println("  ...  host_vars/")
	cobra.CheckErr(ensureDir("host_vars"))

	fmt.Println("  ...    debian.yml")

	file, err := os.Create("host_vars/debian.yml")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte("---\nname: value")

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}

	fmt.Println("  ...    alma.yml")

	file, err = os.Create("host_vars/alma.yml")
	cobra.CheckErr(err)

	defer file.Close()

	content = []byte("---\nname: value")

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func inventory_file() {
	fmt.Println("  ...  hosts.ini")
	file, err := os.Create("hosts.ini")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`[vagrant]
debian ansible_host=192.168.57.5
alma ansible_host=192.168.57.6

[all:vars]
ansible_user=vagrant
ansible_ssh_common_args='-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o CheckHostIP=no'
ansible_port=22
ansible_ssh_private_key_file=~/.ssh/insecure_private_key
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func runbook() {
	fmt.Println("  ...  playbooks/")
	cobra.CheckErr(ensureDir("playbooks"))

	fmt.Println("  ...    runbook.yml")
	file, err := os.Create("playbooks/runbook.yml")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`---
- name: Test Ansible Runbook
  hosts: all
  become: true
  any_errors_fatal: true
  gather_facts: true

  handlers:
  #  - name: [Name of Handler]
  #    ...

  tasks:
  #  - name: [Name of Task]
  #    ...

  #  - # Repeat as necessary

  vars:
    name: value
    # variables needed for runbook
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func yaml_ignore() {
	fmt.Println("  ...  .yamlignore")
	file, err := os.Create(".yamlignore")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`secrets.yml
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func yaml_lint() {
	fmt.Println("  ...  .yamlint")
	file, err := os.Create(".yamllint")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`---
extends: default

rules:
  braces:
    min-spaces-inside: 0
    max-spaces-inside: 1
  comments:
    min-spaces-from-content: 1
  comments-indentation: false
  indentation:
    indent-sequences: true
  line-length:
    max: 125
  new-lines: disable
  octal-values:
    forbid-implicit-octal: true
    forbid-explicit-octal: true

ignore-from-file:
  - .gitignore
  - .yamlignore
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func vagrant_file() {
	fmt.Println("  ...  Vagrantfile")
	file, err := os.Create("Vagrantfile")
	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`Vagrant.configure("2") do |config|
  config.ssh.insert_key = false
  if Vagrant.has_plugin?("vagrant-vbguest")
    config.vbguest.auto_update = false
  end
  config.vm.box_check_update = true
  config.vm.provision "ping", type: "shell", inline: "ping -c 1 192.168.57.1", run: "always"
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.provider "virtualbox" do |vb|
    vb.gui = false
    vb.cpus = 2
    vb.memory = 4096
    vb.check_guest_additions = false
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
    vb.customize [ "modifyvm", :id, "--graphicscontroller", "vmsvga"]
    vb.customize [ "modifyvm", :id, "--ioapic", "on"]
  end
  config.vm.define "debian" do |c|
    c.vm.box = "debian/bookworm64"
    c.vm.hostname = "debian.dev"
    c.vm.network "private_network", ip: "192.168.57.5"
  end
  config.vm.define "alma" do |c|
    c.vm.box = "almalinux/10"
    c.vm.hostname = "alma.dev"
    c.vm.network "private_network", ip: "192.168.57.6"
  end
end
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}
