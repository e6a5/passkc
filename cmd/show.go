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
	"encoding/csv"
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/e6a5/passkc/kc"
	"github.com/spf13/cobra"
)

type showCmdRunner struct {
	kcManager KeychainManager
}

func (r *showCmdRunner) run(cmd *cobra.Command, args []string) {
	outputFormat, _ := cmd.Flags().GetString("output")
	pattern, _ := cmd.Flags().GetString("pattern")
	sortBy, _ := cmd.Flags().GetString("sort")
	quiet, _ := cmd.Flags().GetBool("quiet")

	creds, err := r.kcManager.ListData()
	if err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}

	// Filter credentials if pattern is provided
	if pattern != "" {
		filtered := make([]kc.Credential, 0)
		for _, cred := range creds {
			if strings.Contains(cred.Domain, pattern) || strings.Contains(cred.Username, pattern) {
				filtered = append(filtered, cred)
			}
		}
		creds = filtered
	}

	// Sort credentials
	switch sortBy {
	case "domain":
		sort.Slice(creds, func(i, j int) bool {
			return creds[i].Domain < creds[j].Domain
		})
	case "username":
		sort.Slice(creds, func(i, j int) bool {
			return creds[i].Username < creds[j].Username
		})
	}

	// Output in requested format
	switch outputFormat {
	case "json":
		// Ensure we output a valid JSON array even if creds is nil
		if creds == nil {
			creds = make([]kc.Credential, 0)
		}
		json.NewEncoder(cmd.OutOrStdout()).Encode(creds)
	case "csv":
		w := csv.NewWriter(cmd.OutOrStdout())
		w.Write([]string{"Domain", "Username"})
		for _, cred := range creds {
			w.Write([]string{cred.Domain, cred.Username})
		}
		w.Flush()
	default:
		if !quiet {
			cmd.Println("List of credentials:")
		}
		for _, cred := range creds {
			cmd.Printf("%s (%s)\n", cred.Domain, cred.Username)
		}
	}
}

func newShowCmd(kcManager KeychainManager) *cobra.Command {
	runner := &showCmdRunner{
		kcManager: kcManager,
	}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "List all stored credentials",
		Long: `The passkc show command lists all stored credentials.
It supports filtering, sorting, and multiple output formats.

Examples:
  # List all credentials
  passkc show

  # List credentials in JSON format
  passkc show -o json

  # Filter credentials by pattern
  passkc show --pattern "*.com"

  # Sort credentials by domain
  passkc show --sort domain

  # Combine filtering and sorting
  passkc show --pattern "google" --sort username`,
		Run: runner.run,
	}
	cmd.Flags().String("pattern", "", "Filter credentials by pattern")
	cmd.Flags().String("sort", "", "Sort by field (domain|username)")
	return cmd
}

func init() {
	// The real command uses the live keychain manager.
	rootCmd.AddCommand(newShowCmd(&LiveKeychainManager{}))
}
