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
