package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func secretKeys(account string) ([]string, error) {
	// No -d or -r: dump attributes only, never password data.
	out, err := exec.Command("/usr/bin/security", "dump-keychain").Output()
	if err != nil {
		return nil, fmt.Errorf("reading keychain metadata: %w", err)
	}
	return parseKeychainKeys(bytes.NewReader(out), account)
}

func parseKeychainKeys(in io.Reader, account string) ([]string, error) {
	var keys []string
	var class, user, service string
	var hasService bool
	flush := func() {
		if class == "genp" && user == account && hasService {
			keys = append(keys, service)
		}
		class, user, service, hasService = "", "", "", false
	}

	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "keychain:") || strings.HasPrefix(line, "class:") {
			flush()
			if strings.HasPrefix(line, "class:") {
				class = strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "class:")), "\"")
			}
			continue
		}
		if class != "genp" {
			continue
		}
		for _, attr := range []string{"acct", "svce"} {
			prefix := "\"" + attr + "\"<blob>="
			if !strings.HasPrefix(line, prefix) {
				continue
			}
			raw := strings.TrimSpace(strings.TrimPrefix(line, prefix))
			value, err := keychainAttribute(raw)
			if err != nil {
				// Do not include raw keychain output in errors.
				return nil, fmt.Errorf("invalid keychain %s attribute", attr)
			}
			if attr == "acct" {
				user = value
			} else {
				service, hasService = value, raw != "<NULL>"
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading keychain metadata: %w", err)
	}
	flush()
	return keys, nil
}

// security prints text as quoted strings, or hex with an optional text preview.
func keychainAttribute(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "<NULL>" {
		return "", nil
	}
	if strings.HasPrefix(value, "0x") {
		data, err := hex.DecodeString(strings.Fields(value)[0][2:])
		return string(data), err
	}
	return strconv.Unquote(value)
}
