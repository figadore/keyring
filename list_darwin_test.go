package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestParseKeychainKeys(t *testing.T) {
	// Synthetic metadata, not a dump of the user's keychain. Data fields are
	// deliberately included to verify that only account/service attributes count.
	const metadata = `keychain: "/tmp/test.keychain-db"
version: 512
class: "genp"
attributes:
    "acct"<blob>="keyring"
    "svce"<blob>="github-pat"
data:
"NEVER PRINT THIS PASSWORD"
keychain: "/tmp/test.keychain-db"
class: "genp"
attributes:
    "svce"<blob>="other-app"
    "acct"<blob>="someone-else"
keychain: "/tmp/test.keychain-db"
class: "inet"
attributes:
    "acct"<blob>="keyring"
    "svce"<blob>="not-a-generic-password"
keychain: "/tmp/test.keychain-db"
class: "genp"
attributes:
    "svce"<blob>="missing-account"
keychain: "/tmp/test.keychain-db"
class: "genp"
attributes:
    "acct"<blob>="keyring"
keychain: "/tmp/test.keychain-db"
class: "genp"
attributes:
    "acct"<blob>="keyring"
    "svce"<blob>=<NULL>
keychain: "/tmp/test.keychain-db"
class: "genp"
attributes:
    "svce"<blob>=0x636166c3a9  "café"
    "acct"<blob>=0x6b657972696e67  "keyring"
class: "genp"
attributes:
    "acct"<blob>="keyring"
    "svce"<blob>="two words"
`
	got, err := parseKeychainKeys(strings.NewReader(metadata), "keyring")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"github-pat", "café", "two words"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("keys = %q, want %q", got, want)
	}
}

func TestKeychainAttribute(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{input: `"plain text"`, want: "plain text"},
		{input: `"escaped\\"quote"`, wantErr: true},
		{input: `"escaped\"quote"`, want: `escaped"quote`},
		{input: `0x636166c3a9 "preview"`, want: "café"},
		{input: `<NULL>`},
		{input: `0xZZ`, wantErr: true},
		{input: `unquoted`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := keychainAttribute(tt.input)
			if got != tt.want || (err != nil) != tt.wantErr {
				t.Fatalf("keychainAttribute = %q, %v; want %q, error %v", got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestParseKeychainKeysEmpty(t *testing.T) {
	keys, err := parseKeychainKeys(strings.NewReader(""), "keyring")
	if err != nil || len(keys) != 0 {
		t.Fatalf("keys = %q, error = %v", keys, err)
	}
}

func TestParseKeychainKeysMalformed(t *testing.T) {
	const metadata = "class: \"genp\"\n    \"acct\"<blob>=sensitive-invalid-attribute\n"
	_, err := parseKeychainKeys(strings.NewReader(metadata), "keyring")
	if err == nil || strings.Contains(err.Error(), "sensitive-invalid-attribute") {
		t.Fatalf("expected sanitized parse error, got %v", err)
	}
}

type failingReader struct{}

func (failingReader) Read(p []byte) (int, error) { return 0, errors.New("read failed") }

func TestParseKeychainKeysReadError(t *testing.T) {
	if _, err := parseKeychainKeys(failingReader{}, "keyring"); err == nil {
		t.Fatal("expected read error")
	}
}
