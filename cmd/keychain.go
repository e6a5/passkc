package cmd

import "github.com/e6a5/passkc/kc"

// KeychainManager defines the interface for interacting with the keychain.
// This allows for mocking in tests.
type KeychainManager interface {
	ListData() ([]kc.Credential, error)
	GetData(domain string) (*kc.Credential, error)
	SetData(domain, username, password string) error
	RemoveData(domain string) error
}

// LiveKeychainManager is the implementation that uses the real keychain.
type LiveKeychainManager struct{}

func (lkm *LiveKeychainManager) ListData() ([]kc.Credential, error) {
	return kc.ListData()
}

func (lkm *LiveKeychainManager) GetData(domain string) (*kc.Credential, error) {
	return kc.GetData(domain)
}

func (lkm *LiveKeychainManager) SetData(domain, username, password string) error {
	return kc.SetData(domain, username, password)
}

func (lkm *LiveKeychainManager) RemoveData(domain string) error {
	return kc.RemoveData(domain)
}

// newKeychainManager creates a new instance of the live keychain manager.
func newKeychainManager() KeychainManager {
	return &LiveKeychainManager{}
} 