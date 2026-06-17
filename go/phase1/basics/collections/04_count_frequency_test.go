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
		{"nomal", []string{"apple", "banana", "apple", "cherry", "banana", "apple"}, []CountElement{{element: "apple", count: 3}, {element: "banana", count: 2}, {element: "cherry", count: 1}}},
		{"countDuplicated", []string{"apple", "banana", "apple", "cherry", "banana"}, []CountElement{{element: "apple", count: 2}, {element: "banana", count: 2}, {element: "cherry", count: 1}}},
		{"emptyInput", []string{}, []CountElement{}},
		{"nilInput", nil, []CountElement{}},
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
