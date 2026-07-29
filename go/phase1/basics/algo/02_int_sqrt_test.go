package algo

import "testing"

// 02. 答えで二分探索（整数平方根） のテスト
// 入力例: n=26 → 5（5*5=25 <= 26 < 36）
// 境界:
//   - n=0 → 0
//   - n=1 → 1
//   - 完全平方数: n=4 → 2, n=25 → 5
//   - 大きな n: k*k のオーバーフローに注意（k <= n/k へ書き換え可）
//
// TODO: テストを書く
func TestFindLargestInteger(t *testing.T) {
	cases := []struct {
		name  string
		input int // n
		want  int // k
	}{
		{name: "normal", input: 26, want: 5},
		{name: "inputNisZero", input: 0, want: 0},
		{name: "inputNisOne", input: 1, want: 1},
		{name: "perfectSquare", input: 4, want: 2},
		{name: "perfectSquare2", input: 25, want: 5},
		{name: "largeInputN", input: 2000000000, want: 44721}, // int: -2,147,483,648 - 2,147,483,647
		{name: "NisNegativeInteger", input: -1, want: -1},     // Nが負だったら、入力のnをそのまま返す
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := findLargestInteger(c.input)

			if gotValue != c.want {
				t.Errorf("期待する値と違います。gotValue=%d, want=%d", gotValue, c.want)
			}
		})
	}
}
