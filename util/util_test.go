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

	want := `<i>1) twix 79.00 руб</i>
<i>2) Сырок 69.00 руб</i>
<i>3) Печенье 110.00 руб</i>

Лиз заплатила: <b>149.50 руб</b>
Пау заплатил: <b>108.50 руб</b>

Итого: <b>258.00 бублей.</b>`

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

func TestCalculateSessionTotal(t *testing.T) {
	items1 := []models.Item{
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
	items2 := []models.Item{
		{
			CheckId: 1,
			Price:   123,
			Owner:   models.OWNER_LIZ,
		},
		{
			CheckId: 1,
			Price:   73,
			Owner:   models.OWNER_BOTH,
		},
		{
			CheckId: 1,
			Price:   22,
			Owner:   models.OWNER_PAU,
		},
	}
	items3 := []models.Item{
		{
			CheckId: 1,
			Price:   23,
			Owner:   models.OWNER_LIZ,
		},
		{
			CheckId: 1,
			Price:   22,
			Owner:   models.OWNER_BOTH,
		},
		{
			CheckId: 1,
			Price:   321,
			Owner:   models.OWNER_PAU,
		},
	}
	items4 := []models.Item{
		{
			CheckId: 1,
			Price:   123,
			Owner:   models.OWNER_PAU,
		},
		{
			CheckId: 1,
			Price:   23,
			Owner:   models.OWNER_PAU,
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
			Owner:   models.OWNER_LIZ,
		},
		{
			CheckId: 1,
			Price:   50,
			Owner:   models.OWNER_LIZ,
		},
		{
			CheckId: 1,
			Price:   321,
			Owner:   models.OWNER_LIZ,
		},
	}
	// total = 1034
	// totalLiz = 293.5
	// totalPau = 740.5
	check1 := &models.CheckWithItems{Id: 1}
	check1.SetCheck(&models.Check{Name: "test 1", Owner: models.OWNER_LIZ})
	check1.SetItems(items1)

	// total = 218
	// totalLiz = 159.5
	// totalPau = 58.5
	check2 := &models.CheckWithItems{Id: 2}
	check2.SetCheck(&models.Check{Name: "test 2", Owner: models.OWNER_LIZ})
	check2.SetItems(items2)

	// total = 366
	// totalLiz = 34
	// totalPau = 332
	check3 := &models.CheckWithItems{Id: 3}
	check3.SetCheck(&models.Check{Name: "test 3", Owner: models.OWNER_PAU})
	check3.SetItems(items3)

	// total = 1034
	// totalLiz = 740.5
	// totalPau = 293.5
	check4 := &models.CheckWithItems{Id: 4}
	check4.SetCheck(&models.Check{Name: "test 4", Owner: models.OWNER_PAU})
	check4.SetItems(items4)

	type args struct {
		sessionId int64
		checks    []*models.CheckWithItems
	}
	tests := []struct {
		name string
		args args
		want *models.SessionTotal
	}{
		{
			name: "3 checks. pau -> liz",
			args: args{
				sessionId: 1,
				checks:    []*models.CheckWithItems{check1, check2, check3},
			},
			want: &models.SessionTotal{SessionId: 1, Total: 1618, Recipient: models.OWNER_LIZ, Amount: 765},
		},
		{
			name: "2 checks. pau -> liz",
			args: args{
				sessionId: 2,
				checks:    []*models.CheckWithItems{check2, check3},
			},
			want: &models.SessionTotal{SessionId: 2, Total: 584, Recipient: models.OWNER_LIZ, Amount: 24.5},
		},
		{
			name: "3 checks. liz -> pau",
			args: args{
				sessionId: 3,
				checks:    []*models.CheckWithItems{check2, check3, check4},
			},
			want: &models.SessionTotal{SessionId: 3, Total: 1618, Recipient: models.OWNER_PAU, Amount: 716},
		},
		{
			name: "2 checks. parity",
			args: args{
				sessionId: 4,
				checks:    []*models.CheckWithItems{check1, check4},
			},
			want: &models.SessionTotal{SessionId: 4, Total: 2068, Recipient: models.OWNER_LIZ, Amount: 0},
		},
	}

	validate := func(a, b *models.SessionTotal) bool {
		if a.Total != b.Total {
			return false
		}
		if a.SessionId != b.SessionId {
			return false
		}
		if a.Amount != b.Amount {
			return false
		}
		// if fuunction haven't returned yet, means  a.Amount and b.Amount are equal
		if a.Recipient != b.Recipient && a.Amount != 0 {
			return false
		}
		return true
	}

	for _, tt := range tests {
		got := util.CalculateSessionTotal(tt.args.sessionId, tt.args.checks)
		if !validate(got, tt.want) {
			t.Fatalf("want: %+v\ngot: %+v", *(tt.want), *got)
		}
	}
}
