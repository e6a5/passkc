/*
Copyright © 2023 Hiep Tran <tranhiepqna@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/e6a5/passkc/kc"
	"github.com/spf13/cobra"
)

type removeCmdRunner struct {
	kcManager KeychainManager
}

func (r *removeCmdRunner) run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.PrintErrf("Usage: passkc remove <domain>\n\n")
		cmd.PrintErrf("Examples:\n")
		cmd.PrintErrf("  passkc remove github.com                 # Remove credentials for github.com\n")
		cmd.PrintErrf("  passkc remove github.com -q              # Remove without confirmation\n")
		cmd.PrintErrf("\nFor more help: passkc remove --help\n")
		os.Exit(1)
	}

	domain := args[0]
	quiet, _ := cmd.Flags().GetBool("quiet")
	force, _ := cmd.Flags().GetBool("force")

	// Validate domain
	if err := kc.ValidateDomain(domain); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if credentials exist first
	cred, err := r.kcManager.GetData(domain)
	if err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	// Confirmation prompt (unless forced or quiet)
	if !force && !quiet {
		cmd.Printf("Are you sure you want to remove credentials for '%s' (username: %s)? [y/N]: ", domain, cred.Username)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			response := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if response != "y" && response != "yes" {
				cmd.Printf("Cancelled.\n")
				return
			}
		}
	}

	err = r.kcManager.RemoveData(domain)
	if err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	if !quiet {
		cmd.Printf("✓ Removed credentials for %s\n", domain)
	}
}

func newRemoveCmd(kcManager KeychainManager) *cobra.Command {
	runner := &removeCmdRunner{
		kcManager: kcManager,
	}
	cmd := &cobra.Command{
		Use:   "remove <domain>",
		Short: "Remove credentials for a website or service",
		Long: `Remove stored credentials for a domain from the keychain.

This action permanently deletes the saved username and password.
By default, you'll be asked to confirm the deletion.

Examples:
  passkc remove github.com                 # Remove with confirmation prompt
  passkc remove github.com --force         # Remove without confirmation
  passkc remove github.com -q              # Remove quietly (no output)`,
		Args: cobra.ExactArgs(1),
		Run:  runner.run,
	}
	cmd.Flags().BoolP("force", "f", false, "Remove without confirmation prompt")
	return cmd
}

func init() {
	rootCmd.AddCommand(newRemoveCmd(&LiveKeychainManager{}))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
