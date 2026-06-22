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
	// isMoveもつべき？
}

// １ファイルに対する操作結果
type FileOperationResult struct {
	CurrentFilePath string // 現在のファイルパス
	PlannedPath     string // 計画後のファイルパス
	Result          string // 対象ファイルの操作結果 // complete（完了）| skip（衝突）| error（エラー）| unchange（変更なし）
}

// 計画サマリ
type PlanSummary struct {
	AllFiles      int // 全てのファイル数
	CompleteFiles int // 完了する/した　ファイル数
	SkipFiles     int // スキップしたファイル数（衝突したファイル数）
	ErrorFiles    int // エラーファイル数（衝突以外の何らかの理由でエラー）
	UnchangeFiles int // 変更なしファイル
}

// 実行結果サマリ
type ResultSummary struct {
	AllFiles      int // 全てのファイル数
	CompleteFiles int // 完了する/した　ファイル数
	SkipFiles     int // スキップしたファイル数
	ErrorFiles    int // エラーファイル数
	UnChangeFiles int // 変更なしファイル
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
	if _, err := isExistingDir(inputPath); err != nil {
		return nil, err
	}

	// 追加：空ディレクトリか
	empty, err := isEmptyDir(inputPath)
	if err != nil {
		return nil, err
	}
	if empty {
		return nil, errors.New("空ディレクトリです。")
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
func isExistingDir(path string) (bool, error) {

	pathInfo, err := os.Stat(path) // ファイルもしくはディレクトリが存在するか

	if err != nil { // ファイルまたはディレクトリが存在しない
		return false, errors.New("存在しないパスです。")
	}
	if !pathInfo.IsDir() { // パスは存在するが、ディレクトリではない
		return false, errors.New("ディレクトリを指定してください。")
	}
	return true, nil
}

// ディレクトリが空か判定
func isEmptyDir(path string) (bool, error) {
	empty, err := os.ReadDir(path)
	if err != nil { // Read Error
		return false, err
	}

	if len(empty) == 0 { // 空の場合
		return true, nil
	}
	return false, nil
}

// 対象ディレクトリのファイルパスを全て返す（ディレクトリかどうかを区別しない）
func listFilesPath(targetDir string) ([]FileBaseInfo, error) {
	fileBaseInfo := make([]FileBaseInfo, 0) // len設定はしない。理由：対象ディレクトリのファイル数がわからないため。

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
func listFilesIsMove(targetDir string, recursive bool, fileBaseInfo []FileBaseInfo) []IsFileMove {

	isFilesPathAndMove := make([]IsFileMove, 0, len(fileBaseInfo))

	// recursive: false（サブディレクトリを除外）
	if !recursive {
		for _, file := range fileBaseInfo {
			path := strings.TrimPrefix(file.FilePath, targetDir+"/") //　パスから対象ディレクトリのパスを除外
			if strings.Contains(path, "/") {                         // 最初の文字が「/」だったら、ディレクトリと判断して除外。
				isFilesPathAndMove = append(isFilesPathAndMove, IsFileMove{FileBaseInfo: file, IsMove: false})
				continue
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

// ext: ファイル名と拡張子から移動後のパスを返す関数
func makeExtPathPlan(targetDir string, extFilesMove []ExtFileMove) []FilePathMove {

	extFilePathMove := make([]FilePathMove, 0, len(extFilesMove))

	for _, file := range extFilesMove {
		if !file.FileInfo.IsMove { // 移動しない場合 unchange
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: file.FileInfo.FileBaseInfo.FilePath}) // パス変更なし
			continue                                                                                                                                                 // 忘れ注意
		}
		// 移動する場合
		if file.Extension == ".jpg" {
			planPath := targetDir + "/JPG" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
		}
		if file.Extension == ".pdf" {
			planPath := targetDir + "/PDF" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
		}
		if file.Extension == ".png" {
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

// 計画ロジック（重要）
// 判定層１：冪等かどうかを判定
func isIndempotent(filePathMove []FilePathMove) bool {

	var pathMatchCount int // 移動前と移動後でマッチしたパスの回数をカウント

	if len(filePathMove) == 0 { // 空スライス、nilを判定
		return false
	}

	for _, file := range filePathMove {
		if file.BeforePath == file.AfterPath {
			pathMatchCount++
		}
	}

	if pathMatchCount == len(filePathMove) { // もし、[]FilePathMoveが空スライスだった場合にもtrueになってしまう...
		return true
	}
	return false
}

// 判定層２：衝突を判定
// 移動前パスおよび移動後パスで衝突検知。工夫：recursiveに意識しないような実装にした。
// recursiveを意識すると、recursiveの値で分岐することになりそうだから、依存を無くそうとした。
func hasPathConflict(filePathMove []FilePathMove) []FileOperationResult {

	notChangePathes := make([]FilePathMove, 0, len(filePathMove)) // 変更しないパスグループ
	changePathes := make([]FilePathMove, 0, len(filePathMove))    // 変更予定のパスグループ

	filesOperationResult := make([]FileOperationResult, 0, len(filePathMove)) // 予定が決定したパスを格納するスライス

	for _, file := range filePathMove { //グループ化
		if file.BeforePath == file.AfterPath {
			notChangePathes = append(notChangePathes, file)
		} else {
			changePathes = append(changePathes, file)
		}
	}

	// notChangeを構造体にマッピング（変更しないものは計画がすでに決まっているので、先に格納）
	for _, notChangePath := range notChangePathes {
		filesOperationResult = append(filesOperationResult, FileOperationResult{CurrentFilePath: notChangePath.BeforePath, PlannedPath: notChangePath.AfterPath, Result: "unchange"})
	}

	// 変更しないパスと、変更予定のパスの衝突を判定
	for _, notChangePath := range notChangePathes { // notChangePathは一意。 -> changePathが2回appendされることはない。
		for _, changePath := range changePathes {
			if notChangePath.BeforePath == changePath.AfterPath { // 同じもののグループの現在のパスと、違うものの変更予定のパスが同じ
				filesOperationResult = append(filesOperationResult, FileOperationResult{CurrentFilePath: changePath.BeforePath, PlannedPath: changePath.BeforePath, Result: "skip"})
				break // あるchangePathとnotChangePathが被った時点で、そのchangePathを調べる必要はない。
			}
			// かぶらなければ何もしない（理由：まだ計画が決まっていないため）
		}
	}

	// スライスchangePathからすでにappendされているパスを削除する。
	// すでに計画が決定した現在のパスをの存在を格納する。
	m := make(map[string]struct{}, len(changePathes))

	for _, v := range filesOperationResult { // 計画が決定した
		m[v.CurrentFilePath] = struct{}{}
	}

	// changePathのパスを１つずつ見ていって、mにすでに存在したら、新しいスライスに入れる
	remainingChangePathes := make([]FilePathMove, 0, len(changePathes))

	for _, path := range changePathes {
		if _, ok := m[path.BeforePath]; !ok { // mにキーが存在しなかったら、まだ計画が決まっていないパス -> スライスに格納
			remainingChangePathes = append(remainingChangePathes, path)
		}
	}

	// 計画が決まっていないパス（remainingChangePath）どうしの衝突を判定する。
	// 注意：重複appendは許さないようにする。

	n := make(map[string]struct{}, len(remainingChangePathes)) // remainingChangePathesで計画が決まったもののBeforePathを格納

	for _, i := range remainingChangePathes {
		for _, j := range remainingChangePathes {
			if i == j { // 同じパスの比較はスキップ
				continue
			}
			if i.AfterPath == j.AfterPath { //衝突 // 移動予定のパスで比較
				if _, ok := n[j.BeforePath]; ok { // すでに結果が決まっているかを判定。 -> 決まっていたらappendをスキップ
					continue
				}
				// まだ決まっていない場合 // パスをスライス`filesOperationResult`にappend
				filesOperationResult = append(filesOperationResult, FileOperationResult{CurrentFilePath: j.BeforePath, PlannedPath: j.BeforePath, Result: "skip"})
				// 計画が決まったパスをnに格納
				n[j.BeforePath] = struct{}{}
			}
		}
	}

	// 衝突しなかったパス（計画が決まっていないパス）を抽出 -> complete
	for _, path := range remainingChangePathes {
		if _, ok := n[path.BeforePath]; !ok { // nにパスがない -> まだ計画が決まっていない。(衝突しなかったパス)
			filesOperationResult = append(filesOperationResult, FileOperationResult{CurrentFilePath: path.BeforePath, PlannedPath: path.AfterPath, Result: "complete"})
		}
	}

	fmt.Println(len(filesOperationResult))

	return filesOperationResult
}

// 計画結果を受け取って、サマリを返す関数（dry-run）
func summaryFileOperationPlan(fileOperationResult []FileOperationResult) PlanSummary {

	var completeFilesCount int
	var skipFilesCount int
	var errorFilesCount int
	var unchangeFilesCount int

	allFilesCount := len(fileOperationResult)

	for _, path := range fileOperationResult {
		if path.Result == "complete" {
			completeFilesCount++
		}
		if path.Result == "skip" {
			skipFilesCount++
		}
		if path.Result == "error" {
			errorFilesCount++
		}
		if path.Result == "unchange" { //unchange
			unchangeFilesCount++
		}
	}

	return PlanSummary{AllFiles: allFilesCount, CompleteFiles: completeFilesCount, SkipFiles: skipFilesCount, ErrorFiles: errorFilesCount, UnchangeFiles: unchangeFilesCount}
}

// ファイル操作のサマリを受け取って計画を表示する関数
func displayFileOperationPlan(planSummary PlanSummary) {

	if planSummary.CompleteFiles > 0 { // completeが１件以上で表示
		fmt.Printf("%d件のファイルをサブフォルダに移動します。", planSummary.CompleteFiles)
	}
	if planSummary.SkipFiles > 0 { // skipが1件以上で表示
		fmt.Printf("%d件のファイル移動をスキップします。", planSummary.SkipFiles)
	}
	if planSummary.ErrorFiles > 0 { // errorが1件以上で表示
		fmt.Printf("%d件のファイルが移動できません。", planSummary.SkipFiles)
	}
	if planSummary.UnchangeFiles > 0 { // unchangeが1件以上で表示
		fmt.Printf("%d件のファイルの変更はありません。", planSummary.UnchangeFiles)
	}
}

// 計画を実行するかしないかを表示して結果を返す関数
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

	// コマンドパース
	commands, err := parseCommnads() // コマンドとパスを取得
	if err != nil {                  // エラーの場合、コンソールに出す。
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	allFilesPath, err := listFilesPath(commands.Path)
	filesIsMove := listFilesIsMove(commands.Path, commands.Recursive, allFilesPath)

	// date

	// seq

	if !commands.Apply { // dry-run

		// コマンドごとに計画を立てる
		if commands.By == "ext" { // ext
			extFileMove := extractFilesExtension(filesIsMove)
			filePathMove := makeExtPathPlan(commands.Path, extFileMove)
			planPath := hasPathConflict(filePathMove)
			filOperationSummary := summaryFileOperationPlan(planPath)
			displayFileOperationPlan(filOperationSummary) // サマリーを表示
		}
		if commands.By == "date" { // date

		}
		if commands.By == "seq" { // seq

		}

		isApply, err := isApplyPlan() // ユーザ入力
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
		fmt.Println("実行します。")

		return

	} else { // apply（not dry-run）

		// Aplly
		resultSummary := applyPlanAndaggregateResultSummary(commands.By)
		displayFileOperationResult(resultSummary)

		return
	}
}
