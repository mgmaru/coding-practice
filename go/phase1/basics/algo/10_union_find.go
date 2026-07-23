package algo

// 10. Union-Find（素集合） — README.md「Tier 3 / 10. Union-Find」参照
// n 個の要素の素集合を管理。Find(x)（代表元）/ Union(x,y) / Same(x,y) bool。
// 経路圧縮（Find で通ったノードの親を根へ張り替え）+ ランク/サイズ併合で、ほぼ O(1)。
// TODO: 自分で実装する
