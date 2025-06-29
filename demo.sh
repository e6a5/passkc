#!/bin/bash

# Demo script showing how simple passkc is to use
# This demonstrates the core workflow: store, retrieve, list, and manage passwords

echo "🔐 passkc Demo - Simple Password Management for macOS"
echo "=================================================="
echo ""

# Build the latest version
echo "Building passkc..."
go build -o passkc .
echo "✓ Build complete"
echo ""

# Show help
echo "📚 Quick help:"
./passkc --help
echo ""

# Test showing empty state
echo "📋 Checking current passwords:"
./passkc show
echo ""

# Show examples of usage without running (since they would require keychain access)
echo "🎯 Basic Usage Examples:"
echo ""
echo "Save a password (interactive):"
echo "  passkc set github.com"
echo ""
echo "Save a password (with username):"
echo "  passkc set github.com myusername"
echo ""
echo "Get a password:"
echo "  passkc get github.com"
echo ""
echo "Copy password to clipboard:"
echo "  passkc get github.com -q | pbcopy"
echo ""
echo "List all passwords:"
echo "  passkc show"
echo ""
echo "Search for passwords:"
echo "  passkc show --pattern github"
echo ""
echo "Update credentials:"
echo "  passkc modify github.com newusername"
echo ""
echo "Remove a password:"
echo "  passkc remove github.com"
echo ""

echo "✨ Key Features:"
echo "  • Secure storage in macOS Keychain"
echo "  • Hidden password input (no echo to screen)"
echo "  • Simple, intuitive commands"
echo "  • Helpful error messages and guidance"
echo "  • Works great in scripts and automation"
echo "  • JSON export for backups"
echo ""

echo "🎉 passkc makes password management simple and secure!" 