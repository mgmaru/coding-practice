package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Request struct {
	Ip     string `json:"ip"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Status string `json:"status"`
	Bytes  string `json:"bytes"`
}

type Requests []Request

type Commands struct {
	CommandIndex string
	LogFilePath  string
	By           string
	Top          int // 修正：flag.Intで指定したため、stringからintへ
	Json         bool
}

// 追加:件数を格納する構造体
type Count struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// コマンドを受け取り、解析する関数
func parseCommands() (*Commands, error) {

	// var CommandIndex string
	// var LogFilePath string
	// var By string
	// var Top string
	// var Json string

	by := flag.String("by", "status", "集計軸 status|path|ip")
	top := flag.Int("top", 0, "上位N件（0で全件）")
	asJSON := flag.Bool("json", false, "JSON出力")

	flag.Parse()

	args := flag.Args()

	// ログファイルが指定されていない
	if len(args) < 1 {
		return nil, errors.New("ログファイルを指定してください")
	}

	// byコマンドの引数が不正
	if *by != "status" && *by != "path" && *by != "ip" {
		return nil, errors.New("コマンドが不正です。")
	}
	// 正しいコマンドの場合、コマンドを返す
	commands := Commands{
		LogFilePath: args[0],
		By:          *by,
		Top:         *top,
		Json:        *asJSON,
	}
	return &commands, nil
}

// ファイルを開く
func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	// なんらかの理由でファイルが開けない場合はリターン
	if err != nil {
		defer file.Close()
		return nil, errors.New("ファイルが開けません。")
	}
	return file, nil
}

// ファイルを行単位で読み込んで、有効なリクエストを返す
func extractValidRequestLines(file *os.File, request Request) (Requests, error) {

	scanner := bufio.NewScanner(file)

	var validRequestLines Requests
	var skipped int

	for scanner.Scan() {

		requestline := scanner.Text()

		// 修正
		fields := strings.Fields(requestline)

		// 修正
		// filedsを利用してバリデーション
		// ただし、ここの"5"がマジックナンバーではないのか。。。
		if len(fields) != 5 {
			skipped++
			// 以下の処理を実行しない
			continue
		}

		// 有効なリクエストを構造体にマッピング
		request.Ip = fields[0]
		request.Method = fields[1]
		request.Path = fields[2]
		request.Status = fields[3]
		request.Bytes = fields[4]

		validRequestLines = append(validRequestLines, request)
	}

	if err := scanner.Err(); err != nil {
		// 本当はファイルクローズの責務はここではない（開いた側の責務）
		defer file.Close()
		return nil, err
	}

	defer file.Close()

	// 修正：不正な行を警告
	fmt.Fprintf(os.Stderr, "警告: 不正な行を %d 件スキップしました。\n", skipped)

	return validRequestLines, nil
}

// 追加:--byで指定されたコマンドの件数を集計して、キーとカウントペアで返す関数
func aggregate(reqs Requests, keyOf func(Request) string) []Count {
	m := make(map[string]int)
	for _, r := range reqs {
		m[keyOf(r)]++
	}
	// Countのスライス定義
	counts := make([]Count, 0, len(m))
	for k, v := range m {
		counts = append(counts, Count{Key: k, Count: v})
	}
	return counts
}

// 降順に並べてスライスを返す
// 修正：[]Countを降順に並べ替える
// counts: [{Key:"200", Count:15}, {Key:"401", Count:3}, ... ]
func sortByDescend(counts []Count) []Count {

	sort.Slice(counts, func(i, j int) bool {
		// 比較対象の件数が違えば降順でソート
		if counts[i].Count != counts[j].Count {
			return counts[i].Count > counts[j].Count
		}
		// 比較対象の件数が同数の場合、キーの昇順でソート
		return counts[i].Key < counts[j].Key
	})
	return counts
}

// Top N件のみ出力
func displayTopN(commandTop int, sortValidRequests []Count) ([]Count, error) {

	// 修正：上位N件を格納するスライス
	topSortValidRequests := make([]Count, 0, len(sortValidRequests))

	// バグ修正：topが0の場合、全件出力になる
	if commandTop == 0 {
		return sortValidRequests, nil
	}

	// 考察点：ここで、Topの数に対して、リクエストが足りない場合のバリデーションを行う必要がある
	for i := 0; i < commandTop && i < len(sortValidRequests); i++ {
		topSortValidRequests = append(topSortValidRequests, sortValidRequests[i])
	}
	return topSortValidRequests, nil
}

// jsonで出力する関数
func formatJson(topSortValidRequests []Count) ([]byte, error) {

	// リクエストをJSONに変換
	jsonRequests, err := json.Marshal(topSortValidRequests)

	if err != nil {
		return nil, errors.New("JSONに変換できませんでした。")
	}

	return jsonRequests, nil
}

func main() {

	// 初期化（インスタンス作成）
	var request Request

	// コマンド解析
	commands, err := parseCommands()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// ファイル開く
	file, err := openFile(commands.LogFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 有効なリクエストを抽出
	validRequests, err := extractValidRequestLines(file, request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 追加:有効なコマンドをaggregateに渡す
	keyFns := map[string]func(Request) string{
		"status": func(r Request) string { return r.Status },
		"path":   func(r Request) string { return r.Path },
		"ip":     func(r Request) string { return r.Ip },
	}
	counts := aggregate(validRequests, keyFns[commands.By])

	// ソート
	// sortValidRequests := sortByDescend(commands.By, validRequests)
	sortValidRequests := sortByDescend(counts)

	// 最初のN件
	topSortValidRequests, err := displayTopN(commands.Top, sortValidRequests)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// JSONに変換
	jsonRequests, err := formatJson(topSortValidRequests)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 出力
	fmt.Println(string(jsonRequests))
}
