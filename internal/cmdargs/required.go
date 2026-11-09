/*
Copyright © 2026 Julian Easterling

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

// Package cmdargs holds the positional-argument checks shared by the
// ansible-dev commands, so a missing argument is reported the same way
// everywhere.
package cmdargs

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Required returns a cobra Args check for a command that takes exactly one
// argument, described by what (for example "a role name"):
//
//   - No arguments and no flags: the command was typed bare, most likely to
//     see how it is used, so the help is printed and the exit status is 0.
//   - Flags but no argument: the caller meant to run the command and left the
//     argument out, so that is an error naming what is missing ("expects a
//     role name") with a non-zero exit. Printing the help here instead looked
//     as if the flags were being rejected.
//   - More than one argument: an error saying how many were given.
//
// It runs as the command's Args check on purpose. Cobra validates arguments
// before PersistentPreRunE and PreRunE, so the help appears, and the
// missing-argument error is reported, even outside a development folder where
// a PreRunE such as EnsureAnsibleDirectory would otherwise fail first. The help
// is requested by returning pflag.ErrHelp, which cobra answers by printing the
// help and returning no error, so neither the PreRunE nor RunE runs.
func Required(what string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		switch {
		case len(args) == 1:
			return nil
		case len(args) == 0 && cmd.Flags().NFlag() == 0:
			return pflag.ErrHelp
		case len(args) == 0:
			return fmt.Errorf("expects %s", what)
		default:
			return fmt.Errorf("expects %s, received %d arguments", what, len(args))
		}
	}
}

// AtLeastOne is [Required] for a command whose argument is a list of words
// rather than a single name, such as "shell -- uptime" or "shell -- cat
// /etc/hostname": any number of arguments from one up is accepted. With none,
// it behaves exactly as Required does: the help when the command was typed
// bare (a lone "--" counts as bare, since cobra drops it), and an error naming
// what is missing when flags were given.
func AtLeastOne(what string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		switch {
		case len(args) > 0:
			return nil
		case cmd.Flags().NFlag() == 0:
			return pflag.ErrHelp
		default:
			return fmt.Errorf("expects %s", what)
		}
	}
}
