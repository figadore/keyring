package main

import (
	"flag"
	"io"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		list    bool
		service string
	}{
		{name: "list", args: []string{"-l"}, list: true},
		{name: "fetch", args: []string{"github-pat"}, service: "github-pat"},
		{name: "existence", args: []string{"-e", "github-pat"}, service: "github-pat"},
		{name: "delete", args: []string{"-d", "github-pat"}, service: "github-pat"},
		{name: "missing key", wantErr: true},
		{name: "list with key", args: []string{"-l", "github-pat"}, wantErr: true},
		{name: "list with delete", args: []string{"-l", "-d"}, wantErr: true},
		{name: "list with existence", args: []string{"-l", "-e"}, wantErr: true},
		{name: "list with silent", args: []string{"-l", "-s"}, wantErr: true},
		{name: "list with cleartext", args: []string{"-l", "-c"}, wantErr: true},
		{name: "list with prompt", args: []string{"-l", "-prompt", "secret?"}, wantErr: true},
		{name: "unknown flag", args: []string{"-unknown"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := flag.NewFlagSet("kr", flag.ContinueOnError)
			flags.SetOutput(io.Discard)
			opts, err := parseFlags(flags, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseFlags error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if opts.list != tt.list || opts.service != tt.service || opts.account != "keyring" {
				t.Fatalf("unexpected options: %+v", opts)
			}
			if !opts.list && opts.prompt != "Secret for '"+tt.service+"' not found. Enter it now:" {
				t.Fatalf("unexpected default prompt: %q", opts.prompt)
			}
		})
	}
}
