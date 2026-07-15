package collections

import (
	"reflect"
	"testing"
)

func TestInvertKeyValue(t *testing.T) {
	cases := []struct {
		name  string
		input map[string]int
		want  map[int][]string
	}{
		{name: "NotDuplicatedValue", input: map[string]int{"a": 1, "b": 2, "c": 3}, want: map[int][]string{1: {"a"}, 2: {"b"}, 3: {"c"}}},
		{name: "DuplicatedValue", input: map[string]int{"a": 1, "b": 2, "c": 1}, want: map[int][]string{1: {"a", "c"}, 2: {"b"}}},
		{name: "IncludedEnptykey", input: map[string]int{"": 1, "a": 2, "b": 3}, want: map[int][]string{1: {""}, 2: {"a"}, 3: {"b"}}},
		{name: "AllDuplicatedValues", input: map[string]int{"a": 1, "b": 1, "c": 1}, want: map[int][]string{1: {"a", "b", "c"}}},
		{name: "IncludedUpperCaseKey", input: map[string]int{"a": 1, "b": 2, "A": 1}, want: map[int][]string{1: {"A", "a"}, 2: {"b"}}},
		{name: "IncluededUpperCaseKeyAndEmptyKey", input: map[string]int{"a": 1, "A": 1, "b": 2, "": 1}, want: map[int][]string{1: {"", "A", "a"}, 2: {"b"}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := invertKeyValue(c.input)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値が違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}
