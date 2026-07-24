package algo

// 01. 二分探索 — README.md「Tier 1 / 1. 二分探索」参照
// 昇順ソート済み配列から target の位置を O(log n) で求める。無ければ -1。
// 不変条件: 「答えがあるなら常に探索範囲 [lo, hi) の中に居る」を宣言してから書く。
// TODO: 自分で実装する

func getTargetIndex(input []int, target int) int {

	m := make(map[int]int, len(input))
	// キー：要素、値：indexとしてmapに格納
	for i, v := range input { // 計算量O(n)
		if _, ok := m[v]; !ok { // 要素がmapに含まれていたら、重複なのでスキップ（一番小さいインデックスの要素が適用される）
			m[v] = i
		}
	}

	if _, ok := m[target]; !ok { // targetが含まれていない場合は-1を返す
		return -1
	}
	return m[target] // 計算量O(1)
}
