package collections

// var input = [][]int{{1, 2}, {3}, {4, 5, 6}} // []intにnilが入ってくる場合があることに注意

func flattenList(input [][]int) []int {

	// 修正：平滑化後のスライスのサイズを計算
	max_flatten_list_size := 0
	for _, l := range input {
		max_flatten_list_size += len(l)
	}

	oneDimList := make([]int, 0, max_flatten_list_size) // capはlen(input)では足りない -> 修正：最初に平滑化後のリストを計算して、capに指定

	for _, numList := range input {
		oneDimList = append(oneDimList, numList...) // 修正：リストを展開してappend
	}
	return oneDimList
}
