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
