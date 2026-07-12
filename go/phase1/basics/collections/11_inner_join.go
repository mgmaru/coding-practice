package collections

type UserInfo struct {
	id   int
	name string
}

type OrderInfo struct {
	userId int
	item   string
}

type UserOrderInfo struct {
	name string
	item string
}

func innerJoinById(users []UserInfo, orders []OrderInfo) []UserOrderInfo {

	userOrderInfo := make([]UserOrderInfo, 0, len(orders)) // 出力スライス // 修正：userOrderInfoの容量はusersのサイズでは足りない（１ユーザが複数の注文をしていることがあるため）。
	itemsByUserId := make(map[int][]string, len(orders))   // 一時マップ

	// OrderInfoに関して、useridをキーとして、itemのリストを作る
	for _, order := range orders {
		itemsByUserId[order.userId] = append(itemsByUserId[order.userId], order.item)
	}
	// UserInfoからidをループで取り出して、キーuserIdと一致を確認
	// 一致した場合はitemリストから１つずつ取り出して、[]UserOrderInfoにappend
	for _, user := range users { // usersはスライスなので並びも保証
		for _, item := range itemsByUserId[user.id] { // useridのitemリストが１個以上あればappend。空の場合は何もしない。
			userOrderInfo = append(userOrderInfo, UserOrderInfo{name: user.name, item: item})
		}
	}
	return userOrderInfo
}
