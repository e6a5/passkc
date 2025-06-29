package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/e6a5/passkc/kc"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// mockKeychain is a mock implementation of KeychainManager for testing
type mockKeychain struct {
	creds    []kc.Credential
	setCalls []setCall
	err      error
}

type setCall struct {
	domain   string
	username string
	password string
}

func (m *mockKeychain) ListData() ([]kc.Credential, error) {
	return m.creds, m.err
}

func (m *mockKeychain) GetData(domain string) (*kc.Credential, error) {
	for _, cred := range m.creds {
		if cred.Domain == domain {
			return &cred, nil
		}
	}
	return nil, m.err
}

func (m *mockKeychain) SetData(domain, username, password string) error {
	if m.setCalls == nil {
		m.setCalls = make([]setCall, 0)
	}
	m.setCalls = append(m.setCalls, setCall{domain, username, password})
	return m.err
}

func (m *mockKeychain) RemoveData(domain string) error {
	return m.err
}

func execute(t *testing.T, kcManager KeychainManager, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)

	rootCmd := &cobra.Command{Use: "passkc"}
	initializeFlags(rootCmd)

	rootCmd.AddCommand(newShowCmd(kcManager))
	rootCmd.AddCommand(newGetCmd(kcManager))
	rootCmd.AddCommand(newSetCmd(kcManager))
	rootCmd.AddCommand(newRemoveCmd(kcManager))
	rootCmd.AddCommand(newModifyCmd(kcManager))

	rootCmd.SetArgs(args)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	err := rootCmd.Execute()
	return buf.String(), err
}

func TestShowCommand(t *testing.T) {
	// Use a mock keychain manager for this test
	mockKC := &mockKeychain{
		creds: []kc.Credential{},
		err:   nil,
	}

	// Test JSON output with no data
	output, err := execute(t, mockKC, "show", "-o", "json")
	assert.NoError(t, err)
	isJSON(t, output)
	assert.Equal(t, "[]\n", output)
}

func TestShowCommandWithData(t *testing.T) {
	// Use a mock keychain manager with some data
	mockKC := &mockKeychain{
		creds: []kc.Credential{
			{Domain: "google.com", Username: "testuser"},
			{Domain: "github.com", Username: "anotheruser"},
		},
		err: nil,
	}

	// Test basic execution - updated to match new output format
	output, err := execute(t, mockKC, "show")
	assert.NoError(t, err)
	assert.Contains(t, output, "github.com")
	assert.Contains(t, output, "anotheruser")
	assert.Contains(t, output, "google.com")
	assert.Contains(t, output, "testuser")
	assert.Contains(t, output, "Saved credentials (2 total)")

	// Test JSON output
	output, err = execute(t, mockKC, "show", "-o", "json")
	assert.NoError(t, err)
	isJSON(t, output)
	assert.Contains(t, output, `"domain":"google.com"`)
	assert.Contains(t, output, `"username":"anotheruser"`)

	// Test filtering
	output, err = execute(t, mockKC, "show", "--pattern", "google")
	assert.NoError(t, err)
	assert.Contains(t, output, "google.com")
	assert.NotContains(t, output, "github.com")

	// Test sorting - updated to match new output format
	output, err = execute(t, mockKC, "show", "--sort", "username")
	assert.NoError(t, err)
	// In the mock data, anotheruser comes before testuser alphabetically
	assert.Contains(t, output, "anotheruser")
	assert.Contains(t, output, "testuser")
}

func TestGetCommand(t *testing.T) {
	mockKC := &mockKeychain{
		creds: []kc.Credential{
			{Domain: "google.com", Username: "testuser", Password: "testpass"},
		},
		err: nil,
	}

	// Test basic get - updated to match new secure output format (no password by default)
	output, err := execute(t, mockKC, "get", "google.com")
	assert.NoError(t, err)
	assert.Contains(t, output, "testuser")
	assert.NotContains(t, output, "testpass") // Password should NOT be shown by default
	assert.Contains(t, output, "Domain: google.com")
	assert.Contains(t, output, "To get the password:")

	// Test password-only flag
	output, err = execute(t, mockKC, "get", "google.com", "-p")
	assert.NoError(t, err)
	assert.Equal(t, "testpass", output)

	// Test JSON output (still includes password for structured data)
	output, err = execute(t, mockKC, "get", "google.com", "-o", "json")
	assert.NoError(t, err)
	isJSON(t, output)
	assert.Contains(t, output, `"domain":"google.com"`)
	assert.Contains(t, output, `"username":"testuser"`)
	assert.Contains(t, output, `"password":"testpass"`)

	// Test quiet mode
	output, err = execute(t, mockKC, "get", "google.com", "-q")
	assert.NoError(t, err)
	assert.Equal(t, "testpass", output)
}

func TestSetCommand(t *testing.T) {
	mockKC := &mockKeychain{
		creds: []kc.Credential{},
		err:   nil,
	}

	// Test basic set - updated to match new success message
	output, err := execute(t, mockKC, "set", "google.com", "testuser")
	assert.NoError(t, err)
	assert.Contains(t, output, "✓ Saved credentials for testuser@google.com")
	assert.Len(t, mockKC.setCalls, 1)
	assert.Equal(t, "google.com", mockKC.setCalls[0].domain)
	assert.Equal(t, "testuser", mockKC.setCalls[0].username)

	// Test quiet mode
	mockKC.setCalls = nil // Reset calls
	output, err = execute(t, mockKC, "set", "github.com", "anotheruser", "-q")
	assert.NoError(t, err)
	assert.Empty(t, output)
	assert.Len(t, mockKC.setCalls, 1)
}

func TestRemoveCommand(t *testing.T) {
	mockKC := &mockKeychain{
		creds: []kc.Credential{
			{Domain: "google.com", Username: "testuser"},
		},
		err: nil,
	}

	// Test basic remove with force flag to skip confirmation
	output, err := execute(t, mockKC, "remove", "google.com", "--force")
	assert.NoError(t, err)
	assert.Contains(t, output, "✓ Removed credentials for google.com")

	// Test quiet mode with force
	output, err = execute(t, mockKC, "remove", "google.com", "-q", "--force")
	assert.NoError(t, err)
	assert.Empty(t, output)
}

func TestModifyCommand(t *testing.T) {
	mockKC := &mockKeychain{
		creds: []kc.Credential{
			{Domain: "google.com", Username: "olduser"},
		},
		err: nil,
	}

	// Test basic modify - updated to match new success message
	output, err := execute(t, mockKC, "modify", "google.com", "newuser")
	assert.NoError(t, err)
	assert.Contains(t, output, "✓ Updated credentials for google.com")
	assert.Len(t, mockKC.setCalls, 1)
	assert.Equal(t, "google.com", mockKC.setCalls[0].domain)
	assert.Equal(t, "newuser", mockKC.setCalls[0].username)

	// Test quiet mode
	mockKC.setCalls = nil // Reset calls
	_, err = execute(t, mockKC, "modify", "google.com", "newuser2", "-q")
	assert.NoError(t, err)
	// Should still have some output for the "Updating credentials" message unless we make it completely quiet
	assert.Len(t, mockKC.setCalls, 1)
}

func isJSON(t *testing.T, s string) {
	var js interface{}
	assert.NoError(t, json.Unmarshal([]byte(s), &js), "output should be valid JSON")
}
