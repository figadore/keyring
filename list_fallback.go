//go:build !darwin && !linux && !windows
// +build !darwin,!linux,!windows

package main

import "fmt"

func secretKeys(account string) ([]string, error) {
	return nil, fmt.Errorf("listing secret keys is not supported on this platform")
}
