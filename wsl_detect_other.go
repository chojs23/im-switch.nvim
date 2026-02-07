//go:build !windows

package main

func isWSLFromWindowsProcess() bool {
	return false
}
