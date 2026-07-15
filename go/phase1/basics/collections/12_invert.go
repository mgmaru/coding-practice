package collections

import (
	"slices"
)

func invertKeyValue(input map[string]int) map[int][]string {
	inverted := make(map[int][]string, len(input))
	for key, val := range input { // mapからキーを取り出しているので、順序を保証しない。
		inverted[val] = append(inverted[val], key)
	}
	// mapで出力するので順序は保証しない。
	// ただし、キー自体の並びは保証しないが、値のスライスに関してはスライスなので、並び替えて順序を保証できる。
	for _, value := range inverted {
		slices.Sort(value) // 修正
	}
	return inverted
}
