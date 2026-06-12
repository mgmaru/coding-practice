package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
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
	Top          string
	Json         string
}

// コマンドを受け取り、解析する関数
func parseCommands() (*Commands, error) {

	var CommandIndex string
	var LogFilePath string
	var By string
	var Top string
	var Json string

	// コマンド入力受付
	if _, err := fmt.Scanf("%s %s --by %s --top %s %s", &CommandIndex, &LogFilePath, &By, &Top, &Json); err != nil {
		// 引数の不足を検出した場合
		return nil, errors.New("引数が不足しています")
	}
	// コマンドが不正
	if CommandIndex != "logstat" {
		return nil, errors.New("コマンドが不正です。")
	}
	// コマンドが不正
	if By != "status" && By != "path" && By != "ip" {
		return nil, errors.New("コマンドが不正です。")
	}
	// コマンド不正（topが数字ではない）
	if _, err := strconv.Atoi(Top); err != nil {
		return nil, errors.New("コマンドが不正です。")
	}
	// コマンド不正
	if Json != "--json" {
		return nil, errors.New("コマンドが不正です。")
	}
	// ログファイルが見つからない場合
	if !fileExits(LogFilePath) {
		return nil, errors.New("ファイルが存在しません")
	}
	// 正しいコマンドの場合、コマンドを返す
	commands := Commands{
		CommandIndex: CommandIndex,
		LogFilePath:  LogFilePath,
		By:           By,
		Top:          Top,
		Json:         Json,
	}
	return &commands, nil
}

// ファイルが存在するか判定する関数
func fileExits(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return false
	}
	return true
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

	for scanner.Scan() {

		requestline := scanner.Text()

		// 有効な行かを判断（フィールド数が合わない場合はスキップ）
		if isValidLine(requestline) {

			// 有効なリクエストを構造体にマッピング
			request.Ip = strings.Fields(requestline)[0]
			request.Method = strings.Fields(requestline)[1]
			request.Path = strings.Fields(requestline)[2]
			request.Status = strings.Fields(requestline)[3]
			request.Bytes = strings.Fields(requestline)[4]

			validRequestLines = append(validRequestLines, request)
		}
	}

	if err := scanner.Err(); err != nil {
		defer file.Close()
		return nil, err
	}

	defer file.Close()
	return validRequestLines, nil
}

// 有効な行かを判断する関数
func isValidLine(line string) bool {
	const VALIDREQUESTLENGTH = 5
	requestLength := strings.Fields(line)
	if len(requestLength) == VALIDREQUESTLENGTH {
		return true
	}
	return false
}

// 降順に並べてスライスを返す
func sortByDescend(commandBy string, validRequests Requests) Requests {

	m := make(map[string]int)

	// statusのソート
	if commandBy == "status" {
		//statusのソート
		for i := 0; i < len(validRequests); i++ {
			m[validRequests[i].Status]++
		}

		// キーで配列をソートする
		sort.Slice(validRequests, func(i, j int) bool { return m[validRequests[i].Status] > m[validRequests[j].Status] })
	}

	// pathのソート
	if commandBy == "path" {
		for i := 0; i < len(validRequests); i++ {
			m[validRequests[i].Path]++
		}
		sort.Slice(validRequests, func(i, j int) bool { return m[validRequests[i].Path] > m[validRequests[j].Path] })
	}

	// ipのソート
	if commandBy == "ip" {
		for i := 0; i < len(validRequests); i++ {
			m[validRequests[i].Ip]++
		}
		sort.Slice(validRequests, func(i, j int) bool { return m[validRequests[i].Ip] > m[validRequests[j].Ip] })
	}
	return validRequests
}

// Top N件のみ出力
func displayTopN(commandTop string, sortValidRequests Requests) (Requests, error) {

	commandTopInt, err := strconv.Atoi(commandTop)

	if err != nil {
		return nil, errors.New("コマンドが不正です。")
	}

	var topSortValidRequests Requests

	if commandTopInt > len(sortValidRequests) {
		for i := 0; i < len(sortValidRequests); i++ {
			topSortValidRequests = append(topSortValidRequests, sortValidRequests[i])
		}
		return topSortValidRequests, nil
	}

	for i := 0; i < commandTopInt; i++ {
		topSortValidRequests = append(topSortValidRequests, sortValidRequests[i])

	}
	return topSortValidRequests, nil
}

// jsonで出力する関数
func formatJson(topSortValidRequests Requests) ([]byte, error) {

	// リクエストをJSONに変換
	jsonRequests, err := json.Marshal(topSortValidRequests)

	if err != nil {
		return nil, errors.New("Jsonに変換できませんでした。")
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
	// ソート
	sortValidRequests := sortByDescend(commands.By, validRequests)

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
