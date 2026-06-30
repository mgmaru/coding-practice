package collections

import "sort"

type Person struct {
	Last string
	Age  int
	// 修正：LastとAgeが同じ場合に、安定を見るために加えたフィールド
	Order int
}

func sortStable(input []Person) []Person {
	// 修正：非破壊に統一する。
	// 理由：今回のinputスライスはそこまでデータ量が多くないので、コピーして非破壊にする（入力のinputを直接操作しない）。

	out := make([]Person, len(input)) // 注意：make([]Person, 0 , len(input))と、初期サイズに0を設定してしまうと、だめ。
	// 理由：copy はコピー先の長さ分しかコピーしない

	// 修正：返すスライスをmakeで定義したことにより、nilチェックが不要

	copy(out, input)
	// 修正：sort.Sliceからsort.SliceStableに修正
	sort.SliceStable(out, func(i, j int) bool { // 姓を昇順でソート
		if out[i].Last == out[j].Last { // 姓が同じだったら、年齢の降順でソート
			return out[i].Age > out[j].Age
		}
		return out[i].Last < out[j].Last // それ以外は姓の昇順でソート
	})
	return out
}
