package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

// 対象ディレクトリのファイルの基本情報（ファイル名とファイルパス）
type FileBaseInfo struct {
	FileName string
	FilePath string
}

// 移動前のファイルパスと移動するかの判定結果
// ext, date共通に使用するファイル情報
type IsFileMove struct {
	FileBaseInfo FileBaseInfo
	IsMove       bool
}

// ext用のファイル情報
type ExtFileMove struct {
	FileInfo  IsFileMove
	Extension string
}

// date用のファイル情報
type DateFileMove struct {
	FileInfo IsFileMove
	Date     string
}

// 1ファイルの移動前のパスと移動後のパス
type FilePathMove struct {
	BeforePath string
	AfterPath  string
}

// １ファイルに対する操作結果
type FlieOperationResult struct {
	FileName string // 操作対象のファイル
	Result   string // 対象ファイルの操作結果
}

// ファイルの操作の計画サマリー(計画サマリ)
type PlanSummary struct {
	AllFiles      int // 全てのファイル数
	CompleteFiles int // 完了する/した　ファイル数
	SkipFiles     int // スキップしたファイル数（衝突したファイル数）
	ErrorFiles    int // エラーファイル数（衝突以外の何らかの理由でエラー）
}

type ResultSummary struct {
	AllFiles      int // 全てのファイル数
	CompleteFiles int // 完了する/した　ファイル数
	SkipFiles     int // スキップしたファイル数
	ErrorFiles    int // エラーファイル数
}

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

// 対象ディレクトリのファイルパスを全て返す（ディレクトリかどうかを区別しない）
func listAllFilesPath(targetDir string) ([]FileBaseInfo, error) {
	fileBaseInfo := make([]FileBaseInfo, 0) // len設定は？

	err := filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {

		if !d.IsDir() {
			fileBaseInfo = append(fileBaseInfo, FileBaseInfo{FileName: d.Name(), FilePath: path})
		}
		return nil // ここのエラーは返す？
	})

	// 探索失敗
	if err != nil {
		return fileBaseInfo, errors.New("File Search Error")
	}
	return fileBaseInfo, nil
}

// ファイル移動の対象かどうかを判定してファイル名とパスと判定結果を返す。（--recursiveに対応）
func listAllFilesIsMove(targetDir string, recursive bool, fileBaseInfo []FileBaseInfo) []IsFileMove {

	isFilesPathAndMove := make([]IsFileMove, 0, len(fileBaseInfo))

	// recursive: false（サブディレクトリを除外）
	if !recursive {
		for _, file := range fileBaseInfo {
			path := strings.TrimPrefix(file.FilePath, targetDir+"/") //　パスから対象ディレクトリのパスを除外
			if strings.Contains(path, "/") {                         // 最初の文字が「/」だったら、ディレクトリと判断して除外。
				isFilesPathAndMove = append(isFilesPathAndMove, IsFileMove{FileBaseInfo: file, IsMove: false})
			}
			isFilesPathAndMove = append(isFilesPathAndMove, IsFileMove{FileBaseInfo: file, IsMove: true})

		}
		return isFilesPathAndMove
	}
	// recursive: true（サブディレクトリも対象）
	for _, file := range fileBaseInfo {
		isFilesPathAndMove = append(isFilesPathAndMove, IsFileMove{FileBaseInfo: file, IsMove: true})
	}
	return isFilesPathAndMove
}

// ext: ファイルの拡張子を抽出
func extractFilesExtension(isFilesdPathAndMove []IsFileMove) []ExtFileMove {

	extFilesMove := make([]ExtFileMove, 0, len(isFilesdPathAndMove))

	for _, file := range isFilesdPathAndMove {
		ex := filepath.Ext(file.FileBaseInfo.FilePath)
		extFilesMove = append(extFilesMove, ExtFileMove{FileInfo: file, Extension: ex})
	}
	return extFilesMove
}

// ext: ファイル名と拡張子から移動後のパス（計画）を作成する。
func makeExtPathPlan(targetDir string, extFilesMove []ExtFileMove) []FilePathMove {

	extFilePathMove := make([]FilePathMove, 0, len(extFilesMove))

	for _, file := range extFilesMove {
		if !file.FileInfo.IsMove { // 移動しない場合
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: file.FileInfo.FileBaseInfo.FilePath}) // パス変更なし
		}

		// 移動する場合
		if file.Extension == ".jpg" { // jpg
			planPath := targetDir + "/JPG" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
		}
		if file.Extension == ".pdf" { // pdf
			planPath := targetDir + "/PDF" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
		}
		if file.Extension == ".png" { //png
			planPath := targetDir + "/PNG" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})

		}
		if file.Extension == "" { // no extension
			planPath := targetDir + "/" + file.FileInfo.FileBaseInfo.FileName + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})

		}
		// その他の拡張子があれば追加していく
	}

	return extFilePathMove
}

// ext:現在のパスと移動後のパスを比較して、衝突と冪等を判定する関数
func compareCurrentPathAndPlannedPath(currentPath []string, plannedPath []string) {
	// 全く同じ場合
}

// コマンドとパスを受け取って、ファイル操作を計画してサマリを返す関数（dry-run）
func summaryFileOperationPlan(byCommand string) PlanSummary {

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
	fmt.Printf("%d件のファイルをサブフォルダに移動します。", planSummary.CompleteFiles)
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

	// dry-run
	if !commands.Apply {

		// Plan
		filOperationSummary := summaryFileOperationPlan(commands.By)

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

		allFilesPath, err := listAllFilesPath(commands.Path)
		fmt.Println(makeExtPathPlan(commands.Path, extractFilesExtension(listAllFilesIsMove(commands.Path, commands.Recursive, allFilesPath))))

		return

		// apply（not dry-run）
	} else {

		// Aplly
		resultSummary := applyPlanAndaggregateResultSummary(commands.By)
		displayFileOperationResult(resultSummary)

		return
	}
}
