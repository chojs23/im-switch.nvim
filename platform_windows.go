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
	user32 = syscall.NewLazyDLL("user32.dll")
	imm32  = syscall.NewLazyDLL("imm32.dll")

	getForegroundWindow      = user32.NewProc("GetForegroundWindow")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	getKeyboardLayout        = user32.NewProc("GetKeyboardLayout")
	getKeyboardLayoutList    = user32.NewProc("GetKeyboardLayoutList")
	immGetDefaultIMEWnd      = imm32.NewProc("ImmGetDefaultIMEWnd")
	sendMessageA             = user32.NewProc("SendMessageA")
)

const (
	WM_IME_CONTROL            = 0x0283
	WM_INPUTLANGCHANGEREQUEST = 0x0050
	IMC_GETCONVERSIONMODE     = 0x0001
	IMC_SETCONVERSIONMODE     = 0x0002
	IMC_GETOPENSTATUS         = 0x0005
	IMC_SETOPENSTATUS         = 0x0006
	IME_CMODE_NATIVE          = 0x0001
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
	currentInputSource := getCurrentInputSource()
	if !isRequestSupportedForLayout(currentInputSource, sourceID) {
		return false
	}

	if !switchLayoutIfAvailable(sourceID) {
		return false
	}

	// Windows request maps to layout and IME mode control.
	return setIMEStatus(sourceID)
}

func switchLayoutIfAvailable(sourceID string) bool {
	targetLayout, hkl, found := findInstalledLayout(sourceID)
	if !found {
		return true
	}

	foregroundWnd, _, _ := getForegroundWindow.Call()
	if foregroundWnd == 0 {
		return false
	}

	sendMessageA.Call(
		foregroundWnd,
		WM_INPUTLANGCHANGEREQUEST,
		0,
		uintptr(hkl),
	)

	for i := 0; i < 5; i++ {
		if strings.EqualFold(getCurrentInputSource(), targetLayout) {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}

	return false
}

func findInstalledLayout(sourceID string) (string, HKL, bool) {
	normalized := strings.TrimSpace(sourceID)
	if normalized == "" {
		return "", 0, false
	}

	count, _, _ := getKeyboardLayoutList.Call(0, 0)
	if count == 0 {
		return "", 0, false
	}

	layouts := make([]HKL, count)
	ret, _, _ := getKeyboardLayoutList.Call(uintptr(count), uintptr(unsafe.Pointer(&layouts[0])))
	if ret == 0 {
		return "", 0, false
	}

	for _, hkl := range layouts[:ret] {
		layout := getLayoutName(hkl)
		if isSourceMatchLayout(normalized, layout) {
			return layout, hkl, true
		}
	}

	return "", 0, false
}

func isSourceMatchLayout(sourceID string, layout string) bool {
	if strings.EqualFold(sourceID, layout) {
		return true
	}

	sourceParts := strings.SplitN(strings.ToLower(strings.TrimSpace(sourceID)), "-", 2)
	layoutParts := strings.SplitN(strings.ToLower(strings.TrimSpace(layout)), "-", 2)
	if len(sourceParts) > 0 && len(layoutParts) > 0 && sourceParts[0] != "" {
		return sourceParts[0] == layoutParts[0]
	}

	return false
}

func isRequestSupportedForLayout(currentLayout string, requestedSource string) bool {
	if !strings.EqualFold(strings.TrimSpace(currentLayout), "ko-KR") {
		return true
	}

	normalizedRequest := strings.TrimSpace(requestedSource)
	return strings.EqualFold(normalizedRequest, "ko-KR") || strings.EqualFold(normalizedRequest, "en-US")
}

func shouldOpenIMEForSourceID(sourceID string) bool {
	normalized := strings.TrimSpace(sourceID)
	isCJK := cjkLanguages[normalized]

	if !isCJK {
		for lang := range cjkLanguages {
			if strings.HasPrefix(normalized, strings.Split(lang, "-")[0]) {
				isCJK = true
				break
			}
		}
	}

	return isCJK
}

func setIMEStatus(sourceID string) bool {
	open := shouldOpenIMEForSourceID(sourceID)
	if !setIMEOpenStatus(open) {
		return false
	}

	if shouldForceNativeMode(sourceID) {
		return setNativeConversionMode(true)
	}

	return true
}

func shouldForceNativeMode(sourceID string) bool {
	normalized := strings.ToLower(strings.TrimSpace(sourceID))
	return normalized == "ko-kr" || normalized == "zh-cn" || normalized == "zh-tw"
}

func setIMEOpenStatus(open bool) bool {
	for i := 0; i < 5; i++ {
		imeWnd, ok := getIMEWindow()
		if !ok {
			return false
		}

		if !setIMEOpenStatusForWindow(imeWnd, open) {
			return false
		}

		if isIMEOpenForWindow(imeWnd) == open {
			return true
		}

		time.Sleep(20 * time.Millisecond)
	}

	return false
}

func getIMEWindow() (uintptr, bool) {
	foregroundWnd, _, _ := getForegroundWindow.Call()
	if foregroundWnd == 0 {
		return 0, false
	}

	imeWnd, _, _ := immGetDefaultIMEWnd.Call(foregroundWnd)
	if imeWnd == 0 {
		return 0, false
	}

	return imeWnd, true
}

func setIMEOpenStatusForWindow(imeWnd uintptr, open bool) bool {
	var target uintptr
	if open {
		target = 1
	}

	ret, _, _ := sendMessageA.Call(
		imeWnd,
		WM_IME_CONTROL,
		IMC_SETOPENSTATUS,
		target,
	)

	// Some IMEs return 0 even when set succeeds, so verify via GETOPENSTATUS.
	_ = ret
	return true
}

func isIMEOpenForWindow(imeWnd uintptr) bool {
	status, _, _ := sendMessageA.Call(
		imeWnd,
		WM_IME_CONTROL,
		IMC_GETOPENSTATUS,
		0,
	)
	return status != 0
}

func setNativeConversionMode(enabled bool) bool {
	imeWnd, ok := getIMEWindow()
	if !ok {
		return false
	}

	var target uintptr
	if enabled {
		target = IME_CMODE_NATIVE
	}

	sendMessageA.Call(
		imeWnd,
		WM_IME_CONTROL,
		IMC_SETCONVERSIONMODE,
		target,
	)

	for i := 0; i < 5; i++ {
		current, _, _ := sendMessageA.Call(
			imeWnd,
			WM_IME_CONTROL,
			IMC_GETCONVERSIONMODE,
			0,
		)
		nativeOn := (current & IME_CMODE_NATIVE) != 0
		if nativeOn == enabled {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}

	return false
}

func getDetailedIMEStatus() string {
	imeWnd, ok := getIMEWindow()
	if !ok {
		return "Unknown"
	}

	if isIMEOpenForWindow(imeWnd) {
		return "open"
	}
	return "closed"
}

func getBackendStatus() string {
	return "windows-ime"
}
