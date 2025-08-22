//go:build windows

package main

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32      = syscall.NewLazyDLL("user32.dll")
	imm32       = syscall.NewLazyDLL("imm32.dll")
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	keybd_event = user32.NewProc("keybd_event")

	getForegroundWindow      = user32.NewProc("GetForegroundWindow")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	getKeyboardLayout        = user32.NewProc("GetKeyboardLayout")
	getKeyboardLayoutList    = user32.NewProc("GetKeyboardLayoutList")
	getLocaleInfoW           = kernel32.NewProc("GetLocaleInfoW")
	immGetDefaultIMEWnd      = imm32.NewProc("ImmGetDefaultIMEWnd")
	immGetOpenStatus         = imm32.NewProc("ImmGetOpenStatus")
	immSetOpenStatus         = imm32.NewProc("ImmSetOpenStatus")
	immGetContext            = imm32.NewProc("ImmGetContext")
	immReleaseContext        = imm32.NewProc("ImmReleaseContext")
	sendMessageW             = user32.NewProc("SendMessageW")
)

const (
	LOCALE_SENGLANGUAGE = 0x1001
)

const (
	WM_IME_CONTROL    = 0x0283
	IMC_GETOPENSTATUS = 0x0005
)

const (
	VK_HANGUL = 0x15

	VK_KANJI = 0x19
	VK_KANA  = 0x15

	VK_SHIFT = 0x10

	KEYEVENTF_KEYUP = 0x0002
)

// HKL represents a handle to keyboard layout
type HKL uintptr

// Common keyboard layout identifiers for Windows
var layoutNames = map[uint32]string{
	0x0409: "en-US", // US English
	0x0809: "en-GB", // UK English
	0x0407: "de-DE", // German
	0x040c: "fr-FR", // French
	0x0410: "it-IT", // Italian
	0x040a: "es-ES", // Spanish
	0x0411: "ja-JP", // Japanese
	0x0412: "ko-KR", // Korean
	0x0804: "zh-CN", // Chinese (Simplified)
	0x0404: "zh-TW", // Chinese (Traditional)
	0x0419: "ru-RU", // Russian
}

// CJK language codes that should have IME enabled
var cjkLanguages = map[string]bool{
	"ja-JP": true, // Japanese
	"ko-KR": true, // Korean
	"zh-CN": true, // Chinese (Simplified)
	"zh-TW": true, // Chinese (Traditional)
}

func getCurrentInputSource() string {
	foregroundWnd, _, _ := getForegroundWindow.Call()
	if foregroundWnd == 0 {
		return ""
	}

	var processId uintptr
	threadId, _, _ := getWindowThreadProcessId.Call(foregroundWnd, uintptr(unsafe.Pointer(&processId)))
	if threadId == 0 {
		return ""
	}

	keyboardLayout, _, _ := getKeyboardLayout.Call(threadId)
	if keyboardLayout == 0 {
		return ""
	}

	hkl := HKL(keyboardLayout)
	return getLayoutName(hkl)
}

func getLayoutName(hkl HKL) string {
	// Extract the low 16 bits for the primary language identifier
	langId := uint32(hkl) & 0xFFFF

	if name, exists := layoutNames[langId]; exists {
		return name
	}

	// Return the hex representation as fallback
	return fmt.Sprintf("%04X", langId)
}

func getAllInputSources() []string {
	count, _, _ := getKeyboardLayoutList.Call(0, 0)
	if count == 0 {
		return nil
	}

	layouts := make([]HKL, count)
	ret, _, _ := getKeyboardLayoutList.Call(
		uintptr(count),
		uintptr(unsafe.Pointer(&layouts[0])),
	)

	if ret == 0 {
		return nil
	}

	var sources []string
	seen := make(map[string]bool)

	for _, hkl := range layouts[:ret] {
		name := getLayoutName(hkl)
		if name != "" && !seen[name] {
			sources = append(sources, name)
			seen[name] = true
		}
	}

	return sources
}

func setInputSource(sourceID string) bool {
	// On Windows, we don't change the keyboard layout
	// Instead, we control the IME status based on the source language
	return setIMEStatus(sourceID)
}

func setIMEStatus(sourceID string) bool {
	// Determine if the source is a CJK language
	isCJK := cjkLanguages[sourceID]

	// Also check for partial matches (e.g., "ja" for Japanese)
	if !isCJK {
		for lang := range cjkLanguages {
			if strings.HasPrefix(sourceID, strings.Split(lang, "-")[0]) {
				isCJK = true
				break
			}
		}
	}

	return setIMEOpenStatus(isCJK)
}

func setIMEOpenStatus(open bool) bool {
	isOpen := getDetailedIMEStatus() == "open"

	currentInputSource := getCurrentInputSource()

	if (open && !isOpen) || (!open && isOpen) {
		if currentInputSource == "ja-JP" {
			toggleJapanese()
		}
		if currentInputSource == "ko-KR" {
			toggleKorean()
		}
		if currentInputSource == "zh-CN" || currentInputSource == "zh-TW" {
			toggleChinese()
		}
	}

	return isOpen == open
}

func getDetailedIMEStatus() string {
	foregroundWnd, _, _ := getForegroundWindow.Call()
	if foregroundWnd == 0 {
		return "Unknown"
	}

	var processId uintptr
	threadId, _, _ := getWindowThreadProcessId.Call(foregroundWnd, uintptr(unsafe.Pointer(&processId)))
	if threadId == 0 {
		return "Unknown"
	}

	keyboardLayout, _, _ := getKeyboardLayout.Call(threadId)
	if keyboardLayout == 0 {
		return "Unknown"
	}

	imeWnd, _, _ := immGetDefaultIMEWnd.Call(foregroundWnd)
	imeOpen := false
	if imeWnd != 0 {
		status, _, _ := sendMessageW.Call(
			imeWnd,
			WM_IME_CONTROL,
			IMC_GETOPENSTATUS,
			0,
		)
		imeOpen = (status != 0)
	}

	if imeOpen {
		return "open"
	} else {
		return "closed"
	}
}

func toggleKorean() {
	pressKey(VK_HANGUL)
}

func toggleJapanese() {
	pressKey(VK_KANJI)

	// If fails try using Kana key
	// pressAltKey(0xC0)
}

func toggleChinese() {
	pressKey(VK_SHIFT)
}

func pressKey(vkCode uint32) {
	keybd_event.Call(
		uintptr(vkCode),
		0,
		0,
		0,
	)

	time.Sleep(10 * time.Millisecond)

	keybd_event.Call(
		uintptr(vkCode),
		0,
		uintptr(KEYEVENTF_KEYUP),
		0,
	)
}

func pressAltKey(vkCode uint32) {
	const VK_MENU = 0x12

	keybd_event.Call(uintptr(VK_MENU), 0, 0, 0)
	time.Sleep(10 * time.Millisecond)

	keybd_event.Call(uintptr(vkCode), 0, 0, 0)
	time.Sleep(10 * time.Millisecond)

	keybd_event.Call(uintptr(vkCode), 0, uintptr(KEYEVENTF_KEYUP), 0)
	time.Sleep(10 * time.Millisecond)

	keybd_event.Call(uintptr(VK_MENU), 0, uintptr(KEYEVENTF_KEYUP), 0)
}

func pressCtrlKey(vkCode uint32) {
	const VK_CONTROL = 0x11

	keybd_event.Call(uintptr(VK_CONTROL), 0, 0, 0)
	time.Sleep(10 * time.Millisecond)

	keybd_event.Call(uintptr(vkCode), 0, 0, 0)
	time.Sleep(10 * time.Millisecond)

	keybd_event.Call(uintptr(vkCode), 0, uintptr(KEYEVENTF_KEYUP), 0)
	time.Sleep(10 * time.Millisecond)

	keybd_event.Call(uintptr(VK_CONTROL), 0, uintptr(KEYEVENTF_KEYUP), 0)
}

