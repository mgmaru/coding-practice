package collections

// 課題：Aの要素がBに含まれるかどうかを判定し、含まれるものを返す。

// 線形探索
// 設計意思：課題では、[]int型の例を示しているが、スライスの中にstringなどの他の型が入ってきても良いように、any型にした。
// 設計意思：今回は、あくまで線形探索およびsetによる探索の比較であるので、[]any型とした。しかし、本来の実装においては、any型は使用しない方が良いと考える。

// 修正（観点２）：線形探索もset探索も重複を許さない実装をする。
// 根拠：課題では、「Aの各要素がBに含まれるか」ということなので、要素の重複は無視して良いと解釈。
// Aにある要素が複数含まれていようが、Bにある要素が複数福間荒れていようが、その要素が含まれるのには変わりはない（個数は関係ない）。

// 修正（観点２）：スライスの重複要素を削除する汎用関数を追加 / 計算量 O(2n) -> O(n)
// 実装工夫：map[要素]空のstructで抜き出してあげて、その要素を新しいスライスにappendして重複しないスライスを作成する（mapのキーが重複しないことを利用）。
// 補足：deleteDuplicatedElementsについてはテストは実施しない（今回の課題の趣旨とは離れる）。

func deleteDuplicateElements(input []any) []any {

	// 重複する要素を除いたスライス, 最大容量：len(input)
	noIncludeDuplicatedElements := make([]any, 0, len(input))

	// 要素と空のstructを持つmapを定義
	m := make(map[any]struct{}, 0)

	// バグ修正：
	for _, e := range input {
		// もし、要素（キー）がmに存在しない場合は、mに要素を追加し、スライスに要素を追加する。
		if _, ok := m[e]; !ok {
			m[e] = struct{}{}
			noIncludeDuplicatedElements = append(noIncludeDuplicatedElements, e)
		}
		// mに要素が存在する場合はスキップ（何もしない）
	}

	return noIncludeDuplicatedElements
}

func linearSearchElements(inputA []any, inputB []any) []any {

	// 修正（観点２）：AもしくはBの重複によって出力結果の要素が重複してしまうので、あらかじめ入力AおよびBの中で重複を取り除いておく。　-> 重複をガード
	// 線形探索は、inputAおよびinputBの重複によって、出力するスライスに要素の重複が現れる。
	noDuplicatedinputA := deleteDuplicateElements(inputA)
	noDuplicatedinputB := deleteDuplicateElements(inputB)

	// 一致する要素を格納する配列
	containingElements := make([]any, 0, len(inputB))
	// 設計意思：入力がnilもしくは空スライスの場合のガードはしない。
	// 根拠：make([]any, 0)でスライスを初期化しているので、入力にnilもしくは空スライスが入ってきた時に返されるのは、空スライスが返る。これは想定した挙動である。
	// また、入力Aの要素がBに含まれていなかった場合にも、空のスライスが返る。これも想定した挙動。

	// 線形探索について、前の実装だと、要素aを取り出してbの要素と一致するかを反映していたので、一致する要素に対して、aの要素の並び順で出力される。
	// しかしset探索は、bの要素をsetと比較していたので、bの要素の並び順で出力されている。　->　同じ入力に対して、線形探索とset探索で出力結果が違う。
	// 修正（観点３）：線形探索について、bの要素を取り出して、aの要素と比較することで、set探索の出力の並びと合わせることとした。
	for _, b := range noDuplicatedinputB {
		// 修正（観点2）：containingElementsの中に要素aが含まれていた場合、スキップ(重複は許さない)
		for _, a := range noDuplicatedinputA {
			if b == a {
				containingElements = append(containingElements, b) //ここ、aで書くべきかbで書くべきか...今回は、bのスライスでループを回しているのでbにした。
			}
		}
	}
	// 計算量：n個の要素に対してn個の操作　-> O(n^2)

	return containingElements
}

// set探索
func setSearchElements(inputA []any, inputB []any) []any {

	// 修正（観点２）：set探索の場合は、inputBの重複によって結果のスライスが重複する可能性 -> 重複をガード
	noDuplicatedinputB := deleteDuplicateElements(inputB)

	// 一致する要素を格納する配列
	containingElements := make([]any, 0, len(inputB))

	// キーにはany、値には空のstructが入るmapを定義
	set := make(map[any]struct{})

	// setのキーにAの要素、値に空のstructを格納（Aの要素の集合を格納）
	// inputAの重複をsetによって畳まれる
	for _, a := range inputA {
		set[a] = struct{}{}
	}

	// inputBの要素１つ１つに対して、setに含まれるか判定
	for _, b := range noDuplicatedinputB {
		if _, ok := set[b]; ok {
			containingElements = append(containingElements, b)
		}
	}
	// 計算量：n個の要素に対してn+n(=2n)回の操作　-> O(2n) -> O(n) （定数倍無視）

	return containingElements
}
