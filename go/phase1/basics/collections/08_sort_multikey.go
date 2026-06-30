package collections

import "sort"

type Person struct {
	Last string
	Age  int
}

func sortStable(input []Person) []Person {

	// inputが空もしくはnilの時
	if len(input) == 0 {
		return []Person{} // inputがいk空もしくはnilの場合、共通で空スライスを返す
		// 理由：この関数を呼び出し側でnilか空スライスかを意識する必要をなくすため。
	}

	sort.Slice(input, func(i, j int) bool { // 姓を昇順でソート
		if input[i].Last == input[j].Last { // 姓が同じだったら、年齢の降順でソート
			return input[i].Age > input[j].Age
		}
		return input[i].Last < input[j].Last // それ以外は姓の昇順でソート
	})
	return input
}
