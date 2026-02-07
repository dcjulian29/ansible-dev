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
	"strings"

	"github.com/spf13/cobra"
)

type Play struct {
	Name          string
	Tags          []string
	AskVaultPass  bool
	AskBecomePass bool
	FlushCache    bool
	Step          bool
	Verbose       bool
}

var (
	playFromFlags Play

	playCmd = &cobra.Command{
		Use:   "play <role>",
		Args:  cobra.ExactArgs(1),
		Short: "Provision ansible role to Ansible development vagrant environment",
		Long:  "Provision ansible role to Ansible development vagrant environment",
		Run: func(cmd *cobra.Command, args []string) {
			playFromFlags.Name = args[0]
			playFromFlags.Tags, _ = cmd.Flags().GetStringSlice("tags")
			playFromFlags.AskBecomePass, _ = cmd.Flags().GetBool("becomepass")
			playFromFlags.AskVaultPass, _ = cmd.Flags().GetBool("vaultpass")
			playFromFlags.FlushCache, _ = cmd.Flags().GetBool("flushcache")
			playFromFlags.Step, _ = cmd.Flags().GetBool("step")
			playFromFlags.Verbose, _ = cmd.Flags().GetBool("verbose")

			generate_play(playFromFlags.Name)
			execute_play(playFromFlags)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			ensureAnsibleDirectory()
			ensureVagrantfile()
		},
	}
)

func init() {
	rootCmd.AddCommand(playCmd)

	playCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	playCmd.Flags().Bool("ask-vault-password", false, "ask for vault password")
	playCmd.Flags().Bool("ask-become-password", false, "ask for privilege escalation password")
	playCmd.Flags().Bool("flush-cache", false, "clear the fact cache for every host in inventory")
	playCmd.Flags().BoolP("step", "s", false, "one-step-at-a-time: confirm each task before running")
	playCmd.Flags().StringSlice("tags", []string{}, "only plays and task tagged with these values")
}

func ensurePlayFile() (*os.File, error) {
	if !dirExists(".tmp") {
		cobra.CheckErr(os.Mkdir(".tmp", 0755))
	}

	if fileExists(".tmp/play.yml") {
		cobra.CheckErr(os.Remove(".tmp/play.yml"))
	}

	return os.Create(".tmp/play.yml")
}

func generate_play(roleName string) {
	file, err := ensurePlayFile()
	cobra.CheckErr(err)

	defer file.Close()

	content := "---\n- hosts: all\n  any_errors_fatal: true\n  become: true\n  roles:\n"
	content = fmt.Sprintf("%s%s", content, fmt.Sprintf("    - %s\n", roleName))

	if _, err = file.WriteString(content); err != nil {
		cobra.CheckErr(err)
	}
}

func execute_play(play Play) {
	var param []string

	if len(play.Limit) > 0 {
		param = append(param, "--limit")
		param = append(param, play.Limit)
	}

	if len(play.Tags) > 0 {
		param = append(param, "--tags")
		param = append(param, strings.Join(play.Tags, ","))
	}

	if play.FlushCache {
		param = append(param, "--flush-cache")
	}

	if play.AskVaultPass {
		param = append(param, "--ask-vault-password")
	}

	if play.Verbose {
		param = append(param, "-v")
	}

	if play.Step {
		param = append(param, "--step")
	}

	param = append(param, ".tmp/play.yml")

	if fileExists(".tmp/play.yml") {
		executeExternalProgram("ansible-playbook", param...)
	}
}
