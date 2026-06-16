package collections

import (
	"reflect"
	"testing"
)

// 設計意思：グローバル変数でテストケースを定義（2つの関数で使いまわすため）
type Case struct {
	Name   string
	inputA []any
	inputB []any
	Want   []any
}

var cases = []Case{
	{"IntSliceAndAisLonger", []any{1, 2, 3, 4, 5}, []any{2, 4, 6}, []any{2, 4}},
	{"IntSliceAndBisLonger", []any{2, 4, 6}, []any{1, 2, 3, 4, 5}, []any{2, 4}},
	{"StringSliceAndAisLonger", []any{"apple", "banana", "grape", "lemon"}, []any{"grape"}, []any{"grape"}},
	{"StringSlicendBisLonger", []any{"grape"}, []any{"apple", "banana", "grape", "lemon"}, []any{"grape"}},
	{"AnySlice", []any{1, "apple", 2, "lemon"}, []any{2, "lemon", 4, "banana"}, []any{2, "lemon"}},
	{"AisEmptySlice", []any{}, []any{2, 4, 6}, []any{}},
	{"BisEmptySlice", []any{1, 2, 3, 4, 5}, []any{}, []any{}},
	{"AisNilSlice", nil, []any{2, 4, 6}, []any{}},
	{"BisNilSlice", []any{1, 2, 3, 4, 5}, nil, []any{}},
	{"notIcluded", []any{1, 2, 3, 4, 5}, []any{6, 7, 8}, []any{}},
}

func TestLinearSearchElements(t *testing.T) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			gotValue := linearSearchElements(c.inputA, c.inputB)
			if !reflect.DeepEqual(gotValue, c.Want) {
				t.Errorf("期待する値が違います。 gotValue=%v , Want=%v", gotValue, c.Want)
			}
		})
	}
}

func TestSetSearchElements(t *testing.T) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			gotValue := linearSearchElements(c.inputA, c.inputB)
			if !reflect.DeepEqual(gotValue, c.Want) {
				t.Errorf("期待する値が違います。 gotValue=%v , Want=%v", gotValue, c.Want)
			}
		})
	}
}
