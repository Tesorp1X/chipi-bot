package util_test

import (
	"slices"
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

func TestExtractAdminsIDs(t *testing.T) {
	t.Run("line like [id, id]", func(t *testing.T) {
		s := "[123, 234, 456]"
		got := util.ExtractAdminsIDs(s)
		want := []int64{123, 234, 456}
		if !slices.Equal(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
	t.Run("line like id, id", func(t *testing.T) {
		s := "123, 234, 456"
		got := util.ExtractAdminsIDs(s)
		want := []int64{123, 234, 456}
		if !slices.Equal(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
	t.Run("line without spaces", func(t *testing.T) {
		s := "[123,234,456]"
		got := util.ExtractAdminsIDs(s)
		want := []int64{123, 234, 456}
		if !slices.Equal(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
}

func TestCalculateCheckTotal(t *testing.T) {
	items := []models.Item{
		{
			CheckId: 1,
			Price:   123,
			Owner:   models.OWNER_LIZ,
		},
		{
			CheckId: 1,
			Price:   23,
			Owner:   models.OWNER_LIZ,
		},
		{
			CheckId: 1,
			Price:   222,
			Owner:   models.OWNER_BOTH,
		},
		{
			CheckId: 1,
			Price:   73,
			Owner:   models.OWNER_BOTH,
		},
		{
			CheckId: 1,
			Price:   222,
			Owner:   models.OWNER_PAU,
		},
		{
			CheckId: 1,
			Price:   50,
			Owner:   models.OWNER_PAU,
		},
		{
			CheckId: 1,
			Price:   321,
			Owner:   models.OWNER_PAU,
		},
	}

	check := &models.CheckWithItems{Id: 1}
	check.SetCheck(&models.Check{Name: "test 1", Owner: models.OWNER_PAU})
	check.SetItems(items)

	wantTotal := models.CheckTotal{
		Id:          1,
		OwnerId:     models.OWNER_PAU,
		OwnerTotal:  740.5,
		DebtorTotal: 293.5,
		Total:       1034,
	}

	gotTotal := util.CalculateCheckTotal(check)

	if gotTotal.Id != wantTotal.Id {
		t.Fatalf("wanted id: %d got id: %d", wantTotal.Id, gotTotal.Id)
	}

	if gotTotal.OwnerId != wantTotal.OwnerId {
		t.Fatalf("wanted OwnerId: %s got OwnerId: %s", wantTotal.OwnerId, gotTotal.OwnerId)
	}

	if gotTotal.OwnerTotal != wantTotal.OwnerTotal {
		t.Fatalf("wanted OwnerTotal: %f got OwnerTotal: %f", wantTotal.OwnerTotal, gotTotal.OwnerTotal)
	}

	if gotTotal.DebtorTotal != wantTotal.DebtorTotal {
		t.Fatalf("wanted OwnerDebtorTotal: %f got DebtorTotal: %f", wantTotal.DebtorTotal, gotTotal.DebtorTotal)
	}

	if gotTotal.Total != wantTotal.Total {
		t.Fatalf("wanted Total: %f got Total: %f", wantTotal.Total, gotTotal.Total)
	}
}
