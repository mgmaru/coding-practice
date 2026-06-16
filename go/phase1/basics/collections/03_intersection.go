package collections

// 課題：Aの要素がBに含まれるかどうかを判定し、含まれるものを返す。

// 線形探索
// 設計意思：課題では、[]int型の例を示しているが、スライスの中にstringなどの他の型が入ってきても良いように、any型にした。
// 設計意思：今回は、あくまで線形探索およびsetによる探索の比較であるので、[]any型とした。しかし、本来の実装においては、any型は使用しない方が良いと考える。

func linearSearchElements(inputA []any, inputB []any) []any {

	// 一致する要素を格納する配列
	containingElements := make([]any, 0)
	// 設計意思：入力がnilもしくは空スライスの場合のガードはしない。
	// 根拠：make([]any, 0)でスライスを初期化しているので、入力にnilもしくは空スライスが入ってきた時に返されるのは、空スライスが返る。これは想定した挙動である。
	// また、入力Aの要素がBに含まれていなかった場合にも、空のスライスが返る。これも想定した挙動。

	for _, a := range inputA {
		for _, b := range inputB {
			if a == b {
				containingElements = append(containingElements, b) //ここ、aで書くべきかbで書くべきか...今回は、bのスライスでループを回しているのでbにした。
			}
		}
	}
	// 計算量：n個の要素に対してn個の操作　-> O(n^2)

	return containingElements
}

// set探索
func setSearchElements(inputA []any, inputB []any) []any {

	// 一致する要素を格納する配列
	containingElements := make([]any, 0)

	// キーにはany、値には空のstructが入るmapを定義
	set := make(map[any]struct{})

	// setのキーにAの要素、値に空のstructを格納（Aの要素の集合を格納）
	for _, a := range inputA {
		set[a] = struct{}{}
	}

	// inputBの要素１つ１つに対して、setに含まれるか判定
	for _, b := range inputB {
		if _, ok := set[b]; ok {
			containingElements = append(containingElements, b)
		}
	}
	// 計算量：n個の要素に対してn+n(=2n)回の操作　-> O(2n) -> O(n) （定数倍無視）

	return containingElements
}
