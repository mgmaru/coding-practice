package collections

func calcMovingAverage(input []int, number int) []int {

	movingAverageResult := make([]int, 0, len(input))

	if number > len(input) { // error
		return movingAverageResult
	}
	if number <= 0 { // number >= 1
		return movingAverageResult
	}
	for i, _ := range input { // iを1個ずつずらしていく
		if i+number > len(input) { // 平均の計算対象となる数が要素を超えた場合は平均計算終わり
			return movingAverageResult
		}
		// number分、inputから対象となる数を取り出して平均を計算する
		targetNumsSum := 0
		for j := range number { // number分要素を取り出す
			targetNumsSum += input[i+j]
		}
		movingAverage := targetNumsSum / number
		movingAverageResult = append(movingAverageResult, movingAverage)
	}
	return movingAverageResult
}
