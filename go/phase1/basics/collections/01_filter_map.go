package collections

// Goには標準のフィルター関数がないらしい...
// 標準の関数のみで実装する

var input = []int{1, 2, 3, 4, 5, 6}

func filterSquaredEvens([]int) []int {

	var squaredEven []int

	for i := range input {
		if i%2 == 0 {
			squaredEven = append(squaredEven, i)
		}
	}
	return squaredEven
}
