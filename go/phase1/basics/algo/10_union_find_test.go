package algo

// 10. Union-Find のテスト（要素 0..4、操作列で検証）
// Union(0,1),Union(2,3),Union(1,2) の後: Same(0,3) → true, Same(0,4) → false
// 境界:
//   - 自分自身との Union(1,1)（何も起きない）
//   - 既に同じグループの Union（何も起きない）
//   - 経路圧縮の前後で Same の結果が変わらない
// TODO: テストを書く
