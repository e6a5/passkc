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

	"github.com/spf13/cobra"
)

type setCmdRunner struct {
	kcManager KeychainManager
}

func (r *setCmdRunner) run(cmd *cobra.Command, args []string) {
	filePath, _ := cmd.Flags().GetString("file")
	quiet, _ := cmd.Flags().GetBool("quiet")

	if filePath != "" {
		// Read credentials from file
		file, err := os.Open(filePath)
		if err != nil {
			cmd.PrintErrf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			parts := strings.Fields(scanner.Text())
			if len(parts) >= 2 {
				domain := parts[0]
				username := parts[1]
				password := ""
				if len(parts) > 2 {
					password = parts[2]
				}
				if err := r.kcManager.SetData(domain, username, password); err != nil {
					cmd.PrintErrf("Error setting credentials for %s: %v\n", domain, err)
				} else if !quiet {
					cmd.Printf("Saved credentials for %s\n", domain)
				}
			}
		}
	} else if len(args) == 2 {
		// Set credentials from command line arguments
		domain := args[0]
		username := args[1]
		if err := r.kcManager.SetData(domain, username, ""); err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			os.Exit(1)
		}
		if !quiet {
			cmd.Println("Saved successfully")
		}
	} else {
		// Read from stdin
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
					cmd.Println("Saved successfully")
				}
			} else {
				cmd.PrintErrf("Error: invalid input format\n")
				os.Exit(1)
			}
		}
	}
}

func newSetCmd(kcManager KeychainManager) *cobra.Command {
	runner := &setCmdRunner{
		kcManager: kcManager,
	}
	cmd := &cobra.Command{
		Use:   "set [domain] [username]",
		Short: "Store credentials for a domain in the Keychain",
		Long: `The passkc set command stores credentials for a domain in the Keychain.
It supports input from file and stdin.

Examples:
  # Set credentials interactively
  passkc set domain.com username

  # Set credentials from file
  passkc set -f credentials.txt

  # Set credentials from stdin
  echo "domain.com username password" | passkc set

  # Set credentials in quiet mode
  passkc set domain.com username -q`,
		Args: cobra.MaximumNArgs(2),
		Run:  runner.run,
	}
	cmd.Flags().StringP("file", "f", "", "Read credentials from file")
	return cmd
}

func init() {
	rootCmd.AddCommand(newSetCmd(&LiveKeychainManager{}))
}
