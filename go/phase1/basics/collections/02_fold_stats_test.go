package collections

import (
	"bytes"
	"testing"
)

func TestCalcSumAndMaxAndMin(t *testing.T) {

	// バグ：ここでbufを定義してしまうと、テストケースを検証するごとにbufに文字列が書き込まれてしまう。
	// var buf bytes.Buffer

	cases := []struct {
		Name  string
		Input []int
		Want  string
	}{
		{"normal_1", []int{3, 1, 4, 1, 5, 9, 2, 6}, "合計=31, 最大=9, 最小=1"},
		// 修正：スライス要素が負の場合
		{"normal_2", []int{-3, -1, -5}, "合計=-9, 最大=-1, 最小=-5"},
		// 修正：要素が１つの場合
		{"normal_3", []int{42}, "合計=42, 最大=42, 最小=42"},
		{"nilSlice", nil, "入力されたスライスがnilです。"},
		{"emptySlice", []int{}, "入力されたスライスが空です。"},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {

			//修正：bufはテストケースごとに初期化（バグ：初期化しないと次々と書き込まれていく）
			var buf bytes.Buffer

			calcSumAndMaxAndMin(&buf, c.Input) // bufで文字列のbyteを受け取る。
			gotValue := buf.String()           // bufを文字列に変換

			// 比較(文字列を比較)
			if gotValue != c.Want {
				t.Errorf("期待する結果が違います。%s, %s\n", gotValue, c.Want)
			}

		})
	}
}
