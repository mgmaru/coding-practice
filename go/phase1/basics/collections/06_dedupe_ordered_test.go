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
		// 修正：型が違う要素に関するテストを追加
		{name: "sameNumberDifferentType", input: []any{1, int64(1), 1.0}, want: []any{1, int64(1), 1.0}}, // 3 つとも残る
		// 修正：全要素同じ場合のテストを追加
		{name: "allElementSame", input: []any{3, 3, 3, 3}, want: []any{3}},
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

// 修正：anyではない、具体的な型に関するスライスのテストを追加
func TestDeleteDuplicate_Concrete(t *testing.T) {
	gotValueInt := deleteDuplicate([]int{3, 1, 3, 2, 4})
	if !reflect.DeepEqual(gotValueInt, []int{3, 1, 2, 4}) {
		t.Errorf("期待する値が違います。 gotValue= %v, want= %v", gotValueInt, []int{3, 1, 2, 4})
	}

	gotValueString := deleteDuplicate([]string{"apple", "banana", "grape", "banana"})
	if !reflect.DeepEqual(gotValueString, []string{"apple", "banana", "grape"}) {
		t.Errorf("期待する値が違います。 gotValue= %v, want= %v", gotValueString, []string{"apple", "banana", "grape"})
	}
}
