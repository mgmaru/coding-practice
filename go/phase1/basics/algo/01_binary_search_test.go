package algo

import (
	"reflect"
	"testing"
)

// 01. 二分探索 のテスト（TDD: まず境界ケースを書く）
// 入力例: a=[1,3,5,7,9,11], target=7 → 3
// 境界:
//   - 見つからない: target=4 → -1
//   - 先頭/末尾: target=1 → 0, target=11 → 5
//   - 空配列: [] → -1
//   - 1要素: a=[5],target=5 → 0 / target=3 → -1
// TODO: テストを書く

func TestGetTargetIndexOriginal(t *testing.T) {
	cases := []struct {
		name   string
		input  []int
		target int
		want   int
	}{
		{name: "TargetIncludedInput", input: []int{1, 3, 5, 7, 9, 11}, target: 7, want: 3},
		{name: "TargetNotIncludedInput", input: []int{1, 3, 5, 7, 9, 11}, target: 4, want: -1},
		{name: "TargetBeginningOfInput", input: []int{1, 3, 5, 7, 9, 11}, target: 1, want: 0},
		{name: "TargetEndgOfInput", input: []int{1, 3, 5, 7, 9, 11}, target: 11, want: 5},
		{name: "EmptyInput", input: []int{}, target: 7, want: -1},
		{name: "NilInput", input: nil, target: 7, want: -1},
		{name: "InputIsOneElementAndTargetIncludedInput", input: []int{7}, target: 7, want: 0},
		{name: "InputIsOneElementAndTargetNotIncludedInput", input: []int{7}, target: 3, want: -1},
		{name: "InputContainsDuplicateElementsAndTargetNotduplicated", input: []int{1, 3, 7, 7, 9, 11}, target: 3, want: 1},
		{name: "InputContainsDuplicateElementsAndTargetduplicated", input: []int{1, 3, 7, 7, 9, 11}, target: 7, want: 2}, // ターゲットが重複した場合、インデックスが小さいものを返す
		{name: "AllInputSameElementAndIncludedTarget", input: []int{3, 3, 3, 3, 3, 3}, target: 3, want: 0},
		{name: "AllInputSameElementAndNotIncludedTarget", input: []int{3, 3, 3, 3, 3, 3}, target: 1, want: -1}, // 不要なような気がするが、Includedを確かめたら、Not Includedも確かめるべきかと思った
		// {name: "InputContainsNegativeInteger", input: []int{1, -3, 5, -7, 9, -11}, target: -7, want: 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := getTargetIndexOriginal(c.input, c.target)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}

func TestGetTargetIndexBinarySearch(t *testing.T) {
	cases := []struct {
		name   string
		input  []int
		target int
		want   int
	}{
		{name: "TargetIncludedInput", input: []int{1, 3, 5, 7, 9, 11}, target: 7, want: 3},
		{name: "TargetNotIncludedInput", input: []int{1, 3, 5, 7, 9, 11}, target: 4, want: -1},
		{name: "TargetBeginningOfInput", input: []int{1, 3, 5, 7, 9, 11}, target: 1, want: 0},
		{name: "TargetEndgOfInput", input: []int{1, 3, 5, 7, 9, 11}, target: 11, want: 5},
		{name: "EmptyInput", input: []int{}, target: 7, want: -1},
		{name: "NilInput", input: nil, target: 7, want: -1},
		{name: "InputIsOneElementAndTargetIncludedInput", input: []int{7}, target: 7, want: 0},
		{name: "InputIsOneElementAndTargetNotIncludedInput", input: []int{7}, target: 3, want: -1},
		{name: "InputContainsDuplicateElementsAndTargetNotduplicated", input: []int{1, 3, 7, 7, 9, 11}, target: 3, want: 1},
		{name: "InputContainsDuplicateElementsAndTargetduplicated", input: []int{1, 3, 7, 7, 9, 11}, target: 7, want: 2}, // ターゲットが重複した場合、インデックスが小さいものを返す
		{name: "AllInputSameElementAndIncludedTarget", input: []int{3, 3, 3, 3, 3, 3}, target: 3, want: 0},
		{name: "AllInputSameElementAndNotIncludedTarget", input: []int{3, 3, 3, 3, 3, 3}, target: 1, want: -1}, // 不要なような気がするが、Includedを確かめたら、Not Includedも確かめるべきかと思った
		// {name: "InputContainsNegativeInteger", input: []int{1, -3, 5, -7, 9, -11}, target: -7, want: 3}, // 修正：ソートできていないのでそもそも二分探索の利用条件に合っていないので削除
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := getTargetIndexBinarySearch(c.input, c.target)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}
