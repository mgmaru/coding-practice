package algo

// 08. 二分ヒープ → Top-K 再実装 — README.md「Tier 3 / 8. 二分ヒープ」参照
// 配列表現の二分ヒープ（push/pop O(log n)）を実装し、collections #15 Top-N をヒープで解き直す。
// サイズ K の min-heap を保ち全要素を流し込む。K 超過なら根（=保持中の最小）を捨てる → O(n log K)。
// 親 (i-1)/2、子 2i+1 / 2i+2。比較関数を出力順（頻度→名前）と整合させる（でないと同頻度境界で取りこぼす）。
// TODO: 自分で実装する
