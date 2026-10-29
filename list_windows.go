package main

import (
	"strings"

	"github.com/danieljoos/wincred"
)

func secretKeys(account string) ([]string, error) {
	// Unlike List, FilteredList returns target names without namespace prefixes.
	credentials, err := wincred.FilteredList("*")
	if err != nil {
		return nil, err
	}
	var keys []string
	suffix := ":" + account
	for _, credential := range credentials {
		// go-keyring stores targets as service:account.
		if credential.UserName == account && strings.HasSuffix(credential.TargetName, suffix) {
			keys = append(keys, strings.TrimSuffix(credential.TargetName, suffix))
		}
	}
	return keys, nil
}
