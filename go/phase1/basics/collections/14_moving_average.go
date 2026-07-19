package collections

func calcMovingAverage(input []int, windowSize int) []int { //修正：平均はintで計算する。
	// 修正：変数名をwindowSizeからwindowSizeに変更

	if windowSize <= 0 || windowSize > len(input) { // windowSize >= 1
		return []int{}
	}
	result := make([]int, 0, len(input)-windowSize+1) // 修正：平均の計算の数は、入力のサイズ-窓の要素数+1

	// 修正：最初の窓の合計を計算する
	windowSum := 0
	for i := range windowSize {
		windowSum += input[i]
	}
	result = append(result, windowSum/windowSize)

	for i := windowSize; i < len(input); i++ { // 修正 i：入ってくる数 i-windosSize：出ていく数
		// 補足：もし、windowSizeがlenより大きかった場合、ループ処理には入らない。
		windowSum += input[i] - input[i-windowSize]
		movingAverage := windowSum / windowSize
		result = append(result, movingAverage)
	}
	return result
}
