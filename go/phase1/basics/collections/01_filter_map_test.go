package collections

import (
	"reflect"
	"testing"
)

func TestFilterSquaredEvens(t *testing.T) {
	cases := []struct {
		Name      string
		Input     []int
		WantValue []int
		WantError error
	}{
		{"normal", []int{1, 2, 3, 4, 5, 6}, []int{4, 16, 36}, nil},
		{"noneEvens", []int{3, 5, 7}, []int{}, nil},
		{"emptyValue", []int{}, []int{}, nil},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {

			gotValue, gotError := filterSquaredEvens(c.Input)

			// 値が違う場合
			if !reflect.DeepEqual(gotValue, c.WantValue) {
				t.Errorf("値が違います。")
			}
			// エラーが違う場合
			if gotError != c.WantError {
				t.Errorf("エラーが違います。")
			}
		})
	}
}
