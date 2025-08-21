//go:build windows

package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                     = windows.NewLazySystemDLL("user32.dll")
	procGetKeyboardLayout      = user32.NewProc("GetKeyboardLayout")
	procActivateKeyboardLayout = user32.NewProc("ActivateKeyboardLayout")
	procGetKeyboardLayoutNameW = user32.NewProc("GetKeyboardLayoutNameW")
	procGetKeyboardLayoutList  = user32.NewProc("GetKeyboardLayoutList")
	procLoadKeyboardLayoutW    = user32.NewProc("LoadKeyboardLayoutW")
)

// HKL represents a handle to keyboard layout
type HKL uintptr

// Common keyboard layout identifiers for Windows
var commonLayouts = map[string]string{
	"00000409": "en-US", // US English
	"00000809": "en-GB", // UK English
	"00000407": "de-DE", // German
	"0000040c": "fr-FR", // French
	"00000410": "it-IT", // Italian
	"0000040a": "es-ES", // Spanish
	"00000411": "ja-JP", // Japanese
	"00000412": "ko-KR", // Korean
	"00000804": "zh-CN", // Chinese (Simplified)
	"00000404": "zh-TW", // Chinese (Traditional)
	"00000419": "ru-RU", // Russian
}

func getCurrentInputSource() string {
	ret, _, _ := procGetKeyboardLayout.Call(0) // 0 means current thread
	if ret == 0 {
		return ""
	}

	hkl := HKL(ret)
	return getLayoutName(hkl)
}

func getLayoutName(hkl HKL) string {
	// Try to get the layout name using GetKeyboardLayoutNameW
	buf := make([]uint16, 9) // KL_NAMELENGTH is 9 characters
	ret, _, _ := procGetKeyboardLayoutNameW.Call(uintptr(unsafe.Pointer(&buf[0])))
	if ret != 0 {
		return windows.UTF16ToString(buf)
	}

	// Fallback: convert HKL to hex string
	layoutId := fmt.Sprintf("%08X", uint32(hkl))
	if name, exists := commonLayouts[layoutId]; exists {
		return name
	}

	return layoutId
}

func getAllInputSources() []string {
	// Get the number of keyboard layouts
	count, _, _ := procGetKeyboardLayoutList.Call(0, 0)
	if count == 0 {
		return nil
	}

	// Allocate buffer for layout handles
	layouts := make([]HKL, count)
	ret, _, _ := procGetKeyboardLayoutList.Call(
		uintptr(count),
		uintptr(unsafe.Pointer(&layouts[0])),
	)

	if ret == 0 {
		return nil
	}

	var sources []string
	for _, hkl := range layouts[:ret] {
		name := getLayoutName(hkl)
		if name != "" {
			sources = append(sources, name)
		}
	}

	return sources
}

func setInputSource(sourceID string) bool {
	// First, try to find the layout by name
	var targetHKL HKL
	found := false

	// Get all available layouts to find matching HKL
	count, _, _ := procGetKeyboardLayoutList.Call(0, 0)
	if count == 0 {
		return false
	}

	hklList := make([]HKL, count)
	ret, _, _ := procGetKeyboardLayoutList.Call(
		uintptr(count),
		uintptr(unsafe.Pointer(&hklList[0])),
	)

	if ret == 0 {
		return false
	}

	// Find the HKL that matches our source ID
	for _, hkl := range hklList[:ret] {
		name := getLayoutName(hkl)
		if name == sourceID {
			targetHKL = hkl
			found = true
			break
		}
	}

	// If not found by name, try to parse as hex layout ID
	if !found {
		// Try to load the keyboard layout if it's not already loaded
		sourcePtr, _ := windows.UTF16PtrFromString(sourceID)
		ret, _, _ := procLoadKeyboardLayoutW.Call(
			uintptr(unsafe.Pointer(sourcePtr)),
			0, // KLF_ACTIVATE flag
		)
		if ret != 0 {
			targetHKL = HKL(ret)
			found = true
		}
	}

	if !found {
		return false
	}

	// Activate the keyboard layout
	ret, _, _ = procActivateKeyboardLayout.Call(
		uintptr(targetHKL),
		0, // KLF_SETFORPROCESS flag
	)

	return ret != 0
}
