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
	"net"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Starts and provision the Ansible development vagrant environment",
		Long:  "Starts and provision the Ansible development vagrant environment",
		Run: func(cmd *cobra.Command, args []string) {
			vagrant_up(cmd)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			ensureAnsibleDirectory()
			ensureVagrantfile()
		},
	}
)

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().Bool("base", true, "provision the VMs with the base role minimal tag")
	upCmd.Flags().String("role", "", "provision the VMs with the specified role")
	upCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
}

func vagrant_up(cmd *cobra.Command) {
	inv, err := ini.Load("hosts.ini")
	cobra.CheckErr(err)

	section, err := inv.GetSection("vagrant")
	cobra.CheckErr(err)

	for _, vm := range section.KeyStrings() {
		name := strings.Split(vm, " ")[0]
		addr := section.Key(vm).String()

		fmt.Printf(Yellow("\nBringing '%s' online...\n\n"), name)

		executeExternalProgram("vagrant", "up", name)

		fmt.Printf(Yellow("\nSearching for '%s' at %s..."), name, addr)

		found := false
		count := 0

		for !found {
			found = ping(addr)

			if found {
				fmt.Println(Green(" [Found]"))
			} else {
				if count < 20 {
					fmt.Print(".")
					count++
				} else {
					fmt.Println(Red(" [NotFound]"))
					return
				}
			}
		}
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	if r, _ := cmd.Flags().GetBool("base"); r {
		fmt.Println(Yellow("\nApplying the base role with the minimal tag..."))
		generate_play("dcjulian29.base")
		execute_play(Play{
			Tags:       []string{"minimal"},
			FlushCache: true,
			Verbose:    verbose,
		})
	}

	role, _ := cmd.Flags().GetString("role")

	if len(role) > 0 {
		fmt.Printf(Teal("\nApplying the '%s' role...\n"), role)
		generate_play(role)
		execute_play(Play{
			Name:    role,
			Verbose: verbose,
		})
	}
}

func ping(address string) bool {
	addr, err := net.ResolveIPAddr("ip", address)
	cobra.CheckErr(err)

	var output []byte

	if runtime.GOOS == "windows" {
		output, _ = exec.Command("ping", "-w", "1000", "-n", "1", addr.String()).CombinedOutput()
	} else {
		output, _ = exec.Command("ping", "-c", "1", addr.String()).CombinedOutput()
	}

	if strings.Contains(string(output[:]), "TTL") {
		return true
	}

	return false
}
