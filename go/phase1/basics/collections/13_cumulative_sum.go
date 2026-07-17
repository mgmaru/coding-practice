package collections

// 狙い: 前の状態を持ち越しながら走査する。
// 入力: [1, 2, 3, 4]
// 出力: [1, 3, 6, 10]

func calcCumulativeSum(input []int) []int {
	output := make([]int, 0, len(input))
	sum := 0 // 修正：一時変数をおいて、デバッグ容易性を向上
	for _, num := range input {
		sum += num
		output = append(output, sum)
	}
	return output
}
