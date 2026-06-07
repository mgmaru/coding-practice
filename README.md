# coding-practice

設計経験を実装力に変えるための学習リポジトリ。約6ヶ月のロードマップに沿って学習する。

## 背景と目的

- **背景**: 設計経験はあるが実装経験が薄い4年目エンジニア。「設計はできるが、自分でコードを書ききれない」を解消する。
- **目標**: ①自分でロジックを組んで実装できる ②AIや他人のコードを読んで品質・セキュリティを判断できる
- **完了の定義（6ヶ月後）**: 第三者（転職先の面接官を想定）に見せられる、テスト・型・CI・CodeQLの揃ったリポジトリが1本あり、AIや他人のコードに対して品質・セキュリティの観点で改善点を指摘できる

## 学習方針（要約）

- 言語はまず **Python に集中**（2言語目以降は Phase 4 から）
- コードは **まず自分で書き、AIはレビュアーとして使う**（先にAIに書かせない）
- 題材（何を作るか）は **自分** が決め、機能チケットと受け入れ条件は **AIにプロダクトオーナー役として決めさせる**（実務の「仕様が与えられる」感覚を再現）→ 詳細: 下記「役割分担」と [docs/ai-workflow.md](docs/ai-workflow.md)
- テスト・型・CI・セキュリティ検証（CodeQL）を学習の早い段階から組み込む
- AIレビューは観点を固定する（ツールが見ない層に集中）→ [docs/review-rubric.md](docs/review-rubric.md)

## 役割分担（誰が何を決めるか）

「何を作るか（題材）」は自分、「次にどの機能を作るか（チケット）」と受け入れ条件は **AI（PO役）**、実装と検証は自分、レビューもAIが担う。学習の課題は2系統ある:

| トラック | 場所 | 題材を決める人 | サイクル |
|---|---|---|---|
| ドリル | `exercises/` | 外部教材（Exercism等）または自分 | 短（模範解答と比較）|
| プロジェクト/ツール | `tools/` | **自分**（機能チケットはAIがPO役で発行）| 長（Issue → 実装 → 自己検収）|

役割分担の詳細とフロー: [docs/ai-workflow.md](docs/ai-workflow.md) ／ レビュー観点: [docs/review-rubric.md](docs/review-rubric.md)

## 進捗

- [ ] **Phase 0**: 環境構築と土台（Git / GitHub / Python環境） — 1週目
- [ ] **Phase 1**: ロジックを組む基礎体力（小さなツールの量産） — 1〜2ヶ月目
- [ ] **Phase 2**: 実務で通用するコード（可読性 / SRP / pytest / 型 / CI） — 2〜3ヶ月目
- [ ] **Phase 3**: セキュリティ（OWASP Top 10 / Bandit / CodeQL） — 3〜4ヶ月目
- [ ] **Phase 4**: 統合プロジェクト（別リポジトリで開発） — 4〜6ヶ月目

各Phaseの詳細とマイルストーン: [docs/roadmap.md](docs/roadmap.md)

## ディレクトリ構成

```
coding-practice/
├── README.md           # このファイル（全体像と進捗）
├── docs/
│   ├── roadmap.md       # 学習ロードマップ本体
│   ├── ai-workflow.md   # AIとの役割分担（AIプロダクトオーナー方式）
│   ├── review-rubric.md # AIレビュアーのレビュー観点
│   ├── learning-log/    # 週次の学習ログ（YYYY-WW.md）
│   └── notes/           # 概念メモ（srp.md, sql-injection.md など）
├── exercises/          # 演習・ドリル
│   ├── 01_basics/
│   ├── 02_refactoring/
│   ├── 03_testing/
│   └── 04_security/    # 脆弱版と修正版をペアで残す（意図的に脆弱なコードを含む）
├── tools/              # 実用ミニツール群（各ツールにtests/を持つ）
└── .github/workflows/  # ci.yml（Phase 2〜）, CodeQL（Phase 3〜）
```

## 運用ルール

1. 作業はブランチ + Pull Request で行い、CI・CodeQLの指摘に対応してから自分でマージする
2. ツール開発はIssue駆動（AIがPO役として発行したチケットをIssueに登録し、仕様の質疑もIssueコメントに残す）
3. `exercises/04_security/` 配下のコードは学習目的で**意図的に脆弱**。実環境では使用しないこと
4. 週1回 `docs/learning-log/` に「やったこと・学んだこと・次にやること」を記録する
5. AIが書いたコードを採用した場合は、その旨と自分が検証した内容をコミットメッセージに残す

## セットアップ

```bash
git clone https://github.com/<username>/coding-practice.git
cd coding-practice
# Python環境（uv使用の場合）
uv sync
```

---

学習開始日: 2026-06-XX / 目標達成日: 2026-12-XX
