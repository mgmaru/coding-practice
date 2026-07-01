package collections

type Status struct { // ログ集計結果
	Date  string
	Level string
}

// ロジック設計
// 1. 入力から日付ごとのレベルをスライスにまとめて抽出する
// 2. 日付ごとにレベルのカウントをする
// 3. map of mapで出力（キー：日付 値：mapレベル:カウント）
func aggregateDateStatus(input []Status) map[string]map[string]int {

	m := make(map[string][]string, len(input)) // キー：日付 値：レベルのスライス
	for _, v := range input {
		m[v.Date] = append(m[v.Date], v.Level) //一旦、日付ごとのレベルを格納する // {06-01:[INFO, ERROR, INFO, WARNING, ...], 06-02:[INFO, WARNIG, ...]}
	}

	// 日付ごとのレベルの出現回数をカウントする
	n := make(map[string]int) // キー：レベル　値：出現回数
	out := make(map[string]map[string]int)
	for key, vals := range m { // key:日付　vals:レベルのスライス // 注意：map mからforで取り出しているので、並び順が保証されていない。
		for _, v := range vals {
			n[v]++
		}
		out[key] = n
	}
	return out
}
