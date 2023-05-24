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
	"fmt"
	"os"

	"github.com/keybase/go-keychain"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve username and password for a domain from the Keychain.",
	Long:  `The hiepass get command retrieves the stored username and password for a specific domain`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		accounts, err := keychain.GetAccountsForService(service)
		if err != nil {
			fmt.Printf("Failed to get data for <%s>.\n Error: <%s>.\n", service, err.Error())
			os.Exit(0)
		}
		if len(accounts) > 1 {
			fmt.Printf("Too many accounts for <%s>.\n", service)
			os.Exit(0)
		}
		if len(accounts) == 1 {
			label := getLabel(service, accounts[0])
			password, err := keychain.GetGenericPassword(service, accounts[0], label, "")
			if err == nil {
				clipboard.Init()
				clipboard.Write(clipboard.FmtText, password)
				fmt.Printf("Copied password for account <%s> in service <%s> to clipboard.\n", accounts[0], service)
				return
			}
		}
		fmt.Printf("No information for service <%s>.\n", service)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
