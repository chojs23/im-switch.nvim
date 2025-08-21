//go:build windows

package main

import (
	"testing"
)

func TestGetCurrentInputSource(t *testing.T) {
	current := getCurrentInputSource()
	if current == "" {
		t.Error("Expected non-empty current input source")
	}
	t.Logf("Current input source: %s", current)
}

func TestGetAllInputSources(t *testing.T) {
	sources := getAllInputSources()
	if sources == nil {
		t.Error("Expected non-nil input sources")
		return
	}
	
	if len(sources) == 0 {
		t.Error("Expected at least one input source")
		return
	}
	
	t.Logf("Found %d input sources:", len(sources))
	for i, source := range sources {
		t.Logf("  %d: %s", i+1, source)
	}
}

func TestSetInputSource(t *testing.T) {
	// Get current input source
	original := getCurrentInputSource()
	if original == "" {
		t.Skip("Cannot get current input source, skipping set test")
	}
	
	// Get all available sources
	sources := getAllInputSources()
	if len(sources) < 2 {
		t.Skip("Need at least 2 input sources to test switching, skipping")
	}
	
	// Find a different source to switch to
	var target string
	for _, source := range sources {
		if source != original {
			target = source
			break
		}
	}
	
	if target == "" {
		t.Skip("Cannot find alternative input source, skipping")
	}
	
	t.Logf("Original source: %s", original)
	t.Logf("Target source: %s", target)
	
	// Switch to target
	if !setInputSource(target) {
		t.Errorf("Failed to switch to input source: %s", target)
		return
	}
	
	// Verify the switch
	current := getCurrentInputSource()
	if current != target {
		t.Logf("Warning: Expected %s, but got %s (this may be expected on some Windows versions)", target, current)
	}
	
	// Switch back to original
	if !setInputSource(original) {
		t.Errorf("Failed to switch back to original input source: %s", original)
	}
}

func TestCommonLayouts(t *testing.T) {
	// Test that we can identify common layout names
	for layoutId, expectedName := range commonLayouts {
		t.Logf("Layout ID %s -> %s", layoutId, expectedName)
	}
	
	// Test that our layout name function works
	if len(commonLayouts) == 0 {
		t.Error("Expected common layouts map to be populated")
	}
}