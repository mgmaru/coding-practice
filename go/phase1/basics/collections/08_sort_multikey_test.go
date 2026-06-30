package collections

import (
	"reflect"
	"testing"
)

func TestSortStable(t *testing.T) {
	cases := []struct {
		name  string
		input []Person
		want  []Person
	}{
		{name: "Normal", input: []Person{{Last: "Sato", Age: 30}, {Last: "Ito", Age: 25}, {Last: "Suzuki", Age: 40}}, want: []Person{{Last: "Ito", Age: 25}, {Last: "Sato", Age: 30}, {Last: "Suzuki", Age: 40}}},
		{name: "DuplicateLast", input: []Person{{Last: "Sato", Age: 30}, {Last: "Ito", Age: 25}, {Last: "Sato", Age: 40}}, want: []Person{{Last: "Ito", Age: 25}, {Last: "Sato", Age: 40}, {Last: "Sato", Age: 30}}},
		{name: "DuplicateAge", input: []Person{{Last: "Suzuki", Age: 40}, {Last: "Ito", Age: 40}, {Last: "Sato", Age: 30}}, want: []Person{{Last: "Ito", Age: 40}, {Last: "Sato", Age: 30}, {Last: "Suzuki", Age: 40}}},
		{name: "AllDuplicate", input: []Person{{Last: "Sato", Age: 30}, {Last: "Sato", Age: 30}, {Last: "Sato", Age: 30}}, want: []Person{{Last: "Sato", Age: 30}, {Last: "Sato", Age: 30}, {Last: "Sato", Age: 30}}},
		{name: "NilSlice", input: nil, want: []Person{}},
		{name: "EmptySlice", input: []Person{}, want: []Person{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := sortStable(c.input)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}
