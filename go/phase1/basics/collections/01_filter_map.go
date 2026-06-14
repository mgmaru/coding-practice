package collections

import "errors"

var ErrNilInput = errors.New("入力がnilです。")

// Goには標準のフィルター関数がないらしい...
// 標準の関数のみで実装する
func filterSquaredEvens(input []int) ([]int, error) {

	// 修正：inputがnilの時をガード
	if input == nil {
		return nil, ErrNilInput
	}

	// caption!: var squaredEvens []int -> nil slice
	// 修正：makeで容量を指定。
	squaredEvens := make([]int, 0, len(input))

	for _, v := range input {
		if v%2 == 0 {
			squaredEvens = append(squaredEvens, v*v)
		}
	}
	return squaredEvens, nil
}
