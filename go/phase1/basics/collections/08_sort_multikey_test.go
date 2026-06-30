package collections

import (
	"reflect"
	"testing"
)

func TestSortStable(t *testing.T) {
	cases := []struct {
		name  string
		input []Person
		want  []Person
	}{
		{name: "Normal", input: []Person{{Last: "Sato", Age: 30, Order: 1}, {Last: "Ito", Age: 25, Order: 2}, {Last: "Suzuki", Age: 40, Order: 3}}, want: []Person{{Last: "Ito", Age: 25, Order: 2}, {Last: "Sato", Age: 30, Order: 1}, {Last: "Suzuki", Age: 40, Order: 3}}},
		{name: "DuplicateLast", input: []Person{{Last: "Sato", Age: 30, Order: 1}, {Last: "Ito", Age: 25, Order: 2}, {Last: "Sato", Age: 40, Order: 3}}, want: []Person{{Last: "Ito", Age: 25, Order: 2}, {Last: "Sato", Age: 40, Order: 3}, {Last: "Sato", Age: 30, Order: 1}}},
		{name: "DuplicateAge", input: []Person{{Last: "Suzuki", Age: 40, Order: 1}, {Last: "Ito", Age: 40, Order: 2}, {Last: "Sato", Age: 30, Order: 3}}, want: []Person{{Last: "Ito", Age: 40, Order: 2}, {Last: "Sato", Age: 30, Order: 3}, {Last: "Suzuki", Age: 40, Order: 1}}},
		{name: "AllDuplicate", input: []Person{{Last: "Sato", Age: 30, Order: 1}, {Last: "Sato", Age: 30, Order: 2}, {Last: "Sato", Age: 30, Order: 3}}, want: []Person{{Last: "Sato", Age: 30, Order: 1}, {Last: "Sato", Age: 30, Order: 2}, {Last: "Sato", Age: 30, Order: 3}}},
		{name: "NilSlice", input: nil, want: []Person{}},
		{name: "EmptySlice", input: []Person{}, want: []Person{}},
		// 修正：２つがかぶっている場合のテストを追加
		{name: "TwoDuplicate", input: []Person{{Last: "Sato", Age: 30, Order: 3}, {Last: "Ito", Age: 40, Order: 1}, {Last: "Sato", Age: 30, Order: 2}}, want: []Person{{Last: "Ito", Age: 40, Order: 1}, {Last: "Sato", Age: 30, Order: 3}, {Last: "Sato", Age: 30, Order: 2}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := sortStable(c.input)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}

// 退行検知テスト：安定性が壊れたら（例: 実装を sort.SliceStable → sort.Slice に戻す退行）赤になる見張り番。
// 上の TestSortStable（3要素）は「安定な結果を主張」できるが、小さい入力では sort.Slice でも偶然同じ結果になり
// 退行を捕まえられない。タイ（同 Last・同 Age）を多数（>12）並べると Go の sort は pdqsort 経路に入り、
// 非安定だとタイの入力順が壊れる。安定なら各群内で入力順（Order 昇順）が保たれる。
// litmus:「実装を sort.Slice に戻すと落ち、sort.SliceStable で通る」を満たす。
func TestSortStableRegression(t *testing.T) {
	const n = 30 // n>12 で挿入ソート経路を抜け、pdqsort 経路に入る
	var input []Person
	for i := 0; i < n; i++ {
		last := "B"
		if i%2 == 0 {
			last = "A"
		}
		// Last は A/B を交互、Age は全件同じ＝タイ。Order に入力順を記録（ソートキーには含めない）。
		input = append(input, Person{Last: last, Age: 10, Order: i})
	}

	// 期待値：Last 昇順で A 群 → B 群。各群内は入力順（Order 昇順）を保つ＝安定。
	var want []Person
	for _, last := range []string{"A", "B"} {
		for _, p := range input {
			if p.Last == last {
				want = append(want, p)
			}
		}
	}

	got := sortStable(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("安定性が壊れています（sort.Slice への退行の疑い）。\ngot = %v\nwant= %v", got, want)
	}
}
