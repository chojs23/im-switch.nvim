package main

import (
	"fmt"
	"os"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Carbon

#include <Carbon/Carbon.h>
#include <CoreFoundation/CoreFoundation.h>

// Get current input source
CFStringRef getCurrentInputSource() {
    TISInputSourceRef currentSource = TISCopyCurrentKeyboardInputSource();
    if (currentSource == NULL) {
        return NULL;
    }

    CFStringRef sourceID = (CFStringRef)TISGetInputSourceProperty(currentSource, kTISPropertyInputSourceID);
    CFRetain(sourceID);
    CFRelease(currentSource);

    return sourceID;
}

// Get all input sources
CFArrayRef getAllInputSources() {
    CFStringRef keys[] = {kTISPropertyInputSourceCategory};
    CFStringRef values[] = {kTISCategoryKeyboardInputSource};
    CFDictionaryRef filter = CFDictionaryCreate(
        kCFAllocatorDefault,
        (const void**)keys,
        (const void**)values,
        1,
        &kCFTypeDictionaryKeyCallBacks,
        &kCFTypeDictionaryValueCallBacks
    );

    CFArrayRef inputSources = TISCreateInputSourceList(filter, false);
    CFRelease(filter);

    return inputSources;
}

// Set input source by ID
bool setInputSource(CFStringRef sourceID) {
    CFStringRef keys[] = {kTISPropertyInputSourceID};
    CFStringRef values[] = {sourceID};
    CFDictionaryRef filter = CFDictionaryCreate(
        kCFAllocatorDefault,
        (const void**)keys,
        (const void**)values,
        1,
        &kCFTypeDictionaryKeyCallBacks,
        &kCFTypeDictionaryValueCallBacks
    );

    CFArrayRef inputSources = TISCreateInputSourceList(filter, false);
    CFRelease(filter);

    if (inputSources == NULL || CFArrayGetCount(inputSources) == 0) {
        if (inputSources) CFRelease(inputSources);
        return false;
    }

    TISInputSourceRef source = (TISInputSourceRef)CFArrayGetValueAtIndex(inputSources, 0);
    OSStatus status = TISSelectInputSource(source);
    CFRelease(inputSources);

    return status == noErr;
}

// Convert CFString to C string
char* cfStringToCString(CFStringRef cfStr) {
    if (cfStr == NULL) return NULL;

    CFIndex length = CFStringGetLength(cfStr);
    CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
    char* buffer = malloc(maxSize);
    if (buffer == NULL) return NULL;

    if (CFStringGetCString(cfStr, buffer, maxSize, kCFStringEncodingUTF8)) {
        return buffer;
    }

    free(buffer);
    return NULL;
}
*/
import "C"

func getCurrentInputSource() string {
	cfStr := C.getCurrentInputSource()
	if cfStr == C.CFStringRef(unsafe.Pointer(nil)) {
		return ""
	}
	defer C.CFRelease(C.CFTypeRef(cfStr))

	cStr := C.cfStringToCString(cfStr)
	if cStr == (*C.char)(unsafe.Pointer(nil)) {
		return ""
	}
	defer C.free(unsafe.Pointer(cStr))

	return C.GoString(cStr)
}

func getAllInputSources() []string {
	sources := C.getAllInputSources()
	if sources == C.CFArrayRef(unsafe.Pointer(nil)) {
		return nil
	}
	defer C.CFRelease(C.CFTypeRef(sources))

	count := C.CFArrayGetCount(sources)
	result := make([]string, 0, count)

	for i := C.CFIndex(0); i < count; i++ {
		source := C.TISInputSourceRef(C.CFArrayGetValueAtIndex(sources, i))
		if source == C.TISInputSourceRef(unsafe.Pointer(nil)) {
			continue
		}

		sourceID := C.CFStringRef(C.TISGetInputSourceProperty(source, C.kTISPropertyInputSourceID))
		if sourceID == C.CFStringRef(unsafe.Pointer(nil)) {
			continue
		}

		cStr := C.cfStringToCString(sourceID)
		if cStr != (*C.char)(unsafe.Pointer(nil)) {
			result = append(result, C.GoString(cStr))
			C.free(unsafe.Pointer(cStr))
		}
	}

	return result
}

func setInputSource(sourceID string) bool {
	cStr := C.CString(sourceID)
	defer C.free(unsafe.Pointer(cStr))

	cfStr := C.CFStringCreateWithCString(C.kCFAllocatorDefault, cStr, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfStr))

	return bool(C.setInputSource(cfStr))
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  im-switch                    # Show current input source")
	fmt.Println("  im-switch -l                 # List all input sources")
	fmt.Println("  im-switch [input-source-id]  # Switch to input source")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  im-switch com.apple.keylayout.ABC")
	fmt.Println("  im-switch com.apple.inputmethod.Korean.2SetKorean")
}

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		// Show current input source
		current := getCurrentInputSource()
		if current == "" {
			fmt.Fprintf(os.Stderr, "Error: Could not get current input source\n")
			os.Exit(1)
		}
		fmt.Println(current)

	case 1:
		arg := args[0]
		if arg == "-l" || arg == "--list" {
			// List all input sources
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
			// Set input source
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
