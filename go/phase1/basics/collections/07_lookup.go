package collections

type User struct {
	ID   int
	Name string
}

// T：keyの型 K：値の型
// 修正：値の型Kは、mapのキーには入らないので、comparableとしなくて良い。 -> K anyに修正
func findUsersByKey[T comparable, K any](users []User, input []T, keyOf func(User) T, valueOf func(User) K) []K {

	// Userのキーと値を格納するmap
	m := make(map[T]K, len(users)) // 追記：mapにはusersのキーと値を入れていくので、容量はlen(users)で足りる。
	// 値を格納するスライス
	// 修正：容量は最大input （usersだと足りない） -> inputがIDで[2, 1, 1, 2, 3]の場合、usersの容量を超える
	result := make([]K, 0, len(input))

	// キーと値を格納（ハッシュ化）
	for _, user := range users {
		m[keyOf(user)] = valueOf(user) // keyとvalueは呼び出し元で注入
	}
	// keyで値を引いて、スライスに格納
	for _, v := range input { // inputはスライスから取り出しているので、順序も保証
		// 修正：mapの中にキーがある場合のみ、append
		// 理由：mapにキーがない場合、ゼロ値をappendしてしまうため。
		if val, ok := m[v]; ok {
			result = append(result, val)
		}
		// キーがない場合は何もしない
	}
	return result
}
