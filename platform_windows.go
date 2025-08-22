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
	ret, _, _ := procGetKeyboardLayout.Call(0)
	if ret == 0 {
		return ""
	}

	hkl := HKL(ret)
	return getLayoutName(hkl)
}

func getLayoutName(hkl HKL) string {
	layoutId := fmt.Sprintf("%08X", uint32(hkl))
	if name, exists := commonLayouts[layoutId]; exists {
		return name
	}

	return layoutId
}

func getAllInputSources() []string {
	count, _, _ := procGetKeyboardLayoutList.Call(0, 0)
	if count == 0 {
		return nil
	}

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
	var targetHKL HKL
	found := false

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

	for _, hkl := range hklList[:ret] {
		name := getLayoutName(hkl)
		if name == sourceID {
			targetHKL = hkl
			found = true
			break
		}
	}

	if !found {
		sourcePtr, _ := windows.UTF16PtrFromString(sourceID)
		ret, _, _ := procLoadKeyboardLayoutW.Call(
			uintptr(unsafe.Pointer(sourcePtr)),
			0,
		)
		if ret != 0 {
			targetHKL = HKL(ret)
			found = true
		}
	}

	if !found {
		return false
	}

	ret, _, _ = procActivateKeyboardLayout.Call(
		uintptr(targetHKL),
		0,
	)

	return ret != 0
}
