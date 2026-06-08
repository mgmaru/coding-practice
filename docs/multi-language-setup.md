# 多言語対応への移行ガイド（環境構築の学習メモ）

このリポジトリを Python 専用から **Python / TypeScript / Go / Rust の4言語対応** に作り替えるための手順書。
ただの作業ログではなく、**各ステップが「何を」「なぜ」しているのか**を理解することを目的にする。
コマンドを写経して終わりにせず、「このコマンドは何のファイルを作り、何を解決しているのか」を一つずつ説明できる状態を目指す。

> このドキュメント自体は学習用メモ。形式は厳密でなくてよいが、**概念の理解**を最優先にしている。

---

## 1. このドキュメントの位置づけ（前提）

ロードマップの言語戦略は「**まず1言語を深く／2言語目は Phase 4 以降**」。
今回やるのは *4言語を並行して学ぶこと* ではなく、**リポジトリの「器」を多言語対応にしておく（インフラ準備）** だけ。

> **構成（器）は多言語対応にする。学習計画（中身）は “1言語集中 → Phase 4 で2言語目” のまま据え置く。**

今は Python の中身がほぼ空（`main.py` の雛形だけ）なので、作り替えコストがほぼゼロ＝**やるなら今が最適**。

決定事項（議論の結論）:

| 項目 | 決定 | 理由 |
|------|------|------|
| ディレクトリ構成 | **言語をトップ階層** | パッケージ管理を言語ごとに完全分離するため |
| TypeScript のパッケージ管理 | **pnpm** | 高速・省ディスク・厳格。業務でも採用が増えている |
| バージョン管理 | **各言語ネイティブ** | 統一ツール（mise 等）を増やさず、各言語標準の版固定を使う |

---

## 2. 全体像 — 「1言語 = 1ツールチェーンの根」

作り替え後の構成。**言語をトップに置き、その中にフェーズ × トラック（ドリル/ツール）を入れる**。

```
coding-practice/
├── README.md                  # 横断（言語に依存しない）
├── docs/                      # 横断（このファイルもここ）
├── .github/workflows/         # 言語ごとに job を分ける（Phase 2 で着手）
├── .gitignore                 # 各言語の「生成物」を無視する設定を追記
│
├── python/                    # ← Python の世界。ここがツールチェーンの根
│   ├── pyproject.toml         #   依存・ツール設定
│   ├── .python-version        #   使う Python の版
│   └── phase1/{basics, log_analyzer}/   # basics=ドリル / log_analyzer=ツール
│
├── typescript/                # ← TypeScript の世界
│   ├── package.json
│   ├── tsconfig.json
│   ├── .node-version
│   └── phase1/ ...
│
├── go/                        # ← Go の世界
│   ├── go.mod
│   └── phase1/ ...
│
└── rust/                      # ← Rust の世界
    ├── Cargo.toml             #   workspace（複数 crate をまとめる）
    ├── rust-toolchain.toml
    └── phase1/ ...
```

### なぜ「言語トップ」なのか

各言語は **自分専用の設定ファイル**を必要とする（Python は `pyproject.toml`、Node は `package.json`、Go は `go.mod`、Rust は `Cargo.toml`）。
これらを**リポジトリのルートに全部置くと衝突する**し、「どの言語のツールがどこを見るのか」が曖昧になる。

言語をトップに置けば、**「`python/` の中だけが Python の世界」「`go/` の中だけが Go の世界」** と境界がはっきりする。
これが「1言語 = 1ツールチェーンの根」という考え方で、`_prompts` に書いた「パッケージを別管理にする」の具体形。

> 逆に「フェーズをトップ（`phase1/python/`, `phase1/go/`…）」にすると、`go.mod` や Cargo の設定がフェーズをまたいで分断され、管理が散らかる。だから言語トップを選んだ。

---

## 3. 先に押さえる5つの概念（これが分かると全部読める）

言語ごとに名前は違うが、**やっていることは同じ**。まずこの対応関係を頭に入れる。

| 概念 | これは何か | Python | Node(TS) | Go | Rust |
|------|-----------|--------|----------|----|------|
| **パッケージマネージャ** | 依存ライブラリを入れる/管理する道具 | uv | pnpm | go（標準内蔵） | Cargo |
| **依存の宣言ファイル** | 「このプロジェクトは何に依存するか」を書く台帳 | `pyproject.toml` | `package.json` | `go.mod` | `Cargo.toml` |
| **ロックファイル** | 依存の**正確な版**を固定し再現性を担保 | `uv.lock` | `pnpm-lock.yaml` | `go.sum` | `Cargo.lock` |
| **言語自体の版固定** | 「このプロジェクトは○○の何版で動く」 | `.python-version` | `engines`/`.node-version` | `go.mod` の `toolchain` 行 | `rust-toolchain.toml` |
| **依存の隔離場所** | 入れたライブラリが置かれる隔離領域 | `.venv/` | `node_modules/` | モジュールキャッシュ | `target/` |

- **宣言ファイル vs ロックファイルの違い**: 宣言は「だいたいこの版以上」とゆるく書く（人間が編集）。ロックは「**まさにこの版**」を機械が固定する（自動生成）。両方あることで「自分の手元でも他人の手元でも同じものが入る（再現性）」が成立する。
- **隔離場所はコミットしない**: `.venv/` や `node_modules/` や `target/` は「宣言ファイルから再生成できる成果物」なので Git に入れない（`.gitignore` で無視する。後述）。

---

## 4. 「Go と Rust に “仮想環境” は無い」— よくある誤解

Python に慣れていると「Go の venv は？Rust の venv は？」と探してしまうが、**Go と Rust に Python の `venv` に相当する“仮想環境”は無い**。
代わりに、各言語が別の仕組みで同じ目的（＝プロジェクトごとの依存の隔離と再現性）を達成している。

| | Python | Go | Rust |
|---|--------|----|------|
| 依存の隔離 | `.venv/` を**明示的に作る** | モジュールキャッシュ（共有）＋ `go.mod` が版を固定 | `target/` に**自動で**隔離 |
| 「環境を有効化」する操作 | `source .venv/bin/activate` 等が要る | **不要**（`go` コマンドが go.mod を見て自動解決） | **不要**（`cargo` が Cargo.toml を見て自動解決） |
| 言語本体の版切替 | `.python-version`（uv が読む） | `go.mod` の `toolchain` 行（必要版を自動DL） | `rust-toolchain.toml`（rustup が読む） |

ポイント: **Go と Rust は「環境を有効化する」という手順が要らない**。ディレクトリに置かれた宣言ファイル（`go.mod` / `Cargo.toml`）をコマンドが自動で読むので、Python の `activate` のような儀式が無い。
だから「仮想環境を用意する」ではなく、**「宣言ファイルを正しく置く」** ことが Go/Rust での等価な作業になる。

> 補足: 「ネイティブなバージョンマネージャ」が公式に無いのは **Node だけ**。Python は `.python-version`（uv）、Go は `toolchain` 行、Rust は `rust-toolchain.toml` が公式にあるが、Node は `package.json` に「必要な版」を**宣言**するだけで、複数版の切替自体は手動（Homebrew 等で入れる）になる。今は1言語集中なので実害は小さい。

---

## 5. 手順（各ステップが何をしているか）

> 実際の作業は **ブランチ上**で行う（`docs/pr-workflow.md` のルール: main に直接コミットしない）。

> **前提（macOS）**: 各言語ステップの先頭に「**インストール（初回のみ）**」を付けた。**そのツールがまだ入っていなければ**実行する（`xxx --version` が表示されれば導入済み＝skip 可）。例は Homebrew を使う。Homebrew 自体が無ければ先に https://brew.sh の手順で入れる（`brew --version` で確認）。

### ステップ 0: 作業ブランチを切る

```bash
git checkout -b chore/multi-language-layout
```

- **何をしている**: main から枝分かれした作業用の枝を作って、そこに移る。
- **なぜ**: main を壊さずに作業し、後で PR（変更のまとまり）としてレビュー → マージするため。`chore/` は「雑務・設定変更」を表す慣習的な接頭辞。

---

### ステップ 1: Python を `python/` へ引っ越す

#### インストール（初回のみ）: uv

```bash
brew install uv          # uv（Python のパッケージ/プロジェクト管理ツール）を入れる
uv --version             # 入ったか確認
```

- **uv とは**: 依存管理・仮想環境・Python 本体の導入までまとめて行う高速ツール（Rust 製）。`pip` + `venv` + `pyenv` をまとめた立ち位置。
- **Python 本体は別途入れなくてよい**: uv は `.python-version` を見て、**必要な版の Python を自分で入れてくれる**（`uv python install 3.12` 相当を自動で行う）。だから「Python をどう入れるか」で悩まずに済む。
- 公式インストーラ派なら: `curl -LsSf https://astral.sh/uv/install.sh | sh`

#### 引っ越し作業

```bash
mkdir python
git mv pyproject.toml uv.lock main.py .python-version python/
```

- **`mkdir python`**: Python 用の部屋（フォルダ）を作る。
- **`git mv <file> python/`**: ファイルを `python/` へ移動する。**`mv` ではなく `git mv` を使う**のがポイント。
  - **なぜ git mv か**: 普通の `mv` だと Git からは「古いファイルが消えて、新しいファイルが現れた」ように見え、変更履歴が途切れやすい。`git mv` は「移動した（rename）」として記録するので、**ファイルの履歴が引っ越し先に引き継がれる**。
  - **何を移動しているか**:
    - `pyproject.toml` … Python の依存・ツール設定の台帳
    - `uv.lock` … 依存の正確な版を固定したロックファイル
    - `main.py` … 現状の雛形コード
    - `.python-version` … 使う Python の版（uv がこれを読む）
- **`.venv/` は移動しない**: `.gitignore` で無視されている「再生成できる成果物」なので、引っ越し先で作り直す。

引っ越し後、Python の世界を作り直す:

```bash
cd python
uv sync
```

- **`uv sync`**: `pyproject.toml` と `uv.lock` を読み、`.venv/`（隔離場所）を作って、宣言通りの依存を入れる。
- **結果**: `python/.venv/` が再生成され、`uv run python main.py` が動く状態になる。これが「環境が再現できた」状態。

---

### ステップ 2: TypeScript（pnpm）の世界を作る

#### インストール（初回のみ）: Node.js + corepack

```bash
brew install node        # Node.js 本体を入れる（npm が同梱される）
brew install corepack    # corepack を入れる（最近の Node/Homebrew では別 formula）
corepack enable          # corepack を有効化（pnpm の版を自動で揃えられるようになる）
node -v && corepack --version   # 入ったか確認
```

- **なぜ Node を入れるのか**: TypeScript は最終的に JavaScript として **Node.js（JS の実行環境）** の上で動く。pnpm も Node の上で動くので、まず Node が要る。
- **corepack を別に入れる理由**: 以前は corepack が Node に同梱されていたが、**最近の Node（Homebrew 配布の v26 など）は corepack を同梱・PATH 配置しない**。そのため `corepack enable` だけでは `command not found` になる。Homebrew では `brew install corepack` で別途入れる。
- **pnpm を直接入れず corepack を使う理由**: corepack は `package.json` の `packageManager` 欄に書いた pnpm の版を自動で用意する。だから将来の自分とも版が揃う。
- **もっと手数を減らしたいなら**: corepack を使わず `brew install pnpm` で pnpm を直接入れてもよい。その場合 `packageManager` 欄によるプロジェクト単位の版固定は効かず、pnpm の版は Homebrew 管理になる（学習用途なら実害は小さい）。
- **Node の版固定は手動**（セクション4の通り）: Node だけネイティブな版マネージャが無い。複数版を切り替えたくなったら Homebrew で入れ直すか、将来 mise 等の導入を再検討する。

#### プロジェクトを作る

```bash
mkdir typescript && cd typescript
pnpm init
```

- **`pnpm init`**: `package.json`（依存・スクリプトの台帳）を作る。

作った `package.json` に、版固定の宣言を加える（手で編集）:

```jsonc
{
  "packageManager": "pnpm@9.0.0",      // 使う pnpm の版を固定（corepack が読む）
  "engines": { "node": ">=22" }          // 必要な Node の版を宣言
}
```

次に、TypeScript 開発に必要な道具を入れる:

```bash
pnpm add -D typescript vitest eslint prettier
pnpm exec tsc --init
echo "22" > .node-version
```

- **`pnpm add -D <パッケージ>`**: ライブラリを入れる。`-D` は **devDependencies**（開発時だけ必要で、成果物には同梱しないもの）の意味。
  - `typescript` … TS を JS に変換するコンパイラ（`tsc`）。**型チェッカも兼ねる**。
  - `vitest` … テストランナー（pytest に相当）。
  - `eslint` … lint（コードの問題検出）。
  - `prettier` … フォーマッタ（整形）。
- **`pnpm exec tsc --init`**: `tsconfig.json`（TypeScript コンパイラの設定）を作る。
- **`echo "22" > .node-version`**: 使う Node の版をファイルに明記（人や一部ツールが読む“目印”）。
- **生成物**: `package.json` / `pnpm-lock.yaml`（ロック）/ `node_modules/`（隔離場所）/ `tsconfig.json`。

---

### ステップ 3: Go（modules）の世界を作る

#### インストール（初回のみ）: Go

```bash
brew install go          # Go 本体（コンパイラ＋go コマンド）を入れる
go version               # 入ったか確認
```

- **Go は1つで完結**: `go` コマンドにビルド・テスト・依存管理（modules）が全部入っている。別途パッケージマネージャを入れる必要が無い。
- 公式 pkg 派なら: https://go.dev/dl から macOS 用インストーラを入れる。

#### モジュールを作る

```bash
mkdir go && cd go
go mod init github.com/<github-username>/coding-practice/go
```

- **Go にはパッケージマネージャが別に無い**。`go` コマンド自体が依存解決もビルドもテストも行う。
- **`go mod init <module-path>`**: `go.mod`（依存の台帳）を作る。引数の **module path** は「このモジュールの住所（import するときの接頭辞）」。
  - 慣習として **リポジトリの URL** にする: `github.com/<github-username>/coding-practice/go`。
  - これは import 文の先頭になり、外部から取得されるときの識別子にもなる。だから一意である必要がある（公開予定のリポジトリ URL に合わせるのが安全）。
  - ⚠️ `<github-username>` は **GitHub のユーザー名**に置き換える（git の表示名ではない）。要確認。
- **版固定**: 生成された `go.mod` に `go 1.xx`（最低版）と必要なら `toolchain go1.xx.y`（使うツールチェーン）が入る。`toolchain` 行があると、その版が無ければ Go が**自動でダウンロード**する。これが Go ネイティブの版固定。
- **生成物**: `go.mod`。依存を足すと `go.sum`（ロック相当）も増える。

動作確認:

```bash
go build ./...   # go/ 以下の全パッケージをビルド（... は再帰）
```

---

### ステップ 4: Rust（Cargo workspace）の世界を作る

#### インストール（初回のみ）: Rust（rustup 経由）

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh   # rustup を入れる（画面の指示に従う／既定で可）
source "$HOME/.cargo/env"                                        # 今のシェルに PATH を通す（新しいシェルでは不要）
rustc --version && cargo --version                               # 入ったか確認
```

- **なぜ Homebrew でなく rustup か**: Rust は **rustup**（ツールチェーン管理ツール）経由で入れるのが公式推奨。`rust-toolchain.toml`（版固定）を読んでプロジェクトごとに適切な版・コンポーネント（clippy/rustfmt）へ自動で切り替えるのは rustup の役割だから。
- **入るもの**: `rustc`（コンパイラ）/ `cargo`（ビルド＋パッケージ管理）/ `rustup`（版管理）一式。
- brew 派なら `brew install rustup` 後に `rustup-init` でも可。

#### ワークスペースを作る

```bash
mkdir rust && cd rust
```

`rust/Cargo.toml` を手で作り、**workspace（複数 crate の入れ物）**として宣言する:

```toml
[workspace]
resolver = "2"
members = ["phase1/*"]   # phase1 配下の各ディレクトリを 1 crate として扱う
```

- **crate（クレート）** = Rust のパッケージ1個（ライブラリ or 実行ファイル）。
- **workspace** = 複数の crate を束ねる仕組み。トップに1つ `Cargo.toml` を置き、ドリル/ツール1個 = 1 crate にする。
  - **なぜ workspace か**: ドリルは小さい crate が大量にできる。workspace なら `target/`（隔離場所）と `Cargo.lock`（ロック）が**1か所に集約**され、`cargo test` で全部まとめて回せる。

版固定のため `rust/rust-toolchain.toml` を作る:

```toml
[toolchain]
channel = "stable"                    # 使う Rust の系列（安定版）
components = ["clippy", "rustfmt"]     # lint(clippy) と整形(rustfmt) も一緒に入れる
```

- **`rust-toolchain.toml`**: このディレクトリに入ると rustup が自動でこの版・道具に切り替える。これが Rust ネイティブの版固定。

動作確認:

```bash
cargo check   # 実際のバイナリは作らず、コンパイルが通るかだけ高速に確認
```

---

### ステップ 5: `.gitignore` に「生成物」を追記する

```gitignore
node_modules/
typescript/dist/
rust/target/
python/.venv/
```

- **なぜ無視するのか**: これらは全て **宣言ファイルから再生成できる成果物**（隔離場所やビルド出力）。Git に入れると巨大になり、再現性にも寄与しない。
  - `node_modules/` … pnpm が入れた依存（`package.json` + `pnpm-lock.yaml` から再生成できる）
  - `typescript/dist/` … TS をコンパイルした出力
  - `rust/target/` … Cargo のビルド出力
  - `python/.venv/` … uv が作る仮想環境（移動に伴いパスが変わるので明記）

> 逆に、**宣言ファイルとロックファイルは必ずコミットする**（`pyproject.toml`/`uv.lock`, `package.json`/`pnpm-lock.yaml`, `go.mod`/`go.sum`, `Cargo.toml`/`Cargo.lock`）。これらが「再現の素」だから。

---

## 6. 検収（受け入れ条件）

各言語の世界が独立して動くことを確認する:

| 言語 | 確認コマンド（その言語のルートで） | 通れば何が言える |
|------|------------------------------------|------------------|
| Python | `uv sync` → `uv run python main.py` | 依存が再現でき、コードが動く |
| TypeScript | `pnpm install` → `pnpm exec tsc --noEmit` | 依存が入り、型チェックが通る |
| Go | `go build ./...` | 全パッケージがビルドできる |
| Rust | `cargo check` | 全 crate がコンパイルできる |

全部通ったら、PR を作って自己レビュー → AI レビュー（観点: `docs/review-rubric.md`）→ 自分でマージ（手順: `docs/pr-workflow.md`）。

---

## 7. 移行に伴って更新する他のドキュメント（実装後でよい）

- **README「ディレクトリ構成」**: いまの *フォルダ = フェーズ* から、*言語 / フェーズ / トラック* の3軸に書き換える。
- **roadmap「言語戦略」**: 「**器は多言語対応・学習は1言語集中のまま**」という注記を足す（並行学習の許可ではないと明記）。
- **review-rubric**: すでに「各言語の同等物に読み替える」と書いてあるので追記は少なめ。
- **（Phase 3 着手時）CodeQL × Rust**: Rust の CodeQL 対応はプレビュー段階。Rust では `cargo-audit` + `clippy` を SAST の主役に読み替える、と注記する。

---

## 8. 用語ミニ辞典

| 用語 | 意味 |
|------|------|
| パッケージマネージャ | 依存ライブラリを入れ、版を管理する道具（uv / pnpm / go / Cargo） |
| 依存（dependency） | 自分のコードが利用する外部のライブラリ |
| 宣言ファイル | 「何に依存するか」を書く台帳（pyproject.toml / package.json / go.mod / Cargo.toml） |
| ロックファイル | 依存の**正確な版**を固定し再現性を保証する自動生成ファイル |
| 仮想環境（venv） | Python でプロジェクトごとに依存を隔離する仕組み。**Go/Rust には無い**（別方式で隔離） |
| module path（Go） | モジュールの住所。import の接頭辞。慣習的にリポジトリ URL |
| crate（Rust） | Rust のパッケージ1個（ライブラリ or 実行ファイル） |
| workspace（Rust） | 複数 crate を束ねて1か所で管理する仕組み |
| corepack | Node 同梱の、パッケージマネージャ（pnpm 等）の版を揃える仕組み |
| rustup | Rust のツールチェーン（コンパイラ・cargo・clippy 等）を管理・切替する公式ツール |
| Homebrew | macOS のパッケージ管理ツール（`brew`）。各言語のツール導入に使う |
| devDependencies | 開発時だけ必要で成果物には同梱しない依存（`pnpm add -D`） |
| `./...`（Go） | 「このディレクトリ以下の全パッケージ」を表す指定 |
