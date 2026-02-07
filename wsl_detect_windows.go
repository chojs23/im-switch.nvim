//go:build windows

package main

import (
	"os"
	"strings"
)

func isWSLFromWindowsProcess() bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	lower := strings.ToLower(wd)
	return strings.HasPrefix(lower, `\\wsl$\`) || strings.HasPrefix(lower, `\\wsl.localhost\`)
}
