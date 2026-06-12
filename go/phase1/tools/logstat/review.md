# コードレビュー結果（Issue #25 / logstat）

> 対象: `src/main.go`
> 観点: `/rubric-review`（Phase 1 × ツール、観点 ①仕様適合 ②計算量/設計 ③重複/可読性 ＋ 層2の問い）
> 補足: 「直したか」の進捗は **Issue #25** で、「どう直したか・なぜ」は `notes.md` で管理する。ここは **レビュー所見** を残す。
>
> - **1回目レビュー（修正前 / 2026-06-12）**: 以下「1回目レビュー」。
> - **2回目レビュー（修正後 / 2026-06-12）**: 末尾「2回目レビュー（修正後の確認結果）」を参照。

---

# 1回目レビュー（修正前）

## 全体評価

よく書けている。ガード節・早期リターン、構造体へのマッピング、関数分割（parse / open / extract / sort / topN / json）はきれい。stderr 化と `Atoi` エラー処理の修正も入っており、DoD の半分は満たしている。

ただし **仕様の核（「件数を集計してランキング出力」）に対して、出力が実は集計になっていない**。これが、notes.md で「わからなかった」と書かれた must の 1 本化・第 2 ソート、層2 問2 を貫く**震源**になっている。ここを直すと残りが連鎖的に解ける。

---

## 🔴 最重要：出力が「集計結果」になっていない（層2 問2 の答え）

notes.md 101 行目で感じていた違和感は **正しい**。これが本丸。

`--by status --top 3` の実際の出力:

```json
[{"ip":"192.168.0.3",...,"status":"200",...},
 {"ip":"192.168.0.4",...,"status":"200",...},
 {"ip":"192.168.0.1",...,"status":"200",...}]
```

これは「status 200 の生ログ行が 3 本」であって、仕様が求める「集計結果のランキング」ではない。
README は「**件数を集計し、ランキング形式で出力**」、DoD は「`--by status` で**件数の降順に集計結果**を」。期待される出力はこの形のはず:

```json
[{"key":"200","count":38},{"key":"404","count":5},{"key":"401","count":3}]
```

`Request` に件数フィールドが無いので件数を出しようがなく、「件数の多い軸の生ログを並べた」状態になっている。`--top 3` も「上位 3 つの軸」ではなく「3 行」になってしまう。

**ここを直すと残りが全部ほどける。** 集計を「軸の値 → 件数」の `(key, count)` ペアの一覧にする:

```go
type Count struct {
    Key   string `json:"key"`
    Count int    `json:"count"`
}

// 軸の値を取り出す関数を1つ渡すだけ（3ブロックの重複が消える）
func aggregate(reqs Requests, keyOf func(Request) string) []Count {
    m := make(map[string]int)
    for _, r := range reqs {
        m[keyOf(r)]++
    }
    counts := make([]Count, 0, len(m))
    for k, v := range m {
        counts = append(counts, Count{Key: k, Count: v})
    }
    return counts
}
```

呼び出し側で軸を選ぶ:

```go
keyFns := map[string]func(Request) string{
    "status": func(r Request) string { return r.Status },
    "path":   func(r Request) string { return r.Path },
    "ip":     func(r Request) string { return r.Ip },
}
counts := aggregate(validRequests, keyFns[commands.By])
```

これで **must①（3 ブロック 1 本化）が解決**する。

---

## 🔴 must：3 ブロックの重複 1 本化（notes「わからなかった #1」への答え）

notes.md 82・100 行目「インターフェースを使えば良いのか？」→ **インターフェースは不要**。
notes.md 54-55 行目で自分で出した結論「出現回数でソートする共通関数として切り出す方が良い」が正解だった。自分の直感を信じてよかったところ。

なぜインターフェースが過剰か: 3 つの分岐の違いは「**どのフィールドを見るか**」だけ。振る舞いの種類が違うわけではないので、必要なのは「フィールドを取り出す関数（上の `keyOf`）」を 1 つ差し替えること。Go ではこれを **関数値（`func(Request) string`）** で渡すのが最も軽量。

> インターフェースが要るのは「実装ごとに**中身の処理自体**が変わる」とき。今回は中身は同じで入口（取り出すフィールド）だけ違うので、関数 1 個で十分——という判断を notes に書ければ層2 問1 の満点回答。

---

## 🔴 must：同件数の第 2 ソート（notes「わからなかった #2」への答え）

集計を `[]Count` にすれば一気に簡単になる。比較関数に「件数が同じなら**キー名の昇順**」を足すだけ:

```go
sort.Slice(counts, func(i, j int) bool {
    if counts[i].Count != counts[j].Count {
        return counts[i].Count > counts[j].Count // 件数の降順
    }
    return counts[i].Key < counts[j].Key          // 同数ならキー名の昇順
})
```

ポイントは「**比較に第 2 の基準を足す**」だけ。今の生ログ方式だと「同じ status の行が大量」で第 2 ソートのしようがなかった（行のどれを優先？が決められない）——これも震源が同じ。

補足: `sort.Slice` は不安定ソートだが、比較関数が同値（件数もキーも同じ）を返さない限り（= キーが一意なので返さない）、結果は完全に決まる。`sort.SliceStable` でなくても安定する。

---

## 🟠 should：引数を flag へ（notes「わからなかった #1」: `--by <変数>` の検出）

notes.md 10-13 行目「`--by <変数>` をまとめて処理したい / コマンド不正を検出できない」は、**入力経路が根本的にズレている**のが原因。

今は `fmt.Scanf` で **標準入力から** `logstat ... --by ... --top ...` を読んでいる。だが仕様は CLI コマンド:

```
logstat <logfile> [--by status|path|ip] [--top N] [--json]
```

これは **`os.Args`（コマンドライン引数）** で渡されるもの。`logstat` というプログラム名を stdin から読んでいるのが不自然さの正体。`flag` パッケージを使うと「やりたかったこと」がそのまま実現する:

```go
func parseCommands() (*Commands, error) {
    by := flag.String("by", "status", "集計軸 status|path|ip") // ← 既定値 status
    top := flag.Int("top", 0, "上位N件（0で全件）")              // ← int で受けられる
    asJSON := flag.Bool("json", false, "JSON出力")
    flag.Parse()

    args := flag.Args() // フラグ以外の残り = logfile
    if len(args) < 1 {
        return nil, errors.New("ログファイルを指定してください")
    }
    // by の妥当性チェックは今のままでOK
    ...
}
```

これで:

- **`--by` のデフォルト（status）・`--top` 省略（全件）・順不同**が自動で効く（DoD/README の「デフォルト」列を満たす）
- `--top` を最初から `int` で受けられる（notes.md 94 行目の「最初のパースで int に」がここで実現。`displayTopN` の `Atoi` 二重化も消える）
- `--by abc` のような不正値も検出できる

---

## 🟠 should：不正行スキップの「件数の警告」が未実装（DoD 未達）

`isValidLine` でスキップ自体はできている（sample.log の "this line is broken" / "GET /no-ip 200 100" / "...extra-field-here" は弾けている）が、DoD「**件数を警告する**」が出ていない。スキップ数を数えて最後に stderr へ:

```go
fmt.Fprintf(os.Stderr, "警告: 不正な行を %d 件スキップしました\n", skipped)
```

終了コードは 0 のまま（エラーではなく警告だから）——これが「エラーと警告を出力先で区別」という設計ポイントの実演になる。

---

## 🟡 should：`strings.Fields` の多重実行（notes「わからなかった #2」: 5 回問題）

notes.md 85-87 行目の通り。今 `extractValidRequestLines` で 1 行につき `strings.Fields` を **6 回**（`isValidLine` 内 1 回 + マッピング 5 回）呼んでいる。1 回にまとめる:

```go
for scanner.Scan() {
    fields := strings.Fields(scanner.Text()) // 1回だけ
    if len(fields) != 5 {
        skipped++
        continue
    }
    validRequestLines = append(validRequestLines, Request{
        Ip: fields[0], Method: fields[1], Path: fields[2],
        Status: fields[3], Bytes: fields[4],
    })
}
```

ついでに `isValidLine` は不要になり、**マジックナンバー添字も 1 か所に閉じる**（nice 項目）。

---

## 🟡 細かい点

- **`openFile` の `defer file.Close()`（main.go:90）**: `err != nil` の枝、つまり `file` が `nil` か無効なときに Close している。意味がなく、ここは `defer` 不要（`return nil, err` だけ）。Close は呼び出し側で開いた直後に `defer` するのが定石。
- **`extractValidRequestLines` が file を閉じている**: 開いた場所（呼び出し側）で閉じる方が責務が揃う。「開けた者が閉じる」。
- **`Status` / `Bytes` を `string` で持つ件（層2 問2 の後半）**: 集計の軸（status）として使うだけなら string で問題なし。ただし発展課題「バイト数の合計・平均」をやるなら `Bytes` は `int` が必要。「今は集計キーとしてしか使わないので string、数値演算が必要になったら int に変える」と notes に書ければ OK。
- **stderr にした理由（notes.md 97 行目「わからなかった」）**: 標準出力(stdout)は「**プログラムの成果物**」を流す管。`logstat ... | jq` のようにパイプで次に渡る。エラー/警告メッセージをここに混ぜると JSON が壊れる。標準エラー出力(stderr)は「**人間向けの連絡**」専用の別の管で、パイプに乗らない。だから成果物(JSON)は stdout、警告・エラーは stderr に分ける——これが「出力先で区別」の意味。

---

## まとめ（おすすめ着手順）

「お手上げ」だった 3 つは、**①集計を `[]Count` に変える**を最初にやれば連鎖的に解ける:

1. **集計を `(key, count)` 化**（震源）→ 出力が仕様通りに / 1 本化・第 2 ソートが自然に解ける
2. `flag` パッケージ化 → `--by` 検出・デフォルト・top の int 化
3. スキップ件数を stderr 警告（DoD 最後の未達）
4. `strings.Fields` 1 回化 + マッピング整理

断念はもったいなくない。むしろ「件数フィールドが無いのが変かも（層2 問2）」「インターフェースより共通関数では（54 行目）」と**自力で震源に近いところまで気づけていた**のが大きい。あとは「集計 = map の (key→count) を一覧にする」という型を一つ覚えれば、3 つまとめて崩せる。

---

## Issue #25 項目との対応

| Issue #25 の項目 | 区分 | 状態 | 本レビューの該当節 |
| --- | --- | --- | --- |
| `sortByDescend` 3 ブロックの 1 本化 | must | 未 | 最重要 / must 1本化 |
| 同件数時の第 2 ソート | must | 未 | must 第2ソート |
| エラー出力を stderr へ | must | **済** | （対応済み） |
| 不正行スキップの件数警告 | should | 未 | should 件数警告 |
| 引数を `flag` へ | should | 未 | should flag |
| `strings.Fields` の 1 回化 | should | 未 | should Fields多重 |
| フィールド添字のマジックナンバー解消 | nice | 未 | should Fields多重（同時に解消） |
| `displayTopN` の全件コピー / `Atoi` | nice | 一部 | should flag（int 化で不要に） |

> 加えて本レビュー独自の最重要指摘: **出力が集計結果になっていない**（Issue には未記載。層2 問2 の核）。

---

# 2回目レビュー（修正後の確認結果）

> 対象: 1回目レビューを受けて修正された `src/main.go`（2026-06-12）
> 確認方法: `go build` ＋ `sample.log` での実行（出力先・安定性・順不同・デフォルト値を実測）。

## 総評

**震源（出力が集計結果になっていない）を直しきった、いい修正。** 出力が `(key, count)` のランキングになり、`aggregate` + `keyFns` で 3 ブロックの重複も解消。第 2 ソートも安定動作を実測で確認した（`--by path` で `count:7` の `/api/login` と `/api/users` が、キー昇順で login → users と固定。3 回実行とも同一の並び）。stdout に JSON / stderr に警告、の分離もできている。

**残課題は実質 1 件（`--top` 省略時のバグ）のみ。** ここを直せば must/should は全クリア。

## 実測ログ（要点）

```
$ logstat --by status --top 3 --json sample.log
警告: 不正な行を 3 件スキップしました。      # ← stderr
[{"key":"200","count":30},{"key":"404","count":5},{"key":"401","count":3}]   # ← stdout, 件数ランキング ✓

$ logstat --by path --top 3 sample.log        # 3回とも同一 = 第2ソート安定 ✓
[{"key":"/index.html","count":10},{"key":"/api/login","count":7},{"key":"/api/users","count":7}]

$ logstat --top 2 --by path sample.log        # 順不同（logfileが末尾）でも動く ✓
[{"key":"/index.html","count":10},{"key":"/api/login","count":7}]

$ logstat                                     # 引数なし
ログファイルを指定してください   (exit=1) ✓

$ logstat --by xxx sample.log                 # 不正な --by
コマンドが不正です。            (exit=1) ✓
```

## 🔴 残バグ：`--top` 省略（=0）で「全件」にならず空になる

`--top` の help は「上位N件（**0で全件**）」だが、実際は省略時に `[]` が返る（DoD「`--top` 省略で全件」未達）。

```
$ logstat sample.log          # --top 省略 → 全件のはずが…
[]
$ logstat --by ip --json sample.log
[]
```

原因は `displayTopN`（main.go:170）。`commandTop == 0` でループ 0 回 → 空スライス。help の宣言と実装が食い違っている。

**直し方**（先頭で 0 以下を全件扱い。`error` 戻り値はもう常に nil なので撤去し、`sorted[:n]` でスライス）:

```go
func displayTopN(commandTop int, sorted []Count) []Count {
    if commandTop <= 0 || commandTop > len(sorted) {
        return sorted // 0以下=全件 / 件数より多い指定も全件
    }
    return sorted[:commandTop]
}
```

## 🟡 細かい点（任意）

- **スキップ 0 件でも警告が出る**（main.go:128）: `警告: 不正な行を 0 件スキップしました。` が毎回出る。`if skipped > 0 { ... }` で囲むと、警告すべき時だけ出る（「警告は必要な時だけ」という設計判断の練習）。
- **`len(fields) != 5` の `5`（main.go:103）**: notes でも自問していた通り。`const fieldsPerLine = 5` と名付けると「5 が何か」がコード上で説明される（DRY というより*自己文書化*）。
- **`extractValidRequestLines(file, request Request)` の `request` 引数（main.go:86）**: 外から渡す必要はなく、ループ内ローカル変数で十分。`main` 側の `var request Request`（192行）も不要に。`file.Close()` の責務もコメント自覚どおり「開いた側（main）で `defer`」が筋。
- **使わなくなった残骸の削除**: `Commands.CommandIndex`（もうセットされない）、コメントアウトの `var` 群（41-45行）、`sortByDescend` 旧呼び出しコメント（224行）。`Request.Method`/`Bytes` は 5 フィールド整合チェックの役目があるので残してよい。

## Issue #25 進捗（2回目レビュー時点）

| 項目 | 区分 | 状態 |
| --- | --- | --- |
| `sortByDescend` 3 ブロックの 1 本化（`aggregate` + `keyFns`） | must | ✅ |
| 同件数時の第 2 ソート（同数でキー昇順） | must | ✅ 安定を実測確認 |
| エラー出力を stderr へ | must | ✅ |
| 不正行スキップの件数警告 | should | ✅（0件時の抑制は任意） |
| 引数を `flag` へ（デフォルト/順不同/int 化） | should | ✅ 順不同も実測確認 |
| `strings.Fields` の 1 回化 | should | ✅ |
| フィールド添字のマジックナンバー解消 | nice | ✅（`5` の定数化が残るのみ） |
| `displayTopN` の全件コピー / `Atoi` | nice | △ コピーは残（→ 残バグ修正と同時に `sorted[:n]` で解消可） |
| 出力が集計結果（層2 問2） | — | ✅ 解決 |

> **残タスクは実質「`--top 0` = 全件のバグ修正」1 本。** ここを直せば must/should 全クリア。
