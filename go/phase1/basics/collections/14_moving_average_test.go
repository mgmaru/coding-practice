package collections

import (
	"reflect"
	"testing"
)

func TestCalcMovingAverage(t *testing.T) {
	cases := []struct {
		name       string
		input      []int
		windowSize int
		want       []int
	}{
		{name: "AllPositiveIntegers", input: []int{1, 2, 3, 4, 5}, windowSize: 3, want: []int{2, 3, 4}},
		{name: "IncludeBothPositiveAndNegativeIntegers", input: []int{1, -2, 5, 10, -8}, windowSize: 3, want: []int{1, 4, 2}},
		{name: "AllNegativeIntegers", input: []int{-1, -2, -3, -4, -5}, windowSize: 3, want: []int{-2, -3, -4}},
		{name: "AllZeros", input: []int{0, 0, 0, 0, 0}, windowSize: 3, want: []int{0, 0, 0}},
		{name: "InputSizeIsSameNumberInput", input: []int{1, 2, 3, 4, 5}, windowSize: 5, want: []int{3}},
		{name: "EmptyInput", input: []int{}, windowSize: 3, want: []int{}},
		{name: "NilInput", input: nil, windowSize: 3, want: []int{}},
		{name: "windowSizeIsNegativeInteger", input: []int{1, 2, 3, 4, 5}, windowSize: -1, want: []int{}},
		{name: "windowSizeIsZero", input: []int{1, 2, 3, 4, 5}, windowSize: 0, want: []int{}},
		{name: "windowSizeExceedsInputSize", input: []int{1, 2, 3, 4, 5}, windowSize: 6, want: []int{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := calcMovingAverage(c.input, c.windowSize)

			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値とは違います。gotVlaue= %v, want= %v", gotValue, c.want)
			}
		})
	}

}
