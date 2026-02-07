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
		{name: "korean layout allows official klid code", currentLayout: "ko-KR", requestedSource: "00000804", want: true},
		{name: "japanese layout allows official 0x code", currentLayout: "ja-JP", requestedSource: "0x412", want: true},
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

func TestIsOfficialLayoutCode(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		want     bool
	}{
		{name: "8 digit klid", sourceID: "00000804", want: true},
		{name: "0x klid", sourceID: "0x411", want: true},
		{name: "locale tag", sourceID: "ja-JP", want: false},
		{name: "random", sourceID: "abc", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOfficialLayoutCode(tt.sourceID)
			if got != tt.want {
				t.Fatalf("isOfficialLayoutCode(%q) = %v, want %v", tt.sourceID, got, tt.want)
			}
		})
	}
}

func TestShouldSwitchLayoutForSourceID(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		want     bool
	}{
		{name: "klid should switch layout", sourceID: "00000804", want: true},
		{name: "0x code should switch layout", sourceID: "0x411", want: true},
		{name: "language tag should not switch layout", sourceID: "en-US", want: false},
		{name: "cjk language tag should not switch layout", sourceID: "zh-CN", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSwitchLayoutForSourceID(tt.sourceID)
			if got != tt.want {
				t.Fatalf("shouldSwitchLayoutForSourceID(%q) = %v, want %v", tt.sourceID, got, tt.want)
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

func TestParseRequestedLangID(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		want     uint32
		ok       bool
	}{
		{name: "locale tag", sourceID: "ja-JP", want: 0x0411, ok: true},
		{name: "klid hex", sourceID: "00000412", want: 0x0412, ok: true},
		{name: "0x value", sourceID: "0x409", want: 0x0409, ok: true},
		{name: "invalid", sourceID: "not-a-layout", want: 0, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseRequestedLangID(tt.sourceID)
			if ok != tt.ok {
				t.Fatalf("parseRequestedLangID(%q) ok = %v, want %v", tt.sourceID, ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("parseRequestedLangID(%q) = %#x, want %#x", tt.sourceID, got, tt.want)
			}
		})
	}
}

func TestNormalizeKLID(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		want     string
		ok       bool
	}{
		{name: "existing klid", sourceID: "00000411", want: "00000411", ok: true},
		{name: "0x format", sourceID: "0x412", want: "00000412", ok: true},
		{name: "locale tag", sourceID: "ko-KR", want: "00000412", ok: true},
		{name: "invalid", sourceID: "xx", want: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := normalizeKLID(tt.sourceID)
			if ok != tt.ok {
				t.Fatalf("normalizeKLID(%q) ok = %v, want %v", tt.sourceID, ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("normalizeKLID(%q) = %q, want %q", tt.sourceID, got, tt.want)
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

func TestShouldDisableNativeMode(t *testing.T) {
	tests := []struct {
		name     string
		sourceID string
		want     bool
	}{
		{name: "english disables native mode", sourceID: "en-US", want: true},
		{name: "english case-insensitive", sourceID: "EN-us", want: true},
		{name: "korean does not disable native mode", sourceID: "ko-KR", want: false},
		{name: "chinese does not disable native mode", sourceID: "zh-CN", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldDisableNativeMode(tt.sourceID)
			if got != tt.want {
				t.Fatalf("shouldDisableNativeMode(%q) = %v, want %v", tt.sourceID, got, tt.want)
			}
		})
	}
}
