package collections

// 仕様：スライスを受け取って、その中の要素の重複を削除する関数
// 仕様：出現順は保つ
func deleteDuplicate(input []any) []any {

	noDuplicateSlice := make([]any, 0, len(input)) // 重複していない要素を格納する
	m := make(map[any]struct{}, len(input))        // 要素の存在を判定するmap

	for _, i := range input {
		if _, ok := m[i]; !ok { // 要素が存在しない
			m[i] = struct{}{} // 要素の存在を格納
			noDuplicateSlice = append(noDuplicateSlice, i)
		}
		// 重複している場合は何もしない
	}
	return noDuplicateSlice
}
