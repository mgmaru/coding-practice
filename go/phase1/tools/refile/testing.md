# テスト手引き: どのディレクトリでどのテストケースを実行するか

`testdata/` 配下の各 fixture と、それで検証する DoD（受け入れ条件）の対応表。
README の DoD 1つ1つを「どのフォルダで・どのコマンドで・何を確認するか」に落としたもの。

- 実行場所: **`refile/` ディレクトリから**実行する（`<dir>` 引数は `testdata/<case>` の相対パスで渡す）。
- 実行形式: `go run main.go [オプション] <dir>`（`<dir>` は flag の仕様上**最後**に置く）。

> ⚠️ **`--apply` は破壊的。fixture の原本を壊さないこと。**
> `testdata/` 配下は「書き換えない原本（テンプレート）」として扱う。`--apply` を試すときは必ずコピー先で行う:
> ```
> cp -r testdata/conflict /tmp/refile-x && go run main.go --by ext --recursive --apply /tmp/refile-x
> ```
> 冪等は前回の計画を保存せず FS の現状を見て判断する方針（README 参照）なので、原本を保てば何度でも同じ初期状態から検証できる。

---

## fixture 別の対応表

| fixture | 構成 | 主に検証する DoD | 代表コマンド |
| --- | --- | --- | --- |
| `testdata/sample/` | 基本セット（pdf/png/jpg/拡張子なし/nested） | dry-run 無変更・`ext`/`date`/`seq` の計画・`--recursive`・並び順安定・件数サマリ | `go run main.go --by ext testdata/sample` |
| `testdata/conflict/` | `a/report.pdf`, `b/report.pdf`（同名・別階層） | 衝突を**計画段階で**検知・`--on-conflict error|skip` | `go run main.go --by ext --recursive testdata/conflict` |
| `testdata/idempotent/` | `pdf/report.pdf`, `png/cat.png`（既に行き先にいる） | 冪等（現在パス==変更後パス → 変更なし） | `go run main.go --by ext --recursive testdata/idempotent` |
| （fixture 不要） | 引数なし / 非ディレクトリ / `seq` で `--prefix` 無し | 引数・オプションのバリデーション | 下記参照 |
| （`t.TempDir()`） | 空ディレクトリ | 0件サマリ | Go テストで生成（git は空 dir を追跡できないため fixture を置かない） |

---

## 1. `testdata/sample/` — 基本動作

構成:
```
sample/
  report.pdf  invoice.pdf  cat.png  noext  IMG_001.jpg  IMG_002.jpg
  nested/photo.jpg
```

### DoD: `--apply` なしはファイルを変更せず、計画（現在 → 変更後）のみ出力
```
go run main.go --by ext testdata/sample
```
- 期待: 計画だけ表示され、`sample/` の中身は1つも動かない（実行後に `find testdata/sample` で確認）。

### DoD: `--by ext` の計画を出力できる
```
go run main.go --by ext testdata/sample
```
- 期待（トップレベルのみ。`--recursive` 無しなので `nested/photo.jpg` は対象外）:
  - `report.pdf`  → `testdata/sample/pdf/report.pdf`
  - `invoice.pdf` → `testdata/sample/pdf/invoice.pdf`
  - `cat.png`     → `testdata/sample/png/cat.png`
  - `noext`       → `testdata/sample/noext/noext`
  - `IMG_001.jpg` → `testdata/sample/jpg/IMG_001.jpg`
  - `IMG_002.jpg` → `testdata/sample/jpg/IMG_002.jpg`

### `--by date`（更新月で振り分け）
```
go run main.go --by date testdata/sample
```
- 期待: 各ファイルの mtime（更新月 `YYYY-MM`）のサブフォルダ行きの計画。月は実行時のファイル mtime 依存。

### `--by seq`（連番リネーム。サブフォルダは作らない）
```
go run main.go --by seq --prefix img testdata/sample
```
- 期待: トップレベルのファイルが並び順（元パスの辞書順）で `img-001`, `img-002`, … にリネームされる計画。

### DoD: `--recursive` でサブディレクトリも対象
```
go run main.go --by ext --recursive testdata/sample
```
- 期待: 上記に加えて `nested/photo.jpg` → `testdata/sample/jpg/photo.jpg` が計画に入る。

### DoD: 変更計画の並び順が安定（元パスの辞書順など）
```
go run main.go --by ext --recursive testdata/sample
```
- 期待: 何度実行しても計画の行順が同じ（`IMG_001.jpg` < `IMG_002.jpg` < `cat.png` < … の安定順）。

### DoD: `--apply` 時のみ移動・改名し、件数サマリ（移動/スキップ/エラー）を出力
```
cp -r testdata/sample /tmp/refile-sample
go run main.go --by ext --recursive --apply /tmp/refile-sample
```
- 期待: 実際に振り分けが行われ、「移動 N 件 / スキップ 0 / エラー 0」のサマリが出る。

---

## 2. `testdata/conflict/` — 衝突検知

構成:
```
conflict/
  a/report.pdf   # 中身: report in a/
  b/report.pdf   # 中身: report in b/（別ファイル・同名）
```
`--by ext --recursive` で `a/report.pdf` と `b/report.pdf` が**両方** `conflict/pdf/report.pdf` 行きになり衝突する。

### DoD: 複数ファイルが同じ変更後パスになる衝突を、計画段階で検知して報告
```
go run main.go --by ext --recursive testdata/conflict
```
- 期待（`--on-conflict` 既定 = `error`）: 計画段階で衝突を検知し、エラーとして報告して **apply 前に止まる**（exit code 1）。ファイルは1つも動かない。

### DoD: `--on-conflict skip` は衝突ファイルをスキップして続行
```
go run main.go --by ext --recursive --on-conflict skip testdata/conflict
```
- 期待: 衝突する `report.pdf` 群はスキップ（サマリの「スキップ」に計上）、衝突しないファイルがあれば続行。

> 実際に apply してスキップ件数を見たい場合はコピー先で:
> ```
> cp -r testdata/conflict /tmp/refile-conflict
> go run main.go --by ext --recursive --on-conflict skip --apply /tmp/refile-conflict
> ```

---

## 3. `testdata/idempotent/` — 冪等

構成:
```
idempotent/
  pdf/report.pdf   # 既に pdf/ にいる
  png/cat.png      # 既に png/ にいる
```
`--by ext` で計算した変更後パスが現在のパスと一致する（`pdf/report.pdf` の拡張子は `pdf` → 行き先も `pdf/report.pdf`）。

### DoD: 同じコマンドを2回適用しても破綻しない（2回目は「変更なし」）＝冪等
```
go run main.go --by ext --recursive testdata/idempotent
```
- 期待: すべて「現在パス == 変更後パス」となり、計画は「変更なし」。apply しても何も動かない。

別の確認方法（sample を2回 apply して2回目が無変更になることを見る）:
```
cp -r testdata/sample /tmp/refile-idem
go run main.go --by ext --recursive --apply /tmp/refile-idem   # 1回目: 移動が起きる
go run main.go --by ext --recursive --apply /tmp/refile-idem   # 2回目: 変更なし
```

---

## 4. fixture 不要のケース（引数・オプションのバリデーション）

### DoD: 引数なし → わかりやすいエラーを stderr、exit code 1
```
go run main.go --by ext
```
- 期待: 「ディレクトリを指定してください。」を stderr に出して exit 1。

### DoD: 対象がディレクトリでない／存在しない → エラー、exit code 1
```
go run main.go --by ext testdata/sample/report.pdf   # ファイルを指定
go run main.go --by ext testdata/does-not-exist      # 存在しないパス
```
- 期待: それぞれ「ディレクトリを指定してください。」「存在しないパスです。」を stderr に出して exit 1。

### DoD: `--by seq` で `--prefix` が無い → エラー
```
go run main.go --by seq testdata/sample
```
- 期待: 「prefixを指定してください。」を stderr に出して exit 1。

---

## 5. 空ディレクトリ（fixture を置かない）

git は空ディレクトリを追跡できず、`.gitkeep` を置くとツールからは「拡張子なしファイル1件」に見えて**真の空でなくなる**。
そのため空ディレクトリのケースは fixture にせず、Go テスト内で `t.TempDir()`（最初から空）を使って検証する。
- 期待: 「ファイルが移動できません。」等、0件サマリ（操作対象なし）になること。

---

## 6. 部分失敗（apply 途中の個別失敗）

### DoD: apply 中に個別操作が失敗しても全体を中断せず続行し、失敗はエラー件数に計上（ロールバックしない）

静的 fixture では作りにくい（権限なし・実行中に元ファイルが消えた、等の動的な状況）。
手動で再現する例（コピー先で行う / macOS）:
```
cp -r testdata/sample /tmp/refile-fail
# 行き先サブフォルダを先に作って書き込み不可にする等で、一部だけ失敗させる
chmod 555 /tmp/refile-fail        # 例: 親を読み取り専用にして移動を失敗させる
go run main.go --by ext --recursive --apply /tmp/refile-fail
chmod 755 /tmp/refile-fail        # 後始末
```
- 期待: 失敗したファイルはエラー件数に計上され、残りのファイル処理は続行される（途中で全体が落ちない）。
