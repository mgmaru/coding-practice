package collections

import (
	"reflect"
	"testing"
)

// 設計意思：グローバル変数でテストケースを定義（2つの関数で使いまわすため）
// 修正（観点６）：Caseのフィールドの最初の大文字小文字を統一（すべて大文字に統一）
type Case struct {
	Name   string
	InputA []any
	InputB []any
	Want   []any
}

var Cases = []Case{
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
	// 修正（観点２）：inputAおよびinputBに重複する要素が含まれている場合のテストケースを追加
	{"DupInA", []any{2, 2, 1, 3, 4}, []any{2}, []any{2}}, // 修正：重複要素は比較対象から外す（線形探索、set探索ともに）　-> 重複は許さない仕様とする
	{"DupInB", []any{2}, []any{2, 2, 1, 3, 4}, []any{2}}, // 修正：重複要素は比較対象から外す（線形探索、set探索ともに）　-> 重複は許さない仕様とする
	// 修正（観点3）：順序についてはbの要素の並びとする。
	{"Orderdiff", []any{4, 2}, []any{2, 4}, []any{2, 4}},
}

func TestLinearSearchElements(t *testing.T) {
	for _, c := range Cases {
		t.Run(c.Name, func(t *testing.T) {
			gotValue := linearSearchElements(c.InputA, c.InputB)
			if !reflect.DeepEqual(gotValue, c.Want) {
				t.Errorf("期待する値が違います。 gotValue=%v , Want=%v", gotValue, c.Want)
			}
		})
	}
}

func TestSetSearchElements(t *testing.T) {
	for _, c := range Cases {
		t.Run(c.Name, func(t *testing.T) {
			gotValue := setSearchElements(c.InputA, c.InputB)
			if !reflect.DeepEqual(gotValue, c.Want) {
				t.Errorf("期待する値が違います。 gotValue=%v , Want=%v", gotValue, c.Want)
			}
		})
	}
}
