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
	Short: "A Unix-style password manager using macOS Keychain",
	Long: `passkc is a command-line tool for managing credentials in macOS Keychain.
It follows Unix philosophy by:
- Doing one thing well: managing credentials
- Working with text streams
- Being composable with other tools
- Using plain text interfaces
- Avoiding captive user interfaces

Examples:
  # Get credentials and pipe to clipboard
  passkc get domain.com | pbcopy

  # List credentials and filter
  passkc show | grep "google"

  # Set credentials from file
  passkc set -f credentials.txt

  # Output in JSON format
  passkc show -o json | jq '.[] | select(.domain == "google.com")'`,
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
