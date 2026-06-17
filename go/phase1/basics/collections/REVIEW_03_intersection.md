# レビュー: `03_intersection` —「テストが別の関数を叩く」と「重複・順序の契約」

`03_intersection`（Phase 1 ドリル）のレビュー。
線形探索と set 探索の **2 実装を比較する** のが狙い。だからこそ「2 つが本当に同じ仕様か」「両方が本当にテストされているか」が肝になる。

> 大原則（review-rubric の ②道具の良し悪し）: イディオムか・計算量に無駄はないか・データ構造の選択は妥当か。
> 今回は型は通り、`go test` も**全部緑**。だが緑が嘘をついている。そこをどう潰すか。

前提: `go test -run SearchElements -v` は 20 サブテスト全 PASS。**しかしこの緑は信用できない**（観点1）。

---

## 観点1【重大】テストが「テスト対象の関数」を呼んでいない

`03_intersection_test.go:43`、`TestSetSearchElements` の中身がこう:

```go
func TestSetSearchElements(t *testing.T) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			gotValue := linearSearchElements(c.inputA, c.inputB) // ← setSearchElements のはず
			...
		})
	}
}
```

`TestLinearSearchElements` をコピペして関数名を**差し替え忘れた**。
結果、`setSearchElements` は **一度も実行されていない**。`linearSearchElements` を 2 回テストしているだけ。

### なぜ致命的か
- このドリルの主役は set 実装。その set 実装の**カバレッジが実質ゼロ**。
- 「全部 PASS」という最も安心する signal が、**主役を素通り**して出ている。テストの存在価値が反転している。
- `go test` の緑は「書いたコードが正しい」ではなく「**呼んだコードが正しい**」しか保証しない。何を呼んだかを毎回確認する。

### 直し方
```go
gotValue := setSearchElements(c.inputA, c.inputB)
```
直したうえで、**わざと `setSearchElements` を壊して赤くなることを確認**する（テストがその関数を握っている証明）。

> 合言葉: **全部緑のときこそ「何を呼んだか」を疑う。コピペテストは関数名の差し替えが本体。**

---

## 観点2【重大】2 実装の「重複の扱い」が食い違っている（仕様が未定）

線形と set は、**入力に重複があると違う答えを返す**。比較対象が同じ仕様でないと、そもそも比較にならない。

| 入力 | `linearSearchElements` | `setSearchElements` |
|---|---|---|
| A=`[2,2]`, B=`[2]` | `[2, 2]`（A の重複ぶん二重カウント） | `[2]`（A は set に畳まれる） |
| A=`[2]`, B=`[2,2]` | `[2, 2]`（B の重複ぶん） | `[2, 2]`（B の重複ぶん） |

- 線形: `for a in A { for b in B { if a==b append(b) } }` → **A の重複 × B の重複**ぶん出る。
- set : A を set 化（重複が消える）→ B を走査 → **B の重複ぶんだけ**出る。

つまり「A 側の重複」で挙動が割れる。
notes.md には **「set とは重複しない要素の集まり」** と書いたのに、出力スライスは**どちらもデデュープしていない**。狙いと実装がねじれている。

### 判断の指針 ―― まず「intersection の契約」を一文で決める
1. **集合的 intersection**（一般的）: 結果は**ユニーク**。`[2,2]∩[2]=[2]`。
   → 両実装とも「結果 set」を持って append 前に存在チェックし、デデュープして揃える。
2. **B のうち A に含まれる要素を全部**（重複保持）: なら線形の二重ループは**バグ**（A の重複で増殖する）。
   線形側も「B を走査し、A に含まれるか」に統一する:
   ```go
   for _, b := range inputB {
       for _, a := range inputA {
           if a == b { containingElements = append(containingElements, b); break } // break で A 重複の増殖を止める
       }
   }
   ```
どちらでもよいが、**2 実装が同じ契約**になっていること。そして**重複ケースをテストに足す**（観点5）。

> 合言葉: **2 つの実装を比べるなら、まず「同じ答えを返すべき仕様」を 1 つ決める。重複は最初に割れる。**

---

## 観点3 出力の「順序」も 2 実装で割れうる（テスト未カバー）

- 線形: **A の順**（外ループ A）にグルーピングされて出る。
- set : **B の順**（B を走査）で出る。

例: A=`[4,2]`, B=`[2,4]` → 線形 `[4,2]` / set `[2,4]`。
今のテストケースは**たまたま順序が一致する並び**ばかりで、この差が露見しない。
契約として「結果は B の出現順」等を決め、**順序が割れる入力**をテストに入れて固定する。

> 合言葉: **「順序は問わない」も立派な契約。だが言語化してテストで固定しないと、片方が静かに別物になる。**

---

## 観点4 `any` キー / `any` 比較の落とし穴（設計意思の補足として）

コメント通り「線形 vs set の比較」目的で `[]any` にしたのは妥当。ただし `any` には地雷がある:

- `a == b`（線形）も `set[a]`（set のキー）も、**動的型が comparable でないと実行時 panic**。
  例: 要素にスライスや map が入ると落ちる（コンパイルは通る）。
- 本番では `comparable` 制約のジェネリクスが本筋。コメントの「本来 any は使わない方が良い」は正しい。
  一歩進めるなら **「なぜ」=「comparable でない要素で panic する/型安全が消えるから」** まで書けると強い。

```go
// 将来の到達点（メモ）: 型安全 & デデュープ込み
func Intersection[T comparable](a, b []T) []T { ... }
```

これは今直す必要はない。**「any にした代償（panic 可能性）」を一文持っておく**だけでよい。

---

## 観点5 nil / 空スライスは OK。だが「重複」ケースが丸ごと無い

- `make([]any, 0)`（非 nil 空）で統一し、空・nil 入力で `[]any{}` を返す設計は **`01` の学びを踏襲できていて良い**。テストの `Want` も `[]any{}` で揃っている。✓
- 一方、テストケース 10 個は **すべて要素ユニーク**。観点2・3 の核心（重複・順序差）を**1 つも踏んでいない**。
  比較ドリルの検証として穴。最低限これを足す:
  - `DupInA`: A=`[2,2,3]`, B=`[2,3]` → 契約に応じた `Want`
  - `DupInB`: A=`[2,3]`, B=`[2,2,3]` → 同上
  - `OrderDiff`: A=`[4,2]`, B=`[2,4]` → 契約どおりの順で `Want`
  - （集合契約にするなら）`Want` は両実装で**同一**になるはず。ならない=仕様未統一。

---

## 観点6 軽微（イディオム・命名）

- **容量先取り**（②nice）: 結果は最大でも走査側の長さ。`make([]any, 0, len(inputB))` で再確保を減らせる。
- テスト構造体のフィールド: `Name`/`Want` は公開、`inputA`/`inputB` は非公開で**ケース内の大小が不揃い**。同パッケージなので動くが、揃えると読みやすい。
- 名前のタイポ: ケース名 `StringSlicend...`（`A` 欠落）、`notIcluded`（`Included`）。notes.md のファイル名 `03_intersectiontest.go`（正しくは `03_intersection_test.go`）。
- 計算量コメント `O(2n) -> O(n)` は方向性 OK。厳密には set 構築・`any` のハッシュ/インターフェース比較のコストが乗るが、ドリルの粒度では十分。

---

## 修正チェックリスト（`03_intersection`）

### テスト `03_intersection_test.go`
- [x] **`TestSetSearchElements` が `setSearchElements` を呼ぶ**（観点1・最優先）
- [x] 重複ケース `DupInA` / `DupInB` を追加（観点5）
- [x] 順序が割れる `Orderdiff` を追加（観点3）
- [x] 追加ケースの `Want` が**決めた契約どおり**（`Orderdiff` は B 順の `[2,4]`）

### 実装 `03_intersection.go`
- [x] intersection の契約を一文で決めた（**重複保持しない＝デデュープ**）（観点2）
- [x] 2 実装がその契約で**同じ答え**を返す（両方デデュープ＋B 順で決定的）
- [x] 結果の**順序の契約**を決め、テストで固定した（B の出現順）（観点3）
- [x] `any` の代償（comparable でない要素で panic）を一文で説明できる（観点4）
- [x] （②nice）`make([]any, 0, len(inputB))` で容量先取り

### 動作確認
- [x] `go test -run SearchElements -count=5 ./phase1/basics/collections` が **5 回連続で緑**（PASS 140 件）
- [x] `[no tests to run]` が出ていない / `gofmt`・`go vet` クリーン

---

## 再レビュー依頼の前に（層2: 自分の言葉で）

1. `TestSetSearchElements` は今、**どの関数**を呼んでいるか。直した後、壊して赤を見たか。
2. intersection の契約は「ユニーク」か「重複保持」か。2 実装はその契約で**一致**するか。
3. 結果の順序は何順と決めたか。それが割れる入力で、テストは固定できているか。
4. `any` にした代償を一文で言えるか（このドリルの狙いは線形 vs set の**速度比較**。型安全はその対価）。

> 進め方は rubric の大原則③: **①自己レビュー → ②AIレビュー → ③差分を読む**。
> まず観点1（呼び間違い）を潰し、観点2 の契約を決めてから `/rubric-review` を再依頼する。

---

# 修正サイクルの記録（観点7）—「緑の嘘」から「フレーキー」へ

ここからは、上のチェックリストを潰す**修正のやり取りで実際にハマった**こと。
このドリルは見た目より難しい。**`go test` の緑／赤が、3 回連続で「嘘」をついた**。その嘘の質が毎回違う。

## 修正1回目: 観点1〜6 に着手 → だが「赤がフレーキー」になった

やったこと（良い前進）:
- 観点1: `TestSetSearchElements` が `setSearchElements` を呼ぶよう修正。主役がやっとテストされた。✓
- 観点6: フィールドを `InputA`/`InputB` に統一。✓
- 観点2/3/5: `DupInA`/`DupInB`/`Orderdiff` を追加し、「重複は許さない／順序は問わない」と契約を**言語化**。✓
- 観点2: 重複除去の汎用関数 `deleteDuplicateElements` を新設。✓

ところが**テストがほぼ毎回 FAIL** するようになった。前回（嘘の緑）から**嘘の赤**へ。

### 根本原因: `map` の走査順は「ランダム」

```go
func deleteDuplicateElements(input []any) []any {
    m := make(map[any]struct{})
    for _, e := range input { m[e] = struct{}{} }
    out := make([]any, 0, len(input))
    for e := range m {                 // ★ Go の map 走査は実行ごとに順不同（仕様）
        out = append(out, e)
    }
    return out
}
```

`map` から要素を取り出して作った `out` は、**実行のたびに並びが変わる**。
両関数ともこの戻り値を起点に結果を組むので、出力スライスの並びが毎回バラける。

### そこに `reflect.DeepEqual`（順序に厳密）がぶつかる

```
TestSetSearchElements/IntSliceAndAisLonger  gotValue=[4 2] , Want=[2 4]   ← Orderdiff以外も落ちる
TestSetSearchElements/Orderdiff             gotValue=[2 4] , Want=[4 2]
```

`Orderdiff` だけでなく、**2 要素以上一致する全ケース**が確率的に落ちる。`-count=1` を 8 回回して全 FAIL。

> ここでの学び: **`map` 由来のスライスは順不同。`reflect.DeepEqual` は順序に厳密。この 2 つが出会うと、ロジックが正しくてもテストが踊る。**

### さらに「契約」と「テスト」が逆を向いていた

コメントには「順序は問わない」と書いたのに、テストは `reflect.DeepEqual(got, []any{4,2})` で**順序を問うていた**。
宣言と実装が逆。**この矛盾を放置したまま並びだけ直しても直らない**。先に契約を 1 つに倒す必要があった。

## 修正2回目: コメントだけ変えてコードは据え置き（実体ゼロ更新）

`git diff` の実体は**テストのコメント 1 行だけ**:

```diff
- 修正（観点3）：順序については問わない仕様とする。
+ 修正（観点3）：順序についてはbの要素の並びとする。
```

`deleteDuplicateElements` の `for e := range m` は**手つかず**。当然まだ毎回 FAIL。

> 学び: **コメントは契約の宣言にすぎない。コードが `map` 走査のままなら順序は決まらない。「直した」とは差分のことで、文言のことではない。**

## 修正3回目: テスト側だけ直った（`Want=[2,4]`）→ 「たまに緑」になった

`Orderdiff` の `Want` を B 順の `[2,4]` に修正（テスト側 ✓）。だが実装側 `deleteDuplicateElements` は**まだ `map` 走査**。
結果、出力が `[2,4]` になった時だけ通る = **8 回中 1 回だけ `ok`、残り FAIL**。

```
FAIL / FAIL / ok / FAIL / FAIL / FAIL / FAIL / FAIL
```

> 学び（最重要）: **「たまに緑」は赤より危険**。1 回の `go test` で判断すると「直った」と誤認する。
> フレーキーは**反復実行（`-count=5` など）でしか見えない**。緑を 1 回見て信じない。

## 決着: 直すのは実装 1 か所 ——「`map` は判定、順序は入力スライスで持つ」

```go
func deleteDuplicateElements(input []any) []any {
    seen := make(map[any]struct{}, len(input)) // ★ map は「既出か」の判定だけに使う
    out  := make([]any, 0, len(input))
    for _, e := range input {
        if _, ok := seen[e]; !ok {
            seen[e] = struct{}{}
            out = append(out, e)               // ★ append は入力順 → 順序が決定的に
        }
    }
    return out
}
```

`deleteDuplicateElements(inputB)` が **B の出現順を保つ**ので、両関数とも結果が「B 順」で決定的になり、`Want=[2,4]` と安定一致する。
ポイントは **「重複判定（`map`）」と「順序（スライス）」を別の道具に分担**させたこと。`map` 一個に両方やらせると順序が死ぬ。

> 合言葉: **`map` に順序を期待しない。順序が要るなら、順序はスライスが持つ。`map` は存在判定の係。**

---

## このドリルが「難しかった」理由（メタな振り返り）

1 つのバグではなく、**観点が連鎖**していた:

1. **観点1（呼び間違い）** を直すと、隠れていた **観点2/3（重複・順序の仕様未定）** が露出する。
2. 観点2 を直すために入れた `deleteDuplicateElements` が、今度は **観点7（map 順不同＝フレーキー）** を生む。
3. しかも各段階で `go test` の出力が **緑の嘘 → 赤の嘘 → たまに緑** と表情を変え、**毎回ちがう読み方**を要求してくる。

> 一番の教訓: **テストの色（緑/赤）はそのまま信じる対象ではなく、「何を呼び」「何順で」「何回回したか」とセットで読む signal。**
> 特に `map`・`reflect.DeepEqual`・フレーキーの三点セットは、初見で全部踏むと迷宮になる。今回ここを一通り踏めたのは大きい。

---

## 修正チェックリスト（観点7 追補）

### 実装 `03_intersection.go`
- [x] `deleteDuplicateElements` を **入力順保持**に変更（`map` は存在判定のみ、`append` は入力順）
- [x] 両関数の結果順が「B の出現順」で**決定的**になっている

### 動作確認（フレーキー対策）
- [x] `go test -run SearchElements -count=5 ./phase1/basics/collections` が **5 回連続で緑**
- [x] 1 回だけの緑で「直った」と判断していない（`-count=1` を 10 回回して全 `ok` も確認）

### 再レビュー前の自問（観点7）
1. `map` を走査して作ったスライスの並びは保証されるか？（→ されない）
2. 「順序を問う／問わない」を決めたか。テストの比較方法（`DeepEqual` か、ソート/集合比較か）はその契約と**同じ向き**か。
3. その修正は**コードに入っているか**（コメントだけになっていないか）。`git diff` で実体を見たか。
4. フレーキーを**反復実行で**潰したか（1 回の緑を信じていないか）。

---

# 決着の記録 — 安定緑まで到達

最終形の `deleteDuplicateElements`（コミット時点）:

```go
func deleteDuplicateElements(input []any) []any {
	noIncludeDuplicatedElements := make([]any, 0, len(input))
	m := make(map[any]struct{}, 0)
	for _, e := range input {
		if _, ok := m[e]; !ok {        // map は「既出か」の判定だけ
			m[e] = struct{}{}
			noIncludeDuplicatedElements = append(noIncludeDuplicatedElements, e) // append は入力順
		}
	}
	return noIncludeDuplicatedElements
}
```

- `map` 走査をやめ、**入力スライス順を保ったデデュープ**に変更 → 両関数の出力が「B の出現順」で決定的に。
- 結果: `go test -count=5` で **PASS 140 件・FAIL 0**、`-count=1` を 10 回連続で全 `ok`。**フレーキー撲滅**。
- `gofmt` 差分なし・`go vet` クリーン。途中にあった `else { continue }`（無意味な分岐）も除去済み。

## 全観点の最終ステータス

| 観点 | 内容 | 状態 |
|---|---|---|
| 1 | set 探索がテストされる（呼び間違い） | ✅ |
| 2 | 重複の契約＝デデュープで統一、2 実装が一致 | ✅ |
| 3 | 順序の契約＝「B の出現順」で決定的・テストで固定 | ✅ |
| 4 | `any` の代償（非 comparable で panic）を理解 | ✅（設計意思として明記） |
| 5 | 重複・順序のテストケース追加 | ✅ |
| 6 | フィールド命名統一 | ✅ |
| 7 | **map 順不同によるフレーキー撲滅**（本ドリル最大の山場） | ✅ |

## 残タスク（任意・動作に無関係）

- コメントのタイポ: `:11`「複数福間荒れ」→「複数含まれ」、`linearSearch` の「containingElements の中に…スキップ」コメント（実体のないロジック説明）の整理。
- ケース名 `notIcluded` → `notIncluded`、`StringSlicend` → `StringSliceAnd`。
- `make(map[any]struct{}, 0)` の `, 0` は無害（省略 or `len` でもよい）。

> 総括: 1 個のバグではなく**観点が連鎖する**問題だった。`go test` の表情が
> **「緑の嘘（呼び間違い）→ 赤の嘘（仕様矛盾）→ たまに緑（フレーキー）→ 安定緑」**と 4 段階で変わり、
> 各段階で**違う読み方**を要求してきた。`map`×`reflect.DeepEqual`×フレーキーの三点を一通り踏み抜けたのが本ドリルの収穫。
