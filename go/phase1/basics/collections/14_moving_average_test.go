package collections

import (
	"reflect"
	"testing"
)

func TestCalcMovingAverage(t *testing.T) {
	cases := []struct {
		name        string
		input       []int
		inputNumber int
		want        []int
	}{
		{name: "AllPositiveIntegers", input: []int{1, 2, 3, 4, 5}, inputNumber: 3, want: []int{2, 3, 4}},
		{name: "IncludeBothPositiveAndNegativeIntegers", input: []int{1, -2, 5, 10, -8}, inputNumber: 3, want: []int{1, 4, 2}},
		{name: "AllPositiveIntegers", input: []int{-1, -2, -3, -4, -5}, inputNumber: 3, want: []int{-2, -3, -4}},
		{name: "AllZeros", input: []int{0, 0, 0, 0, 0}, inputNumber: 3, want: []int{0, 0, 0}},
		{name: "InputSizeIsSameNumberInput", input: []int{1, 2, 3, 4, 5}, inputNumber: 5, want: []int{3}},
		{name: "EmptyInput", input: []int{}, inputNumber: 3, want: []int{}},
		{name: "NilInput", input: nil, inputNumber: 3, want: []int{}},
		{name: "InputNumberIsNegativeInteger", input: []int{1, 2, 3, 4, 5}, inputNumber: -1, want: []int{}},
		{name: "InputNumberIsZero", input: []int{1, 2, 3, 4, 5}, inputNumber: 0, want: []int{}},
		{name: "InputNumberExceedsInputSize", input: []int{1, 2, 3, 4, 5}, inputNumber: 6, want: []int{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := calcMovingAverage(c.input, c.inputNumber)

			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値とは違います。gotVlaue= %v, want= %v", gotValue, c.want)
			}
		})
	}

}
