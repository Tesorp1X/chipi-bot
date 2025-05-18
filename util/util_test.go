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

func TestCreateItemsListResponse(t *testing.T) {
	items := []models.Item{}
	items = append(items,
		models.Item{Name: "twix", Price: 79, Owner: models.OWNER_BOTH},
		models.Item{Name: "Сырок", Price: 69, Owner: models.OWNER_PAU},
		models.Item{Name: "Печенье", Price: 110, Owner: models.OWNER_LIZ},
	)

	want := `1) twix 79 руб
2) Сырок 69 руб
3) Печенье 110 руб
Лиз заплатила: 149.50 руб
Пау заплатил: 108.50 руб
Итого: 258 бублей.`

	got := util.CreateItemsListResponse(items...)

	if got != want {
		t.Errorf("Got:\n%s\n Wanted:\n%s", got, want)
	}
}
