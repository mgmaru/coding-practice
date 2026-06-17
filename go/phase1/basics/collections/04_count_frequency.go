package collections

import (
	"sort"
)

// 課題：要素ごとの出現回数を数える。

// 要素の数をカウントする関数
// 仕様：string型のスライスを受け取って、その中の要素と出現回数をペアで出力する。
// 仕様：出現回数の降順でソートする。同数の場合、キーの昇順でソートする。
// 仕様：入力はstringのスライス。出力はCountElementのスライス。

// 修正：構造体が公開になっていたので、それに対応してフィールドも公開する形へ修正
type CountElement struct {
	Element string // 要素
	Count   int    // 出現回数
}

func CountElements(input []string) []CountElement {

	// 修正：inputが空スライスおよびnilスライスの場合のガードは削除した。
	// 理由：makeでスライスは初期化しているために、入力がnilであっても、全て空配列で返る。

	// 修正：長さ（容量）を指定
	m := make(map[string]int, len(input))

	// 要素と出現回数をペアで格納
	for _, v := range input {
		m[v]++
	}

	counts := make([]CountElement, 0, len(m))

	// 構造体CountElementにマッピング
	for k, v := range m {
		counts = append(counts, CountElement{Element: k, Count: v})
	}

	// countを降順でソート。countが同数だった場合は、要素を昇順でソート
	sort.Slice(counts, func(i, j int) bool {

		// 出現回数が同数だった場合、要素を昇順でソート
		if counts[i].Count == counts[j].Count {
			return counts[i].Element < counts[j].Element
		}

		return counts[i].Count > counts[j].Count
	})

	return counts
}
