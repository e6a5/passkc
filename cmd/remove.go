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

type removeCmdRunner struct {
	kcManager KeychainManager
}

func (r *removeCmdRunner) run(cmd *cobra.Command, args []string) {
	domain := args[0]
	quiet, _ := cmd.Flags().GetBool("quiet")

	err := r.kcManager.RemoveData(domain)
	if err != nil {
		cmd.PrintErrf("Error removing credentials for %s: %v\n", domain, err)
		os.Exit(1)
	}

	if !quiet {
		cmd.Printf("Removed credentials for %s successfully\n", domain)
	}
}

func newRemoveCmd(kcManager KeychainManager) *cobra.Command {
	runner := &removeCmdRunner{
		kcManager: kcManager,
	}
	return &cobra.Command{
		Use:   "remove [domain]",
		Short: "Remove credentials for a domain from the Keychain",
		Long: `The passkc remove command removes the stored credentials for a specific domain.
This action is irreversible.

Examples:
  # Remove credentials for a domain
  passkc remove domain.com`,
		Args: cobra.ExactArgs(1),
		Run:  runner.run,
	}
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
