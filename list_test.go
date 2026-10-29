package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestWriteSecretKeys(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want string
	}{
		{name: "empty"},
		{name: "sorted and unique", keys: []string{"zebra", "github-pat", "alpha", "github-pat"}, want: "alpha\ngithub-pat\nzebra\n"},
		{name: "spaces and unicode", keys: []string{"two words", "café"}, want: "café\ntwo words\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := writeSecretKeys(&out, tt.keys); err != nil {
				t.Fatal(err)
			}
			if got := out.String(); got != tt.want {
				t.Fatalf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

type failingWriter struct{ err error }

func (w failingWriter) Write(p []byte) (int, error) { return 0, w.err }

func TestWriteSecretKeysError(t *testing.T) {
	want := errors.New("write failed")
	if err := writeSecretKeys(failingWriter{want}, []string{"github-pat"}); !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}
