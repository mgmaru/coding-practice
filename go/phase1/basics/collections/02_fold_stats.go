package collections

import (
	"fmt"
	"io"
)

// スライス要素の合計と最大値、最小値を計算して表示する関数
// 文字列のbyteを格納する変数w io.Writerも渡してあげる。
func calcSumAndMaxAndMin(w io.Writer, input []int) {

	var sum int
	var max int
	var min int

	// 入力inputがnilの場合は、何も返さない。
	if input == nil {
		fmt.Fprintf(w, "入力されたスライスがnilです。")
		return
	}

	for i, num := range input {

		// 要素1個目は初期値を代入
		if i == 0 {
			max = num
			min = num
		}

		// 要素2個目以降

		sum += num

		// 要素が最大値を超えた場合、最大値を入れ替え
		if max < num {
			max = num
		}
		// 要素が最小値よりも小さい場合、最小値を入れ替え
		if min > num {
			min = num
		}
	}
	fmt.Fprintf(w, "合計=%d, 最大=%d, 最小=%d", sum, max, min)
}
