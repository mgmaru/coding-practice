package collections

// 仕様：スライスを受け取って、その中の要素の重複を削除する関数
// 仕様：出現順は保つ
// 修正：Go では値を map のキーに使うには comparable でなければならない -> ジェネリクスでTをcomparableの型として定義した。
func deleteDuplicate[T comparable](input []T) []T {

	noDuplicateSlice := make([]T, 0, len(input)) // 重複していない要素を格納する
	// 修正：存在するかどうかを判定するmapの変数名を「m」から「seen」に変更
	// 理由：変数名から責務が見られるようにするため。
	seen := make(map[T]struct{}, len(input)) // 要素の存在を判定するmap

	// 修正：「i」だとindexと勘違い指定しまうので、「v」に変更
	for _, v := range input {
		if _, ok := seen[v]; !ok { // 要素が存在しない
			seen[v] = struct{}{} // 要素の存在を格納
			noDuplicateSlice = append(noDuplicateSlice, v)
		}
		// 重複している場合は何もしない
	}
	return noDuplicateSlice
}
