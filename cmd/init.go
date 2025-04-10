/*
Copyright © 2023 Julian Easterling <julian@julianscorner.com>

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
		Long: `Initialize an development environment for Ansible development by creating the folder
structure and generating the needed files to quickly set up a virtual environment
ready for development. Vagrant can be used to manage the environment and connect
to troubleshoot and/or validate.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing development environment...")

			init_env()
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite an existing development environment")
}

func init_env() {
	if workingDirectory != folderPath {
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			fmt.Println("Creating development environment folder...")
			if err := os.MkdirAll(folderPath, 0755); err != nil {
				fmt.Println("unable to create development environment folder")
				os.Exit(1)
			}
		}

		ensureAnsibleDirectory()
	}

	if fileExists("ansible.cfg") && !force {
		fmt.Println("folder for the ansible environment already contains a context and force was not provided")
		ensureWorkingDirectoryAndExit()
	}

	fmt.Println("    ...   roles/")
	if err := ensureDir("roles"); err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	fmt.Println("    ...   playbooks/")
	if err := ensureDir("playbooks"); err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	ansible_cfg()
	ansible_lint()
	inventory_file()
	vagrant_file()

	ensureWorkingDirectoryAndExit()
}

func ansible_cfg() {
	fmt.Println("    ...   ansible.cfg")

	file, err := os.Create("ansible.cfg")

	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`[defaults]
any_errors_fatal            = true
collections_path            = ./collections
duplicate_dict_key          = error
error_on_undefined_vars     = true
gathering                   = smart
host_key_checking           = false
inventory                   = ./hosts.ini
log_path                    = ./ansible.log
roles_path                  = ./roles
callback_result_format      = yaml
verbosity                   = 1

[diff]
always                      = true
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func inventory_file() {
	fmt.Println("    ...   hosts.ini")
	file, err := os.Create("hosts.ini")

	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`[ansibledev]
debian ansible_host=192.168.57.5

[vagrant]
debian ansible_host=192.168.57.5
rocky ansible_host=192.168.57.6

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

func vagrant_file() {
	fmt.Println("    ...   Vagrantfile")
	file, err := os.Create("Vagrantfile")

	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`Vagrant.configure("2") do |config|
  config.ssh.insert_key = false
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.network :forwarded_port, guest: 22, host: 2220, id: "ssh", disabled: true
  config.vm.provision "ping", type: "shell", inline: "ping -c 1 192.168.57.1", run: "always"
  config.ssh.extra_args = "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null"

  config.vm.provider "virtualbox" do |vb|
    vb.gui = false
    vb.memory = 4096
    vb.cpus = 2
    vb.check_guest_additions = false
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
    vb.customize [ "modifyvm", :id, "--graphicscontroller", "vmsvga"]
    vb.customize [ "modifyvm", :id, "--nestedpaging", "on"]
    vb.customize [ "modifyvm", :id, "--largepages", "on"]
    vb.customize [ "modifyvm", :id, "--ioapic", "on"]
  end

  config.vm.define "debian" do |c|
    c.vm.box = "debian/bullseye64"
    c.vm.network "private_network", ip: "192.168.57.5"
    c.vm.network :forwarded_port, guest: 22, host: 8005, id: 'ssh'
  end

  config.vm.define "rocky" do |c|
    c.vm.box = "rockylinux/9"
    c.vm.network "private_network", ip: "192.168.57.6"
    c.vm.network :forwarded_port, guest: 22, host: 8006, id: 'ssh'
  end

end
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}

func ansible_lint() {
	fmt.Println("    ...   .ansible-lint")
	file, err := os.Create(".ansible-lint")

	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`skip_list:
  - var-naming[no-role-prefix]

warn_list:
  - internal-error
  - no-handler
  - experimental

kinds:
  - playbook: "playbooks/*.{yml,yaml}"
  - yaml: "**/*.{yml,yaml}"

verbosity: 1
`)

	if _, err = file.Write(content); err != nil {
		cobra.CheckErr(err)
	}
}
