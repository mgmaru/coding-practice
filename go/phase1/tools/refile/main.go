package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

// 移動前のファイルパスと移動（リネーム）するかの判定結果
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

// seq用のファイル情報
type SeqFileMove struct {
	Extension string
	Files     []IsFileMove
}

// 1ファイルの移動前のパスと移動後のパス // 衝突の判定に使用する
type FilePathMove struct {
	BeforePath string
	AfterPath  string
	// isMoveもつべき？ -> いらない。理由：移動前のパスと移動後のパスで衝突は判定できる。
}

// ファイルパスの衝突結果
type FilePathHasConflict struct {
	FilePathMove FilePathMove // 移動前のパスと移動後のパス
	HasConflict  bool         // 衝突判定結果
}

// 移動計画サマリ
type PlanSummary struct {
	CompleteFiles     int // 移動できるファイル数
	SkipFiles         int // 衝突してスキップしたファイル数
	ErrorFiles        int // エラーファイル数（衝突によって実行できなかったファイル数）
	UnchangeFiles     int // 変更なしファイル数
	CannotChangeFiles int // エラーによって変更できないファイル数
}

// 最終的な計画サマリ
type PlanedFilePath struct {
	FilePath FilePathMove
	Status   FileStatus
}

// 実行後のサマリ
type ApplyedFilePath struct {
	FilePath FilePathMove
	Status   FileStatus
}

// リネーム計画サマリ
type RenamePlanSummary struct {
	CompleteFiles     int // 完了する/した　ファイル数
	SkipFiles         int // スキップしたファイル数（衝突したファイル数）
	ErrorFiles        int // エラーファイル数（衝突以外の何らかの理由でエラー）
	UnchangeFiles     int // 変更なしファイル
	CannotChangeFiles int // エラーによって変更できないファイル数
}

// 実行結果サマリ
type ResultSummary struct {
	CompleteFiles     int // 完了する/した　ファイル数
	SkipFiles         int // スキップしたファイル数
	ErrorFiles        int // エラーファイル数
	UnChangeFiles     int // 変更なしファイル
	CannotChangeFiles int // エラーによって変更できないファイル数
}

type FileStatus string

// 計画完了後、適用後のファイルステータス
const (
	CompleteFileStatus FileStatus = "complete"
	SkipFileStatus     FileStatus = "skip"
	ErrorFileStatus    FileStatus = "error"
	UnchangeFileStatus FileStatus = "unchange"
	CannotFileStatus   FileStatus = "cannot"
)

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

	args := flag.Args()

	// ディレクトリが指定されていない場合
	if len(args) < 1 {
		return nil, errors.New("ディレクトリを指定してください。")
	}

	inputPath := args[0]

	// 入力されたパスがファイルかディレクトリか判定
	// ディレクトリだった場合、存在するかどうか
	if _, err := isExistingDir(inputPath); err != nil {
		return nil, err
	}

	// 空ディレクトリか
	isEmpty, err := isEmptyDir(inputPath)
	if err != nil {
		return nil, err
	}
	if isEmpty {
		return nil, errors.New("空ディレクトリです。")
	}

	// --by seqの時に、--prefixが指定されているか
	if *by == "seq" && *prefix == "" {
		return nil, errors.New("prefixを指定してください。")
	}

	// コマンドを返却
	return &Commands{
		Path:      inputPath,
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

// 共通：対象ディレクトリのファイルパスを全て返す（ディレクトリかどうかを区別しない）。
// 設計意思：共通で仕様する関数なので、どのコマンドであるかを意識しない。
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

// 共通：ファイル移動の対象かどうかを判定してファイル名とパスと判定結果を返す（--recursiveに対応）
// 設計意思：共通で仕様する関数なので、どのコマンドであるかを意識しない。
// 注意：seqのみ、recursiveは全てfalse
func listFilesIsMove(targetDir string, recursive bool, fileBaseInfo []FileBaseInfo) []IsFileMove {

	isFilesPathAndMove := make([]IsFileMove, 0, len(fileBaseInfo))

	// recursive: false（サブディレクトリを除外）
	if !recursive {
		for _, file := range fileBaseInfo {
			path := strings.TrimPrefix(file.FilePath, targetDir+"/") //　パスから対象ディレクトリのパスを除外
			if strings.Contains(path, "/") {                         // 残りのパスに「/」が含まれていたらディレクトリと判断（ファイルではない）。
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

// ext / seq : ファイルの拡張子を抽出
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
			planPath := targetDir + "/jpg" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
		}
		if file.Extension == ".pdf" {
			planPath := targetDir + "/pdf" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
		}
		if file.Extension == ".png" {
			planPath := targetDir + "/png" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
		}
		if file.Extension == "" { // no extension
			planPath := targetDir + "/" + "noext" + "/" + file.FileInfo.FileBaseInfo.FileName
			extFilePathMove = append(extFilePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
		}
		// その他の拡張子があれば追加していく
	}
	return extFilePathMove
}

// date:対象ファイルからファイルの更新月を抽出する。
func extractFilesDate(isFilesdPathAndMove []IsFileMove) ([]DateFileMove, error) {

	dateFilesInfo := make([]DateFileMove, 0, len(isFilesdPathAndMove))

	for _, file := range isFilesdPathAndMove {
		fileInfo, err := os.Stat(file.FileBaseInfo.FilePath)
		if err != nil {
			return nil, err
		}
		dateFilesInfo = append(dateFilesInfo, DateFileMove{FileInfo: file, Date: fileInfo.ModTime().Format("2006-01-02")})
	}
	return dateFilesInfo, nil
}

// date: ファイル情報および拡張子から移動計画を返す関数
func makeDatePathPlan(targetDir string, dateFilesMove []DateFileMove) []FilePathMove {

	filePathMove := make([]FilePathMove, 0, len(dateFilesMove))

	for _, file := range dateFilesMove {
		if !file.FileInfo.IsMove { // 移動しない場合
			filePathMove = append(filePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: file.FileInfo.FileBaseInfo.FilePath})
			continue
		}
		// 移動する場合 パスを作る
		date := file.Date
		planPath := targetDir + "/" + string(date[:7]) + "/" + file.FileInfo.FileBaseInfo.FileName // Date: 2026-06形式
		filePathMove = append(filePathMove, FilePathMove{BeforePath: file.FileInfo.FileBaseInfo.FilePath, AfterPath: planPath})
	}
	return filePathMove
}

// seq仕様：seqは--recursiveコマンドは効かない（サブディレクトリは対象にしない。）
// seq: ファイル情報と拡張子から拡張子ごとに分類して、mapで返す（ここでソートはしない）。
func classifyExt(extFilesMove []ExtFileMove) []SeqFileMove {

	// key： Extension value: files
	m := make(map[string][]IsFileMove, len(extFilesMove))
	seqFilesMove := make([]SeqFileMove, 0, len(extFilesMove))

	for _, file := range extFilesMove { // 拡張子ごとにファイルを分類
		m[file.Extension] = append(m[file.Extension], file.FileInfo)
	}
	// 構造体SeqFileMoveにマッピング
	for ex, files := range m {
		seqFilesMove = append(seqFilesMove, SeqFileMove{Extension: ex, Files: files})
	}

	return seqFilesMove
}

// seq: パスを作成する。（1.拡張子を昇順でソート -> 2.拡張子の中でファイル名でソート -> 3.パスを作る）
func makeSeqPathPlan(targetDir string, seqFilesMove []SeqFileMove, prefix string) []FilePathMove {

	filesPathMove := make([]FilePathMove, 0, len(seqFilesMove))

	// 1. キー：拡張子について、昇順ソート
	sort.Slice(seqFilesMove, func(i, j int) bool {
		return seqFilesMove[i].Extension > seqFilesMove[j].Extension
	})

	// 2. 拡張子ごとのファイルのリストの中のファイル名でソート
	for _, files := range seqFilesMove {
		sort.Slice(files.Files, func(i, j int) bool {
			return files.Files[i].FileBaseInfo.FileName > files.Files[j].FileBaseInfo.FileName
		})
	}

	// パスを作成する
	sortedFiles := make([]IsFileMove, 0, len(seqFilesMove))
	for _, files := range seqFilesMove {
		sortedFiles = append(sortedFiles, files.Files...) // ソートされたファイルを順番にappend
	}

	count := 1                         // ファイルの連番
	for _, file := range sortedFiles { // 順番にファイルを取り出してパスを作成する。
		if !file.IsMove { // ファイルリネームをしない場合はスキップ // ディレクトリのファイルは無視。
			continue
		}
		// ファイルのリネームをする場合はパスを作成する
		afterPath := targetDir + "/" + prefix + "-" + intToSerialNumberString(count) + filepath.Ext(file.FileBaseInfo.FilePath)
		filesPathMove = append(filesPathMove, FilePathMove{BeforePath: file.FileBaseInfo.FilePath, AfterPath: afterPath})

		count++
	}

	return filesPathMove
}

// seq: intから連番の文字列を作成する関数
// 仕様:連番2桁までは、桁の前に「0」を付与。3桁以上はそのまま。
func intToSerialNumberString(num int) string {

	strNum := strconv.Itoa(num)

	if len(strNum) == 0 { //0バイト（一応）
		serialNumber := "000"
		return serialNumber
	}
	if len(strNum) == 1 { //1バイト
		serialNumber := "00" + strNum
		return serialNumber
	}
	if len(strNum) == 2 { //2バイト
		serialNumber := "0" + strNum
		return serialNumber
	} else { // 3バイト以上
		serialNumber := strNum
		return serialNumber
	}
}

// 計画ロジック（重要）
// 判定層１：冪等かどうかを判定
func isIndempotent(filePathMove []FilePathMove) bool {

	var pathMatchCount int // 移動前と移動後でマッチしたパスの回数をカウント

	if len(filePathMove) == 0 { // 空スライス、nilを判定
		return true // 空スライスの場合も冪等（変更なし）としてカウント
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
// 移動前パスおよび移動後パスで衝突検知。
// 注意：入力スライスの要素の並びを保証しない。
func hasPathConflict(filePathMove []FilePathMove) []FilePathHasConflict {

	notChangePathes := make([]FilePathMove, 0, len(filePathMove)) // 変更しないパスグループ
	changePathes := make([]FilePathMove, 0, len(filePathMove))    // 変更予定のパスグループ

	filesOperationResult := make([]FilePathHasConflict, 0, len(filePathMove)) // 予定が決定したパスを格納するスライス

	for _, file := range filePathMove { //グループ化
		if file.BeforePath == file.AfterPath {
			notChangePathes = append(notChangePathes, file)
		} else {
			changePathes = append(changePathes, file)
		}
	}

	// notChangeを構造体にマッピング（変更しないものは計画がすでに決まっているので、先に格納）
	for _, notChangePath := range notChangePathes {
		filesOperationResult = append(filesOperationResult, FilePathHasConflict{FilePathMove: notChangePath, HasConflict: false})
	}

	// 変更しないパスと、変更予定のパスの衝突を判定
	for _, notChangePath := range notChangePathes { // notChangePathは一意。 -> changePathが2回appendされることはない。
		for _, changePath := range changePathes {
			if notChangePath.BeforePath == changePath.AfterPath { // 同じもののグループの現在のパスと、違うものの変更予定のパスが同じ
				filesOperationResult = append(filesOperationResult, FilePathHasConflict{FilePathMove: changePath, HasConflict: true}) //衝突
				break                                                                                                                 // あるchangePathとnotChangePathが被った時点で、そのchangePathを調べる必要はない。
			}
			// かぶらなければ何もしない（理由：まだ計画が決まっていないため）
		}
	}

	// スライスchangePathからすでにappendされているパスを削除する。
	// すでに計画が決定した現在のパスをの存在を格納する。
	m := make(map[string]struct{}, len(changePathes))

	for _, v := range filesOperationResult { // 衝突判定を終えたパスをmapに格納
		m[v.FilePathMove.BeforePath] = struct{}{}
	}

	// changePathのパスを１つずつ見ていって、mにすでに存在したら、新しいスライスに入れる
	remainingChangePathes := make([]FilePathMove, 0, len(changePathes))

	for _, path := range changePathes {
		if _, ok := m[path.BeforePath]; !ok { // mにキーが存在しなかったら、まだ衝突判定結果が確定していないパス -> スライスに格納
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
			if i.AfterPath == j.AfterPath { //衝突した
				if _, ok := n[j.BeforePath]; ok { // すでに結果が決まっているかを判定。 -> 決まっていたらappendをスキップ
					continue
				}
				// まだ決まっていない場合 // パスをスライス`filesOperationResult`にappend
				filesOperationResult = append(filesOperationResult, FilePathHasConflict{FilePathMove: j, HasConflict: true})
				// 計画が決まったパスをnに格納
				n[j.BeforePath] = struct{}{}
			}
		}
	}

	// 衝突しなかったパス（計画が決まっていないパス）を抽出 -> complete
	for _, path := range remainingChangePathes {
		if _, ok := n[path.BeforePath]; !ok { // nにパスがない -> まだ計画が決まっていない。(衝突しなかったパス)
			filesOperationResult = append(filesOperationResult, FilePathHasConflict{FilePathMove: path, HasConflict: false})
		}
	}

	fmt.Println(len(filesOperationResult))

	return filesOperationResult
}

// 衝突判定結果を終えたパスを移動前のファイルパスのファイル名を見てソートする関数（ext / date） -> ファイル操作の実行順番で返す
// hasConflictはファイルの実行順を保証しないので、最後にソートする。
func sortFilesNameAcend(planedFilesPath []PlanedFilePath) []PlanedFilePath {
	sort.Slice(planedFilesPath, func(i, j int) bool {
		return filepath.Base(planedFilesPath[i].FilePath.BeforePath) < filepath.Base(planedFilesPath[j].FilePath.BeforePath) // ファイル名でソート
	})
	return planedFilesPath
}

// 予定を受け取って結果を返す関数
// --on-conflictに基づいて、衝突したものをスキップするかエラーとして操作を終了するのかを決めて、結果を返す。
// 仕様：--on-conflict：skip -> 衝突しているファイルをスキップしてファイル操作を実行する。
// 仕様：--on-conflict；error -> 衝突しているファイルを検出した時点で、ファイル操作を中断する。
// 2026-07-13追記：[]PlanedFilePathのファイル順序は保証しない
// conflict skipの場合：スライスからファイルパスを取り出して衝突判定をしているので、順序が保証されている。
// conflict errorの場合：cannotについて、マップからパスを取り出して[]PlanedFilePathにappendしているので、順序が保証されない
func summaryFilesMovePlan(hasPathConflict []FilePathHasConflict, conflict string) ([]PlanedFilePath, PlanSummary) {

	var completeFilesCount int     // ファイル操作できるファイル数
	var skipFilesCount int         // 衝突のためにスキップしたファイル数
	var errorFilesCount int        // エラーファイル数
	var unchangeFilesCount int     // パスの変更がないファイル数
	var cannotChangeFilesCount int //エラーによって操作を中断したファイル数

	planedFilePath := make([]PlanedFilePath, 0, len(hasPathConflict))

	if conflict == "skip" {
		for _, path := range hasPathConflict {
			if path.HasConflict {
				skipFilesCount++
				planedFilePath = append(planedFilePath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.FilePathMove.BeforePath, AfterPath: path.FilePathMove.BeforePath}, Status: SkipFileStatus})
				continue // スキップ
			}
			if path.FilePathMove.BeforePath == path.FilePathMove.AfterPath {
				planedFilePath = append(planedFilePath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.FilePathMove.BeforePath, AfterPath: path.FilePathMove.AfterPath}, Status: UnchangeFileStatus})
				unchangeFilesCount++
				continue // スキップ
			}
			// 衝突していない場合、実行できるファイルパス
			planedFilePath = append(planedFilePath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.FilePathMove.BeforePath, AfterPath: path.FilePathMove.AfterPath}, Status: CompleteFileStatus})
			completeFilesCount++
		}
		return planedFilePath, PlanSummary{CompleteFiles: completeFilesCount, SkipFiles: skipFilesCount, ErrorFiles: errorFilesCount, UnchangeFiles: unchangeFilesCount, CannotChangeFiles: cannotChangeFilesCount}

	} else { // error（default）

		// まだ計画が決定していないファイルパス（FilePathMove）をキー、値を空の構造体として格納する。
		// 計画が決定したファイルはマップから削除する。
		// 最後まで残っているファイルはcannotとして計画を決定する
		m := make(map[FilePathMove]struct{}, len(hasPathConflict))

		// mapにファイルを格納
		for _, fileMovePath := range hasPathConflict {
			m[fileMovePath.FilePathMove] = struct{}{}
		}

		for _, path := range hasPathConflict {
			if path.HasConflict {
				planedFilePath = append(planedFilePath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.FilePathMove.BeforePath, AfterPath: path.FilePathMove.BeforePath}, Status: ErrorFileStatus})
				errorFilesCount++
				delete(m, path.FilePathMove)
				break // 衝突ファイルを検知した時点で操作をやめる
			}
			if path.FilePathMove.BeforePath == path.FilePathMove.AfterPath { // パスの変更がないファイル
				planedFilePath = append(planedFilePath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.FilePathMove.BeforePath, AfterPath: path.FilePathMove.BeforePath}, Status: UnchangeFileStatus})
				unchangeFilesCount++
				delete(m, path.FilePathMove)
				continue
			}
			// 衝突していない場合、実行できるファイルパス
			planedFilePath = append(planedFilePath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.FilePathMove.BeforePath, AfterPath: path.FilePathMove.AfterPath}, Status: CompleteFileStatus})
			completeFilesCount++
			delete(m, path.FilePathMove)
		}

		// もし、衝突していたら残りのファイルは全て変更できない
		cannotChangeFilesCount = len(hasPathConflict) - completeFilesCount - errorFilesCount - unchangeFilesCount // エラーによって実行できないファイル数
		// マップの中で残っているファイルをplanedFilePathに格納していく
		// cannot（移動できない）ものとして決定
		for path, _ := range m {
			planedFilePath = append(planedFilePath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.BeforePath, AfterPath: path.BeforePath}, Status: CannotFileStatus})
		}

		return planedFilePath, PlanSummary{CompleteFiles: completeFilesCount, SkipFiles: skipFilesCount, ErrorFiles: errorFilesCount, UnchangeFiles: unchangeFilesCount, CannotChangeFiles: cannotChangeFilesCount}
	}
}

// リネーム計画を受け取ってサマリを返す関数
func summaryFilesRenamePlan(planPath []FilePathMove) ([]PlanedFilePath, RenamePlanSummary) {
	var completeFilesCount int         // 必要
	var skipFilesCount int = 0         // 不要？ 理由：衝突は起きないから（一時リネームを導入するため）
	var errorFilesCount int = 0        // 不要？　理由：衝突は起きないから（一時リネームを導入するため）
	var unchangeFilesCount int         // 必要
	var cannotChangeFilesCount int = 0 //　不要？　理由：衝突によるエラーは起こらないため（一時リネームを導入するため）

	planedFilesPath := make([]PlanedFilePath, 0, len(planPath))

	for _, path := range planPath {
		if path.BeforePath == path.AfterPath {
			planedFilesPath = append(planedFilesPath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.BeforePath, AfterPath: path.AfterPath}, Status: UnchangeFileStatus})
			unchangeFilesCount++
			continue
		}
		planedFilesPath = append(planedFilesPath, PlanedFilePath{FilePath: FilePathMove{BeforePath: path.BeforePath, AfterPath: path.AfterPath}, Status: CompleteFileStatus})
		completeFilesCount++
	}
	return planedFilesPath, RenamePlanSummary{CompleteFiles: completeFilesCount, SkipFiles: skipFilesCount, ErrorFiles: errorFilesCount, UnchangeFiles: unchangeFilesCount, CannotChangeFiles: cannotChangeFilesCount}
}

// ファイル操作のサマリを受け取って計画を表示する関数（ext, date）
func displayFilesMovePlan(planSummary PlanSummary) {

	if planSummary.CompleteFiles > 0 { // completeが１件以上で表示
		fmt.Printf("%d件のファイルをサブフォルダに移動します。\n", planSummary.CompleteFiles)
	}
	if planSummary.SkipFiles > 0 { // skipが1件以上で表示
		fmt.Printf("%d件のファイルのパスが衝突しています。衝突したファイルの移動をスキップします。\n", planSummary.SkipFiles)
	}
	if planSummary.ErrorFiles > 0 { // errorが1件以上で表示
		fmt.Printf("%d件のファイルがエラーです。\n", planSummary.ErrorFiles)
	}
	if planSummary.UnchangeFiles > 0 { // unchangeが1件以上で表示
		fmt.Printf("%d件のファイルの変更はありません。\n", planSummary.UnchangeFiles)
	}
	if planSummary.CannotChangeFiles > 0 { // cannotChangeが1件以上で表示
		fmt.Printf("%d件のファイルが移動できません。\n", planSummary.CannotChangeFiles)
	}
}

// ファイルのリネームのサマリを受け取って計画を表示する関数
// 備考：リネームは衝突が起きないので、skipCountはない。
func displayFilesRenamePlan(renamePlanSummary RenamePlanSummary) {
	if renamePlanSummary.CompleteFiles > 0 { // completeが１件以上で表示
		fmt.Printf("%d件のファイルをリネームします。\n", renamePlanSummary.CompleteFiles)
	}
	if renamePlanSummary.ErrorFiles > 0 { // errorが1件以上で表示
		fmt.Printf("%d件のファイルがエラーです。リネームできません。\n", renamePlanSummary.ErrorFiles)
	}
	if renamePlanSummary.UnchangeFiles > 0 { // unchangeが1件以上で表示
		fmt.Printf("%d件のファイルの変更はありません。\n", renamePlanSummary.UnchangeFiles)
	}
}

// 計画を実行するかしないかを表示して結果を返す関数
func doApplyPlan() (bool, error) {

	var doApply string

	fmt.Print("計画を実行しますか？ Y/n：")
	_, err := fmt.Scan(&doApply)
	if err != nil {
		// scan error
		return false, err // 実行しない
	}

	if doApply == "Y" {
		return true, nil //実行する
	}
	return false, nil // 実行しない
}

func main() {

	// コマンドパース
	commands, err := parseCommnads() // コマンドとパスを取得
	if err != nil {                  // エラーの場合、コンソールに出す。
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 共通処理
	allFilesPath, err := listFilesPath(commands.Path)

	// dry-run
	if !commands.Apply {

		// ext
		if commands.By == "ext" {
			filesIsMove := listFilesIsMove(commands.Path, commands.Recursive, allFilesPath)
			extFileMove := extractFilesExtension(filesIsMove)
			filePathMove := makeExtPathPlan(commands.Path, extFileMove)

			if isIndempotent(filePathMove) { // 冪等判定
				fmt.Println("ファイルの変更はありません。")
				return
			}

			planPath := hasPathConflict(filePathMove)
			planedFilesPath, planedSummary := summaryFilesMovePlan(planPath, commands.Conflict)
			sortPlanedFilesPath := sortFilesNameAcend(planedFilesPath)
			displayFilesMovePlan(planedSummary) // サマリーを表示

			fmt.Println(sortPlanedFilesPath)
		}

		// date
		if commands.By == "date" {

			filesIsMove := listFilesIsMove(commands.Path, commands.Recursive, allFilesPath)
			dateFileMove, err := extractFilesDate(filesIsMove)
			if err != nil {
				fmt.Println(err)
				return
			}
			filePathMove := makeDatePathPlan(commands.Path, dateFileMove)

			if isIndempotent(filePathMove) {
				fmt.Println("ファイルの変更はありません。")
				return
			}

			planPath := hasPathConflict(filePathMove)
			planedFilesPath, planedSummary := summaryFilesMovePlan(planPath, commands.Conflict)
			sortPlanedFilesPath := sortFilesNameAcend(planedFilesPath)
			displayFilesMovePlan(planedSummary) // サマリーを表示

			fmt.Println(sortPlanedFilesPath)
		}

		// seq
		if commands.By == "seq" {
			filesIsMove := listFilesIsMove(commands.Path, false, allFilesPath)       // Case Seq is automatically recursive false.
			extFileMove := extractFilesExtension(filesIsMove)                        // ファイル情報と拡張子
			seqFileMove := classifyExt(extFileMove)                                  // 拡張子ごとに分類
			planPath := makeSeqPathPlan(commands.Path, seqFileMove, commands.Prefix) // ソートしてパスを作成

			if isIndempotent(planPath) {
				fmt.Println("ファイルの変更はありません。")
				return
			}

			// 冪等でなければ、planPathがそのまま決定計画となる。
			// 以下追記：2026/7/12
			// ただし、「planPathがそのまま決定計画」について注意が必要。これは実行時に一時リネームによって、上書きを防ぐ前提だからこそ成り立つ。
			// seqコマンドは、衝突は起こらないが上書きは起こる。これは別のものとして考えた方が良さそう。
			planedFilesPath, renameFileSummary := summaryFilesRenamePlan(planPath)
			displayFilesRenamePlan(renameFileSummary)
			sortPlanedFilesPath := sortFilesNameAcend(planedFilesPath)

			fmt.Println(sortPlanedFilesPath)
		}

		// ユーザ入力
		isApply, err := doApplyPlan()
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

		// apply（not dry-run）
	} else {

		// Aplly

		return
	}
}
