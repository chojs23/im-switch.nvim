package main

import (
	"runtime"
	"testing"
)

func TestGetCurrentInputSource(t *testing.T) {
	current := getCurrentInputSource()
	if current == "" {
		t.Error("getCurrentInputSource() returned empty string")
	}
}

func TestGetAllInputSources(t *testing.T) {
	sources := getAllInputSources()
	if sources == nil {
		t.Error("getAllInputSources() returned nil")
		return
	}
	if len(sources) == 0 {
		t.Error("getAllInputSources() returned empty slice")
	}
}

func TestSetInputSourceWithInvalidID(t *testing.T) {
	invalidID := "invalid-input-source-id-that-should-not-exist"
	result := setInputSource(invalidID)
	if result {
		t.Error("setInputSource() should return false for invalid input source")
	}
}

func TestSetInputSourceWithCurrentID(t *testing.T) {
	current := getCurrentInputSource()
	if current == "" {
		t.Skip("Cannot get current input source, skipping test")
	}

	result := setInputSource(current)
	if !result {
		t.Error("setInputSource() should return true when setting to current input source")
	}
}

func TestPlatformSpecificFunctions(t *testing.T) {
	if runtime.GOOS == "darwin" {
		testDarwinFunctions(t)
	} else if runtime.GOOS == "linux" {
		testLinuxFunctions(t)
	} else {
		t.Skip("Unsupported platform")
	}
}

func testDarwinFunctions(t *testing.T) {
	current := getCurrentInputSource()
	if current == "" {
		t.Error("macOS getCurrentInputSource() returned empty string")
	}

	sources := getAllInputSources()
	if len(sources) == 0 {
		t.Error("macOS getAllInputSources() returned no sources")
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

func testLinuxFunctions(t *testing.T) {
	current := getCurrentInputSource()
	if current == "" {
		t.Error("Linux getCurrentInputSource() returned empty string")
	}

	sources := getAllInputSources()
	if len(sources) == 0 {
		t.Error("Linux getAllInputSources() returned no sources")
	}
}

func TestInputSourceToggling(t *testing.T) {
	originalSource := getCurrentInputSource()
	if originalSource == "" {
		t.Skip("Cannot get current input source, skipping toggle test")
	}

	sources := getAllInputSources()
	if len(sources) < 2 {
		t.Skip("Need at least 2 input sources for toggle test")
	}

	var alternativeSource string
	for _, source := range sources {
		if source != originalSource {
			alternativeSource = source
			break
		}
	}

	if alternativeSource == "" {
		t.Skip("Could not find alternative input source")
	}

	if !setInputSource(alternativeSource) {
		t.Error("Failed to switch to alternative input source")
		return
	}

	current := getCurrentInputSource()
	if current != alternativeSource {
		t.Errorf("Expected input source %s, got %s", alternativeSource, current)
	}

	if !setInputSource(originalSource) {
		t.Error("Failed to restore original input source")
	}
}

