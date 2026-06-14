# 実行方法 — コレクションドリルの回し方

このドリルは **`go run`（main関数）ではなく `go test` で回す**。
各ファイルは `package collections`（ライブラリパッケージ）で、`func main()` は1つも作らない。

## なぜ main を作らないか

- ドリルの各ファイルは「関数を書いて、テストで検証する」形。実行可能プログラムではないので `func main()` は不要。
- もし各問を `func main()` で書くと、同じパッケージ内で main が重複して **ビルドエラー**になる。
- だから Go ではテスト駆動が自然で、衝突も起きない。

## 基本コマンド

モジュールルートは `go/`（`go.mod` のある場所）。そこを基準に実行する。

```bash
# 全ドリルのテストを実行
go test ./phase1/basics/collections/

# 詳細表示（どのテストが通った/落ちたか1件ずつ）
go test -v ./phase1/basics/collections/

# 1問だけ実行（テスト関数名を正規表現で指定）
go test -run TestFilterSquareEvens ./phase1/basics/collections/

# go/ 配下すべて（logstat等も含む）
go test ./...
```

リポジトリのどこからでも動かしたいなら `-C` でモジュールルートを指定する。

```bash
go -C /Users/hiroaki/Developer/coding-practice/go test ./phase1/basics/collections/
```

## 進め方（TDDの型）

1. `NN_xxx_test.go` に入力/出力例からテストを書く（赤）
2. `NN_xxx.go` に関数を実装する
3. `go test -run TestXxx ./phase1/basics/collections/` で緑にする
4. 次の問題へ

各問のテスト関数名を `TestXxx` に揃えておくと、`-run` で1問だけ素早く回せる。

## Goテストの書き方と注意点

### 守らないと「テストとして認識されない」ルール

- **ファイル名は `_test.go` で終える。** `go test` はこの接尾辞のファイルだけをテストとして拾う。`01_filter_map_test.go` はOK。
- **テスト関数は `Test` で始め、続く文字は大文字。** シグネチャは `func TestXxx(t *testing.T)`。
  - `TestFilterSquareEvens` → 実行される。
  - `testFilter`（小文字始まり）や `Test_filter`（Testの次が小文字）→ **実行されない**（エラーも出ずに黙って無視されるので注意）。
  - これは Go の「大文字始まり＝エクスポート」の規則と、testing パッケージが名前で関数を探す仕組みによる。
- **`import "testing"` を入れ、引数は `*testing.T`。** これが無いとテストにならない。

```go
package collections

import (
	"reflect"
	"testing"
)

func TestFilterSquareEvens(t *testing.T) {
	got := filterSquareEvens([]int{1, 2, 3, 4, 5, 6})
	want := []int{4, 16, 36}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
```

### 失敗の報告: Error 系と Fatal 系

- `t.Errorf(...)` / `t.Error(...)` … 失敗を記録するが**その後の行も実行を続ける**。複数の検証をまとめて見たいとき。
- `t.Fatalf(...)` / `t.Fatal(...)` … 失敗を記録して**そのテスト関数を即中断**。これ以降を続けると panic する前提のときに使う。
- メッセージは `got %v, want %v` の形が定番。「実際 → 期待」の順で書くと読みやすい。

### スライス・マップの比較は `==` が使えない

- スライスやマップは `==` で中身比較できない（コンパイルエラー、または参照比較）。`reflect.DeepEqual(got, want)` を使う。
- Go 1.21+ なら `slices.Equal` / `maps.Equal` も使える（`import "slices"` / `"maps"`）。

### テーブル駆動テスト（複数ケースをまとめる）

```go
func TestFlatten(t *testing.T) {
	cases := []struct {
		name string
		in   [][]int
		want []int
	}{
		{"basic", [][]int{{1, 2}, {3}, {4, 5, 6}}, []int{1, 2, 3, 4, 5, 6}},
		{"empty", [][]int{}, []int{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := flatten(c.in); !reflect.DeepEqual(got, c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}
```

- `t.Run(名前, func)` でサブテストに分けると、`-run TestFlatten/empty` のように個別実行できる。失敗箇所もケース名で分かる。

### このドリル特有のハマりどころ

- **マップのイテレーション順は非決定**（DRILLS.md の罠表のとおり）。順序を含む期待値と比べるときは、結果をソートしてから比較するか、`reflect.DeepEqual` で**マップ同士**を比べる（マップ比較はキー集合で見るので順序非依存）。
- **float の比較は厳密一致を避ける**。移動平均（14）などは `math.Abs(got-want) < 1e-9` のように誤差許容で比較する。
- 空配列・空マップの戻り値（2, 14 など）は自分で仕様を決め、その決めた値をテストに書く。

### 分岐をテストする

`if/else` や早期 return がある関数は、「**分岐ごとに1ケースを用意して、真の経路と偽の経路を両方通す**」のが基本。

#### 1. 分岐ごとにケースを足す（テーブル駆動が最適）

```go
func TestFoldStats(t *testing.T) {
	cases := []struct {
		name   string
		in     []int
		want   Stats
		wantOK bool // ← 空かどうかの分岐を ok で表現する例
	}{
		{"normal", []int{3, 1, 4, 1, 5}, Stats{Sum: 14, Max: 5, Min: 1}, true},
		{"single", []int{7}, Stats{Sum: 7, Max: 7, Min: 7}, true}, // 初期化分岐の境界
		{"empty", []int{}, Stats{}, false},                        // 空の分岐
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := foldStats(c.in)
			if ok != c.wantOK || got != c.want {
				t.Errorf("foldStats(%v) = %v,%v; want %v,%v", c.in, got, ok, c.want, c.wantOK)
			}
		})
	}
}
```

`if len(x)==0` があるなら「空」と「非空」、ループの初回を特別扱いするなら「1要素」と「複数要素」、というように各経路を1ケースずつ持つ。

#### 2. どの分岐が未テストか可視化する（カバレッジ）

```bash
# カバレッジ率を表示
go test -cover ./phase1/basics/collections/

# 行ごとに通った/通ってないをHTMLで可視化
go test -coverprofile=cover.out ./phase1/basics/collections/
go tool cover -html=cover.out
```

HTML では実行された行が緑、通っていない行が赤で出る。`else` 側や `if` の中身が赤いままなら、その分岐を通すケースが足りない。

> `-cover` は厳密には「ステートメント（行）カバレッジ」で、ブランチカバレッジそのものではない。ただし `if` の各ブロックは別の行になるので、赤い行を潰していけば分岐の取りこぼしはほぼ見つかる。

#### 3. 見落としがちな「隠れ分岐」

- **境界**: 0件 / 1件 / N件、窓が満たない先頭（14番）、キーの有無（`v, ok := m[k]` の ok 両方）。
- **早期 return / continue**: そこに到達するケースとしないケースの両方。
- **エラー経路**: 関数が error を返すなら「正常」と「異常」を両方。

## 手で動かして確認したいとき

- テスト内で `t.Log(...)` を使えば、`-v` 付き実行で出力を確認できる。
- `fmt.Println` での動作確認に慣れているなら、別途 `cmd/playground/main.go`（`package main`）を1つ作って
  collections をインポートして呼ぶ手もある。ただしドリルの目的（反射で書く＋テストで即検証）には `go test` が最短。

## 注意

中身が空のうちは `go test` は「no test files」、または `package collections` 未記述でビルドエラーになる。
最初の1問に `package collections` とテストを書いた時点で回り始める。
