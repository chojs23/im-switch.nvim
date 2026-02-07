//go:build windows

package main

import "testing"

func TestShouldOpenIMEForSourceID(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		wantOpen bool
	}{
		{name: "english mode", sourceID: "en-US", wantOpen: false},
		{name: "chinese mode", sourceID: "zh-CN", wantOpen: true},
		{name: "japanese prefix", sourceID: "ja", wantOpen: true},
		{name: "korean mode", sourceID: "ko-KR", wantOpen: true},
		{name: "trim whitespace", sourceID: "  zh-CN  ", wantOpen: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldOpenIMEForSourceID(tt.sourceID)
			if got != tt.wantOpen {
				t.Fatalf("shouldOpenIMEForSourceID(%q) = %v, want %v", tt.sourceID, got, tt.wantOpen)
			}
		})
	}
}

func TestIsRequestSupportedForLayout(t *testing.T) {
	tests := []struct {
		name            string
		currentLayout   string
		requestedSource string
		want            bool
	}{
		{name: "korean layout allows ko-KR", currentLayout: "ko-KR", requestedSource: "ko-KR", want: true},
		{name: "korean layout allows en-US", currentLayout: "ko-KR", requestedSource: "en-US", want: true},
		{name: "korean layout rejects ja-JP", currentLayout: "ko-KR", requestedSource: "ja-JP", want: false},
		{name: "korean layout rejects random id", currentLayout: "ko-KR", requestedSource: "abc", want: false},
		{name: "japanese layout rejects korean request", currentLayout: "ja-JP", requestedSource: "ko-KR", want: false},
		{name: "chinese layout rejects japanese request", currentLayout: "zh-CN", requestedSource: "ja-JP", want: false},
		{name: "chinese layout allows chinese request", currentLayout: "zh-CN", requestedSource: "zh-CN", want: true},
		{name: "english layout allows english variant", currentLayout: "en-US", requestedSource: "en-GB", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRequestSupportedForLayout(tt.currentLayout, tt.requestedSource)
			if got != tt.want {
				t.Fatalf("isRequestSupportedForLayout(%q, %q) = %v, want %v", tt.currentLayout, tt.requestedSource, got, tt.want)
			}
		})
	}
}

func TestPrimaryLanguage(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		want     string
	}{
		{name: "region language tag", sourceID: "ko-KR", want: "ko"},
		{name: "simple language", sourceID: "ja", want: "ja"},
		{name: "trim and lowercase", sourceID: " ZH-cn ", want: "zh"},
		{name: "empty string", sourceID: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := primaryLanguage(tt.sourceID)
			if got != tt.want {
				t.Fatalf("primaryLanguage(%q) = %q, want %q", tt.sourceID, got, tt.want)
			}
		})
	}
}

func TestIsSourceMatchLayout(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		layout   string
		want     bool
	}{
		{name: "exact match", sourceID: "ko-KR", layout: "ko-KR", want: true},
		{name: "short language match", sourceID: "ko", layout: "ko-KR", want: true},
		{name: "case-insensitive match", sourceID: "EN-us", layout: "en-US", want: true},
		{name: "different language", sourceID: "ja-JP", layout: "ko-KR", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSourceMatchLayout(tt.sourceID, tt.layout)
			if got != tt.want {
				t.Fatalf("isSourceMatchLayout(%q, %q) = %v, want %v", tt.sourceID, tt.layout, got, tt.want)
			}
		})
	}
}

func TestShouldForceNativeMode(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		want     bool
	}{
		{name: "korean uses native mode", sourceID: "ko-KR", want: true},
		{name: "korean case-insensitive", sourceID: "KO-kr", want: true},
		{name: "chinese simplified uses native mode", sourceID: "zh-CN", want: true},
		{name: "chinese traditional uses native mode", sourceID: "zh-TW", want: true},
		{name: "english no native mode", sourceID: "en-US", want: false},
		{name: "japanese no forced native mode", sourceID: "ja-JP", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldForceNativeMode(tt.sourceID)
			if got != tt.want {
				t.Fatalf("shouldForceNativeMode(%q) = %v, want %v", tt.sourceID, got, tt.want)
			}
		})
	}
}
