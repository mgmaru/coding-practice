package collections

import (
	"reflect"
	"testing"
)

func TestCalcCumulativeSum(t *testing.T) {
	cases := []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "Normal", input: []int{1, 2, 3, 4}, want: []int{1, 3, 6, 10}},
		{name: "Normal2", input: []int{2, 4, 6, 8}, want: []int{2, 6, 12, 20}},
		{name: "SingleInput", input: []int{1}, want: []int{1}},
		{name: "EmptyInput", input: []int{}, want: []int{}},
		{name: "NilInput", input: nil, want: []int{}},
		// 修正：テスト追加
		{name: "AllNegativeIntegers", input: []int{-1, -3, -6}, want: []int{-1, -4, -10}},
		{name: "IncludesPositiveAndNegativeIntegers", input: []int{-1, 5, -7, 8}, want: []int{-1, 4, -3, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := calcCumulativeSum(c.input)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}
