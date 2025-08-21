package main

import (
	"fmt"
	"os"
	"runtime"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  im-switch                    # Show current input source")
	fmt.Println("  im-switch -l                 # List all input sources")
	fmt.Println("  im-switch [input-source-id]  # Switch to input source")
	fmt.Println("")
	fmt.Println("Examples:")
	if runtime.GOOS == "darwin" {
		fmt.Println("  # macOS")
		fmt.Println("  im-switch com.apple.keylayout.ABC")
		fmt.Println("  im-switch com.apple.inputmethod.Korean.2SetKorean")
	} else {
		fmt.Println("  # Linux")
		fmt.Println("  im-switch us                    # XKB layout")
		fmt.Println("  im-switch xkb:us::eng           # IBus")
		fmt.Println("  im-switch keyboard-us           # Fcitx")
	}
	fmt.Println("")
	fmt.Printf("Platform: %s\n", runtime.GOOS)
}

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		current := getCurrentInputSource()
		if current == "" {
			fmt.Fprintf(os.Stderr, "Error: Could not get current input source\n")
			os.Exit(1)
		}
		fmt.Println(current)

	case 1:
		arg := args[0]
		if arg == "-l" || arg == "--list" {
			sources := getAllInputSources()
			if sources == nil {
				fmt.Fprintf(os.Stderr, "Error: Could not get input sources\n")
				os.Exit(1)
			}
			for _, source := range sources {
				fmt.Println(source)
			}
		} else if arg == "-h" || arg == "--help" {
			printUsage()
		} else {
			if !setInputSource(arg) {
				fmt.Fprintf(os.Stderr, "Error: Could not set input source to '%s'\n", arg)
				fmt.Fprintf(os.Stderr, "Use 'im-switch -l' to see available input sources\n")
				os.Exit(1)
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: Too many arguments\n\n")
		printUsage()
		os.Exit(1)
	}
}
