package collections

// 狙い: 前の状態を持ち越しながら走査する。
// 入力: [1, 2, 3, 4]
// 出力: [1, 3, 6, 10]

func calcCumulativeSum(input []int) []int {
	output := make([]int, 0, len(input))
	for i, num := range input {
		if len(output) == 0 { // i == 0でも良い？
			output = append(output, num)
			continue
		}
		output = append(output, output[i-1]+num)
	}
	return output
}
