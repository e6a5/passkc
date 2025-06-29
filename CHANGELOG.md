# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enhanced security scanning with gosec configuration
- SARIF output format for security scan results  
- Dedicated gosec configuration file (.gosec.json)
- Improved CI workflow with GitHub CodeQL integration

### Changed
- Upgraded Go version requirement from 1.21 to 1.23
- Enhanced golangci-lint configuration with better gosec integration
- Updated all documentation to reflect Go 1.23 requirement
- Improved development workflow documentation

### Fixed
- Fixed gosec integration in CI workflow (corrected package path)
- Fixed misspelling: "Cancelled" → "Canceled"
- Resolved golangci-lint configuration version compatibility

## [1.0.4] - 2024-01-XX

### Added
- Secure password input (hidden from terminal)
- Interactive prompts for missing arguments
- Improved error messages with actionable guidance
- Visual success indicators (✓ symbols)
- Confirmation prompts for destructive actions

### Changed
- Default behavior: passwords hidden by default, use -p flag to show
- Enhanced user experience with visual indicators
- Simplified GitHub Actions workflows using standard Go tooling
- Removed Taskfile dependency in favor of standard Go commands

### Security
- Fixed service naming consistency in keychain operations
- Enhanced security by hiding passwords by default

### Fixed
- Resolved keychain API errors with proper SetReturnAttributes() calls

## [1.0.3] - Previous Release
### Changed
- Upgraded GitHub workflows

## [1.0.2] - Previous Release
### Changed
- Upgraded GitHub workflows

## [1.0.1] - Previous Release
### Changed
- Updated project with improvements

## [1.0.0] - Initial Release
### Added
- Initial release of passkc
- Basic password storage and retrieval
- macOS Keychain integration
- Command-line interface
