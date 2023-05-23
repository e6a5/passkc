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

	"github.com/keybase/go-keychain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			domain   = "default"
			username string
			password string
		)
		if len(args) >= 1 && args[0] != "" {
			domain = args[0]
		}
		accessGroup := viper.GetString("access_group")
		username, _ = cmd.Flags().GetString("username")
		password, _ = cmd.Flags().GetString("password")

		item := keychain.NewItem()
		item.SetSecClass(keychain.SecClassInternetPassword)
		item.SetService(domain)
		item.SetAccount(username)
		item.SetAccessGroup(accessGroup)
		item.SetData([]byte(password))
		item.SetSynchronizable(keychain.SynchronizableNo)
		item.SetAccessible(keychain.AccessibleWhenUnlocked)
		err := keychain.AddItem(item)
		if err != nil {
			fmt.Printf("failed to set data for doamin <%s> error <%s>\n", domain, err.Error())
		} else {
			fmt.Printf("Your information for domain <%s> has been successfully saved\n", domain)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setCmd.PersistentFlags().StringP("username", "u", "", "username for domain")
	setCmd.PersistentFlags().StringP("password", "p", "", "password for domain")
}