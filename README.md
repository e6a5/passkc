# hiepass

**hiepass** is a command-line interface (CLI) tool for storing and retrieving username-password information from the Keychain on macOS. It provides a convenient way to securely manage and access credentials for various domains.

## Installation

## Usage

## Examples

```sh
# set new domain's information
hiepass set google.com e6a5
> Enter password: 
> Saved successfully 

# retrieve data 
hiepass get google.com
> Copied password for account <e6a5> in service <google.com> to clipboard.

# modify information
hiepass modify google.com tranhiep
> Enter password: 
> Updated successfully

# remove domain's information
hiepass remove google.com
> Removed successfully

```

## Contributing
Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue on the GitHub repository.

## License

This project is licensed under the MIT License. See the [LICENSE](#LICENSE) file for details.
