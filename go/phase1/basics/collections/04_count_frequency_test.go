package collections

import (
	"reflect"
	"testing"
)

func TestCountElements(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		want  []CountElement
	}{
		// 修正：タイポ「nomal」を「normal」に修正
		{"normal", []string{"apple", "banana", "apple", "cherry", "banana", "apple"}, []CountElement{{Element: "apple", Count: 3}, {Element: "banana", Count: 2}, {Element: "cherry", Count: 1}}},
		{"countDuplicated", []string{"apple", "banana", "apple", "cherry", "banana"}, []CountElement{{Element: "apple", Count: 2}, {Element: "banana", Count: 2}, {Element: "cherry", Count: 1}}},
		{"emptyInput", []string{}, []CountElement{}},
		{"nilInput", nil, []CountElement{}},
		// 修正：入力が単一の要素の場合、大文字小文字が混在している場合のテストケースを追加
		{"singleInput", []string{"x"}, []CountElement{{Element: "x", Count: 1}}},
		{"mixUpperCaseAndLowerCase", []string{"apple", "Apple", "Banana", "grape", "banana", "apple", "Banana"}, []CountElement{{Element: "Banana", Count: 2}, {Element: "apple", Count: 2}, {Element: "Apple", Count: 1}, {Element: "banana", Count: 1}, {Element: "grape", Count: 1}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			gotValue := CountElements(c.input)

			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する結果が違います。gotValue=%v, want=%v", gotValue, c.want)
			}
		})
	}
}
