package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh/terminal"
)

type options struct {
	account        string
	service        string
	prompt         string
	silent         bool
	checkExistence bool
	clearText      bool
	del            bool
}

// Read secret from OS keyring. If not found, get from stdin, store, and then continue
func main() {
	// Usage: keyring [-d] [-prompt "<prompt>"] [-e] <key>
	// Retrieve bitbucket username with `keyring bitbucket_user`
	// Use -e to check for existence (exit 0 if checkExistence, exit 4 if not)
	// Include -d to delete

	opts, err := parseArgs()
	if err != nil {
		log.Fatalf("Error parsing args: %v", err)
	}
	// Do the rest in a testable function
	checkFetchDel(opts)
}

func checkFetchDel(opts *options) {
	if opts.del {
		deleteSecret(opts)
		return
	}
	secret, err := keyring.Get(opts.service, opts.account)
	if err == keyring.ErrNotFound {
		if opts.checkExistence {
			os.Exit(4)
		}
		secret = promptAndStore(opts)
	} else if err != nil {
		log.Fatalln("Error getting secret:", err)
	}

	if opts.checkExistence {
		os.Exit(0)
	}
	if !opts.silent {
		fmt.Println(secret)
	}
}

// Handle any complicated arg parsing here
func parseArgs() (*options, error) {
	// go-keyring module uses two keys to look up a secret, 'account' and
	// 'service'. This appears to have been done for cross-platform
	// compatibility, so it seems best to hardcode the account, and store the
	// username and password as two separate secrets
	opts := options{}
	opts.account = "keyring"
	flag.BoolVar(&opts.clearText, "c", false, "If set, show text while typing (for non-secrets only)")
	flag.BoolVar(&opts.silent, "s", false, "If set, do not print value (useful for input-only scenarios)")
	flag.BoolVar(&opts.checkExistence, "e", false, "If set, just check for existence of key")
	flag.BoolVar(&opts.del, "d", false, "Delete")
	flag.StringVar(&opts.prompt, "prompt", "", "Prompt string for secret, if it doesn't exist")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return &options{}, fmt.Errorf("expecting at least 1 arg (secret_name), found %d", len(args))
	}
	opts.service = args[0]
	if opts.prompt == "" {
		// No prompt specified, use default
		opts.prompt = fmt.Sprintf("Secret for '%s' not found. Enter it now:", opts.service)
	}
	return &opts, nil
}

func promptAndStore(opts *options) string {
	var secret string
	var e error
	fmt.Print(opts.prompt)
	if opts.clearText {
		_, e = fmt.Scanln(&secret)
	} else {
		var secretBytes []byte
		secretBytes, e = terminal.ReadPassword(int(os.Stdin.Fd()))
		secret = string(secretBytes)
		// Start new line after prompt
		fmt.Println("")
	}

	if e != nil {
		log.Fatalln("Error reading secret from stdin:", e)
	}
	e = keyring.Set(opts.service, opts.account, secret)
	if e != nil {
		// Secret is in memory, but unable to write it to the OS keyring. Continue with warning
		log.Println("Error writing secret:", e)
	}
	return secret
}

func deleteSecret(opts *options) {
	fmt.Println("Deleting secret for", opts.service)
	err := keyring.Delete(opts.service, opts.account)
	if err != nil {
		log.Fatalln("Error deleting secret:", err)
	}
}
