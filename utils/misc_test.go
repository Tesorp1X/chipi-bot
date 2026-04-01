package utils_test

import (
	"fmt"
	"testing"

	"github.com/Tesorp1X/chipi-bot/utils"
	tele "gopkg.in/telebot.v4"
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

const isResponded = "is_callback_query_responded"

func prepareTeleCtx() tele.Context {
	b, err := tele.NewBot(tele.Settings{Offline: true})
	if err != nil {
		fmt.Println("bot fucked up")
		return nil
	}

	return tele.NewContext(b, tele.Update{})
}

func makeTeleCtxWithTrue() tele.Context {
	c := prepareTeleCtx()
	c.Set(isResponded, true)

	return c
}

func makeTeleCtxWithFalse() tele.Context {
	c := prepareTeleCtx()
	c.Set(isResponded, false)

	return c
}

func makeTeleCtxWithNothing() tele.Context {
	c := prepareTeleCtx()
	if val := c.Get(isResponded); val != nil {
		c.Set(isResponded, nil)
	}

	return c
}

func TestIsCbQueryResponded(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		c    tele.Context
		want bool
	}{
		{
			name: "ctx with true-value in 'is_callback_query_responded'",
			c:    makeTeleCtxWithTrue(),
			want: true,
		},
		{
			name: "ctx with false-value in 'is_callback_query_responded'",
			c:    makeTeleCtxWithFalse(),
			want: false,
		},
		{
			name: "ctx with nil-value in 'is_callback_query_responded'",
			c:    makeTeleCtxWithNothing(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IsCbQueryResponded(tt.c)
			if got != tt.want {
				t.Errorf("IsCbQueryResponded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkCbQueryAsResponded(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		c               tele.Context
		wantIsResponded bool
		wantErr         bool
	}{
		{
			name:            "ctx with already a true-value in 'is_callback_query_responded'",
			c:               makeTeleCtxWithTrue(),
			wantIsResponded: true,
			wantErr:         true,
		},
		{
			name:            "ctx with false-value in 'is_callback_query_responded'",
			c:               makeTeleCtxWithFalse(),
			wantIsResponded: true,
			wantErr:         false,
		},
		{
			name:            "ctx with nil-value in 'is_callback_query_responded'",
			c:               makeTeleCtxWithNothing(),
			wantIsResponded: true,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := utils.MarkCbQueryAsResponded(tt.c)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MarkCbQueryAsResponded() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MarkCbQueryAsResponded() succeeded unexpectedly")
			}

			gotIsResponded := tt.c.Get(isResponded)
			if gotIsResponded != tt.wantIsResponded {

				t.Errorf(
					"MarkCbQueryAsResponded(): gotIsResponded %v, wantIsResponded %v",
					gotIsResponded,
					tt.wantIsResponded,
				)

			}
		})
	}
}
