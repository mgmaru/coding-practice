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

## 手で動かして確認したいとき

- テスト内で `t.Log(...)` を使えば、`-v` 付き実行で出力を確認できる。
- `fmt.Println` での動作確認に慣れているなら、別途 `cmd/playground/main.go`（`package main`）を1つ作って
  collections をインポートして呼ぶ手もある。ただしドリルの目的（反射で書く＋テストで即検証）には `go test` が最短。

## 注意

中身が空のうちは `go test` は「no test files」、または `package collections` 未記述でビルドエラーになる。
最初の1問に `package collections` とテストを書いた時点で回り始める。
