# collections

## ドキュメント概要
`collections`の実装で気付いたことを記載する。

## 反省
### 2026/6/14
#### やったこと
- `01_filter_map.go`実装
- `01_filter_map_test.go`実装
#### わかったこと
- `=`でスライス同士の比較はできない。
    - 標準モジュール`reflect.DeepEqual()`で比較することができる。
- スライスの定義に関して、`var squaredEvens []int`よりも`squaredEvens := []int{}`の方が良い。
    - `var squaredEvens []int`だと、`nil slice`が返ってしまう。
    - `squaredEvens := []int{}`にすれば、`empty slice`を返すことができる。
- ただし、`reflect.DeepEqual(int{}, []int(nil))`は、`false`になることに注意。
    - `[]int{}`は、空スライス。
    - `[]int(nil)`は、nilスライス。`nil`と同じ。
- 注意（重要）：
    - `func filterSquaredEvens(input []int) `の場合
        - `input`に、`[1, "apple", ...]`みたいなスライスは入ってこないので、これに対する型チェックはしなくて良い。
        - **ただし、`input`には、`nil`は入ってくるので、これに対してはガードを設ける。**
- `==`での比較は、ポインタで返されている場合に`false`になるので注意。
    - エラーの比較には、`errors.Is`などを使用する。
---
### 2026/6/15
#### やったこと
- `02_fold_stats.go`実装
- `02_fold_stats_test.go`実装
#### わかったこと
- 入力スライスの最初の要素について、`if i == 0 { sum = num, max = num, max = num }`としていたが、余計な分岐でよくない。
    - `sum=input[0], max=input[0], min=input[0]`のように、初期値を代入してあげる。
    - その次の要素に対して、合計、最大値、最小値を計算するループ処理を行えば良い。`for _, num := range input[1:]`
    - 余計な分岐をすると、テストの対象になりうるので注意。
- テストする関数で文字列を返す場合は、`fmt.Println()`だと、テスト側でその文字列を受け取ることができない。
    - `fmt.Fprintf(w, "文字列")`のように、文字列を、変数`w io.Writer`に渡してあげることで、テスト側でその文字列を取り出すことができる。
    - 受け取る側では、`buf byte.Buffer`で受け取り、それを文字列へ変換`buf.String()`で、文字列を受け取り、比較を行うことができる。
    - むやみに、`fmt.Println()`は使用しない。
- 空スライスの判定は、`reflect.DeepEqual(input, []int{})`ではなく、`len(input) == 0`の方が軽量。
    - ただし、**`input`がnilスライスの場合も`len(input) == 0`となる**ので注意。
    - 逆に言えば、空スライスとnilスライスを区別しないでガードするならば、`len(input) == 0`のガードのみで良い。
- **組み込み関数（`max`,`min`など）を変数名に使用しないようにする**。
#### 実装時間
- 1時間53分
---
### 2026/6/16
#### やったこと
- `03_intersection.go` 実装
- `03_intersectiontest.go` 実装
#### わかったこと
- set型を使用すると、高速で要素が含まれるかどうかを判定することができる。
- **`set`とは、「重複しない要素の集まり」**。
    - 性質1：ある要素が入っているかを高速で判定することができる。
    - 性質2：同じ要素を２個入れても１個分にしかならない。
- Goでset型はない。
    - 上の`set`を満たすGoのデータ構造は`map[キー]値`。**`map`のキーは重複しない。**
    - `map`のキーには要素の集合を実装する。
    - 問題は値に何を実装するかだが、値には**struct{}**つまり、**空のstruct（フィールドが０個）**を実装する。
    - なぜ値に空のstructを実装するのか
        - **空のstructは1バイトも消費しない**。つまり、値自体に意味がない。
    - 以上より、**余計な意味を持たずに、要素の存在自体を表すことができる。**
- Goの`map`のキーの型は、比較可能な値（comparable）であればなんでも良い。
    - 比較可能な値（comparable）：
        - `int`/`int64`/`float64`
        - `bool`
        - `string`
        - ポインタ
        - 比較可能なフィールドだけで構成された`struct`
        - 比較可能な要素型の配列`[3]int`など
        - インターフェース型
    - キーに使用できない型
        - スライス`[]int`
        - `map`
        - 関数
- **`map`の注意：`for key := range map`で`key`を取り出すと、その`key`の順序はランダム**になる。
    - なので、要素の順序を気にする場合は、`for`で`map`のキーを取り出し処理はやめた方が良い。
    - あくまで、`map`のキーは、要素が存在するかどうかの判定に使用した方が良さそう。
    - もし、`for`でキーを取り出す場合は、取り出した後にソートするのが良さそう。
- **スライスの要素の順番も考慮する必要**がある。 ←超重要！！
#### いつか実施したい課題
- `[]any`のスライス（例えば、`[]any{1, "apple", "banana", 5, ...}`）などの様々な型のデータが格納されているスライスを昇順、降順に並べる。
#### 実装時間
- 約3時間
---
### 2026/6/17
#### やったこと
- `04_count_frequency.go` 実装
- `04_count_frequency_test.go` 実装
#### わかったこと
- 入力`["apple", "banana", "apple", "cherry", "banana", "apple"]`に対して、出力`{"apple": 3, "banana": 2, "cherry": 1}`が出力されるには、順序を気にする必要がある。
    - `map`はハッシュ順に格納される（実質ランダム）なので、例えば、以下のような実装だと、順番が保証されない。
```go
func CountElements(input []string) map[string]int {

	m := make(map[string]int, 0)

	for _, v := range input {
		m[v]++ // 順番が保証されない。
	}
	return m
}
```
- ただし、**Goでは出力が`map`のデータ構造かつその順序を保証することはできない。**
    - なので、今回は、`[{element:"apple", count:3}, {element:"banana", count:2}, ...]`のような出力とした。
    - **仕様を決めることが大事。**
- `make`は、スライスとmap、チャネルで指定の仕方が違う。
    - スライス：`make([]T, len, cap)` / `len` ：初期サイズ、　`cap`：容量
    - map：`make(map[K]V, sizeHint)`/ `sizeHint`：容量
    - チャネル：`make(chan T, buffer)`/ `buffer`：容量
- **文字コード(ASCII/Unicodeの数値)の大小**：
    - 大文字`A-Z`：65〜90
    - 小文字`a-z`：97〜122
    - 例えば、`["apple", "Apple", "banana"]`を昇順にソート　→ `["Apple", "apple", "banana"]`になる。
#### 実装時間
- 2時間54分
---
### 2026/6/27-28
#### やったこと
- `05_group_by.go` 実装
- `05_group_by_test.go` 実装
#### わかったこと
1. 構造体のフィールドをキーとして指定する方法がわからなかった。
2. マップの比較の方法がわからなかった。
    - **単純な`==`では比較できないことがわかった。**
    - `reflect.DeepEqual`で比較することにした。
    - ~~ただし、`refrect.DeepEqual`ではスライスの並び順までは判定してくれないみたいなので、ここを改善する必要あり。~~
        - **`reflect.DeepEqual`は、スライスの要素の順序も判定してくれる。**
#### レビュー
- 自分の実装では、`groupBy`関数に`key string`を渡していた。
    - しかし、**関数を`groupBy`関数に渡すのが良い。**
    - 呼び出し元で、キーおよび値を注入してあげる。
#### 総合反省
- 関数を渡してあげるという発想がなかった。これは持っておく必要がある。
    - Goでは文字列から、構造体のフィールドにアクセスするのは難しい。
#### 実装時間（初期実装＆レビュー修正）
- 4時間40分
---
### 2026/6/29
#### やったこと
- `06_dedupe_ordered.go` 実装
- `06_dedupe_ordered_test.go` 実装
#### わかったこと
- **`map`のキーにはcomparable（「`==` と `!=` で比較できる型）が入らないといけない。**
- comparableではないもの
  - slice
  - map
  - 関数
- ジェネリクスを使用した。
- **`[]any{1, int64(1), 1.0},`これらは違う型の要素**であるので、mapのキーは別のものとして扱われる。-> このテストケースも含めるべき。
#### 実装時間
- 2時間31分
---
### 2026/6/29（07_lookup）
#### やったこと
- `07_lookup.go` 実装
- `07_lookup_test.go` 実装
#### わかったこと
- **`m[k]` は、キーが無くてもpanicせず「値型のゼロ値」を黙って返す。**
    - 「在る／無い」を区別したいときは、`val, ok := m[k]` の `ok` で受ける。
    - 単値の `m[k]` だけだと、欠損キーでゼロ値（`string`なら`""`、`int`なら`0`）が混ざる。
    - 今回これを見落として、`usersEmpty`（usersが空）のテストが赤になった（ゼロ値`["","",""]`が返っていた）。
    - **`DRILLS.md` のGo罠表の `v, ok := m[k]` がまさにこれ。** このドリルの主眼。
- **欠損キーの「契約」を先に決める。** 引くキーが表に無いとき、どうするか:
    - ① skip（在ったものだけ積む。今回採用）
    - ② ゼロ値を積む（位置対応を保つ。ただしテストの`want`もゼロ値に揃える必要あり）
    - ③ エラー/panic
    - **どれを選んでも良いが、実装・テスト・コメントを同じ契約に揃えること。** 今回は実装(ゼロ値)とテスト(skip)が食い違っていたのが赤の原因。
    - skipを選ぶと **出力の長さが入力の長さと一致しなくなる**（欠損ぶん短くなる）ことは承知しておく。
- **mapの「値の型」にcomparableは要らない。** 必要なのはキーの型だけ。ジェネリクスは `[T comparable, K any]` でよい（Kまでcomparableにすると不必要に型を狭める）。
- **結果スライスのcapは `len(input)`。** `len(users)` だと、引くキー列inputに重複があると足りない（例 `input=[2,1,1,2,3]`）。
- **バグを「二通り」で直すと、片方が冗長分岐になりやすい。** comma-okを入れたのに `if len(users)==0` の早期returnも残してしまい、後者が冗長だった（comma-okだけで同じ結果）。`04`/`05` と同じ「無くても同じ結果の分岐は消す」。直したら「片方で足りないか」を問う。
- **索引化（このドリルの狙い）**：繰り返し引くなら先にmapへ索引化する。線形探索 O(n×m) → 索引化 O(n+m)。
#### 実装時間
- 約3時間
---
### 2026/6/30
#### わかったこと
- スライスのソートには、`sort.Slice`と`sort.SliceStable`の２種類がある。
  - `sort.SliceStable`：`sort.Slice`の安定ソート版で、`less` が等しいと判定した要素同士の元の順序を保持します。一方`sort.Slice` は安定性を保証しないので、等しい要素の並びが入れ替わる可能性があります。
    - 例：`people := []Person{{"Alice", 30},{"Bob", 25},{"Carol", 30}}`を年齢の降順でソートする。
    - `sort.Slice`:ソート後 -> `people := []Person{{"Bob", 25},{"Alice", 30},{"Carol", 30}}` or `[]Person{{"Bob", 25},{"Carol", 30},{"Alice", 30}}` -> 年齢が同じ場合に元のスライスの並びにならない可能性がある。
    - `sort.SliceStable`:ソート後 -> `people := []Person{{"Bob", 25},{"Alice", 30},{"Carol", 30}}` -> 年齢が同じ場合に、元のスライスの順番である`Alice` -> `Carol`の順番が保たれる。
- **(重要)「破壊」と「非破壊」の意識をしてコーディングすることが大事。**
  - **破壊的**：元のデータを変更する。
  - **非破壊的** ：元のデータにはを加えず、新しいデータを結果として返す。
  - 選択の軸
    1. 元データを使用するかどうか：非破壊を選ぶのは、元のデータをこの先も別の用途で使う場合。破壊を選ぶのは、元のデータをもう使わない、あるいは「変換後の状態こそが欲しい結果」の場合。
    2. パフォーマンスとデータサイズ：非破壊はコピーを作るので、メモリと時間のコストがかかります。要素数が多い、頻繁に呼ばれる、といった場面では無視できません。
      - 迷ったらまず、非破壊で書いて、パフォーマンスが問題になったら、破壊で書く感じで良い。
    3. その関数の「立場」：自分が**他人に使われる関数（特に公開API）を書いているときは、デフォルトを非破壊に寄せるのが安全です。理由は、引数で渡したスライスが勝手に書き換わるのは、呼び出す側にとって予想外の挙動（驚き）**だから。
- 退行テストを理解するべき（今後で良い）。
#### 実装時間
- 2時間36分
---
### 2026/7/1
#### わかったこと
- `map of map`のキーでソートする方法がわからなかった。
  - でも、そもそも`map`構造は順序を保証できないので、キーでソートすることはできない。
  - もし、ソートしたいのであれば、スライスにする必要がある。
- `map`同士を比較する時に、`==`では比較できないので、`reflect.DeepEqual`を使用。
  - `DeepEqual`で`map[string]int{"a":1, "b":2}`と`map[string]int{"b":2, "a":1}`を比較すると、`true`になる。
  - **`DeepEqual`は`map`の並び順までは見ない。**
- `map`のネストが多くなった時に、`{}`の数がわからなくなる。 -> ただし、書いてるとなれる。
- `map`のキーに何が入りうるかは知っておくべき。
#### バグ
- 修正前コード
```go
	n := make(map[string]int, len(input)) // 日付ごと　-> キー：レベル　値：出現回数
	out := make(map[string]map[string]int, len(input))
	for key, vals := range m { // key:日付　vals:レベルのスライス // 注意：map mからforで取り出しているので、並び順が保証されていない。 -> ソートが必要
		for _, v := range vals {
			if _, ok := n[v]; ok { // もし、mの中にLevelがあったら、出現回数をカウントアップ
				n[v]++
			}
			// もしなかったら、キー(日付)と初期値を加える
			n[v] = 1 // 初期値
		}
		out[key] = n // キーを日付として、値にnを代入
	}
```
- 修正後コード
```go
	out := make(map[string]map[string]int, len(input))
	for key, vals := range m { // key:日付　vals:レベルのスライス // 注意：map mからforで取り出しているので、並び順が保証されていない。 -> ソートが必要
		n := make(map[string]int, len(input)) // 修正：初期化する
		for _, v := range vals {
			if _, ok := n[v]; ok { // もし、nの中にLevelがあったら、出現回数をカウントアップ
				n[v]++
				continue // 修正：ここを忘れると、ifの後のn[v]=1が実行されて、レベルのカウントが毎回1になる...
			}
			// もしなかったら、キー(レベル)と初期値を加える
			n[v] = 1 // 初期値
		}
		out[key] = n // キーを日付として、値にnを代入
	}
```
- バグ理由①：`for key, vals`ごとに、つまり`key`ごとに変数`n`を初期化しないと、`key`を跨いでレベルをカウントしてしまう。 
- バグ理由②：`continue`がないと、レベルをカウントしてから、`n[v]= 1`が実行されてしまい、レベルのカウントが全て1になる。
#### レビュー修正（わかったこと）
1. *`map`にキーと値を代入する時に、キーが存在するかどうかはいらない。**
- 冗長なコード
```go
for _, v := range vals {
    if _, ok := n[v]; ok { // もし、nの中にLevelがあったら、出現回数をカウントアップ
        n[v]++
        continue // 忘れずに
    }
    // もしなかったら、キー(レベル)と初期値を加える
    n[v] = 1 // 初期値
}
```
- 正しいコード
```go
for _, v := range vals {
    n[v]++
}
```
2. `if _, ok := m[e]; !ok { ... } `は存在判定が目的の時のみ使う。 -> キーの出現回数をカウントする時には不要。
3. コードの簡素化
- 自分の実装
```go
m := make(map[string][]string, len(input)) // 第1パス：日付ごとにレベルを "スライスに溜める"（＝05 の group_by）
for _, s := range input {
    m[s.Date] = append(m[s.Date], s.Level)
}
out := make(map[string]map[string]int, len(input)) // 第2パス：溜めたスライスを数える（＝04 の count）
for key, vals := range m { ... }
```
- Claudeの実装
```go
func aggregateDateStatus(input []Status) map[string]map[string]int {
    out := make(map[string]map[string]int, len(input))
    for _, s := range input {
        if out[s.Date] == nil {          // 内側 map は書き込む前に必ず初期化（下の注意）
            out[s.Date] = make(map[string]int)
        }
        out[s.Date][s.Level]++           // 欠損レベルはゼロ値 0 → ++ で 1
    }
    return out
}
```
- **`out[s.Date][s.Level]++`のように書くことで、`map`のネストした値にアクセスできる。** -> 初めて知った！！
- **(重要)`map`の書き込みに対しては、必ず初期化しないと`nil`マップへの書き込みとなって`panic`になる。**
  - ここでは、キー`Date`に対する値が`map`なので、必ず初期化する。
#### 知識
##### `Entry`（エントリ）とは
- `entry` は「記入・項目・登録されたもの」という普通の英単語。ソフトウェアでは **「並んでいるデータの1件分」** を指すごく一般的な言葉。
  - ログの1行 = a log **entry**
  - 辞書の1項目 = a dictionary **entry**
  - 名簿の1件 = an **entry**
- 「たくさん並んでいるものの、そのうちの1つ」というニュアンス。
- 今回の `Status` 構造体は「ログが何行もある中の**1行分**（日付＋レベル）」＝まさに **log entry（ログ1件）**。
  - なので `Status`（状態）より `LogEntry`（ログ1件）のほうが「これは入力の1レコードだ」と正確に伝わる。
- 似た用語の使い分け:

| 語 | 意味 | 例 |
|---|---|---|
| **entry** | 並んだデータの1件（記入されたもの） | ログ1行、辞書の1項目 |
| **record** | 1件のデータ（DB由来。ほぼ同義） | DBの1行、`LogRecord` |
| **item** | 一覧・配列の1要素（もっと汎用） | リストの1個 |
| **element** | 集合・配列の1要素（数学寄り） | `04` の `CountElement` |

- `LogEntry` でも `LogRecord` でも意味はほぼ同じ。Go／ログ処理の界隈では **`entry`** がよく使われる（標準ライブラリの `log` など）。
- 補足: `map` の文脈では `entry` に**もう一つ**の意味がある —— 「mapの中の**キー＋値のペア1組**」も entry と呼ぶ（`for k, v := range m` で回している1組がそれ）。今回の改名案は前者（ログ1件）の意味。
#### 反省
- `Go`における`map`は順序を保証しないことは常に頭に入れておく。
- **データ構造ごとの性質（何が得意で不得意か）理解したほうが良い**かも。
  - 例１：`map`はキーが重複しないので、要素の存在の判断に使える。また、要素の数をカウントできる。等
  - 例２：`slice`は、並びを保証しやすい。ソートしやすいデータ。等
#### 実装時間
- 合計：3時間54分
  - 初期実装：3時間
  - 修正：54分
---
### 2026/7/10
#### やったこと
- 10_flatten.go実装
- 10_flatten_test.go実装
#### わかったこと
- リストの展開`...`を使用することで、`for`で要素を１つずつ取り出して`append`しなくても、１発で`append`できる。
1. 自分の実装
```go
for _, numList := range input {
		for _, num := range numList { // numListが空の場合、appendは実行されない
			oneDimList = append(oneDimList, num)
		}
	}
```
2. 展開`...`を使用した実装
```go
for _, numList := range input {
	oneDimList = append(oneDimList, numList...)
	}
```
- 平滑化後の容量の設定
1. 最初の自分の実装では、平滑後のリストの容量を指定していなかった。
```go
func smoothTheList(input [][]int) []int {
	oneDimList := make([]int, 0) // capはlen(input)では足りない ←容量指定なし
	for _, numList := range input {
			oneDimList = append(oneDimList, numList...)
```
2. しかし、冷静に考えると、平滑化後のリストの最大サイズは計算できる。
```go
func flattenList(input [][]int) []int {
	// 修正：平滑化後のスライスのサイズを計算 ←ここ！
	maxFlattenListSize := 0
	for _, l := range input {
		maxFlattenListSize += len(l)
	}
	oneDimList := make([]int, 0, maxFlattenListSize) // capはlen(input)では足りない -> 修正：最初に平滑化後のリストを計算して、capに指定
	for _, numList := range input {
		oneDimList = append(oneDimList, numList...) // 修正：リストを展開してappend
	}
```
- 容量を設定することで、無駄な容量を確保しなくてよくなり、メモリを効率的に使用できる。
#### 実装時間
- 合計：1時間30分
---
### 2026/7/11
#### やったこと
- 11_inner_join.go実装
- 11_inner_join_test.go実装
#### 間違えたところ
- スライスの容量の見積もりを誤った。
#### 実装時間
- 合計：1時間27分
---
### 2026/7/15
#### やったこと
- 12_invert.go実装
- 12_inver_test.go実装
#### わかったこと
- 今回、私の実装では、`[]string`の要素に対して`sort.SliceStable`を使用してソートした。
  - しかし、これは過剰だという指摘を受けた。
  - 考えてみると、今回はそもそも1つの要素だけを比べるだけだから、同じ場合は同じ場合で良い。
  - 例えば、`people := []Person{{"Alice", 30},{"Bob", 25},{"Carol", 30}}`を年齢の降順でソートする場合、以下の２通りがソート後に考えられる（ソートが安定しない）。
    - `[]Person{{"Bob", 25},{"Alice", 30},{"Carol", 30}}` 
    - `[]Person{{"Bob", 25},{"Carol", 30},{"Alice", 30}}` 
    - 年齢でソートすると、AliceとCarolのソートが安定しない。
    - このような場合に、`sort.SliceStable`を使用することで、ソート前の並びのAlice->Carolに安定させることができる。
  - 今回の1つの要素をソートする場合では過剰。
- `slices.Sort`と`sort.StableSlice`では、`sort.StableSlice`の方が計算量が多い（lognだけ多い）。
  - なので、安定性を求める必要がない場面では、`slices.Sort`を使用することが良い。
- `sort.Slice`と`slices.Sort`は違う！！
  - `slices.Sort`：数値・文字列を昇順でソート
  - `sort.Slice`：`func Slice(x any, less func(i, j int) bool)`なので、なんでも渡せる。
    - less は「インデックス i の要素を j より前に置くべきか」を真偽で返す関数
    - less(i, j) が true を返す ＝「インデックス i の要素を j より前に置く」
      - return s[i] < s[j] … 「i の値が小さいとき i を前に」→ 小さいものが先 → 昇順（小→大）
      - return s[i] > s[j] … 「i の値が大きいとき i を前に」→ 大きいものが先 → 降順（大→小）
- 文字列の昇順 ：`<` -> 大文字（A,B,...）、小文字（a,b,...）の順にASCIIの数字が小さくなっていく。
- `input := map[string]int{}`の時
```go
for key, val := range input { // mapからキーを取り出しているので、順序を保証しない。
		inverted[val] = append(inverted[val], key)
	}
```
  - 上の`inverted[val]`でゼロ値を返すのでは？と一瞬疑った。
  - しかし考えてみると、そもそもinvertedのキーとなる**valが空（nil）であるから、ゼロ値は返らない。** -> **修正：ループが0回だから、空のmapが返る（append処理は発生しない）。**
  - もし、valがnilでなくて存在する値であったらゼロ値を返す。
#### 反省
- ASCIIに慣れた方が良い。
#### 実装時間
- 合計：約2時間
---
### 2026/7/17
#### やったこと
- 13_cumulative_sum.go実装
- 13_cumulative_sum_test.go実装
#### わかったこと
- 当初、課題に書かれていた「前の状態を持ち越しながら走査する。」の意味がわからなかった。
  - 聞いたところ、前の状態は、その時点での和であることがわかった。
  - なので、修正では、計算結果を格納する変数`sum` を定義した。
#### 議論
- 自分の実装
```go
func calcCumulativeSum(input []int) []int {
	output := make([]int, 0, len(input))
	for i, num := range input {
		if len(output) == 0 { // i == 0でも良い？
			output = append(output, num)
			continue
		}
		output = append(output, output[i-1]+num)
	}
	return output
}
```
- 修正後の実装
```go
func calcCumulativeSum(input []int) []int {
	output := make([]int, 0, len(input))
	sum := 0
	for _, num := range input {
		sum += num
		output = append(output, sum)
	}
	return output
}
```
- 自分の実装は、`sum`という変数を置いてないので、デバッグ容易性にかけていると感じた。
  - 実際、AIに聞いてみても同意する回答だった。
  - **`sum`を置くことでデバッグ容易性が向上。**
    - 理由：`sum`に「その時点の累積」という**名前と居場所**ができたので、途中経過を1変数で追える。旧実装は途中状態が`output`スライスの中に散っていて、`output[i-1]`という**式**を目で追わないと「今の累積」がわからなかった。
- あと、分岐を減らすことができている。
- **ただし、この修正の芯は「デバッグ容易性」ではない。** それは副次的な効果で、上流に本当の効き目が2つある。
  - ① **狙い「前の状態を持ち越す」に一致した**（`sum`が持ち越す状態そのもの＝このドリルの正解の型）。
  - ② **`output[i-1]`の潜在バグを除去した**。旧実装は入力添字`i`で出力位置を引いており、今は1:1だからたまたま一致するが、filter/skipが1つ入ると`i≠len(output)`で壊れる。`sum`ならそもそも出力位置を計算しないので癒着が起きない。
  - 優先順位：① 狙い一致 → ② 潜在バグ除去 → ③ 分岐消去 → ④ デバッグ容易性。
- ~~`sum`という一時変数を使用することで、メモリを消費するが所詮`int`なのでそこまでメモリの消費は大きくないと感じた。~~ → **訂正：これはトレードオフではない。** 主メモリは旧実装と同じ`output`（`len(input)`個）で、`sum`は`int`1個がスタックに載るだけ（追加ヒープはゼロ）。むしろ`output[i-1]`の毎回の境界チェック＋間接参照が消えるぶん、わずかに速くなる可能性すらある。**メモリを払う代わりの改善ではなく、ほぼ純粋な改善。**
#### 実装時間
- 合計：約1時間
---
### 2026/7/19
#### やったこと
- 14_moving_average.go実装
- 14_moving_average_test.go実装
#### わかったこと
1. 繰り返し処理の書き方
- `for i := 0; i < 10; i++` -> この書き方は古いらしい。
- `for i := range 10` -> この書き方が推奨されているらしい（Go 1.22〜）
2. スライディングウィンドウの実装
  1. 自分の実装
```go
	for i, _ := range input { // iを1個ずつずらしていく
		if i+number > len(input) { // 平均の計算対象となる数が要素を超えた場合は平均計算終わり
			return movingAverageResult
		}
		// number分、inputから対象となる数を取り出して平均を計算する
		targetNumsSum := 0
		for j := range number { // number分要素を取り出す
			targetNumsSum += input[i+j]
		}
		movingAverage := targetNumsSum / number
		movingAverageResult = append(movingAverageResult, movingAverage)
	}
```
- 自分の実装では、内側のループで毎回合計を計算しているので、計算量は`O(n)` -> n(窓の要素数)が大きい場合、計算量が膨大になる。
  2. AIの回答
```go
	windowSum := 0
	for i := range windowSize {
		windowSum += input[i]
	}
	result = append(result, windowSum/windowSize)

	for i := windowSize; i < len(input); i++ { // 修正 i：入ってくる数 i-windosSize：出ていく数
		// 補足：もし、windowSizeがlenより大きかった場合、ループ処理には入らない。
		windowSum += input[i] - input[i-windowSize]
		movingAverage := windowSum / windowSize
		result = append(result, movingAverage)
	}
	return result
```
- AIの回答では、どんなに窓の要素が多くても、計算量は`O(1)`
#### 反省
- 今回の問題は、数学的な要素が大きかった。
- やはり、数学的なアルゴリズムの処理が苦手なことがわかった。
- 基本的なアルゴリズムや、競プロを練習する必要がある。
#### 実装時間
- 合計：3時間40分
  - 初期実装：1時間53分
  - 修正：1時間47分
---
