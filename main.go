package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  im-switch                    # Show current input source")
	fmt.Println("  im-switch --status           # Show OS, binary path, backend, and current input source")
	fmt.Println("  im-switch -l                 # List all input sources")
	fmt.Println("  im-switch --capslock-off     # Turn Caps Lock off if it is on")
	fmt.Println("  im-switch [input-source-id]  # Switch to input source")
	fmt.Println("")
	fmt.Println("Examples:")
	if runtime.GOOS == "darwin" {
		fmt.Println("  # macOS")
		fmt.Println("  im-switch com.apple.keylayout.ABC")
		fmt.Println("  im-switch com.apple.inputmethod.Korean.2SetKorean")
	} else if runtime.GOOS == "windows" {
		fmt.Println("  # Windows")
		fmt.Println("  im-switch en-US                 # English (United States)")
		fmt.Println("  im-switch 00000409              # English (United States) by layout ID")
		fmt.Println("  im-switch de-DE                 # German (Germany)")
		fmt.Println("  im-switch ja-JP                 # Japanese")
		fmt.Println("  # Note: Windows can use IDs as IME mode requests (en-US=IME off, zh-CN/ja-JP/ko-KR=IME on)")
	} else {
		fmt.Println("  # Linux")
		fmt.Println("  im-switch us                    # XKB layout")
		fmt.Println("  im-switch xkb:us::eng           # IBus")
		fmt.Println("  im-switch keyboard-us           # Fcitx")
	}
	fmt.Println("")
	fmt.Printf("Platform: %s\n", runtime.GOOS)
}

func printStatus() int {
	binaryPath, err := os.Executable()
	if err != nil {
		binaryPath = "unavailable"
	}

	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("WSL: %t\n", isWSL())
	fmt.Printf("Binary: %s\n", binaryPath)
	fmt.Printf("Backend: %s\n", getBackendStatus())

	current := getCurrentInputSource()
	if current == "" {
		fmt.Println("Current input source: unavailable")
		if runtime.GOOS == "linux" {
			fmt.Println("Hint: install ibus, fcitx, fcitx5, or ensure setxkbmap is available")
		}
		return 1
	}

	fmt.Printf("Current input source: %s\n", current)
	return 0
}

func isWSL() bool {
	return detectWSL(runtime.GOOS, os.LookupEnv, os.ReadFile)
}

func detectWSL(goos string, lookupEnv func(string) (string, bool), readFile func(string) ([]byte, error)) bool {
	if _, ok := lookupEnv("WSL_DISTRO_NAME"); ok {
		return true
	}

	if _, ok := lookupEnv("WSL_INTEROP"); ok {
		return true
	}

	if goos == "windows" && isWSLFromWindowsProcess() {
		return true
	}

	if goos != "linux" {
		return false
	}

	content, err := readFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return false
	}

	osrelease := strings.ToLower(string(content))
	return strings.Contains(osrelease, "microsoft")
}

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		current := getCurrentInputSource()
		if current == "" {
			if runtime.GOOS == "linux" {
				fmt.Fprintf(os.Stderr, "Error: No input method framework detected\n")
				fmt.Fprintf(os.Stderr, "Please install one of: ibus, fcitx, fcitx5, or ensure setxkbmap is available\n")
			} else {
				fmt.Fprintf(os.Stderr, "Error: Could not get current input source\n")
			}
			os.Exit(1)
		}
		fmt.Println(current)

	case 1:
		arg := args[0]
		if arg == "--status" {
			if code := printStatus(); code != 0 {
				os.Exit(code)
			}
		} else if arg == "-l" || arg == "--list" {
			sources := getAllInputSources()
			if sources == nil {
				if runtime.GOOS == "linux" {
					fmt.Fprintf(os.Stderr, "Error: No input method framework detected\n")
					fmt.Fprintf(os.Stderr, "Please install one of: ibus, fcitx, fcitx5, or ensure setxkbmap is available\n")
				} else {
					fmt.Fprintf(os.Stderr, "Error: Could not get input sources\n")
				}
				os.Exit(1)
			}
			for _, source := range sources {
				fmt.Println(source)
			}
		} else if arg == "-h" || arg == "--help" {
			printUsage()
		} else if arg == "--capslock-off" {
			if !turnOffCapsLock() {
				fmt.Fprintf(os.Stderr, "Error: Could not turn Caps Lock off\n")
				os.Exit(1)
			}
		} else {
			if !setInputSource(arg) {
				fmt.Fprintf(os.Stderr, "Error: Could not set input source to '%s'\n", arg)
				if runtime.GOOS == "windows" {
					fmt.Fprintf(os.Stderr, "On Windows, this controls IME mode as well as input IDs\n")
					fmt.Fprintf(os.Stderr, "Only the active layout language and 'en-US' are supported\n")
					fmt.Fprintf(os.Stderr, "Official layout codes (for example '00000804' or '0x804') request keyboard layout switching\n")
					fmt.Fprintf(os.Stderr, "Try 'en-US' for English mode (IME off) and 'zh-CN'/'ja-JP'/'ko-KR' for IME on\n")
				} else {
					fmt.Fprintf(os.Stderr, "Use 'im-switch -l' to see available input sources\n")
				}
				os.Exit(1)
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: Too many arguments\n\n")
		printUsage()
		os.Exit(1)
	}
}
