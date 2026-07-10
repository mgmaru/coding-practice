package collections

// var input = [][]int{{1, 2}, {3}, {4, 5, 6}} // []intにnilが入ってくる場合があることに注意

func smoothTheList(input [][]int) []int {
	oneDimList := make([]int, 0) // capはlen(input)では足りない

	for _, numList := range input {
		for _, num := range numList { // numListが空の場合、appendは実行されない
			oneDimList = append(oneDimList, num)
		}
	}
	return oneDimList
}
