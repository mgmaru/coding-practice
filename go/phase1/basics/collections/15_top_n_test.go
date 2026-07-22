package collections

import (
	"reflect"
	"testing"
)

func TestDisplayTopN(t *testing.T) {
	cases := []struct {
		name  string
		input []Item
		topN  int
		want  []Item
	}{
		{name: "normal", input: []Item{{Name: "a", Frequency: 5}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: 8}, {Name: "d", Frequency: 1}}, topN: 2, want: []Item{{Name: "c", Frequency: 8}, {Name: "a", Frequency: 5}}},
		{name: "TwoItemsSameFrequency", input: []Item{{Name: "a", Frequency: 5}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: 5}, {Name: "d", Frequency: 1}}, topN: 2, want: []Item{{Name: "a", Frequency: 5}, {Name: "c", Frequency: 5}}}, // 同じ頻度の場合、降順
		{name: "TopNIsOne", input: []Item{{Name: "a", Frequency: 5}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: 8}, {Name: "d", Frequency: 1}}, topN: 1, want: []Item{{Name: "c", Frequency: 8}}},
		{name: "TopNisZero", input: []Item{{Name: "a", Frequency: 5}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: 8}, {Name: "d", Frequency: 1}}, topN: 0, want: []Item{}}, // topN <=0 の時、空のスライスを返す
		// 以下のテスト名とテストケースがあっていなかったので修正
		{name: "TopNIsNegativeInteger", input: []Item{{Name: "a", Frequency: 5}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: 8}, {Name: "d", Frequency: 1}}, topN: -1, want: []Item{}},
		{name: "EmptyInput", input: []Item{}, topN: 2, want: []Item{}},
		{name: "NilInput", input: nil, topN: 2, want: []Item{}},
		{name: "AllItemsTheSameFrequency", input: []Item{{Name: "a", Frequency: 2}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: 2}, {Name: "d", Frequency: 2}}, topN: 2, want: []Item{{Name: "a", Frequency: 2}, {Name: "b", Frequency: 2}}},
		{name: "AnElementIncludeNegativeFrequency", input: []Item{{Name: "a", Frequency: -3}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: 8}, {Name: "d", Frequency: 1}}, topN: 2, want: []Item{{Name: "c", Frequency: 8}, {Name: "b", Frequency: 2}}}, // 負の頻度の要素は件数ソートから除外する。
		{name: "ThereAreMoreNegativeFrequencyElementsThanTopN", input: []Item{{Name: "a", Frequency: -1}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: -3}, {Name: "d", Frequency: -5}}, topN: 2, want: []Item{{Name: "b", Frequency: 2}}},
		{name: "AllElementsIsNegativeFrequency", input: []Item{{Name: "a", Frequency: -1}, {Name: "b", Frequency: -2}, {Name: "c", Frequency: -3}, {Name: "d", Frequency: -5}}, topN: 2, want: []Item{}},
		{name: "InCludeEmptyItem", input: []Item{{Name: "a", Frequency: 5}, {Name: "b", Frequency: 2}, {Name: "c", Frequency: 8}, {Name: "", Frequency: 7}}, topN: 2, want: []Item{{Name: "c", Frequency: 8}, {Name: "", Frequency: 7}}},
		// 修正：Frequencyの境界値テストを追加
		{name: "IncludeZeroFrecuency", input: []Item{{Name: "a", Frequency: 5}, {Name: "b", Frequency: 0}, {Name: "c", Frequency: 8}, {Name: "d", Frequency: 0}}, topN: 3, want: []Item{{Name: "c", Frequency: 8}, {Name: "a", Frequency: 5}, {Name: "b", Frequency: 0}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := displayTopN(c.input, c.topN)

			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値ではありません。 gotValue=%v, want= %v", gotValue, c.want)
			}
		})
	}

}
