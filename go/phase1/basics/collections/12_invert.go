package collections

import "sort"

func invertKeyValue(input map[string]int) map[int][]string {

	invert_map := make(map[int][]string, len(input))
	for key, val := range input { // mapからキーを取り出しているので、順序を保証しない。
		invert_map[val] = append(invert_map[val], key)
	}
	// mapで出力するので順序は保証しない。
	// ただし、キー自体の並びは保証しないが、値のスライスに関してはスライスなので、並び替えて順序を保証できる。
	for _, value := range invert_map {
		sort.SliceStable(value, func(i, j int) bool {
			return value[i] < value[j] // 文字列の昇順でソート（大文字先小文字後（大文字の方がasciiの値が小さい） / a,b,c... / 空文字は最小）
		})
	}
	return invert_map
}
