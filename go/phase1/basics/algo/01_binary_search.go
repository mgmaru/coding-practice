package algo

// 01. 二分探索 — README.md「Tier 1 / 1. 二分探索」参照
// 昇順ソート済み配列から target の位置を O(log n) で求める。無ければ -1。
// 不変条件: 「答えがあるなら常に探索範囲 [lo, hi) の中に居る」を宣言してから書く。
// TODO: 自分で実装する

func getTargetIndexOriginal(input []int, target int) int {

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

// ソート済みのスライスに対して二分探索を実行する関数
func getTargetIndexBinarySearch(input []int, target int) int {

	n := len(input) // nは固定値なので、あらかじめ変数に格納しておく（nが大キックなった時にいちいち、lenでサイズを計算すると計算量が爆増するため）
	if n == 0 {
		return -1
	}

	// 初期値
	l := 0
	h := n - 1
	for l <= h {
		m := (l + h) / 2        // 修正：mの更新処理は共通処理なので最初にまとめた。// 追記：1回のループ処理でmは固定なので、違う分岐に入ることはない。
		if target == input[m] { // targetがinput[m]と同じ場合はその時点で探索終了
			// 初期の探索範囲
			lo := 0
			hi := m
			mid := lo + (hi-lo)/2
			for lo < hi { // 継続条件は loとhiが違う場合 -> lo==hiとなった場合はその時点で終了
				mid = lo + (hi-lo)/2
				if input[mid] == target { // midがtargetと一致する場合は、midよりも小さい範囲に境界がある。
					hi = mid - 1 //探索範囲を更新
				}
				if input[mid] != target { // midがtargetと一致する場合は、midよりも大きい範囲に教会がある。
					lo = mid + 1 // 探索範囲を更新
				}
			}
			return lo
		}
		if target > input[m] { // もし、targetがmidより大きい場合、midを含む左側の要素を削除し、最小要素と真ん中のインデックスを更新
			l = m + 1 // 最小を更新(midの１つ右の要素が最小要素)
		}
		if target < input[m] { // もし、targetがmidよりも小さい場合、midを含む右側の要素を削除し、最大要素と真ん中のインデックスを更新
			h = m - 1 // 最大を更新（midの１つ左の要素が最大要素）
		}
	}
	return -1
}
