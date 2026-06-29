package collections

import (
	"reflect"
	"testing"
)

// IDで引いた場合のテスト
func TestFindUsersByKey_ID(t *testing.T) {
	cases := []struct {
		name    string
		users   []User
		input   []int
		keyOf   func(User) int
		valueOf func(User) string
		want    []string
	}{
		{name: "Normal", users: []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}, input: []int{2, 1, 2}, keyOf: func(u User) int { return u.ID }, valueOf: func(u User) string { return u.Name }, want: []string{"B", "A", "B"}},
		{name: "NameDuplicate(valueDuplicate)", users: []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}, {ID: 3, Name: "A"}}, input: []int{2, 1, 3}, keyOf: func(u User) int { return u.ID }, valueOf: func(u User) string { return u.Name }, want: []string{"B", "A", "A"}},
		{name: "IDDuplicate(keyDuplicate)", users: []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}, {ID: 1, Name: "B"}}, input: []int{1, 2, 1}, keyOf: func(u User) int { return u.ID }, valueOf: func(u User) string { return u.Name }, want: []string{"B", "B", "B"}},
		{name: "usersAndInputEmpty", users: []User{}, input: []int{}, keyOf: func(u User) int { return u.ID }, valueOf: func(u User) string { return u.Name }, want: []string{}},
		{name: "usersEmpty", users: []User{}, input: []int{2, 1, 2}, keyOf: func(u User) int { return u.ID }, valueOf: func(u User) string { return u.Name }, want: []string{}},
		{name: "inputEmpty", users: []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}, input: []int{}, keyOf: func(u User) int { return u.ID }, valueOf: func(u User) string { return u.Name }, want: []string{}},
		// 修正：IDが欠損している場合を追加
		{name: "IDDeficiency", users: []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}, input: []int{1, 3, 2}, keyOf: func(u User) int { return u.ID }, valueOf: func(u User) string { return u.Name }, want: []string{"A", "B"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := findUsersByKey(c.users, c.input, c.keyOf, c.valueOf)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}

// Nameで引いた場合のテスト
func TestFindUsersByKey_Name(t *testing.T) {
	cases := []struct {
		name    string
		users   []User
		input   []string
		keyOf   func(User) string
		valueOf func(User) int
		want    []int
	}{
		{name: "Normal", users: []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}, {ID: 3, Name: "C"}}, input: []string{"C", "A", "B"}, keyOf: func(u User) string { return u.Name }, valueOf: func(u User) int { return u.ID }, want: []int{3, 1, 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := findUsersByKey(c.users, c.input, c.keyOf, c.valueOf)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}
