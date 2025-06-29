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

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "passkc",
	Short: "Simple password manager using macOS Keychain",
	Long: `passkc is a simple command-line password manager for macOS.
Store and retrieve passwords securely using the macOS Keychain.

Common usage:
  passkc set github.com myusername     # Save a password
  passkc get github.com                # Retrieve a password
  passkc show                          # List all saved passwords
  passkc remove github.com             # Delete a password

Advanced usage:
  passkc get github.com -q | pbcopy    # Copy password to clipboard
  passkc show | grep google            # Search for specific sites`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initializeFlags(rootCmd)
}

func initializeFlags(cmd *cobra.Command) {
	// Global flags
	cmd.PersistentFlags().StringP("output", "o", "text", "Output format (text|json|csv)")
	cmd.PersistentFlags().StringP("config", "c", "", "Config file (default is $HOME/.passkc.yaml)")
	cmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress prompts and non-essential output")

	// Environment variable support
	if domain := os.Getenv("PASSKC_DEFAULT_DOMAIN"); domain != "" {
		cmd.PersistentFlags().String("domain", domain, "Default domain to use")
	}
}
