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
	"encoding/csv"
	"encoding/json"
	"os"
	"strings"

	"github.com/e6a5/passkc/kc"
	"github.com/spf13/cobra"
)

type getCmdRunner struct {
	kcManager KeychainManager
}

func (r *getCmdRunner) run(cmd *cobra.Command, args []string) {
	var domain string

	if len(args) > 0 {
		domain = args[0]
	} else {
		// Read from stdin if no domain provided
		if r.hasStdinInput() {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				domain = strings.TrimSpace(scanner.Text())
			}
		}
	}

	if domain == "" {
		cmd.PrintErrf("Usage: passkc get <domain>\n\n")
		cmd.PrintErrf("Examples:\n")
		cmd.PrintErrf("  passkc get github.com                    # Show domain and username only\n")
		cmd.PrintErrf("  passkc get github.com -p                 # Show password only\n")
		cmd.PrintErrf("  passkc get github.com -q                 # Quiet mode (password only)\n")
		cmd.PrintErrf("  passkc get github.com -q | pbcopy        # Copy password to clipboard\n")
		cmd.PrintErrf("  echo \"github.com\" | passkc get          # Read domain from pipe\n")
		cmd.PrintErrf("\nFor more help: passkc get --help\n")
		os.Exit(1)
	}

	// Validate domain
	if err := kc.ValidateDomain(domain); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, _ := cmd.Flags().GetString("output")
	quiet, _ := cmd.Flags().GetBool("quiet")
	passwordOnly, _ := cmd.Flags().GetBool("password-only")

	cred, err := r.kcManager.GetData(domain)
	if err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	switch outputFormat {
	case "json":
		json.NewEncoder(cmd.OutOrStdout()).Encode(cred)
	case "csv":
		w := csv.NewWriter(cmd.OutOrStdout())
		w.Write([]string{cred.Domain, cred.Username, cred.Password})
		w.Flush()
	default:
		if passwordOnly || quiet {
			// Just output the password
			cmd.Print(cred.Password)
		} else {
			// SECURITY FIX: By default, only show domain and username
			// Never show password in plain text unless explicitly requested
			cmd.Printf("Domain: %s\n", cred.Domain)
			cmd.Printf("Username: %s\n", cred.Username)
			cmd.Printf("\nTo get the password:\n")
			cmd.Printf("  passkc get %s -p                 # Show password\n", domain)
			cmd.Printf("  passkc get %s -q | pbcopy        # Copy to clipboard\n", domain)
		}
	}
}

func (r *getCmdRunner) hasStdinInput() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func newGetCmd(kcManager KeychainManager) *cobra.Command {
	runner := &getCmdRunner{
		kcManager: kcManager,
	}
	cmd := &cobra.Command{
		Use:   "get <domain>",
		Short: "Retrieve credentials for a website or service",
		Long: `Retrieve your saved username and password for a domain.

SECURITY: By default, only shows domain and username (password is hidden).
Use the -p flag to show the password, or -q to output only the password.

Examples:
  passkc get github.com                    # Show domain and username only (secure)
  passkc get github.com -p                 # Show password only  
  passkc get github.com -q                 # Quiet mode (password only)
  passkc get github.com -o json            # Output as JSON (includes password)
  passkc get github.com -q | pbcopy        # Copy password to clipboard (recommended)
  echo "github.com" | passkc get           # Read domain from pipe`,
		Args: cobra.MaximumNArgs(1),
		Run:  runner.run,
	}
	cmd.Flags().BoolP("password-only", "p", false, "Output only the password")
	return cmd
}

func init() {
	rootCmd.AddCommand(newGetCmd(&LiveKeychainManager{}))
}
