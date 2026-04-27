//go:build darwin

package main

import (
	"testing"
)

func withDarwinCapsLockMocks(t *testing.T, isOn func() bool, setState func(bool) bool, toggle func() bool) {
	t.Helper()

	originalIsOn := darwinIsCapsLockOn
	originalSetState := darwinSetCapsLockState
	originalToggle := darwinToggleCapsLockWithEvent

	darwinIsCapsLockOn = isOn
	darwinSetCapsLockState = setState
	darwinToggleCapsLockWithEvent = toggle

	t.Cleanup(func() {
		darwinIsCapsLockOn = originalIsOn
		darwinSetCapsLockState = originalSetState
		darwinToggleCapsLockWithEvent = originalToggle
	})
}

func TestDarwinTurnOffCapsLockUsesIOKitWhenStateChanges(t *testing.T) {
	checks := 0
	setCalls := 0
	toggleCalls := 0

	withDarwinCapsLockMocks(t,
		func() bool {
			checks++
			return checks == 1
		},
		func(state bool) bool {
			setCalls++
			if state {
				t.Fatal("turnOffCapsLock should request Caps Lock off")
			}
			return true
		},
		func() bool {
			toggleCalls++
			return true
		},
	)

	if !turnOffCapsLock() {
		t.Fatal("turnOffCapsLock should succeed when IOKit changes state")
	}
	if setCalls != 1 {
		t.Fatalf("darwinSetCapsLockState calls = %d, want 1", setCalls)
	}
	if toggleCalls != 0 {
		t.Fatalf("darwinToggleCapsLockWithEvent calls = %d, want 0", toggleCalls)
	}
}

func TestDarwinTurnOffCapsLockFallsBackWhenIOKitDoesNotChangeState(t *testing.T) {
	setCalls := 0
	toggleCalls := 0

	withDarwinCapsLockMocks(t,
		func() bool { return true },
		func(state bool) bool {
			setCalls++
			if state {
				t.Fatal("turnOffCapsLock should request Caps Lock off")
			}
			return true
		},
		func() bool {
			toggleCalls++
			return true
		},
	)

	if !turnOffCapsLock() {
		t.Fatal("turnOffCapsLock should return fallback success")
	}
	if setCalls != 1 {
		t.Fatalf("darwinSetCapsLockState calls = %d, want 1", setCalls)
	}
	if toggleCalls != 1 {
		t.Fatalf("darwinToggleCapsLockWithEvent calls = %d, want 1", toggleCalls)
	}
}

func TestDarwinTurnOffCapsLockSkipsWorkWhenAlreadyOff(t *testing.T) {
	setCalls := 0
	toggleCalls := 0

	withDarwinCapsLockMocks(t,
		func() bool { return false },
		func(state bool) bool {
			setCalls++
			return true
		},
		func() bool {
			toggleCalls++
			return true
		},
	)

	if !turnOffCapsLock() {
		t.Fatal("turnOffCapsLock should succeed when Caps Lock is already off")
	}
	if setCalls != 0 {
		t.Fatalf("darwinSetCapsLockState calls = %d, want 0", setCalls)
	}
	if toggleCalls != 0 {
		t.Fatalf("darwinToggleCapsLockWithEvent calls = %d, want 0", toggleCalls)
	}
}

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
