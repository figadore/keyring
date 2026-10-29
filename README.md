# Keyring
A simple cross-platform keychain command line utility based on [zalando/go-keyring](https://github.com/zalando/go-keyring)

## Installation
### Global
go get github.com/figadore/keyring

### From source
go install

## Usage

### Example
The following example shows how to use `kr` in a curl command to access the GitHub API. If the `github-pat` secret does not exist yet, it will prompt for populating its value. Once populated, the secret is read out on subsequent calls without any additional action needed.
```
curl -u octocat:$(kr github-pat) https://api.github.com/users/octocat
```

### List existing keys

```sh
kr -l
```

Prints existing secret key names in alphabetical order, one per line, without
printing their values. Only entries belonging to this utility's `keyring` account
are included. An empty keyring produces no output. Listing cannot be combined
with a secret name or other options.

Supported on macOS (Keychain), Linux (Secret Service), and Windows (Credential
Manager), including secrets stored before this option was added.

### Help
kr -h
```
Usage of kr:
  -c    If set, show text while typing (for non-secrets only)
  -d    Delete
  -e    If set, just check for existence of key
  -l    List existing secret key names (not values)
  -prompt string
        Prompt string for secret, if it doesn't exist
  -s    If set, do not print value (useful for input-only scenarios)
```

