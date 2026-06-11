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

// コマンドを受け取り、解析する関数
func parseCommand() ([]string, error) {

	var commandIndex string //コマンド
	var logFilePath string  //ログファイル
	var by string           //集計軸
	var top string          //上位N件
	var json string         //JSON形式

	// コマンド入力受付
	if _, err := fmt.Scanf("%s %s --by %s --top %s %s", &commandIndex, &logFilePath, &by, &top, &json); err != nil {
		// 引数の不足を検出した場合
		return nil, errors.New("引数が不足しています")
	}
	// コマンドが不正
	if commandIndex != "logstat" {
		return nil, errors.New("コマンドが不正です。")
	}
	// コマンドが不正
	if by != "status" && by != "path" && by != "ip" {
		return nil, errors.New("コマンドが不正です。")
	}
	// コマンド不正（topが数字ではない）
	if _, err := strconv.Atoi(top); err != nil {
		return nil, errors.New("コマンドが不正です。")
	}
	// コマンド不正
	if json != "--json" {
		return nil, errors.New("コマンドが不正です。")
	}
	// ログファイルが見つからない場合
	if !fileExits(logFilePath) {
		return nil, errors.New("ファイルが存在しません")
	}
	// 正しいコマンドの場合、コマンドを返す
	return []string{commandIndex, logFilePath, by, top, json}, nil
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

// ファイルを行単位で読み込んで、有効な行のみを返す関数
func extractValidRequestLines(file *os.File) ([][]string, error) {

	scanner := bufio.NewScanner(file)

	var validRequestLines [][]string

	for scanner.Scan() {

		requestline := scanner.Text()

		// 有効な行かを判断（フィールド数が合わない場合はスキップ）
		if isValidLine(requestline) {
			validRequestLine := strings.Fields(requestline)
			validRequestLines = append(validRequestLines, validRequestLine)
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
func sortByDescend(commandBy string, validRequests [][]string) [][]string {

	m := make(map[string]int)

	// statusのソート
	if commandBy == "status" {
		//statusのソート
		for i := 0; i < len(validRequests); i++ {
			// 初期カウント
			count := 1
			// i番目のスライスの値を抽出
			status := validRequests[i][3]

			// キーがすでに存在する場合
			if _, ok := m[status]; ok {
				m[status] += 1
			} else {
				// キーが存在しない場合
				m[status] = count
			}
		}

		// キーで配列をソートする
		sort.Slice(validRequests, func(i, j int) bool { return m[validRequests[i][3]] > m[validRequests[j][3]] })
	}

	// pathのソート
	if commandBy == "path" {
		for i := 0; i < len(validRequests); i++ {
			// 初期カウント
			count := 1
			// i番目のスライスの値を抽出
			path := validRequests[i][2]

			// キーがすでに存在する場合
			if _, ok := m[path]; ok {
				m[path] += 1
			} else {
				// キーが存在しない場合
				m[path] = count
			}
		}
		sort.Slice(validRequests, func(i, j int) bool { return m[validRequests[i][2]] > m[validRequests[j][2]] })
	}

	// ipのソート
	if commandBy == "ip" {
		for i := 0; i < len(validRequests); i++ {
			// 初期カウント
			count := 1
			// i番目のスライスの値を抽出
			ip := validRequests[i][0]

			// キーがすでに存在する場合
			if _, ok := m[ip]; ok {
				m[ip] += 1
			} else {
				// キーが存在しない場合
				m[ip] = count
			}
		}
		sort.Slice(validRequests, func(i, j int) bool { return m[validRequests[i][0]] > m[validRequests[j][0]] })
	}
	return validRequests
}

// Top N件のみ出力
func displayTopN(commandTop string, sortValidRequests [][]string) [][]string {

	commandTopInt, _ := strconv.Atoi(commandTop)

	var topSortValidRequests [][]string

	if commandTopInt > len(sortValidRequests) {
		for i := 0; i < len(sortValidRequests); i++ {
			topSortValidRequests = append(topSortValidRequests, sortValidRequests[i])
		}
		return topSortValidRequests
	}

	for i := 0; i < commandTopInt; i++ {
		topSortValidRequests = append(topSortValidRequests, sortValidRequests[i])

	}
	return topSortValidRequests
}

// jsonで出力する関数
func formatJson(topSortValidRequests [][]string) ([]byte, error) {

	var requests Requests

	for i := 0; i < len(topSortValidRequests); i++ {
		request := Request{
			Ip:     topSortValidRequests[i][0],
			Method: topSortValidRequests[i][1],
			Path:   topSortValidRequests[i][2],
			Status: topSortValidRequests[i][3],
			Bytes:  topSortValidRequests[i][4],
		}
		requests = append(requests, request)
	}
	// リクエストをJSONに変換
	jsonRequests, err := json.Marshal(requests)

	if err != nil {
		return nil, errors.New("Jsonに変換できませんでした。")
	}

	return jsonRequests, nil
}

func main() {

	// コマンド解析
	commands, err := parseCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logFile := commands[1]

	// ファイル開く
	file, err := openFile(logFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	validRequests, err := extractValidRequestLines(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// ソート
	commandBy := commands[2]
	sortValidRequests := sortByDescend(commandBy, validRequests)

	// 最初のN件
	commandTop := commands[3]
	topSortValidRequests := displayTopN(commandTop, sortValidRequests)

	// JSONに変換
	jsonRequests, err := formatJson(topSortValidRequests)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 出力
	fmt.Println(string(jsonRequests))
}
