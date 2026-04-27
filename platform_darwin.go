//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Carbon -framework ApplicationServices -framework IOKit

#include <ApplicationServices/ApplicationServices.h>
#include <Carbon/Carbon.h>
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOKitLib.h>
#include <IOKit/hidsystem/IOHIDLib.h>
#include <IOKit/hidsystem/IOHIDParameter.h>
#include <unistd.h>

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

bool isCapsLockOn() {
    CGEventFlags flags = CGEventSourceFlagsState(kCGEventSourceStateHIDSystemState);
    return (flags & kCGEventFlagMaskAlphaShift) != 0;
}

bool setCapsLockState(bool state) {
    io_service_t service = IOServiceGetMatchingService(kIOMainPortDefault, IOServiceMatching(kIOHIDSystemClass));
    if (service == MACH_PORT_NULL) {
        return false;
    }

    io_connect_t connect = MACH_PORT_NULL;
    kern_return_t status = IOServiceOpen(service, mach_task_self(), kIOHIDParamConnectType, &connect);
    IOObjectRelease(service);
    if (status != KERN_SUCCESS) {
        return false;
    }

    status = IOHIDSetModifierLockState(connect, kIOHIDCapsLockState, state);

    bool currentState = state;
    bool verified = false;
    if (status == KERN_SUCCESS && IOHIDGetModifierLockState(connect, kIOHIDCapsLockState, &currentState) == KERN_SUCCESS) {
        verified = currentState == state;
    }

    IOServiceClose(connect);
    return status == KERN_SUCCESS && verified;
}

bool toggleCapsLockWithKeyEvent() {
    CGEventRef keyDown = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)kVK_CapsLock, true);
    CGEventRef keyUp = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)kVK_CapsLock, false);
    if (keyDown == NULL || keyUp == NULL) {
        if (keyDown) CFRelease(keyDown);
        if (keyUp) CFRelease(keyUp);
        return false;
    }

    CGEventPost(kCGHIDEventTap, keyDown);
    CGEventPost(kCGHIDEventTap, keyUp);
    CFRelease(keyDown);
    CFRelease(keyUp);
    usleep(20000);

    return !isCapsLockOn();
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
import (
	"time"
	"unsafe"
)

var (
	darwinIsCapsLockOn            = func() bool { return bool(C.isCapsLockOn()) }
	darwinSetCapsLockState        = func(state bool) bool { return bool(C.setCapsLockState(C.bool(state))) }
	darwinToggleCapsLockWithEvent = func() bool { return bool(C.toggleCapsLockWithKeyEvent()) }
)

func waitForDarwinCapsLockState() {
	time.Sleep(20 * time.Millisecond)
}

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

func turnOffCapsLock() bool {
	if !darwinIsCapsLockOn() {
		return true
	}

	darwinSetCapsLockState(false)
	waitForDarwinCapsLockState()
	if !darwinIsCapsLockOn() {
		return true
	}

	return darwinToggleCapsLockWithEvent()
}

func getBackendStatus() string {
	return "macos-tis"
}
