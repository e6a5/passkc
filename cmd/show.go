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

	// Show helpful message if no credentials exist
	if len(creds) == 0 && outputFormat == "text" && !quiet {
		cmd.Printf("No credentials found.\n\n")
		cmd.Printf("To add your first credential:\n")
		cmd.Printf("  passkc set github.com myusername\n\n")
		cmd.Printf("For help: passkc --help\n")
		return
	}

	originalCount := len(creds)

	// Filter credentials if pattern is provided
	if pattern != "" {
		filtered := make([]kc.Credential, 0)
		for _, cred := range creds {
			if strings.Contains(strings.ToLower(cred.Domain), strings.ToLower(pattern)) ||
				strings.Contains(strings.ToLower(cred.Username), strings.ToLower(pattern)) {
				filtered = append(filtered, cred)
			}
		}
		creds = filtered

		// Show message if pattern filtered out all results
		if len(creds) == 0 && outputFormat == "text" && !quiet {
			cmd.Printf("No credentials found matching pattern '%s'.\n", pattern)
			cmd.Printf("Found %d total credentials. Try a different search pattern.\n", originalCount)
			return
		}
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
	default:
		// Default sort by domain
		sort.Slice(creds, func(i, j int) bool {
			return creds[i].Domain < creds[j].Domain
		})
	}

	// Output in requested format
	switch outputFormat {
	case "json":
		// Ensure we output a valid JSON array even if creds is nil
		if creds == nil {
			creds = make([]kc.Credential, 0)
		}
		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(creds); err != nil {
			cmd.PrintErrf("Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
	case "csv":
		w := csv.NewWriter(cmd.OutOrStdout())
		if err := w.Write([]string{"Domain", "Username"}); err != nil {
			cmd.PrintErrf("Error writing CSV header: %v\n", err)
			os.Exit(1)
		}
		for _, cred := range creds {
			if err := w.Write([]string{cred.Domain, cred.Username}); err != nil {
				cmd.PrintErrf("Error writing CSV row: %v\n", err)
				os.Exit(1)
			}
		}
		w.Flush()
	default:
		if !quiet {
			if pattern != "" {
				cmd.Printf("Credentials matching '%s' (%d found):\n\n", pattern, len(creds))
			} else {
				cmd.Printf("Saved credentials (%d total):\n\n", len(creds))
			}
		}

		for i, cred := range creds {
			if quiet {
				cmd.Printf("%s\n", cred.Domain)
			} else {
				cmd.Printf("  %d. %s\n", i+1, cred.Domain)
				cmd.Printf("     Username: %s\n", cred.Username)
				if i < len(creds)-1 {
					cmd.Printf("\n")
				}
			}
		}

		if !quiet && len(creds) > 0 {
			cmd.Printf("\nTip: Use 'passkc get <domain>' to retrieve a password\n")
		}
	}
}

func newShowCmd(kcManager KeychainManager) *cobra.Command {
	runner := &showCmdRunner{
		kcManager: kcManager,
	}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "List all saved credentials",
		Long: `List all your saved credentials with their domains and usernames.

By default, shows a numbered list with domains and usernames.
Use flags to filter, sort, or change the output format.

Examples:
  passkc show                              # List all credentials
  passkc show --pattern github            # Search for credentials containing "github"
  passkc show --sort username             # Sort by username instead of domain
  passkc show -o json                     # Output as JSON
  passkc show -q                          # Quiet mode (domains only)`,
		Run: runner.run,
	}
	cmd.Flags().String("pattern", "", "Filter credentials by domain or username")
	cmd.Flags().String("sort", "", "Sort by field (domain|username)")
	return cmd
}

func init() {
	// The real command uses the live keychain manager.
	rootCmd.AddCommand(newShowCmd(&LiveKeychainManager{}))
}
