package collections

type person struct {
	name string // 名前
	dept string // 部門
}

// keyを受け取ってkeyごとのリストを返す関数
func groupBy(key string, input []person) map[string][]string {

	m := make(map[string][]string, len(input))

	// keyが正しく指定できていない
	if key != "name" && key != "dept" {
		return m // 空
	}

	// inputが空もしくはnil
	if len(input) == 0 {
		return m // 空
	}

	// 以下、keyが正しく指定されている
	for _, person := range input {
		m[person.keyFunc(key)] = append(m[person.keyFunc(key)], person.valueFunc(key))
	}
	return m
}

// キーと同じ名前のフィールド要素を返す関数
func (p person) keyFunc(key string) string {
	if key == "name" { // name
		return p.name
	} else { // dept
		return p.dept
	}
}

// keyに対応する要素を返す関数
func (p person) valueFunc(key string) string {
	if key == "name" { // keyがnameだったらdeptを返す
		return p.dept
	} else { // keyがdeptだったらnameを返す
		return p.name
	}
}
