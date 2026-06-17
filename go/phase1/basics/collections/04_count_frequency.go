package collections

import (
	"sort"
)

// 課題：要素ごとの出現回数を数える。

// 要素の数をカウントする関数
// 仕様：string型のスライスを受け取って、その中の要素と出現回数をペアで出力する。
// 仕様：出現回数の降順でソートする。同数の場合、キーの昇順でソートする。
// 仕様：入力はstringのスライス。出力はCountElementのスライス。

type CountElement struct {
	element string // 要素
	count   int    // 出現回数
}

func CountElements(input []string) []CountElement {

	if len(input) == 0 {
		return []CountElement{}
	}

	m := make(map[string]int, 0)

	// 要素と出現回数をペアで格納
	for _, v := range input {
		m[v]++
	}

	counts := make([]CountElement, 0, len(m))

	// 構造体CountElementにマッピング
	for k, v := range m {
		counts = append(counts, CountElement{element: k, count: v})
	}

	// countを降順でソート。countが同数だった場合は、要素を昇順でソート
	sort.Slice(counts, func(i, j int) bool {

		// 出現回数が同数だった場合、要素を昇順でソート
		if counts[i].count == counts[j].count {
			return counts[i].element < counts[j].element
		}

		return counts[i].count > counts[j].count
	})

	return counts
}
