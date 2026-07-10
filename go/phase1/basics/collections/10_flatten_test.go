package collections

import (
	"reflect"
	"testing"
)

func TestFlattenList(t *testing.T) {
	cases := []struct {
		name  string
		input [][]int
		want  []int
	}{
		{name: "Normal", input: [][]int{{1, 2}, {3}, {4, 5}}, want: []int{1, 2, 3, 4, 5}},
		{name: "SingleElement", input: [][]int{{1}, {2}, {3}, {4}, {5}}, want: []int{1, 2, 3, 4, 5}},
		{name: "ContainsEmptyElement", input: [][]int{{1, 2}, {}, {4, 5}}, want: []int{1, 2, 4, 5}},
		{name: "AllEmptyElements", input: [][]int{{}, {}, {}}, want: []int{}},
		{name: "ContainsNilElement", input: [][]int{{1, 2}, nil, {4, 5}}, want: []int{1, 2, 4, 5}},
		{name: "AllNilElements", input: [][]int{nil, nil, nil}, want: []int{}},
		// 修正：テストケース追加
		{name: "EmptyInput", input: [][]int{}, want: []int{}},
		{name: "NilInput", input: nil, want: []int{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotvalue := flattenList(c.input)

			if !reflect.DeepEqual(gotvalue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotvalue, c.want)
			}
		})
	}
}
