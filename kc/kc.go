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
package kc

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/keybase/go-keychain"
	"golang.design/x/clipboard"
	"golang.org/x/crypto/ssh/terminal"
)

func GetData(service string) {
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
}

func SetData(service, account string) {
	label := getLabel(service, account)
	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}

	item := keychain.NewGenericPassword(service, account, label, bytePassword, "")
	err = keychain.AddItem(item)
	if err != nil {
		fmt.Printf("Failed to set data for <%s>.\n Error: <%s>\n", service, err.Error())
		os.Exit(0)
	}
	writeLabelToFile(label)
	fmt.Println("Saved successfully")
}

func ModifyData(service, account string) {
	accounts, err := keychain.GetAccountsForService(service)
	if err != nil {
		fmt.Printf("Failed to get data for <%s>.\n Error: <%s>.\n", service, err.Error())
		os.Exit(0)
	}
	if len(accounts) > 1 {
		fmt.Printf("Too many accounts for <%s>.\n", service)
		os.Exit(0)
	}
	if len(accounts) == 0 {
		fmt.Printf("No information for service <%s>.\n", service)
		os.Exit(0)
	}
	if len(accounts) == 1 && account == "" {
		account = accounts[0]
	}
	oldLabel := getLabel(service, accounts[0])
	newLabel := getLabel(service, account)
	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	queryItem := keychain.NewGenericPassword(service, accounts[0], oldLabel, []byte{}, "")
	updateItem := keychain.NewItem()
	updateItem.SetSecClass(keychain.SecClassGenericPassword)
	updateItem.SetService(service)
	updateItem.SetAccount(account)
	updateItem.SetLabel(newLabel)
	updateItem.SetData(bytePassword)
	err = keychain.UpdateItem(queryItem, updateItem)
	if err != nil {
		fmt.Printf("Failed to modify data for <%s>.\n Error: <%s>\n", service, err.Error())
		return
	}
	fmt.Println("Updated successfully")
}

func RemoveService(service string) {
	accounts, err := keychain.GetAccountsForService(service)
	if err != nil {
		fmt.Printf("Failed to get data for <%s>.\n Error: <%s>.\n", service, err.Error())
		os.Exit(0)
	}
	if len(accounts) == 0 {
		fmt.Printf("No information for service <%s>.\n", service)
		os.Exit(0)
	}
	if len(accounts) > 1 {
		fmt.Printf("Too many accounts for <%s>.\n", service)
		os.Exit(0)
	}
	label := getLabel(service, accounts[0])
	queryItem := keychain.NewGenericPassword(service, accounts[0], label, []byte{}, "")
	err = keychain.DeleteItem(queryItem)
	if err != nil {
		fmt.Printf("Failed to remove <%s>.\n Error: <%s>\n", service, err.Error())
		os.Exit(0)
	}
	deleteLableInFile(label)
	fmt.Println("Removed successfully")
}

func ShowLabels() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	filePath := filepath.Join(homeDir, LABEL_FILE_NAME)
	// Read the existing file contents
	labels, err := readLabelsFromFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(labels) == 0 {
		fmt.Println("No labels")
		return
	}
	// Print the list of domains
	fmt.Println("List of labels:")
	for _, label := range labels {
		fmt.Println(label)
	}
}

func getLabel(serviceName, account string) string {
	return fmt.Sprintf("%s.%s.%s", LABEL_PREFIX, serviceName, account)
}
