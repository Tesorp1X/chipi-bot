package reader

import (
	"reflect"
	"testing"
	"time"
)

var (
	date1, _ = time.Parse(time.DateTime, "2025-12-12 18:40:00")
	date2, _ = time.Parse(time.DateTime, "2025-12-22 03:48:00")
)

func Test_extractTime(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawLine string
		want    time.Time
		wantErr bool
	}{
		{
			rawLine: "12.12.2025   18:40:00\nФренч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98\n           ",
			want:    date1,
		},
		{
			rawLine: "22.12.2025   03:48:00\nПылесос-робот Dreame D20 [Лидар, 5200 мА*ч,\n13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр,\nкартография, белый]\n1.0 18 699,00 18 699,00\n           ",
			want:    date2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := extractTime(tt.rawLine)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("extractTime() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("extractTime() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !tt.want.Equal(got) {
				t.Errorf("extractTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	magnitRawItems = []string{
		"12.12.2025   18:40:00\nФренч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98\n           ",
		"\nРОССИЯ Шоколад молочн миндаль/вафля 75/82г\nва\n1.0 89,99 89,99\n           ",
		"\nЮБИЛЕЙНОЕ Печенье витамин мол с глаз 232г п/у 1.0 84,99 84,99\n           ",
		"\nЛЕГЕНДА ГОР АРХЫЗ МинВод газ природ 1,5л 1.0 51,99 51,99\n           ",
		"\nЭлектронный платеж 510,94\nЦена, Сумма,\n510,94\nТовар или услуга Количество\nКассовый чек\n \n \nИтого\n \nПриход         Смена №353         Чек №55         СНО: ОСН         ФНС: nalog.ru\nВ т.ч. НДС           10%  7,64o           18%  71,17o\n \nФД: 30614         ФПД: 2852541763         ФН: 7384440800325689         РН ККТ: 0007671471032948\nАКЦИОНЕРНОЕ ОБЩЕСТВО \"ТАНДЕР\"         ИНН 2310031475\n630003, Новосибирская обл, Новосибирск г, Владимировская ул, дом № 23",
	}
	magnitCleanItems = []string{
		"Френч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98",
		"РОССИЯ Шоколад молочн миндаль/вафля 75/82г ва 1.0 89,99 89,99",
		"ЮБИЛЕЙНОЕ Печенье витамин мол с глаз 232г п/у 1.0 84,99 84,99",
		"ЛЕГЕНДА ГОР АРХЫЗ МинВод газ природ 1,5л 1.0 51,99 51,99",
		"Электронный платеж 510,94 Цена, Сумма, 510,94 Товар или услуга Количество Кассовый чек     Итого   Приход         Смена №353         Чек №55         СНО: ОСН         ФНС: nalog.ru В т.ч. НДС           10%  7,64o           18%  71,17o   ФД: 30614         ФПД: 2852541763         ФН: 7384440800325689         РН ККТ: 0007671471032948 АКЦИОНЕРНОЕ ОБЩЕСТВО \"ТАНДЕР\"         ИНН 2310031475 630003, Новосибирская обл, Новосибирск г, Владимировская ул, дом № 23",
	}
	magnitExtractedItems = []CheckItem{
		{
			rawText:  "Френч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98",
			Name:     "Френч-дог с сосиской свиной (Кулинария ММ)",
			Price:    99.99,
			Amount:   2.0,
			SubTotal: 199.98,
		},
		{
			rawText:  "РОССИЯ Шоколад молочн миндаль/вафля 75/82г ва 1.0 89,99 89,99",
			Name:     "РОССИЯ Шоколад молочн миндаль/вафля 75/82г ва",
			Price:    89.99,
			Amount:   1.0,
			SubTotal: 89.99,
		},
		{
			rawText:  "ЮБИЛЕЙНОЕ Печенье витамин мол с глаз 232г п/у 1.0 84,99 84,99",
			Name:     "ЮБИЛЕЙНОЕ Печенье витамин мол с глаз 232г п/у",
			Price:    84.99,
			Amount:   1.0,
			SubTotal: 84.99,
		},
		{
			rawText:  "ЛЕГЕНДА ГОР АРХЫЗ МинВод газ природ 1,5л 1.0 51,99 51,99",
			Name:     "ЛЕГЕНДА ГОР АРХЫЗ МинВод газ природ 1,5л",
			Price:    51.99,
			Amount:   1.0,
			SubTotal: 51.99,
		},
	}
)

func Test_normalizeItems(t *testing.T) {

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawLines []string
		want     []string
		wantErr  bool
	}{
		{
			name:     "magnit",
			rawLines: magnitRawItems,
			want:     magnitCleanItems,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := normalizeItems(tt.rawLines)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("normalizeItems() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("normalizeItems() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("normalizeItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractItems(t *testing.T) {
	verify := func(a, b []CheckItem) bool {
		if len(a) != len(b) {
			return false
		}

		for i, itemA := range a {
			if !itemA.IsEqual(&b[i]) {
				return false
			}
		}

		return true
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		lines   []string
		want    []CheckItem
		wantErr bool
	}{
		{
			name:  "magnit",
			lines: magnitCleanItems,
			want:  magnitExtractedItems,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := extractItems(tt.lines)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("extractItems() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("extractItems() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !verify(tt.want, got) {
				t.Errorf("extractItems() = got:\n%v\n, want:\n%v", got, tt.want)
			}
		})
	}
}

func Test_extractTotal(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawText string
		want    float64
		wantErr bool
	}{
		{
			rawText: "Электронный платеж 510,94 Цена, Сумма, 510,94 Товар или услуга Количество Кассовый чек     Итого   Приход         Смена №353         Чек №55         СНО: ОСН         ФНС: nalog.ru В т.ч. НДС           10%  7,64o           18%  71,17o   ФД: 30614         ФПД: 2852541763         ФН: 7384440800325689         РН ККТ: 0007671471032948 АКЦИОНЕРНОЕ ОБЩЕСТВО \"ТАНДЕР\"         ИНН 2310031475 630003, Новосибирская обл, Новосибирск г, Владимировская ул, дом № 23",
			want:    510.94,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := extractTotal(tt.rawText)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("extractTotal() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("extractTotal() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if tt.want != got {
				t.Errorf("extractTotal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractOrgName(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawText string
		want    string
		wantErr bool
	}{
		{
			rawText: "Электронный платеж 510,94 Цена, Сумма, 510,94 Товар или услуга Количество Кассовый чек     Итого   Приход         Смена №353         Чек №55         СНО: ОСН         ФНС: nalog.ru В т.ч. НДС           10%  7,64o           18%  71,17o   ФД: 30614         ФПД: 2852541763         ФН: 7384440800325689         РН ККТ: 0007671471032948 АКЦИОНЕРНОЕ ОБЩЕСТВО \"ТАНДЕР\"         ИНН 2310031475 630003, Новосибирская обл, Новосибирск г, Владимировская ул, дом № 23",
			want:    "АКЦИОНЕРНОЕ ОБЩЕСТВО \"ТАНДЕР\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := extractOrgName(tt.rawText)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("extractOrgName() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("extractOrgName() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if tt.want != got {
				t.Errorf("extractOrgName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCheckData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawCheckText string
		want         *CheckData
		wantErr      bool
	}{
		{
			name:         "magnit",
			rawCheckText: "12.12.2025   18:40:00\nФренч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98\n           В т.ч. НДС 20%\nРОССИЯ Шоколад молочн миндаль/вафля 75/82г\nва\n1.0 89,99 89,99\n           В т.ч. НДС 20%\nЮБИЛЕЙНОЕ Печенье витамин мол с глаз 232г п/у 1.0 84,99 84,99\n           В т.ч. НДС 20%\nХлеб для сэндвичей классический нарез 350г п/ 1.0 83,99 83,99\n           В т.ч. НДС 10%\nЛЕГЕНДА ГОР АРХЫЗ МинВод газ природ 1,5л 1.0 51,99 51,99\n           В т.ч. НДС 20%\nЭлектронный платеж 510,94\nЦена, Сумма,\n510,94\nТовар или услуга Количество\nКассовый чек\n \n \nИтого\n \nПриход         Смена №353         Чек №55         СНО: ОСН         ФНС: nalog.ru\nВ т.ч. НДС           10%  7,64o           18%  71,17o\n \nФД: 30614         ФПД: 2852541763         ФН: 7384440800325689         РН ККТ: 0007671471032948\nАКЦИОНЕРНОЕ ОБЩЕСТВО \"ТАНДЕР\"         ИНН 2310031475\n630003, Новосибирская обл, Новосибирск г, Владимировская ул, дом № 23",
			want: &CheckData{
				rawText:        "12.12.2025   18:40:00\nФренч-дог с сосиской свиной (Кулинария ММ) 2.0 99,99 199,98\n           В т.ч. НДС 20%\nРОССИЯ Шоколад молочн миндаль/вафля 75/82г\nва\n1.0 89,99 89,99\n           В т.ч. НДС 20%\nЮБИЛЕЙНОЕ Печенье витамин мол с глаз 232г п/у 1.0 84,99 84,99\n           В т.ч. НДС 20%\nХлеб для сэндвичей классический нарез 350г п/ 1.0 83,99 83,99\n           В т.ч. НДС 10%\nЛЕГЕНДА ГОР АРХЫЗ МинВод газ природ 1,5л 1.0 51,99 51,99\n           В т.ч. НДС 20%\nЭлектронный платеж 510,94\nЦена, Сумма,\n510,94\nТовар или услуга Количество\nКассовый чек\n \n \nИтого\n \nПриход         Смена №353         Чек №55         СНО: ОСН         ФНС: nalog.ru\nВ т.ч. НДС           10%  7,64o           18%  71,17o\n \nФД: 30614         ФПД: 2852541763         ФН: 7384440800325689         РН ККТ: 0007671471032948\nАКЦИОНЕРНОЕ ОБЩЕСТВО \"ТАНДЕР\"         ИНН 2310031475\n630003, Новосибирская обл, Новосибирск г, Владимировская ул, дом № 23",
				TimeOfCreation: date1,
				Items:          magnitExtractedItems,
				Total:          510.94,
				OrgName:        "АКЦИОНЕРНОЕ ОБЩЕСТВО \"ТАНДЕР\"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewCheckData(tt.rawCheckText)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewCheckData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewCheckData() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if tt.want.IsEqual(got) {
				t.Errorf("NewCheckData() = %v, want %v", got, tt.want)
			}
		})
	}
}
