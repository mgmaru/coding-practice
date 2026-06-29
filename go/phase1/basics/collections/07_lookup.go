package collections

type User struct {
	ID   int
	Name string
}

// T：keyの型 K：値の型
func findUsersByKey[T, K comparable](users []User, input []T, KeyOf func(User) T, ValueOf func(User) K) []K {

	// Userのキーと値を格納するmap
	m := make(map[T]K, len(users))
	// 値を格納するスライス
	result := make([]K, 0, len(users))

	// キーと値を格納（ハッシュ化）
	for _, user := range users {
		m[KeyOf(user)] = ValueOf(user) // keyとvalueは呼び出し元で注入
	}
	// keyで値を引いて、スライスに格納
	for _, v := range input { // inputはスライスから取り出しているので、順序も保証
		result = append(result, m[v])
	}
	return result
}
