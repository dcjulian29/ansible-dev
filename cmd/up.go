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
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
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
		PostRun: func(cmd *cobra.Command, args []string) {
			ensureWorkingDirectoryAndExit()
		},
	}
)

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().BoolP("development", "d", true, "only start and provision the development VMs")
	upCmd.Flags().BoolP("provision", "p", false, "only start and provision the provision VM")
	upCmd.Flags().BoolP("test", "t", false, "only start and provision the test VMs")
	upCmd.Flags().Bool("base", true, "provision the VMs with the base role minimal tag")
	upCmd.Flags().String("role", "", "provision the VMs with the specified role")
	upCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")

	upCmd.MarkFlagsMutuallyExclusive("development", "provision", "test")
}

func vagrant_up(cmd *cobra.Command) {
	inv, err := ini.Load("hosts.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	sectionName := "ansibledev"

	if r, _ := cmd.Flags().GetBool("provision"); r {
		sectionName = "provisiontest"
	}

	if r, _ := cmd.Flags().GetBool("test"); r {
		sectionName = "vagrant"
	}

	section, err := inv.GetSection(sectionName)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, vm := range section.KeyStrings() {
		name := strings.Split(vm, " ")[0]
		addr := section.Key(vm).String()

		fmt.Printf("\nBringing '%s' online...\n\n", name)

		executeExternalProgram("vagrant", "up", name)

		fmt.Printf("\nSearching for '%s' at %s...", name, addr)

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

	verbose, _ := cmd.Flags().GetBool("base")

	if r, _ := cmd.Flags().GetBool("base"); r {
		fmt.Println("\nApplying the base role with the minimal tag...")
		generate_play("dcjulian29.base")
		execute_play(Play{
			Limit:      sectionName,
			Tags:       []string{"minimal"},
			FlushCache: true,
			Verbose:    verbose,
		})
	}

	role, _ := cmd.Flags().GetString("role")

	if len(role) > 0 {
		fmt.Printf("\nApplying the '%s' role..", role)
		generate_play(role)
		execute_play(Play{
			Name:    role,
			Limit:   sectionName,
			Verbose: verbose,
		})
	}
}

func ping(address string) bool {
	addr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		fmt.Println(err)
		return false
	}

	conn, err := net.DialIP("ip4:icmp", nil, addr)
	if err != nil {
		return false
	}

	defer conn.Close()

	msg := append([]byte{8, 0, 0, 0, 0, 13, 0, 37}, (make([]byte, 48))...)

	sum := 0
	for i := 0; i < len(msg)-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16

	checksum := uint16(^sum)

	msg[2] = byte(checksum >> 8)
	msg[3] = byte(checksum & 0xff)

	if _, err = conn.Write(msg); err != nil {
		return false
	}

	r := make([]byte, 1024)

	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	if _, err = conn.Read(r); err != nil {
		return false
	}

	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), r)
	if err != nil {
		return false
	}

	return rm.Code == 0
}
