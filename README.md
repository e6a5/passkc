# passkc

**passkc** is a simple password manager for macOS that stores your passwords securely in the macOS Keychain.

## Quick Start

```bash
# Save a password
passkc set github.com myusername

# Show credentials (password hidden by default)
passkc get github.com

# Copy password to clipboard (secure)
passkc get github.com -q | pbcopy

# List all saved passwords
passkc show
```

## Installation

### Homebrew

```bash
brew install e6a5/tap/passkc
```

### From Source

```bash
go install github.com/e6a5/passkc@latest
```

### Requirements

- macOS (uses macOS Keychain for secure storage)
- Go 1.23+ (if building from source)

## Basic Usage

### Save a Password

```bash
# Interactive: prompts for username and password
passkc set github.com

# Specify username, prompts for password
passkc set github.com myusername
```

### Get a Password

```bash
# Show domain and username (password hidden for security)
passkc get github.com

# Get the password securely
passkc get github.com -p

# Copy password to clipboard (recommended)
passkc get github.com -q | pbcopy
```

### List Your Passwords

```bash
# List all saved passwords
passkc show

# Search for specific sites
passkc show --pattern github

# Sort by username
passkc show --sort username
```

### Update or Remove

```bash
# Change username and password
passkc modify github.com newusername

# Remove a password
passkc remove github.com
```

## Advanced Features

### Import Multiple Passwords

Create a file with your passwords:

```text
# credentials.txt
github.com myuser mypass123
google.com user@email.com secret456
work-vpn.com employee pass789
```

Import them:

```bash
passkc set -f credentials.txt
```

### JSON Output

```bash
# Get single credential as JSON
passkc get github.com -o json

# Export all credentials as JSON
passkc show -o json > backup.json
```

### Scripting

```bash
# Silent operation for scripts
PASSWORD=$(passkc get github.com -q)

# Check if password exists
if passkc get github.com -q > /dev/null 2>&1; then
    echo "Password found"
else
    echo "No password saved"
fi
```

## Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `passkc set <domain> [username]` | Save a password | `passkc set github.com` |
| `passkc get <domain>` | Show credentials (password hidden) | `passkc get github.com` |
| `passkc get <domain> -p` | Show password only | `passkc get github.com -p` |
| `passkc show` | List all passwords | `passkc show --pattern google` |
| `passkc modify <domain> <username>` | Update credentials | `passkc modify github.com newuser` |
| `passkc remove <domain>` | Delete a password | `passkc remove github.com` |

### Useful Flags

| Flag | Description | Example |
|------|-------------|---------|
| `-q, --quiet` | Silent output | `passkc get github.com -q` |
| `-p, --password-only` | Show only password | `passkc get github.com -p` |
| `-o, --output json` | JSON output | `passkc show -o json` |
| `--pattern <text>` | Filter results | `passkc show --pattern google` |
| `--sort <field>` | Sort by domain/username | `passkc show --sort username` |
| `-f, --force` | Skip confirmations | `passkc remove github.com -f` |

## Security

- 🔐 **Secure Storage**: Uses macOS Keychain, not plain text files
- 🔒 **Hidden Input**: Passwords are entered securely (not visible on screen)
- 🛡️ **System Integration**: Follows macOS security practices
- 🚫 **No Cloud**: Everything stays on your Mac

## Tips

**Secure password access (recommended):**
```bash
passkc get github.com -q | pbcopy    # Copy to clipboard without showing
```

**Show password when needed:**
```bash
passkc get github.com -p             # Display password only
```

**Search for a site:**
```bash
passkc show --pattern work
```

**Backup your passwords:**
```bash
passkc show -o json > ~/passwords-backup.json
```

**Add an alias for convenience:**
```bash
# Add to your ~/.zshrc or ~/.bashrc
alias pc='passkc'
alias pcg='passkc get'
alias pcs='passkc set'
```

## Getting Help

```bash
passkc --help              # General help
passkc set --help          # Help for specific command
```

## Development

```bash
git clone https://github.com/e6a5/passkc.git
cd passkc
task build                  # Build the binary
task test                   # Run tests
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Development

### Prerequisites

- macOS (required for keychain integration)
- Go 1.23+ (automatically detected from go.mod)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/e6a5/passkc.git
cd passkc

# Install dependencies
go mod download

# Build
go build -v ./...

# Run tests
go test -v -race ./...

# Format code
go fmt ./...

# Lint (install golangci-lint first)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run

# Security scanning
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec -conf=.gosec.json ./...
```

### Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

- 🐛 [Report bugs](https://github.com/e6a5/passkc/issues/new?template=bug_report.md)
- 💡 [Request features](https://github.com/e6a5/passkc/issues/new?template=feature_request.md)
- 🔧 [Submit pull requests](https://github.com/e6a5/passkc/pulls)

### Philosophy

passkc follows a simple philosophy: **make password management disappear into practice**.

Like a good tool, it should be:
- **Simple**: Intuitive commands that just work
- **Secure**: Safe by default, no compromises  
- **Fast**: Get your password and get back to work
- **Unobtrusive**: Fits seamlessly into your workflow

We believe the best password manager is one you don't have to think about.

## Community

- 📖 [Code of Conduct](CODE_OF_CONDUCT.md)
- 🔒 [Security Policy](SECURITY.md)
- 📝 [Changelog](CHANGELOG.md)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Thanks to all contributors who help make passkc better! 🔐
