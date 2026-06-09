---
name: rubric-review
description: Review the user's own code changes against THIS repository's rubric (docs/review-rubric.md), selecting criteria by Phase × Track. Returns findings as must/should/nice with reasons, written in Japanese. Skips layer-0 concerns (format/lint/type/SAST/deps) that tools already cover. Use when the user asks to review their implementation in this learning repo (drills or tools/projects) before merging — e.g. "/rubric-review", "rubric でレビューして", "Phase2 のツールをレビューして".
disable-model-invocation: true
argument-hint: "[phase] [track]"
arguments: [phase, track]
---

# rubric-review

Review the user's code against the repository's own rubric in `docs/review-rubric.md`.

This skill is the **procedure**. `docs/review-rubric.md` is the **content and the single
source of truth**. Never hardcode or paraphrase the criteria here — always read the rubric
at run time so that edits to the rubric take effect automatically (no drift between this
skill and the doc).

## 0. Output language (read first — non-negotiable)

- Write the entire review in **Japanese**.
- Use the rubric's vocabulary **verbatim**; do NOT translate it:
  `層0 / 層1 / 層2`, `① 土台`, `② 道具の良し悪し`, `③ 可読性・構造`, `④ テスト`,
  `⑤ セキュリティ`, and the severities `must / should / nice`.
- These English instructions steer *how you work*; they must not leak into the output.

## 1. Decide WHAT to review

- Default: the current branch's changes — uncommitted changes if any, otherwise the diff
  against the default branch (`main`).
- If the user names specific files or directories, review those instead.
- Each language lives in its own root (`python/`, `typescript/`, `go/`, `rust/`). Review
  only the changed source files; ignore generated artifacts.

## 2. Decide Phase × Track (this selects which rubric groups apply)

Resolve both, in this order:

1. **Arguments win.** `$phase` (e.g. `1`, `2`, `3`) and `$track`
   (`drill` / `ドリル`, or `tool` / `project` / `ツール`).
2. **Else infer from the path** of the files under review:
   - **Phase** = the `phaseN/` segment (e.g. `python/phase2/...` → Phase 2).
   - **Track**: a drill directory (e.g. `basics/`, or an exercise-style folder that has a
     model answer) → **ドリル**; a named tool/project directory (e.g. `log_analyzer/`)
     → **ツール/プロジェクト**.
3. **Else ask** one short question to confirm Phase and Track. Do not guess silently when
   the path is ambiguous.

State the resolved **Phase** and **Track** at the top of the review so the user can correct you.

## 3. Read the rubric and select criteria

1. Read `docs/review-rubric.md`.
2. Use its **適用マトリクス (Phase × Track)** to pick which groups to apply.
   - Shape reminder (the rubric file is authoritative): tools are a *superset* of drills —
     drills ≈ ①②; tools add ③ and 層2 from Phase 1; ④ unlocks at Phase 2; ⑤ at Phase 3.
3. Apply **only** the selected groups. Applying more — e.g. demanding ④ tests in Phase 1, or
   ⑤ security before Phase 3 — is a *defect* of the review, not thoroughness.

## 4. Stay out of layer 0

Do **not** comment on anything the toolchain already covers (大原則①): formatting, lint,
type-checking, SAST, dependency scanning. If you spot such an issue, assume CI handles it and
say nothing. Spend the whole review budget on what tools cannot see (層1 / 層2).

## 5. Respect the self-review-first discipline

The rubric's 大原則③ order is: ① self-review → ② AI review → ③ read the diff. Before
reviewing, briefly confirm the user has already self-reviewed. If they clearly have not,
gently remind them — but still provide the review (do not block).

## 6. Output format

For each finding:

- **重大度**: `must` / `should` / `nice`
- **観点**: the rubric group it maps to (e.g. `③ 可読性・構造`)
- **場所**: `path:line`
- **指摘と理由**: what + **why** (the "why" is the point — it trains the user's own judgment)
- **直し方のヒント** (optional): a hint, not a full rewrite, and only when useful

Group findings by 重大度 (must → should → nice). Then end with:

- a 1–2 line **総評**, and
- one or two **層2 の問い** for the user to answer themselves (e.g.
  「なぜこの関数分割にしたか自分の言葉で説明できるか」). The final judgment is theirs (層2).

Keep it focused: a few high-signal findings beat an exhaustive list.
