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

type modifyCmdRunner struct {
	kcManager KeychainManager
}

func (r *modifyCmdRunner) run(cmd *cobra.Command, args []string) {
	domain := args[0]
	newUsername := args[1]
	quiet, _ := cmd.Flags().GetBool("quiet")

	// We use SetData which will prompt for a password and update if the item exists
	err := r.kcManager.SetData(domain, newUsername, "")
	if err != nil {
		cmd.PrintErrf("Error modifying credentials for %s: %v\n", domain, err)
		os.Exit(1)
	}

	if !quiet {
		cmd.Printf("Modified credentials for %s successfully\n", domain)
	}
}

func newModifyCmd(kcManager KeychainManager) *cobra.Command {
	runner := &modifyCmdRunner{
		kcManager: kcManager,
	}
	return &cobra.Command{
		Use:   "modify [domain] [new-username]",
		Short: "Modify the username for a domain in the Keychain",
		Long: `The passkc modify command updates the username for a specific domain.
It will prompt for a new password. If you want to change the password, use 'set'.

Examples:
  # Modify the username for a domain
  passkc modify domain.com new-username`,
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
