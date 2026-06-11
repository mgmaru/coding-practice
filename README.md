# coding-practice

設計経験を実装力に変えるための学習リポジトリ。約6ヶ月のロードマップに沿って学習する。

## 背景と目的

- **背景**: 設計経験はあるが実装経験が薄い4年目エンジニア。「設計はできるが、自分でコードを書ききれない」を解消する。
- **目標**:
  - ① 自分でロジックを組んで実装できる
  - ② AIや他人のコードを読んで品質・セキュリティを判断できる
- **実装の姿勢**: 「動けばよい」で止めず、**各実装に根拠（なぜそう書くか）を説明できる**ことを基準にする。ドリルもプロジェクト/ツールも必ずAIレビュー（[docs/review-rubric.md](docs/review-rubric.md)）を通すのは、この「説明できる実装」を担保するため。
- **完了の定義（6ヶ月後）**: 第三者（転職先の面接官を想定）に見せられる、テスト・型・CI・CodeQLの揃ったリポジトリが1本あり、AIや他人のコードに対して品質・セキュリティの観点で改善点を指摘できる

## 学習方針（要約）

- **まず1言語を深く**（言語は問わない。テスト・型・セキュリティのツールが揃う **Python を推奨**。業務等で使う言語があればそれでよい）。複数言語の同時並行は避け、2言語目は1言語が固まってから（Phase 4以降が目安）
- コードは **まず自分で書き、AIはレビュアーとして使う**（先にAIに書かせない）
- 題材（何を作るか）は **自分** が決め、機能チケットと受け入れ条件は **AIにプロダクトオーナー役として決めさせる**（実務の「仕様が与えられる」感覚を再現）→ 詳細: 下記「役割分担」と [docs/ai-workflow.md](docs/ai-workflow.md)
- テスト・型・CI・セキュリティ検証（CodeQL）を学習の早い段階から組み込む
- AIレビューは観点を固定する（ツールが見ない層に集中）→ [docs/review-rubric.md](docs/review-rubric.md)

## 役割分担（誰が何を決めるか）

「何を作るか（題材）」は自分、「次にどの機能を作るか（チケット）」と受け入れ条件は **AI（PO役）**、実装と検証は自分、レビューもAIが担う。学習の課題は2系統ある:

| トラック | 目的 | 場所 | 題材を決める人 | 学習サイクル |
|---|---|---|---|---|
| ドリル | **道具**を単体で習得（部品）| `phaseN/<topic>/` | 外部教材（Exercism等）または自分 | 短: 解く → AIレビュー → 参考コードと比較 |
| プロジェクト/ツール | 道具で**設計込みの実装**を完走（組み立て）| `phaseN/<tool>/`（作成フェーズに置く）| **自分**（チケットはAIがPO役で発行）| 長: チケット → 実装 → AIレビュー → 自己検収 |

> **workflow（PR → AIレビュー → 自分でマージ）は両トラック共通**で、AIレビューは必須（観点: [docs/review-rubric.md](docs/review-rubric.md)）。2つを分けるのは上の **目的** と **題材の出どころ**。「学習サイクル」列は*学習の回し方*で、git運用（PR）とは別レイヤ。目的・サイクルの詳細定義は [docs/ai-workflow.md](docs/ai-workflow.md)。

> **「フェーズ」と「トラック」は別軸。** フェーズ（縦＝いつ/どのレベル）は進んでいくが、上の2トラック（横＝訓練の種類）は**毎フェーズ両方を並行**する。ツールの規模差ではなく、**短サイクルで反復するか（ドリル）／長サイクルで実務を完走するか（プロジェクト）** の違い。配分目安は ドリル3〜4割 / プロジェクト6〜7割。

| フェーズ ＼ トラック | ドリル（短・反復）| プロジェクト/ツール（長・完走）|
|---|---|---|
| **Phase 1** | 文法ドリル | 小ツールを自力で作る（`log_analyzer` 等）|
| **Phase 2** | リファクタ/テストのドリル | ツールにテスト・型・CIを足す（AI-PO方式は後半開始）|
| **Phase 3** | OWASP「壊して直す」ドリル | ツールにセキュリティ点検 |
| **Phase 4** | （圧縮して維持）| 統合プロジェクト（別リポ）|

> 各フェーズフォルダの中に両トラックが同居する。例: `phase1/basics/`（ドリル）＋ `phase1/log_analyzer/`（ツール）。

役割分担の詳細とフロー: [docs/ai-workflow.md](docs/ai-workflow.md) ／ レビュー観点: [docs/review-rubric.md](docs/review-rubric.md)

## 進捗

- [x] **Phase 0**: 環境構築と土台（Git / GitHub / 言語環境。多言語対応の器も整備済み） — 1週目
- [ ] **Phase 1**: ロジックを組む基礎体力（小さなツールの量産） — 1〜2ヶ月目
- [ ] **Phase 2**: 実務で通用するコード（可読性 / SRP / pytest / 型 / CI） — 2〜3ヶ月目
- [ ] **Phase 3**: セキュリティ（OWASP Top 10 / Bandit / CodeQL） — 3〜4ヶ月目
- [ ] **Phase 4**: 統合プロジェクト（別リポジトリで開発） — 4〜6ヶ月目

> 各Phaseのチェック ＝「**フェーズ移行レビュー**」を通過した印（自己検収 ＋ AI検収 → 自分が昇格判断）。承認者を外部に置かず、自己バイアスを防ぐチェックポイントとして運用する。

各Phaseの詳細・完了の定義（DoD）・移行レビューの手順: [docs/roadmap.md](docs/roadmap.md)

## ディレクトリ構成

> **構成（器）は多言語対応、学習計画（中身）は1言語集中のまま。** リポジトリは Python / TypeScript / Go / Rust の4言語を置ける「器」にしてあるが、学習は当面1言語（**Python推奨**）に集中し、2言語目は Phase 4 以降に始める（[docs/roadmap.md](docs/roadmap.md) の言語戦略）。器を多言語にする理由・各言語の環境構築手順は [docs/multi-language-setup.md](docs/multi-language-setup.md) を参照。

**配置は3軸 ＝ 言語 / フェーズ / トラック。** 言語をトップに置き（`python/` の中だけが Python の世界…「1言語 = 1ツールチェーンの根」）、その中に **フェーズ**（`phaseN/`＝いつ着手したか）、さらに **トラック**（`basics/`＝ドリル、`log_analyzer/` 等＝ツール）が入る。

```
coding-practice/
├── README.md                # このファイル（全体像と進捗）
├── docs/                    # ★横断ドキュメント（言語・フェーズに紐づかない）
│   ├── roadmap.md               # 学習ロードマップ本体
│   ├── ai-workflow.md           # AIとの役割分担（AIプロダクトオーナー方式）
│   ├── review-rubric.md         # AIレビュアーのレビュー観点
│   ├── pr-workflow.md           # ブランチ→PR→マージの1サイクル手順
│   ├── branch-naming.md         # ブランチ命名規則（type語彙の正本）
│   ├── issue-writing.md         # Issueの書き方（構造・受け入れ条件）
│   ├── pr-writing.md            # PRの書き方（タイトル=コミット・本文）
│   ├── git-history.md           # 過去のコードは消えない（履歴から取り出す手順）
│   ├── multi-language-setup.md  # 多言語対応への移行ガイド（環境構築メモ）
│   ├── learning-log/            # 週次の学習ログ（YYYY-WW.md）
│   └── notes/                   # 概念メモ（srp.md, sql-injection.md など）
├── .github/workflows/       # 言語ごとにCIを分割（python-ci.yml / go-ci.yml / rust-ci.yml / typescript-ci.yml）。codeql.yml は Phase 3〜
├── .gitignore               # 各言語の生成物（.venv/ node_modules/ target/ dist/）を無視
│
├── python/                  # ← Python の世界（ツールチェーンの根）
│   ├── pyproject.toml           # 依存・Ruff/mypy 設定
│   ├── .python-version          # 使う Python の版（uv が読む）
│   ├── uv.lock                  # 依存の固定（再現性）
│   └── phase1/ …                # basics/（ドリル）＋ log_analyzer/（ツール）…
│
├── typescript/              # ← TypeScript の世界（pnpm）
│   ├── package.json             # 依存・スクリプト（packageManager で pnpm の版を固定）
│   ├── tsconfig.json            # TypeScript コンパイラ設定
│   ├── .node-version            # 使う Node の版
│   ├── pnpm-lock.yaml           # 依存の固定
│   └── phase1/ …
│
├── go/                      # ← Go の世界（modules）
│   ├── go.mod                   # 依存の台帳（依存を足すと go.sum が増える）
│   └── phase1/ …
│
└── rust/                    # ← Rust の世界（Cargo workspace）
    ├── Cargo.toml               # workspace（複数 crate を束ねる）
    ├── rust-toolchain.toml      # 使う Rust の版・道具（clippy / rustfmt）
    ├── Cargo.lock               # 依存の固定
    └── phase1/ …                # 1 ドリル/ツール = 1 crate
```

> **ツールは「作成フェーズ」に置いて以降は動かさない** — 後フェーズの改修（P2でテスト追加・リファクタ、P3でセキュリティ点検）は、そのファイルを**その場で編集**する（別フェーズへコピー/移動しない）。フェーズ番号＝**着手時期**で、経緯はツール直下 README の履歴＋git＋learning-log が表す。
> **Phase 0** はフェーズ用フォルダなし（成果物は各言語ルートの設定と `docs/`）／**Phase 4** は別リポジトリ（成果物リポジトリ）／`docs/` だけは言語・フェーズに属さない横断置き場。
> 各言語の宣言/ロックファイル（`pyproject.toml`・`package.json`・`go.mod`・`Cargo.toml` とそれぞれのロック）は**必ずコミット**、生成物（`.venv/`・`node_modules/`・`target/`・`dist/`）は `.gitignore` で**除外**する（再生成できるため）。

**学習リポジトリは Public 推奨**（CodeQL が無料で使える＋公開前提が機密情報管理の練習になる）。Phase 4 では、これとは別に成果物用の **プロジェクトリポジトリを独立して作る**（例: `log-pilot`）。その構成（参考。下は Python の例で、選んだ言語に読み替える）:

```
log-pilot/
├── README.md                  # 目的・インストール・使い方・設計意図
├── pyproject.toml
├── .github/workflows/         # ci.yml + codeql.yml
├── docs/
│   ├── design.md              # 設計判断の記録（ADR的に）
│   └── threat-model.md        # 簡易脅威モデリング
├── src/
│   └── logpilot/              # src レイアウト（実務で一般的）
├── tests/
└── .env.example               # 必要な環境変数の雛形（実値は置かない）
```

## 運用ルール

1. 作業はブランチ + Pull Request で行い、CI・CodeQLの指摘に対応してから自分でマージする（手順とフロー図: [docs/pr-workflow.md](docs/pr-workflow.md)）
2. ツール開発はIssue駆動（AIがPO役として発行したチケットをIssueに登録し、仕様の質疑もIssueコメントに残す）
3. `<言語>/phase3/security/` 配下のコード（例: `python/phase3/security/`）は学習目的で**意図的に脆弱**。実環境では使用しないこと
4. 週1回 `docs/learning-log/` に「やったこと・学んだこと・次にやること」を記録する
5. AIが書いたコードを採用した場合は、その旨と自分が検証した内容をコミットメッセージに残す

## セットアップ

リポジトリを取得:

```bash
git clone https://github.com/<username>/coding-practice.git
cd coding-practice
```

各言語は独立した「世界」。**使う言語のディレクトリに入って**セットアップする（他言語は触らなくてよい）。各ツール（uv / pnpm / go / rustup）の導入と「なぜそうするのか」は [docs/multi-language-setup.md](docs/multi-language-setup.md) を参照。

| 言語 | 作業ディレクトリ | コマンド（その言語のルートで） | 通れば言えること |
|------|----------------|------------------------------|-----------------|
| Python | `python/` | `uv sync` → `uv run python main.py` | 依存が再現でき、コードが動く |
| TypeScript | `typescript/` | `pnpm install` → `pnpm exec tsc --noEmit` | 依存が入り、型チェックが通る |
| Go | `go/` | `go build ./...` | 全パッケージがビルドできる |
| Rust | `rust/` | `cargo check` | 全 crate がコンパイルできる |

## 週次の回し方（15時間の例）

| 曜日 | 内容 | 時間 |
|------|------|------|
| 週初め | AIにチケットを発行させIssue登録 → 仕様の質疑 | 1h |
| 平日×4 | 実装（ドリル or Issue対応）+ AIレビューの差分読み | 2h×4 |
| 週末1 | インプット（本・OSSリーディング・セキュリティ学習） | 3h |
| 週末2 | 受け入れ条件で自己検収 → PRマージ → 学習ログ記入 | 3h |

## 詰まったときの指針

- 30分悩んだらAIに「答え」ではなく「ヒント」を求める（「答えは言わずに、どこを疑うべきかだけ教えて」）
- エラーは消す前に読む。エラーメッセージを翻訳・解説させるのは良いAIの使い方
- 完璧主義に注意。動く→読みやすく→安全に、の順で良い

> 実装の最中にAIをどこまで使ってよいか（🟢常時OK／🟡自力の後／🔴禁止 の線引き、エディタ補完の扱い）の全体像は [docs/ai-workflow.md の「実装中のAIの使い方」](docs/ai-workflow.md) を参照。

---

学習開始日: 2026-06-XX / 目標達成日: 2026-12-XX
