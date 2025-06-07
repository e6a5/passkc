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
	"strings"

	"github.com/keybase/go-keychain"
)

// Credential holds the data for a keychain entry.
type Credential struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

// GetData retrieves credentials from the Keychain for a given domain.
// It will find the first entry matching the service "com.passkc.<domain>".
func GetData(domain string) (*Credential, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(fmt.Sprintf("com.passkc.%s", domain))
	query.SetMatchLimit(keychain.MatchLimitAll)

	results, err := keychain.QueryItem(query)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no credentials found for domain: %s", domain)
	}

	// Get the first result
	result := results[0]
	username := result.Account
	password := string(result.Data)

	return &Credential{
		Domain:   domain,
		Username: username,
		Password: password,
	}, nil
}

// SetData stores credentials in the Keychain.
// If an entry for the service and account already exists, it will be updated.
// If password is an empty string, the user will be prompted to enter it.
func SetData(domain, username, password string) error {
	service := fmt.Sprintf("com.passkc.%s.%s", domain, username)
	
	if password == "" {
		// Prompt for password if not provided
		fmt.Print("Enter password: ")
		var input string
		fmt.Scanln(&input)
		password = input
	}

	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(username)
	item.SetData([]byte(password))
	item.SetAccessible(keychain.AccessibleWhenUnlocked)
	item.SetSynchronizable(keychain.SynchronizableNo)

	err := keychain.AddItem(item)
	if err == keychain.ErrorDuplicateItem {
		// Update existing item
		query := keychain.NewItem()
		query.SetSecClass(keychain.SecClassGenericPassword)
		query.SetService(service)
		query.SetAccount(username)
		query.SetMatchLimit(keychain.MatchLimitOne)

		attributes := keychain.NewItem()
		attributes.SetData([]byte(password))

		err = keychain.UpdateItem(query, attributes)
	}
	return err
}

// RemoveData removes all credential entries for a given domain.
func RemoveData(domain string) error {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(fmt.Sprintf("com.passkc.%s", domain))
	query.SetMatchLimit(keychain.MatchLimitAll)

	return keychain.DeleteItem(query)
}

func ListData() ([]Credential, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService("com.passkc")
	query.SetMatchLimit(keychain.MatchLimitAll)

	results, err := keychain.QueryItem(query)
	if err == keychain.ErrorItemNotFound {
		return make([]Credential, 0), nil // Not an error, just no items
	}
	if err != nil {
		return nil, err
	}

	creds := make([]Credential, 0)
	for _, result := range results {
		// Parse domain and username from service name
		parts := strings.Split(result.Service, ".")
		if len(parts) >= 4 {
			domain := parts[2]
			username := result.Account
			creds = append(creds, Credential{
				Domain:   domain,
				Username: username,
			})
		}
	}

	return creds, nil
}
