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
	"fmt"
	"os"
	"strings"

	"github.com/e6a5/passkc/kc"
	"github.com/spf13/cobra"
)

type setCmdRunner struct {
	kcManager KeychainManager
}

func (r *setCmdRunner) run(cmd *cobra.Command, args []string) {
	filePath, _ := cmd.Flags().GetString("file")
	quiet, _ := cmd.Flags().GetBool("quiet")

	if filePath != "" {
		r.handleFileInput(cmd, filePath, quiet)
		return
	}

	if len(args) >= 2 {
		r.handleDirectInput(cmd, args[0], args[1], quiet)
		return
	}

	if len(args) == 1 {
		r.handleInteractiveInput(cmd, args[0], quiet)
		return
	}

	// No arguments provided - try stdin or show usage
	if r.hasStdinInput() {
		r.handleStdinInput(cmd, quiet)
		return
	}

	// Show helpful usage message
	cmd.PrintErrf("Usage: passkc set <domain> [username]\n\n")
	cmd.PrintErrf("Examples:\n")
	cmd.PrintErrf("  passkc set github.com                    # Interactive username/password prompt\n")
	cmd.PrintErrf("  passkc set github.com myusername         # Password prompt only\n")
	cmd.PrintErrf("  passkc set -f credentials.txt            # Import from file\n")
	cmd.PrintErrf("\nFor more help: passkc set --help\n")
	os.Exit(1)
}

func (r *setCmdRunner) handleDirectInput(cmd *cobra.Command, domain, username string, quiet bool) {
	// Validate inputs
	if err := kc.ValidateDomain(domain); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}
	if err := kc.ValidateUsername(username); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := r.kcManager.SetData(domain, username, ""); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	if !quiet {
		cmd.Printf("✓ Saved credentials for %s@%s\n", username, domain)
	}
}

func (r *setCmdRunner) handleInteractiveInput(cmd *cobra.Command, domain string, quiet bool) {
	// Validate domain
	if err := kc.ValidateDomain(domain); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	// Interactive username prompt
	fmt.Printf("Username for %s: ", domain)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		cmd.PrintErrf("Error: failed to read username\n")
		os.Exit(1)
	}

	username := strings.TrimSpace(scanner.Text())
	if err := kc.ValidateUsername(username); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := r.kcManager.SetData(domain, username, ""); err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	if !quiet {
		cmd.Printf("✓ Saved credentials for %s@%s\n", username, domain)
	}
}

func (r *setCmdRunner) handleFileInput(cmd *cobra.Command, filePath string, quiet bool) {
	file, err := os.Open(filePath)
	if err != nil {
		cmd.PrintErrf("Error: cannot open file '%s': %v\n", filePath, err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			cmd.PrintErrf("Warning: failed to close file: %v\n", closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	successCount := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			cmd.PrintErrf("Warning: skipping line %d (invalid format): %s\n", lineNum, line)
			continue
		}

		domain := parts[0]
		username := parts[1]
		password := ""
		if len(parts) > 2 {
			password = parts[2]
		}

		if err := r.kcManager.SetData(domain, username, password); err != nil {
			cmd.PrintErrf("Error on line %d: failed to save %s: %v\n", lineNum, domain, err)
		} else {
			successCount++
			if !quiet {
				cmd.Printf("✓ Saved credentials for %s@%s\n", username, domain)
			}
		}
	}

	if !quiet {
		cmd.Printf("\nImported %d credentials successfully\n", successCount)
	}
}

func (r *setCmdRunner) handleStdinInput(cmd *cobra.Command, quiet bool) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) >= 2 {
			domain := parts[0]
			username := parts[1]
			password := ""
			if len(parts) > 2 {
				password = parts[2]
			}
			if err := r.kcManager.SetData(domain, username, password); err != nil {
				cmd.PrintErrf("Error: %v\n", err)
				os.Exit(1)
			}
			if !quiet {
				cmd.Printf("✓ Saved credentials for %s@%s\n", username, domain)
			}
		} else {
			cmd.PrintErrf("Error: invalid input format. Expected: domain username [password]\n")
			os.Exit(1)
		}
	}
}

func (r *setCmdRunner) hasStdinInput() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func newSetCmd(kcManager KeychainManager) *cobra.Command {
	runner := &setCmdRunner{
		kcManager: kcManager,
	}
	cmd := &cobra.Command{
		Use:   "set <domain> [username]",
		Short: "Save credentials for a website or service",
		Long: `Save your username and password for a domain.

If you only provide the domain, you'll be prompted for the username.
You'll always be prompted for the password (for security).

Examples:
  passkc set github.com                    # Interactive: prompts for username and password
  passkc set github.com myusername         # Prompts for password only
  passkc set -f credentials.txt            # Import multiple credentials from file

File format (one per line):
  domain username [password]
  github.com user1 pass123
  google.com user2`,
		Args: cobra.RangeArgs(0, 2),
		Run:  runner.run,
	}
	cmd.Flags().StringP("file", "f", "", "Import credentials from file")
	return cmd
}

func init() {
	rootCmd.AddCommand(newSetCmd(&LiveKeychainManager{}))
}
