package collections

import (
	"reflect"
	"testing"
)

func TestDeleteDuplicate(t *testing.T) {
	cases := []struct {
		name  string
		input []any
		want  []any
	}{
		{name: "duplicateIntSlice", input: []any{3, 1, 3, 2, 4}, want: []any{3, 1, 2, 4}},
		{name: "duplicateStringSlice", input: []any{"apple", "banana", "grape", "banana"}, want: []any{"apple", "banana", "grape"}},
		{name: "duplicateMultiSlice", input: []any{"apple", 3, 1, "banana", "apple", 7}, want: []any{"apple", 3, 1, "banana", 7}},
		{name: "twoDuplicateSlice", input: []any{3, 1, 2, 2, 3, 5}, want: []any{3, 1, 2, 5}},
		{name: "noDuplicateSlice", input: []any{3, 1, 2, 5}, want: []any{3, 1, 2, 5}},
		{name: "emptySlice", input: []any{}, want: []any{}},
		{name: "nilSlice", input: nil, want: []any{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := deleteDuplicate(c.input)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値が違います。 gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}
