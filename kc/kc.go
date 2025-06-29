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
	"syscall"

	"github.com/keybase/go-keychain"
	"golang.org/x/term"
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
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnAttributes(true)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		if err == keychain.ErrorItemNotFound {
			return nil, fmt.Errorf("no credentials found for '%s'. Use 'passkc set %s <username>' to add credentials", domain, domain)
		}
		return nil, fmt.Errorf("failed to access keychain: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no credentials found for '%s'. Use 'passkc set %s <username>' to add credentials", domain, domain)
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
// If password is an empty string, the user will be prompted to enter it securely.
func SetData(domain, username, password string) error {
	// Fixed: Use consistent service naming scheme
	service := fmt.Sprintf("com.passkc.%s", domain)

	if password == "" {
		// Secure password prompt
		fmt.Printf("Enter password for %s@%s: ", username, domain)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		fmt.Println() // Add newline after password input
		password = string(bytePassword)

		if password == "" {
			return fmt.Errorf("password cannot be empty")
		}
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
		if err != nil {
			return fmt.Errorf("failed to update credentials for '%s': %v", domain, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to save credentials for '%s': %v", domain, err)
	}

	return nil
}

// RemoveData removes all credential entries for a given domain.
func RemoveData(domain string) error {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(fmt.Sprintf("com.passkc.%s", domain))
	query.SetMatchLimit(keychain.MatchLimitOne)

	err := keychain.DeleteItem(query)
	if err == keychain.ErrorItemNotFound {
		return fmt.Errorf("no credentials found for '%s'", domain)
	}
	if err != nil {
		return fmt.Errorf("failed to remove credentials for '%s': %v", domain, err)
	}

	return nil
}

func ListData() ([]Credential, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetMatchLimit(keychain.MatchLimitAll)
	query.SetReturnAttributes(true)

	// Search for all items with service names starting with "com.passkc."
	results, err := keychain.QueryItem(query)
	if err == keychain.ErrorItemNotFound {
		return make([]Credential, 0), nil // Not an error, just no items
	}
	if err != nil {
		return nil, fmt.Errorf("failed to access keychain: %v", err)
	}

	creds := make([]Credential, 0)
	for _, result := range results {
		// Parse domain from service name: com.passkc.<domain>
		if strings.HasPrefix(result.Service, "com.passkc.") {
			domain := strings.TrimPrefix(result.Service, "com.passkc.")
			username := result.Account
			creds = append(creds, Credential{
				Domain:   domain,
				Username: username,
			})
		}
	}

	return creds, nil
}

// Helper function for better user input validation
func ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}
	if strings.Contains(domain, " ") {
		return fmt.Errorf("domain cannot contain spaces")
	}
	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	return nil
}
