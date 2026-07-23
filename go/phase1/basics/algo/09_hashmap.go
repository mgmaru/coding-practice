package algo

// 09. ハッシュ表（チェイン法） — README.md「Tier 3 / 9. ハッシュ表」参照
// map の正体を自作。バケット配列 + 各バケットに連結リスト(or slice)で衝突を吸収し、負荷率超過でリサイズ(rehash)。
// API: Put(key string, val int) / Get(key string) (int, bool) / Delete(key string)
// ハッシュ関数は自前の簡単なもの（多項式ローリング / FNV 風）でよい。
// TODO: 自分で実装する
