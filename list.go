package main

import (
	"fmt"
	"io"
	"sort"
)

// Only names reach the output layer.
func writeSecretKeys(out io.Writer, keys []string) error {
	sort.Strings(keys)
	for i, key := range keys {
		if i > 0 && key == keys[i-1] {
			continue
		}
		if _, err := fmt.Fprintln(out, key); err != nil {
			return err
		}
	}
	return nil
}
