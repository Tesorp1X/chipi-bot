package util_test

import (
	"testing"

	"github.com/Tesorp1X/chipi-bot/models"
	"github.com/Tesorp1X/chipi-bot/util"
)

func TestExtractDataFromCallback(t *testing.T) {
	type args struct {
		data   string
		action models.CallbackAction
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "first",
			args: args{
				data:   "\fstart|123",
				action: models.CallbackAction("start"),
			},
			want: "123",
		},
		{
			name: "second",
			args: args{
				data:   "\fcheck_owner|liz",
				action: models.CallbackAction("check_owner"),
			},
			want: "liz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.ExtractDataFromCallback(tt.args.data, tt.args.action); got != tt.want {
				t.Errorf("ExtractDataFromCallback(data: %v, action: %v) = %v, want %v",
					tt.args.data, tt.args.action.String(), got, tt.want)
			}
		})
	}
}
