package collections

import (
	"bytes"
	"testing"
)

func TestCalcSumAndMaxAndMin(t *testing.T) {

	var buf bytes.Buffer

	cases := []struct {
		Name  string
		Input []int
		Want  string
	}{
		{"normal", []int{3, 1, 4, 1, 5, 9, 2, 6}, "合計=31, 最大=9, 最小=1"},
		{"nilSlice", nil, "入力されたスライスがnilです。"},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {

			calcSumAndMaxAndMin(&buf, c.Input) // bufで文字列のbyteを受け取る。
			gotValue := buf.String()           // bufを文字列に変換

			// 比較(文字列を比較)
			if gotValue != c.Want {
				t.Errorf("期待する結果が違います。%s, %s\n", gotValue, c.Want)
			}

		})
	}
}
