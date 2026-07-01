package collections

type Status struct { // ログ集計結果
	Date  string // nilは入らない
	Level string // nilは入らない
}

// ロジック設計
// 1. 入力から日付ごとのレベルをスライスにまとめてmapで抽出する。
// 2. 1で作成したmapから、キーを日付として、スライスを展開してレベルのカウントをする。
// 3. map of mapで出力（キー：日付 値：mapレベル:カウント）
// 注意：mapの構造自体が並びを保証できないため、map of mapのような構造では並びを保証できない仕様とした。
func aggregateDateStatus(input []Status) map[string]map[string]int {

	// inputが空で入ってきても、全て空スライスで返るのでガード不要

	m := make(map[string][]string, len(input)) // キー：日付 値：レベルのスライス
	for _, s := range input {
		m[s.Date] = append(m[s.Date], s.Level) //一旦、日付ごとのレベルを格納する // {06-01:[INFO, ERROR, INFO, WARNING, ...], 06-02:[INFO, WARNIG, ...]}
	}

	n := make(map[string]int, len(input)) // 日付ごと　-> キー：レベル　値：出現回数
	out := make(map[string]map[string]int, len(input))
	for key, vals := range m { // key:日付　vals:レベルのスライス // 注意：map mからforで取り出しているので、並び順が保証されていない。 -> ソートが必要
		for _, v := range vals {
			if _, ok := n[v]; ok { // もし、mの中にLevelがあったら、出現回数をカウントアップ
				n[v]++
			}
			// もしなかったら、キー(日付)と初期値を加える
			n[v] = 1 // 初期値
		}
		out[key] = n // キーを日付として、値にnを代入
	}
	return out
}
