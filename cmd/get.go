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
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			domain = strings.TrimSpace(scanner.Text())
		}
	}

	if domain == "" {
		cmd.PrintErrf("Error: domain is required\n")
		os.Exit(1)
	}

	outputFormat, _ := cmd.Flags().GetString("output")
	quiet, _ := cmd.Flags().GetBool("quiet")

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
		if !quiet {
			cmd.Printf("Domain: %s\nUsername: %s\n", cred.Domain, cred.Username)
		}
		cmd.Print(cred.Password)
	}
}

func newGetCmd(kcManager KeychainManager) *cobra.Command {
	runner := &getCmdRunner{
		kcManager: kcManager,
	}
	cmd := &cobra.Command{
		Use:   "get [domain]",
		Short: "Retrieve username and password for a domain from the Keychain",
		Long: `The passkc get command retrieves the stored username and password for a specific domain.
It supports multiple output formats and can be used in pipelines.

Examples:
  # Get credentials in default format
  passkc get domain.com

  # Get credentials in JSON format
  passkc get domain.com -o json

  # Get credentials and pipe to clipboard
  passkc get domain.com | pbcopy

  # Get credentials for domain from stdin
  echo "domain.com" | passkc get`,
		Args: cobra.MaximumNArgs(1),
		Run:  runner.run,
	}
	cmd.Flags().BoolP("password-only", "p", false, "Output only the password")
	return cmd
}

func init() {
	rootCmd.AddCommand(newGetCmd(&LiveKeychainManager{}))
}
