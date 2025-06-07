# passkc

**passkc** is a Unix-style command-line tool for storing and retrieving credentials from the macOS Keychain. It follows Unix philosophy by doing one thing well, working with text streams, and being composable with other tools.

## Features

🔐 **Secure**: Uses macOS Keychain for secure credential storage
📝 **Text Streams**: Read from stdin, output to stdout, pipe with other commands
🔧 **Multiple Formats**: Support for text, JSON, and CSV output
🔍 **Filtering & Sorting**: Pattern matching and sorting capabilities
📁 **Batch Operations**: Import credentials from files
🤐 **Quiet Mode**: Script-friendly operation with minimal output
🧪 **Well Tested**: Comprehensive test coverage with dependency injection

## Installation

### Homebrew

```sh
brew install e6a5/tap/passkc
```

### From Source

```sh
# Using Go
go install github.com/e6a5/passkc@latest

# Using Task (recommended for development)
git clone https://github.com/e6a5/passkc.git
cd passkc
task build
```

### Requirements

- macOS (uses macOS Keychain)
- Go 1.21+ (for building from source)

## Usage

### Basic Commands

```bash
# Store credentials (interactive password prompt)
passkc set domain.com username

# Retrieve credentials
passkc get domain.com

# List all credentials
passkc show

# Modify credentials
passkc modify domain.com new-username

# Remove credentials
passkc remove domain.com
```

### Unix Philosophy Examples

**Text Stream Processing:**
```bash
# Read domain from stdin
echo "domain.com" | passkc get

# Pipe password to clipboard
passkc get domain.com -q | pbcopy

# Filter and search
passkc show | grep "google"
```

**Structured Output:**
```bash
# JSON output for programmatic use
passkc show -o json | jq '.[] | select(.domain == "google.com")'

# CSV export
passkc show -o csv > credentials.csv

# Quiet mode for scripts
PASSWORD=$(passkc get domain.com -q)
```

**Batch Operations:**
```bash
# Import from file (format: domain username [password])
echo "google.com user1 pass1" > creds.txt
echo "github.com user2 pass2" >> creds.txt
passkc set -f creds.txt

# Bulk export
passkc show -o json > backup.json
```

**Filtering and Sorting:**
```bash
# Filter by pattern
passkc show --pattern "google"

# Sort by domain or username
passkc show --sort domain
passkc show --sort username

# Combine filtering and output formats
passkc show --pattern "*.com" --sort username -o json
```

## Command Reference

### Global Flags

- `-o, --output string`: Output format (`text`|`json`|`csv`) (default: `text`)
- `-q, --quiet`: Suppress prompts and non-essential output
- `-c, --config string`: Config file (default: `$HOME/.passkc.yaml`)

### Commands

#### `passkc show`
List all stored credentials with filtering and sorting options.

**Flags:**
- `--pattern string`: Filter credentials by pattern
- `--sort string`: Sort by field (`domain`|`username`)

**Examples:**
```bash
passkc show
passkc show -o json
passkc show --pattern "google" --sort username
```

#### `passkc get [domain]`
Retrieve credentials for a domain. Can read domain from stdin if not provided.

**Flags:**
- `-p, --password-only`: Output only the password

**Examples:**
```bash
passkc get domain.com
passkc get domain.com -o json
echo "domain.com" | passkc get
```

#### `passkc set [domain] [username]`
Store credentials for a domain. Supports file input and stdin.

**Flags:**
- `-f, --file string`: Read credentials from file

**Examples:**
```bash
passkc set domain.com username
passkc set -f credentials.txt
echo "domain.com username password" | passkc set
```

#### `passkc modify [domain] [new-username]`
Update the username for a domain.

**Examples:**
```bash
passkc modify domain.com new-username
```

#### `passkc remove [domain]`
Remove credentials for a domain.

**Examples:**
```bash
passkc remove domain.com
```

## Development

### Prerequisites

- Go 1.21+
- Task (task runner)
- Git

### Setup

```bash
git clone https://github.com/e6a5/passkc.git
cd passkc
task deps  # Download dependencies
```

### Available Tasks

```bash
task                # List all available tasks
task build          # Build the binary
task test           # Run tests
task test-coverage  # Run tests with coverage report
task lint           # Run linters
task clean          # Clean build artifacts
task install        # Install binary to /usr/local/bin
task demo           # Showcase Unix philosophy features
```

### Testing

The project has comprehensive test coverage using dependency injection for mocking:

```bash
# Run all tests
task test

# Run tests with coverage
task test-coverage

# Run benchmarks
task test-bench
```

### Architecture

```
cmd/           # CLI commands with dependency injection
├── keychain.go    # KeychainManager interface and implementation
├── show.go        # List credentials command
├── get.go         # Retrieve credentials command
├── set.go         # Store credentials command
├── modify.go      # Modify credentials command
├── remove.go      # Remove credentials command
└── cli_test.go    # Comprehensive tests with mocks

kc/            # Keychain abstraction layer
└── kc.go          # Core keychain operations

main.go        # Application entry point
```

## Unix Philosophy Compliance

**passkc** follows Unix philosophy principles:

1. **Do one thing well**: Manage credentials in macOS Keychain
2. **Work together**: Pipe and chain with other Unix tools
3. **Handle text streams**: Read from stdin, write to stdout
4. **Plain text interface**: No captive user interfaces
5. **Composability**: Every command can be a filter

## Examples in Action

**Daily Workflow:**
```bash
# Quick password copy
passkc get work-email -q | pbcopy

# Search and filter
passkc show | grep -i "github\|gitlab" | head -5

# Backup to JSON
passkc show -o json > ~/backups/credentials-$(date +%Y%m%d).json

# Import from CSV
cat exported-passwords.csv | tail -n +2 | while IFS=, read domain user pass; do
  echo "$domain $user $pass" | passkc set
done
```

**Scripting:**
```bash
#!/bin/bash
# Check if credentials exist for multiple domains
DOMAINS=("github.com" "gitlab.com" "work.com")

for domain in "${DOMAINS[@]}"; do
  if passkc get "$domain" -q > /dev/null 2>&1; then
    echo "✓ $domain: credentials found"
  else
    echo "✗ $domain: no credentials"
  fi
done
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`task test`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request
