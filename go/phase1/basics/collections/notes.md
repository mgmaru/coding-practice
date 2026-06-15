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
