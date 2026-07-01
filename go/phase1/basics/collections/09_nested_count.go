package collections

type logEntry struct { // ステータス1件
	date  string // 空文字は入りうる（nilは入らない）
	level string // 空文字は入りうる（nilは入らない）
}

func aggregateDateStatus(input []logEntry) map[string]map[string]int {

	// inputが空で入ってきても、全て空スライスで返るのでガード不要

	// 修正：不要な中間mapであるmを消去し、一本化
	out := make(map[string]map[string]int, len(input))
	for _, s := range input {
		if out[s.date] == nil {
			out[s.date] = make(map[string]int)
		}
		out[s.date][s.level]++
	}
	return out
}
