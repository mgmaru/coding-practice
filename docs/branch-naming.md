# ブランチ命名規則

ブランチに **何の作業かが一目で分かる名前**を付けるための規則。
ブランチを切る・push する **手順** は [pr-workflow.md](./pr-workflow.md)（②④）にある。ここは「**どう名付けるか**」だけを定義する。

> このリポジトリの **type 語彙（feat / fix / docs …）の正本はこのファイル**。PRタイトル（[pr-writing.md](./pr-writing.md)）も同じ語彙を使う。

## 形式

```
<type>/<短い説明>
```

例: `docs/branch-naming`, `feat/csv-export`, `fix/empty-input-crash`

## type 一覧

[Conventional Commits](https://www.conventionalcommits.org/) に揃える。**PRタイトルの先頭・Issueラベルもこの語彙を共有する**（3つで同じ言葉を使うと一貫する）。

| type | 使うとき |
|------|---------|
| `feat` | 機能の追加 |
| `fix` | バグ修正 |
| `docs` | ドキュメントのみの変更 |
| `refactor` | 挙動を変えない内部改善 |
| `test` | テストの追加・修正 |
| `chore` | 設定・依存・雑務（`.gitignore`、ツール設定など）|
| `perf` / `ci` / `build` | 性能改善 / CI設定 / ビルド周り（必要になったら）|

## 命名ルール

- **小文字 + ハイフン区切り（kebab-case）**: `feat/csv-export`（`featCsvExport` や `feat/csv_export` ではない）
- **短く具体的に**（〜50字目安）: `feat/user-auth` は可、`feat/auth` は曖昧、`feat/add-the-new-user-authentication-flow` は冗長
- **スラッシュは type との区切り1つだけ**: `feat/csv-export`（`feat/csv/export` にしない）

## 任意: Issue番号を入れる

このリポジトリは **Issue駆動**（[ai-workflow.md](./ai-workflow.md)）なので、対応するIssueがあれば番号を頭に付けると追跡性が上がる（GitHubがブランチ↔Issueを自動リンクする）。

```
<type>/<Issue番号>-<短い説明>
```

例: `feat/12-csv-export`（Issue #12 に対応）

## つまずきポイント

| 症状 | 原因 | 対処 |
|------|------|------|
| `src refspec ... does not match any` | push時の名前が実在のブランチ名と違う | `git branch --show-current` で確認してから push |

> `git checkout -b` での作成、`git push -u` での upstream 設定など **操作の詳細は [pr-workflow.md](./pr-workflow.md)** を参照。
