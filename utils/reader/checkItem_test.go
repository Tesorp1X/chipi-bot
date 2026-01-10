package reader

import "testing"

func Test_extractPrice(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawItemText string
		pos         int
		want        float64
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:        "robot subtotal",
			rawItemText: "Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый] 1.0 18 699,00 18 699,00",
			pos:         SUBTOTAL,
			want:        18699.00,
			wantErr:     false,
		},
		{
			name:        "robot price",
			rawItemText: "Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый] 1.0 18 699,00 18 699,00",
			pos:         PRICE,
			want:        18699.00,
			wantErr:     false,
		},
		{
			name:        "french-dog subtotal",
			rawItemText: "Френч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98",
			pos:         SUBTOTAL,
			want:        199.98,
			wantErr:     false,
		},
		{
			name:        "french-dog price",
			rawItemText: "Френч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98",
			pos:         PRICE,
			want:        99.99,
			wantErr:     false,
		},
		{
			name:        "red pepper subtotal",
			rawItemText: "ПЕРЕЦ красный 1кг 0.176 364,31 64,12",
			pos:         SUBTOTAL,
			want:        64.12,
			wantErr:     false,
		},
		{
			name:        "red pepper subtotal",
			rawItemText: "ПЕРЕЦ красный 1кг 0.176 364,31 64,12",
			pos:         PRICE,
			want:        364.31,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := extractPrice(tt.rawItemText, tt.pos)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("extractPrice() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("extractPrice() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("extractPrice() = %f, want %f", got, tt.want)
			}
		})
	}
}

func Test_extractItemName(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawItemText string
		want        string
		wantErr     bool
	}{
		{
			rawItemText: "Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый] 1.0 18 699,00 18 699,00",
			want:        "Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый]",
			wantErr:     false,
		},
		{
			rawItemText: "Френч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98",
			want:        "Френч-дог с сосиской свиной (Кулинария ММ)",
			wantErr:     false,
		},
		{
			rawItemText: "ПЕРЕЦ красный 1кг 0.176 364,31 64,12",
			want:        "ПЕРЕЦ красный 1кг",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := extractItemName(tt.rawItemText)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("extractItemName() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("extractItemName() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if tt.want != got {
				t.Errorf("extractItemName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCheckItem(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawItemText string
		want        CheckItem
		wantErr     bool
	}{
		{
			name:        "robot subtotal",
			rawItemText: "Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый] 1.0 18 699,00 18 699,00",

			want: CheckItem{
				rawText:  "Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый] 1.0 18 699,00 18 699,00",
				Name:     "Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый]",
				Price:    18699.00,
				Amount:   1.0,
				SubTotal: 18699.00,
			},
			wantErr: false,
		},
		{
			name:        "french-dog subtotal",
			rawItemText: "Френч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98",
			want: CheckItem{
				rawText:  "Френч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98",
				Name:     "Френч-дог с сосиской свиной (Кулинария ММ)",
				Price:    99.99,
				Amount:   2.0,
				SubTotal: 199.98,
			},
			wantErr: false,
		},
		{
			name:        "red pepper subtotal",
			rawItemText: "ПЕРЕЦ красный 1кг 0.176 364,31 64,12",
			want: CheckItem{
				rawText:  "ПЕРЕЦ красный 1кг 0.176 364,31 64,12",
				Name:     "ПЕРЕЦ красный 1кг",
				Price:    364.31,
				Amount:   0.176,
				SubTotal: 64.12,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewCheckItem(tt.rawItemText)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewCheckItem() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewCheckItem() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !tt.want.IsEqual(&got) {
				t.Errorf("NewCheckItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
