package collections

// Goには標準のフィルター関数がないらしい...
// 標準の関数のみで実装する

func filterSquaredEvens(input []int) ([]int, error) {

	// caption!: var squaredEvens []int -> nil slice
	squaredEvens := []int{}

	for _, v := range input {
		if v%2 == 0 {
			squaredEvens = append(squaredEvens, v*v)
		}
	}
	return squaredEvens, nil
}
