# Keyring
A simple cross-platform keychain command line utility

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

### Help
kr -h
```
Usage of kr:
  -c    If set, show text while typing (for non-secrets only)
  -d    Delete
  -e    If set, just check for existence of key
  -prompt string
        Prompt string for secret, if it doesn't exist
  -s    If set, do not print value (useful for input-only scenarios)
```

