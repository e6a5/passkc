# passkc

**passkc** is a command-line interface (CLI) tool for storing and retrieving username-password information from the Keychain on macOS. It provides a convenient way to securely manage and access credentials for various domains.

## Installation

## Usage

passkc supports the following commands:

- **Show**: Show list of labels.

```bash
  passkc show

```

- **Set**: Store a username for a domain in the Keychain.

```bash
  passkc set <domain> <username>

```
- **Get**: Retrieve the stored username for a domain from the Keychain.

```bash 
passkc get <domain>
```

- **Modify**: Update the username for a domain in the Keychain (optional: include a new username).

```bash
passkc modify <domain> <new_username>
```

- **Remove**: Remove a domain and its associated credentials from the Keychain.

```bash
passkc remove <domain>
```

Make sure to replace `<domain>` and `<username>` with the specific domain and username you want to store, retrieve, modify, or remove in the Keychain.

## Examples

```bash
# set new domain's information
> passkc set google.com e6a5
Enter password: 
Saved successfully 

# retrieve data 
>passkc get google.com
Copied password for account <e6a5> in service <google.com> to clipboard.

# modify information
> passkc modify google.com tranhiep
Enter password: 
Updated successfully

# remove domain's information
> passkc remove google.com
Removed successfully

# show list of labels
> passkc show
List of labels:
com.passkc.google.com.e6a5 #<app_prefix>.<domain>.<username>
```

## Contributing
Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue on the GitHub repository.

## License

This project is licensed under the MIT License. See the [LICENSE](#LICENSE) file for details.
