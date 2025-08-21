//go:build linux

package main

import (
	"testing"
)

func TestDetectInputMethod(t *testing.T) {
	method := detectInputMethod()
	validMethods := []string{"ibus", "fcitx", "fcitx5", "xkb"}

	found := false
	for _, valid := range validMethods {
		if method == valid {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("detectInputMethod() returned invalid method: %s", method)
	}
}

func TestIsProcessRunning(t *testing.T) {
	result := isProcessRunning("nonexistent-process-12345")
	if result {
		t.Error("isProcessRunning() should return false for non-existent process")
	}

	result = isProcessRunning("init")
	if !result {
		t.Error("isProcessRunning() should return true for init process")
	}
}

func TestLinuxInputMethods(t *testing.T) {
	testCases := []struct {
		method string
		getter func() string
		lister func() []string
		setter func(string) bool
	}{
		{"ibus", getCurrentInputSourceIBus, getAllInputSourcesIBus, setInputSourceIBus},
		{"fcitx", getCurrentInputSourceFcitx, getAllInputSourcesFcitx, setInputSourceFcitx},
		{"fcitx5", getCurrentInputSourceFcitx5, getAllInputSourcesFcitx5, setInputSourceFcitx5},
		{"xkb", getCurrentInputSourceXKB, getAllInputSourcesXKB, setInputSourceXKB},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			current := tc.getter()
			sources := tc.lister()

			if sources == nil {
				t.Errorf("%s getAllInputSources returned nil", tc.method)
				return
			}

			if len(sources) == 0 {
				t.Errorf("%s getAllInputSources returned empty slice", tc.method)
				return
			}

			validSource := sources[0]
			result := tc.setter(validSource)
			if tc.method != "xkb" && !result {
				t.Logf("%s setInputSource may have failed (this is expected if %s is not running)", tc.method, tc.method)
			}
		})
	}
}

func TestXKBInputSources(t *testing.T) {
	sources := getAllInputSourcesXKB()
	expectedSources := []string{"us", "gb", "de", "fr", "es", "it", "ru", "cn", "jp", "kr"}

	if len(sources) != len(expectedSources) {
		t.Errorf("Expected %d sources, got %d", len(expectedSources), len(sources))
	}

	for i, expected := range expectedSources {
		if i >= len(sources) || sources[i] != expected {
			t.Errorf("Expected source %s at index %d, got %s", expected, i, sources[i])
		}
	}
}

