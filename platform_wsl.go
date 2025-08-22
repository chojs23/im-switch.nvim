//go:build linux

package main

import (
	"os"
	"os/exec"
	"strings"
)

// WSL-specific implementation that bridges to Windows host
// This file provides WSL detection and Windows host integration

// isWSL detects if we're running in Windows Subsystem for Linux
func isWSL() bool {
	// Check for WSL environment variables
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}

	// Check for WSL in kernel version
	if data, err := os.ReadFile("/proc/version"); err == nil {
		version := strings.ToLower(string(data))
		if strings.Contains(version, "microsoft") && strings.Contains(version, "wsl") {
			return true
		}
	}

	// Check for WSL interop directory
	if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); err == nil {
		return true
	}

	return false
}

// executeWindowsCommand runs a command on the Windows host via PowerShell
func executeWindowsCommand(command string) (string, error) {
	// Use PowerShell to execute Windows commands from WSL
	cmd := exec.Command("powershell.exe", "-Command", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// WSL-specific input source functions

func getCurrentInputSourceWSL() string {
	// Use the more reliable culture info approach
	result, err := executeWindowsCommand("[System.Globalization.CultureInfo]::CurrentCulture.Name")
	if err != nil {
		return ""
	}

	return result
}

func getAllInputSourcesWSL() []string {
	// Try to get installed languages from Windows
	psCommand := `Get-WinUserLanguageList | ForEach-Object { $_.LanguageTag }`
	
	result, err := executeWindowsCommand(psCommand)
	if err != nil || result == "" {
		// Fallback: return common language codes including current
		current := getCurrentInputSourceWSL()
		fallback := []string{
			"en-US", "en-GB", "de-DE", "fr-FR", "es-ES", "it-IT",
			"ja-JP", "ko-KR", "zh-CN", "zh-TW", "ru-RU",
		}
		
		// Add current language if not in fallback list
		if current != "" {
			found := false
			for _, lang := range fallback {
				if lang == current {
					found = true
					break
				}
			}
			if !found {
				fallback = append([]string{current}, fallback...)
			}
		}
		
		return fallback
	}

	sources := strings.Split(result, "\n")
	var cleanSources []string
	for _, source := range sources {
		source = strings.TrimSpace(source)
		if source != "" {
			cleanSources = append(cleanSources, source)
		}
	}

	if len(cleanSources) == 0 {
		return []string{"en-US", "ko-KR"}
	}

	return cleanSources
}

func setInputSourceWSL(sourceID string) bool {
	// First try to set the current culture (simpler approach)
	psCommand := `
		try {
			$culture = [System.Globalization.CultureInfo]::GetCultureInfo("` + sourceID + `")
			[System.Threading.Thread]::CurrentThread.CurrentCulture = $culture
			[System.Threading.Thread]::CurrentThread.CurrentUICulture = $culture
			Write-Output "Success"
		} catch {
			Write-Output "Error: Culture not supported"
		}
	`

	result, err := executeWindowsCommand(psCommand)
	if err == nil && strings.Contains(result, "Success") {
		return true
	}

	// Fallback: try the language list approach for persistent change
	psCommand = `
		$targetLanguage = "` + sourceID + `"
		try {
			$languageList = Get-WinUserLanguageList
			$targetLang = $languageList | Where-Object { $_.LanguageTag -eq $targetLanguage }
			
			if ($targetLang) {
				$newList = @($targetLang) + ($languageList | Where-Object { $_.LanguageTag -ne $targetLanguage })
				Set-WinUserLanguageList -LanguageList $newList -Force
				Write-Output "Success"
			} else {
				$newLang = New-WinUserLanguage -Language $targetLanguage
				$newList = @($newLang) + $languageList
				Set-WinUserLanguageList -LanguageList $newList -Force  
				Write-Output "Success"
			}
		} catch {
			Write-Output "Error"
		}
	`

	result, err = executeWindowsCommand(psCommand)
	return err == nil && strings.Contains(result, "Success")
}

// Override the Linux platform functions when in WSL
func getCurrentInputSourceWSLProxy() string {
	if isWSL() {
		return getCurrentInputSourceWSL()
	}
	// Fall back to original Linux implementation
	return getCurrentInputSourceLinux()
}

func getAllInputSourcesWSLProxy() []string {
	if isWSL() {
		return getAllInputSourcesWSL()
	}
	// Fall back to original Linux implementation
	return getAllInputSourcesLinux()
}

func setInputSourceWSLProxy(sourceID string) bool {
	if isWSL() {
		return setInputSourceWSL(sourceID)
	}
	// Fall back to original Linux implementation
	return setInputSourceLinux(sourceID)
}

// Store original Linux functions
func getCurrentInputSourceLinux() string {
	method := detectInputMethod()

	switch method {
	case "ibus":
		return getCurrentInputSourceIBus()
	case "fcitx5":
		return getCurrentInputSourceFcitx5()
	case "fcitx":
		return getCurrentInputSourceFcitx()
	case "xkb":
		return getCurrentInputSourceXKB()
	default:
		return getCurrentInputSourceXKB()
	}
}

func getAllInputSourcesLinux() []string {
	method := detectInputMethod()

	switch method {
	case "ibus":
		return getAllInputSourcesIBus()
	case "fcitx5":
		return getAllInputSourcesFcitx5()
	case "fcitx":
		return getAllInputSourcesFcitx()
	case "xkb":
		return getAllInputSourcesXKB()
	default:
		return getAllInputSourcesXKB()
	}
}

func setInputSourceLinux(sourceID string) bool {
	method := detectInputMethod()

	switch method {
	case "ibus":
		return setInputSourceIBus(sourceID)
	case "fcitx5":
		return setInputSourceFcitx5(sourceID)
	case "fcitx":
		return setInputSourceFcitx(sourceID)
	case "xkb":
		return setInputSourceXKB(sourceID)
	default:
		return setInputSourceXKB(sourceID)
	}
}