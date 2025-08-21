//go:build linux

package main

import (
	"os"
	"os/exec"
	"strings"
)

// Linux input method switching using multiple backends
// Supports: ibus, fcitx, fcitx5, rime

// detectInputMethod detects which input method framework is running
func detectInputMethod() string {
	if im := os.Getenv("GTK_IM_MODULE"); im != "" {
		if strings.Contains(im, "ibus") {
			return "ibus"
		}
		if strings.Contains(im, "fcitx") {
			return "fcitx"
		}
	}

	if im := os.Getenv("QT_IM_MODULE"); im != "" {
		if strings.Contains(im, "ibus") {
			return "ibus"
		}
		if strings.Contains(im, "fcitx") {
			return "fcitx"
		}
	}

	if isProcessRunning("ibus-daemon") {
		return "ibus"
	}
	if isProcessRunning("fcitx5") {
		return "fcitx5"
	}
	if isProcessRunning("fcitx") {
		return "fcitx"
	}

	// Default fallback
	return "xkb"
}

func isProcessRunning(process string) bool {
	cmd := exec.Command("pgrep", process)
	return cmd.Run() == nil
}

func getCurrentInputSource() string {
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

func getCurrentInputSourceIBus() string {
	cmd := exec.Command("ibus", "engine")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getCurrentInputSourceFcitx5() string {
	cmd := exec.Command("fcitx5-remote", "-n")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getCurrentInputSourceFcitx() string {
	cmd := exec.Command("fcitx-remote", "-n")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getCurrentInputSourceXKB() string {
	cmd := exec.Command("setxkbmap", "-query")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "layout:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}
	return ""
}

func getAllInputSources() []string {
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

func getAllInputSourcesIBus() []string {
	cmd := exec.Command("ibus", "list-engine")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var sources []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "language:") {
			sources = append(sources, line)
		}
	}
	return sources
}

func getAllInputSourcesFcitx5() []string {
	cmd := exec.Command("fcitx5-remote", "-l")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var sources []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			sources = append(sources, line)
		}
	}
	return sources
}

func getAllInputSourcesFcitx() []string {
	cmd := exec.Command("fcitx-remote", "-l")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var sources []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			sources = append(sources, line)
		}
	}
	return sources
}

func getAllInputSourcesXKB() []string {
	// Common keyboard layouts
	return []string{
		"us",
		"gb",
		"de",
		"fr",
		"es",
		"it",
		"ru",
		"cn",
		"jp",
		"kr",
	}
}

func setInputSource(sourceID string) bool {
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

func setInputSourceIBus(sourceID string) bool {
	cmd := exec.Command("ibus", "engine", sourceID)
	return cmd.Run() == nil
}

func setInputSourceFcitx5(sourceID string) bool {
	cmd := exec.Command("fcitx5-remote", "-s", sourceID)
	return cmd.Run() == nil
}

func setInputSourceFcitx(sourceID string) bool {
	cmd := exec.Command("fcitx-remote", "-s", sourceID)
	return cmd.Run() == nil
}

func setInputSourceXKB(sourceID string) bool {
	cmd := exec.Command("setxkbmap", sourceID)
	return cmd.Run() == nil
}

