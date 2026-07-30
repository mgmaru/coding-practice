package algo

import (
	"math"
	"testing"
)

// 02. 答えで二分探索（整数平方根） のテスト
// 入力例: n=26 → 5（5*5=25 <= 26 < 36）
// 境界:
//   - n=0 → 0
//   - n=1 → 1
//   - 完全平方数: n=4 → 2, n=25 → 5
//   - 大きな n: k*k のオーバーフローに注意（k <= n/k へ書き換え可）
//
// TODO: テストを書く
func TestIntSqrt(t *testing.T) {
	cases := []struct {
		name string
		n    int // n : 非負整数、入力
		want int // k :floor(sqrt(n))  // floor(x): floorはx以下の最大の整数を求める
	}{
		{name: "normal", n: 26, want: 5},
		{name: "nIsZero", n: 0, want: 0},
		{name: "nIsOne", n: 1, want: 1},
		{name: "perfectSquare", n: 4, want: 2},
		{name: "perfectSquare2", n: 25, want: 5},
		{name: "largeN", n: 2000000000, want: 44721},  // int: -9,223,372,036,854,775,808 - 9,223,372,036,854,775,807
		{name: "NisNegativeInteger", n: -1, want: -1}, // Nが負だったら、入力のnをそのまま返す
		// 追加
		{name: "nisTwo", n: 2, want: 1},                    // return kに到達するテスト（カバレッジ）
		{name: "maxInt", n: math.MaxInt, want: 3037000499}, // 型の限界ちょうど
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := intSqrt(c.n)

			if gotValue != c.want {
				t.Errorf("期待する値と違います。gotValue=%d, want=%d", gotValue, c.want)
			}
		})
	}
}
