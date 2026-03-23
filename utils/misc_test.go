package utils_test

import (
	"testing"

	"github.com/Tesorp1X/chipi-bot/utils"
)

func TestExtractCallbackData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawData string
		want    string
	}{
		{
			name:    "CallbackActionEditCheck",
			rawData: "\fEditCheck|go_back",
			want:    "go_back",
		},
		{
			name:    "CallbackActionEditItem",
			rawData: "\fEditItem|go_back",
			want:    "go_back",
		},
		{
			name:    "CallbackActionSelector",
			rawData: "\fSelector|go_back",
			want:    "go_back",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.ExtractCallbackData(tt.rawData)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("ExtractCallbackData() = %v, want %v", got, tt.want)
			}
		})
	}
}
