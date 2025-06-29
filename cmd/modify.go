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
	"os"

	"github.com/e6a5/passkc/kc"
	"github.com/spf13/cobra"
)

type modifyCmdRunner struct {
	kcManager KeychainManager
}

func (r *modifyCmdRunner) run(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.PrintErrf("Usage: passkc modify <domain> <new-username>\n\n")
		cmd.PrintErrf("Examples:\n")
		cmd.PrintErrf("  passkc modify github.com newusername     # Change username and password\n")
		cmd.PrintErrf("  passkc modify github.com same-user -q    # Change password only (quiet)\n")
		cmd.PrintErrf("\nNote: This will prompt for a new password.\n")
		cmd.PrintErrf("To keep the same password, use: passkc set <domain> <username>\n")
		cmd.PrintErrf("\nFor more help: passkc modify --help\n")
		os.Exit(1)
	}

	domain := args[0]
	newUsername := args[1]
	quiet, _ := cmd.Flags().GetBool("quiet")

	// Validate inputs
	if err := kc.ValidateDomain(domain); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}
	if err := kc.ValidateUsername(newUsername); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if credentials exist first
	_, err := r.kcManager.GetData(domain)
	if err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	// Show what we're changing
	if !quiet {
		cmd.Printf("Updating credentials for %s with username '%s'\n", domain, newUsername)
	}

	// We use SetData which will prompt for a password and update if the item exists
	err = r.kcManager.SetData(domain, newUsername, "")
	if err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	if !quiet {
		cmd.Printf("✓ Updated credentials for %s\n", domain)
	}
}

func newModifyCmd(kcManager KeychainManager) *cobra.Command {
	runner := &modifyCmdRunner{
		kcManager: kcManager,
	}
	return &cobra.Command{
		Use:   "modify <domain> <new-username>",
		Short: "Update credentials for a website or service",
		Long: `Update the username and/or password for existing credentials.

This command changes both the username and password for a domain.
You'll be prompted to enter a new password securely.

If you only want to change the password but keep the same username,
use the 'set' command instead.

Examples:
  passkc modify github.com newusername     # Change username and password
  passkc modify github.com same-user       # Keep username, change password
  
Tip: To change only the password, use:
  passkc set github.com existing-username`,
		Args: cobra.ExactArgs(2),
		Run:  runner.run,
	}
}

func init() {
	rootCmd.AddCommand(newModifyCmd(&LiveKeychainManager{}))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
