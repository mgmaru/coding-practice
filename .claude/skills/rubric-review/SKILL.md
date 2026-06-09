---
name: rubric-review
description: Reviews the user's own code changes against THIS repository's rubric (docs/review-rubric.md), selecting criteria by Phase × Track. Returns findings as must/should/nice with reasons, written in Japanese. Skips layer-0 concerns (format/lint/type/SAST/deps) that tools already cover. Use when the user asks to review their implementation in this learning repo (drills or tools/projects) before merging — e.g. "/rubric-review", "rubric でレビューして", "Phase2 のツールをレビューして".
disable-model-invocation: true
argument-hint: "[phase] [track]"
arguments: [phase, track]
---

# rubric-review

Review the user's code against the repository's own rubric in `docs/review-rubric.md`.

This skill is the **procedure**. `docs/review-rubric.md` is the **content and the single
source of truth**. The rubric is **injected live at the end of this skill** (see the
"Rubric" section) via `cat docs/review-rubric.md`, so it is the live file — never a copy.
Never hardcode, copy, or paraphrase the criteria; editing the doc updates this skill
automatically (no drift).

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

## 3. Select criteria from the rubric

The full rubric is inlined in the **"Rubric"** section at the end of this skill (live from
`docs/review-rubric.md`). If for some reason it is not present — e.g. the session was not
started from the repo root, so the injection produced nothing — Read `docs/review-rubric.md`
directly before continuing.

1. Use the rubric's **適用マトリクス (Phase × Track)** to pick which groups to apply.
2. Shape reminder (the rubric is authoritative): tools are a *superset* of drills —
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

## Rubric (live single source of truth — injected at load time)

The following is inlined from `docs/review-rubric.md` at skill-load time. This is the live
file, not a copy; do not paraphrase or duplicate it elsewhere. Apply its 適用マトリクス per §2–§3.

!`cat docs/review-rubric.md`
