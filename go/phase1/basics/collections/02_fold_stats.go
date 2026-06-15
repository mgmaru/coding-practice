package collections

import (
	"fmt"
	"io"
)

// スライス要素の合計と最大値、最小値を計算して表示する関数
// 文字列のbyteを格納する変数w io.Writerも渡してあげる。
func calcSumAndMaxAndMin(w io.Writer, input []int) {

	// 修正：max, minは組み込み関数で存在するので、よくない。
	// max -> maxValue, min -> minValueへ修正
	var sum int
	var maxValues int
	var minValues int

	// 入力inputがnilの場合は、何も返さない。
	if input == nil {
		fmt.Fprintf(w, "入力されたスライスがnilです。")
		return
	}

	// 修正1：空スライスの場合もループを通ってしまうのでガードする。
	// ループを通ると、合計0、最大値0、最小値0が返されてしまう。これは正しくない。
	// 修正2： inputが空スライスについて、「refrect.Deepequal()」で比較するよりも、「len(input) == 0」で比較した方が軽量。
	// 注意：ただし、inputがnilの場合、len(input) == 0となるので、inputがnilの場合をガードしてから、len(input) == 0 をしないと、nilスライスか空スライスかが判別できない。
	if len(input) == 0 {
		fmt.Fprintf(w, "入力されたスライスが空です。")
		return
	}

	// 設計意思：len(input) == 0とすれば、nilスライスと空のスライスを同時にガードできるがあえて、２つを分ける実装とした。
	// 理由：入力がnilスライスと空のスライスが渡されたかどうかを区別するため。デバッグ容易性。

	sum = input[0]
	maxValues = input[0]
	minValues = input[0]

	// 2個目の要素からループ処理
	for _, num := range input[1:] {

		sum += num

		// 要素が最大値を超えた場合、最大値を入れ替え
		if maxValues < num {
			maxValues = num
		}
		// 要素が最小値よりも小さい場合、最小値を入れ替え
		if minValues > num {
			minValues = num
		}
	}
	fmt.Fprintf(w, "合計=%d, 最大=%d, 最小=%d", sum, maxValues, minValues)
}
