package algo

// 02. 答えで二分探索（整数平方根） — README.md「Tier 1 / 2. 答えで二分探索」参照
// 非負整数 n に対し k*k <= n を満たす最大の整数 k（= floor(sqrt(n)））を二分探索で。標準ライブラリの sqrt は使わない。
// 核: 判定 k*k <= n が k について単調（境界まで真、その先すべて偽）。
// TODO: 自分で実装する

// アルゴリズム整理
// k*k < n となるような最大の整数kを求める　-> k*k < n と k*k > n と切り替わる境目を見つける
func findLargestInteger(n int) int {
	if n < 0 {
		return n
	}
	if n == 1 {
		return n
	}
	// 初期値計算
	lo := 0
	hi := n
	k := 0 //　求める最大の整数
	for lo < hi {
		mid := lo + (hi-lo)/2 // midを最初に更新して複数の分岐に入らないようにする
		if mid == n/mid {     // 完全平方数
			return mid
		}
		if mid < n/mid { // mid*mid < nだったら、見つける整数はmid以上
			lo = mid + 1 // 最小値を更新
			k = mid      // ただし、このmidは求める整数の候補になりうるので残す。 // ここ重要！！
		}
		if mid > n/mid { // mid*mid > n だったら、見つける整数はmid未満
			hi = mid // 最大値を更新
		}
	}
	return k
}
