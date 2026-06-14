package collections

// Goには標準のフィルター関数がないらしい...
// 標準の関数のみで実装する

func filterSquaredEvens(input []int) []int {

	var squaredEven []int

	for i := range input {
		if i%2 == 0 {
			squaredEven = append(squaredEven, i)
		}
	}
	return squaredEven
}
