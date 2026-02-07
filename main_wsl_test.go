package main

import (
	"errors"
	"testing"
)

func TestDetectWSL(t *testing.T) {
	tests := []struct {
		name      string
		goos      string
		env       map[string]string
		osrelease string
		readErr   error
		want      bool
	}{
		{
			name: "windows process in wsl via env",
			goos: "windows",
			env: map[string]string{
				"WSL_INTEROP": "1",
			},
			want: true,
		},
		{
			name:      "linux wsl via osrelease",
			goos:      "linux",
			osrelease: "5.15.153.1-microsoft-standard-WSL2",
			want:      true,
		},
		{
			name:      "linux non wsl",
			goos:      "linux",
			osrelease: "6.8.0-48-generic",
			want:      false,
		},
		{
			name:    "linux osrelease unavailable",
			goos:    "linux",
			readErr: errors.New("read error"),
			want:    false,
		},
		{
			name: "windows non wsl",
			goos: "windows",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookup := func(key string) (string, bool) {
				v, ok := tt.env[key]
				return v, ok
			}
			read := func(path string) ([]byte, error) {
				if tt.readErr != nil {
					return nil, tt.readErr
				}
				return []byte(tt.osrelease), nil
			}

			got := detectWSL(tt.goos, lookup, read)
			if got != tt.want {
				t.Fatalf("detectWSL(%q) = %v, want %v", tt.goos, got, tt.want)
			}
		})
	}
}
