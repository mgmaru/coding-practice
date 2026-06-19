package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// コマンド
type Commands struct {
	Path      string
	By        string
	Prefix    string
	Apply     bool
	Recursive bool
	Conflict  string
}

// １ファイルに対する操作結果
type FlieOperationResult struct {
	FileName string // 操作対象のファイル
	Result   string // 対象ファイルの操作結果
}

// ファイルの操作の計画サマリー(計画サマリ)
type PlanSummary struct {
	CompleteFiles int // 完了する/した　ファイル数
	SkipFiles     int // スキップしたファイル数
	ErrorFiles    int // エラーファイル数
}

type ResultSummary struct {
	CompleteFiles int // 完了する/した　ファイル数
	SkipFiles     int // スキップしたファイル数
	ErrorFiles    int // エラーファイル数
}

// ファイル操作のサマリー（結果）

// コマンドをパースする関数：入力したコマンド形式が正しいかどうかを判定する。
// コマンドが正しければ、コマンドとパスを返す。
func parseCommnads() (*Commands, error) {

	// コマンド定義
	by := flag.String("by", "ext", "file operations ext|date|seq")
	prefix := flag.String("prefix", "", "prefix when seq is selected")         // 選択：prefixに何も入力されなかった時の値を考える。現状は空文字。
	apply := flag.Bool("apply", false, "plan (dry-run) or apply immediately?") // dry-run: false / apply: true
	recursive := flag.Bool("recursive", false, "whether to include subdirectories")
	conflict := flag.String("conflict", "error", "behavior during collision")

	flag.Parse()

	args := flag.Args() // 対象ディレクトリ（絶対パスも相対パスも入ってくる）

	// ディレクトリが指定されていない
	if len(args) < 1 {
		return nil, errors.New("ディレクトリを指定してください。")
	}

	inputPath := args[0]

	// 入力されたパスがファイルかディレクトリか判定
	// ディレクトリだった場合、存在するかどうか
	if _, err := isDirectoryOrFileAndExist(inputPath); err != nil {
		return nil, err
	}

	// --by seqの時に、--prefixが指定されているか
	if *by == "seq" && *prefix == "" {
		return nil, errors.New("prefixを指定してください。")
	}

	return &Commands{
		Path:      args[0],
		By:        *by,
		Prefix:    *prefix,
		Apply:     *apply,
		Recursive: *recursive,
		Conflict:  *conflict,
	}, nil
}

// ディレクトリかファイルかを判定して、ディレクトリだった場合、存在を判定して、存在する場合はtrueを返す。
// そのほかの場合はfalseを返す。
func isDirectoryOrFileAndExist(inputPath string) (bool, error) {

	pathInfo, err := os.Stat(inputPath) // ファイルもしくはディレクトリが存在するか

	if err != nil { // ファイルまたはディレクトリが存在しない
		return false, errors.New("存在しないパスです。")
	}
	if !pathInfo.IsDir() { // パスは存在するが、ディレクトリではない
		return false, errors.New("ディレクトリを指定してください。")
	}
	return true, nil
}

// コマンドを受け取って、ファイル操作を計画してサマリを返す関数（dry-run）
func aggregateFileOperations(byCommand string) PlanSummary {

	if byCommand == "ext" { // extの計画を立てる
		return PlanSummary{}

	}
	if byCommand == "data" { // dataの計画を立てる
		return PlanSummary{}

	} else { // seqの計画を立てる
		return PlanSummary{}
	}
}

// ファイル操作のサマリを受け取って、計画を表示する関数
func displayFileOperationPlan(planSummary PlanSummary) {

}

// 計画を実行するかしないかを返す関数
func isApplyPlan() (bool, error) {

	var carryOut string

	fmt.Print("計画を実行しますか？ Y/n：")
	_, err := fmt.Scan(&carryOut)
	if err != nil {
		// scan error
		return false, err // 実行しない
	}

	if carryOut == "Y" {
		return true, nil //実行する
	}
	return false, nil // 実行しない
}

// コマンドを受け取って、計画を実行する関数
// 設計意思：ーーby ext|date|seqを１本化
func applyPlanAndaggregateResultSummary(byCommand string) ResultSummary {

	if byCommand == "ext" { // ext
		return ResultSummary{}
	}
	if byCommand == "date" { // date
		return ResultSummary{}
	}
	// seq
	return ResultSummary{}
}

// ファイル操作のサマリを受け取って、操作結果を表示する関数
func displayFileOperationResult(resultSummary ResultSummary) {

}

func main() {

	commands, err := parseCommnads() // コマンドとパスを取得
	if err != nil {                  // エラーの場合、コンソールに出す。
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !commands.Apply { // dry-run

		// Plan
		filOperationSummary := aggregateFileOperations(commands.By)

		displayFileOperationPlan(filOperationSummary) // サマリーを表示
		isApply, err := isApplyPlan()                 // ユーザ入力
		if err != nil {
			// 実行しない
			fmt.Println("scan error")
			return
		}
		if !isApply {
			// 実行しない
			fmt.Println("実行しません。")
			return
		}

		// Apply
		resultSummary := applyPlanAndaggregateResultSummary(commands.By)
		displayFileOperationResult(resultSummary)

		return

	} else { // apply（not dry-run）

		// Aplly
		resultSummary := applyPlanAndaggregateResultSummary(commands.By)
		displayFileOperationResult(resultSummary)

		return
	}
}
