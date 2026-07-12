package collections

import (
	"reflect"
	"testing"
)

func TestInnerJoinById(t *testing.T) {
	cases := []struct {
		name   string
		users  []UserInfo
		orders []OrderInfo
		want   []UserOrderInfo
	}{
		{name: "AllUsersIncludeOrders", users: []UserInfo{{id: 1, name: "A"}, {id: 2, name: "B"}, {id: 3, name: "C"}}, orders: []OrderInfo{{userId: 1, item: "x"}, {userId: 2, item: "y"}, {userId: 3, item: "z"}}, want: []UserOrderInfo{{name: "A", item: "x"}, {name: "B", item: "y"}, {name: "C", item: "z"}}},
		{name: "SomeUserNotIncludeOrders（No order）", users: []UserInfo{{id: 1, name: "A"}, {id: 2, name: "B"}, {id: 3, name: "C"}}, orders: []OrderInfo{{userId: 1, item: "x"}, {userId: 3, item: "y"}}, want: []UserOrderInfo{{name: "A", item: "x"}, {name: "C", item: "y"}}},
		{name: "SomeItemsOrderNotIncludedUser（No user）", users: []UserInfo{{id: 1, name: "A"}, {id: 2, name: "B"}, {id: 3, name: "C"}}, orders: []OrderInfo{{userId: 1, item: "x"}, {userId: 2, item: "y"}, {userId: 3, item: "z"}, {userId: 4, item: "i"}}, want: []UserOrderInfo{{name: "A", item: "x"}, {name: "B", item: "y"}, {name: "C", item: "z"}}},
		{name: "EmptyUser", users: []UserInfo{}, orders: []OrderInfo{{userId: 1, item: "x"}, {userId: 2, item: "y"}, {userId: 3, item: "z"}}, want: []UserOrderInfo{}},
		{name: "EmptyOrder", users: []UserInfo{{id: 1, name: "A"}, {id: 2, name: "B"}, {id: 3, name: "C"}}, orders: []OrderInfo{}, want: []UserOrderInfo{}},
		{name: "EmptyUserAndOrder", users: []UserInfo{}, orders: []OrderInfo{}, want: []UserOrderInfo{}},
		{name: "NilUser", users: nil, orders: []OrderInfo{{userId: 1, item: "x"}, {userId: 2, item: "y"}, {userId: 3, item: "z"}}, want: []UserOrderInfo{}},
		{name: "NilOrder", users: []UserInfo{{id: 1, name: "A"}, {id: 2, name: "B"}, {id: 3, name: "C"}}, orders: nil, want: []UserOrderInfo{}},
		{name: "NilUserAndOrder", users: nil, orders: nil, want: []UserOrderInfo{}},
		{name: "OneUserPlacesMultipleOrders", users: []UserInfo{{id: 1, name: "A"}, {id: 2, name: "B"}, {id: 3, name: "C"}}, orders: []OrderInfo{{userId: 1, item: "x"}, {userId: 2, item: "y"}, {userId: 3, item: "z"}, {userId: 1, item: "i"}}, want: []UserOrderInfo{{name: "A", item: "x"}, {name: "A", item: "i"}, {name: "B", item: "y"}, {name: "C", item: "z"}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := innerJoinById(c.users, c.orders)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}
