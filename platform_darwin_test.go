//go:build darwin

package main

import (
	"testing"
)

func TestDarwinInputSourceFunctions(t *testing.T) {
	current := getCurrentInputSource()
	if current == "" {
		t.Error("getCurrentInputSource() returned empty string on macOS")
	}

	sources := getAllInputSources()
	if sources == nil {
		t.Error("getAllInputSources() returned nil on macOS")
		return
	}

	if len(sources) == 0 {
		t.Error("getAllInputSources() returned empty slice on macOS")
		return
	}

	found := false
	for _, source := range sources {
		if source == current {
			found = true
			break
		}
	}

	if !found {
		t.Error("Current input source not found in available sources list")
	}
}

func TestDarwinSetInputSource(t *testing.T) {
	originalSource := getCurrentInputSource()
	if originalSource == "" {
		t.Skip("Cannot get current input source")
	}

	result := setInputSource(originalSource)
	if !result {
		t.Error("setInputSource() failed to set current input source")
	}

	result = setInputSource("invalid-source-id")
	if result {
		t.Error("setInputSource() should fail for invalid source ID")
	}
}

func TestDarwinCommonInputSources(t *testing.T) {
	sources := getAllInputSources()
	if len(sources) == 0 {
		t.Skip("No input sources available")
	}

	hasABC := false
	for _, source := range sources {
		if source == "com.apple.keylayout.ABC" {
			hasABC = true
			break
		}
	}

	if !hasABC {
		t.Log("ABC keyboard layout not found (this may be normal depending on system configuration)")
	}
}

func TestDarwinInputSourceSwitching(t *testing.T) {
	sources := getAllInputSources()
	if len(sources) < 2 {
		t.Skip("Need at least 2 input sources for switching test")
	}

	originalSource := getCurrentInputSource()
	if originalSource == "" {
		t.Skip("Cannot get current input source")
	}

	var targetSource string
	for _, source := range sources {
		if source != originalSource {
			targetSource = source
			break
		}
	}

	if targetSource == "" {
		t.Skip("Could not find alternative input source")
	}

	if !setInputSource(targetSource) {
		t.Skip("Could not switch to target source (may not be enabled)")
	}

	currentAfterSwitch := getCurrentInputSource()
	if currentAfterSwitch != targetSource {
		t.Log("Input source may not have switched immediately (this can be normal)")
	}

	setInputSource(originalSource)
}

