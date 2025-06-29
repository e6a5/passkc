# Contributing to passkc

Thank you for your interest in contributing to passkc! This document provides guidelines and information for contributors.

## Philosophy

passkc follows a simple philosophy: **make password management disappear into practice**. Like a good tool, it should be:
- Simple and intuitive to use
- Secure by default
- Fast and reliable
- Unobtrusive in daily workflow

When contributing, please keep these principles in mind.

## Code of Conduct

This project adheres to a Code of Conduct that we expect all contributors to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before participating.

## Quick Start for Contributors

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/passkc.git
   cd passkc
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Build and test**:
   ```bash
   go build -v ./...
   go test -v -race ./...
   ```

## Development Workflow

### Prerequisites

- macOS (required for keychain integration)
- Go 1.23+

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards:
   - Follow Go conventions (gofmt, go vet)
   - Write tests for new functionality
   - Update documentation if needed
   - Keep commits focused and atomic

3. **Test your changes**:
   ```bash
   go test -v -race ./...                    # Run all tests
   go test -v -race -coverprofile=coverage.out ./...  # Generate coverage report
   golangci-lint run                         # Run linters
   gosec -conf=.gosec.json ./...            # Security scan
   ```

4. **Test manually**:
   ```bash
   go build -v ./...
   ./passkc --help     # Test CLI functionality
   ```

### Submitting Changes

1. **Push your branch** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** with:
   - Clear title describing the change
   - Detailed description of what and why
   - Link to any related issues
   - Screenshots for UI changes (if applicable)

## Types of Contributions

### 🐛 Bug Reports

Please use the bug report template and include:
- Steps to reproduce
- Expected vs actual behavior
- System information (macOS version, passkc version)
- Error messages or logs

### ✨ Feature Requests

We welcome feature requests, but please:
- Check existing issues first
- Explain the use case and problem being solved
- Consider if it aligns with our philosophy of simplicity
- Be open to alternative solutions

### 🔧 Code Contributions

Areas where contributions are especially welcome:
- Bug fixes
- Performance improvements
- Better error messages
- Additional output formats
- Import/export features
- Documentation improvements

### 📚 Documentation

Help us improve:
- README examples
- Command help text
- Code comments
- Contributing guides
- Blog posts or tutorials

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Write clear, concise comments
- Keep functions focused and small
- Handle errors appropriately

### Testing

- Write tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Mock external dependencies (keychain access)

### User Experience

- Commands should be intuitive and consistent
- Error messages should be helpful and actionable
- Default behavior should be secure
- Consider both interactive and scripting use cases

## Project Structure

```
passkc/
├── cmd/           # CLI commands (cobra)
├── kc/            # Core keychain functionality
├── main.go        # Entry point
├── go.mod         # Go module definition
├── .github/       # GitHub workflows and templates
├── .gosec.json    # Security scanning configuration
├── .golangci.yml  # Linting configuration
└── README.md      # User documentation
```

## Release Process

Releases are automated via GitHub Actions when tags are pushed:

1. **Version tagging**: We use semantic versioning (v1.2.3)
2. **Automated builds**: Cross-platform binaries are built automatically
3. **GitHub Releases**: Release notes are generated from tag annotations

Maintainers handle releases, but contributors can suggest when a release might be appropriate.

## Getting Help

- **Questions**: Open a discussion or issue
- **Chat**: Mention @e6a5 in issues or PRs
- **Security**: See [SECURITY.md](SECURITY.md) for security issues

## Recognition

Contributors are recognized in:
- Release notes
- README acknowledgments
- Git commit history

Thank you for contributing to passkc! 🔐 