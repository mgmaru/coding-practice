package collections

type person struct {
	name string // 名前
	dept string // 部門
}

// keyOf/valueOfでpersonからキー・値を取り出し、キーごとに値をまとめて返す関数
func groupBy(keyOf func(person) string, valueOf func(person) string, input []person) map[string][]string {
	// 修正：keyOf：person型の変数を受け取って、keyとなる要素を返す関数
	// 修正：valueOf：person型の変数を受け取って、valueとなる要素を返す関数
	// 理由：関数をgroupByの引数として渡してあげることで、呼び出し元でkeyとvalueの要素を注入する。

	m := make(map[string][]string, len(input))

	// keyが正しく指定できていない
	// 修正：keyOf関数をgroupByの引数で渡すことで、不正なキーを検証する意味がなくなった。
	// if key != "name" && key != "dept" {
	// 	return m // 空
	// }

	// inputが空もしくはnil
	// 修正：len(input) == 0 を消去
	// 理由：inputが空もしくはnilの時は、以下のfor文が1回も回らず、いずれの場合も空のmが返るので不要。

	// 修正：forの中の変数「person」を「p」へ修正。
	// 理由：personは型と命名が被っているので、forの中でperson型を使おうとしても使えなくなる。
	for _, p := range input {
		k := keyOf(p) // 修正：関数呼び出しを1回にまとめた。
		m[k] = append(m[k], valueOf(p))
	}
	return m
}

// キーと同じ名前のフィールド要素を返す関数
// 修正：不要になったので削除。
// func (p person) keyFunc(key string) string {
// 	if key == "name" { // name
// 		return p.name
// 	} else { // dept
// 		return p.dept
// 	}
// }

// keyに対応する要素を返す関数
// 修正：不要になったので削除
// func (p person) valueFunc(key string) string {
// 	if key == "name" { // keyがnameだったらdeptを返す
// 		return p.dept
// 	} else { // keyがdeptだったらnameを返す
// 		return p.name
// 	}
// }
