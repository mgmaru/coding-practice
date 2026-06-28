package collections

import (
	"reflect"
	"testing"
)

func TestGroupBy(t *testing.T) {
	cases := []struct {
		name         string
		keyOf        func(person) string
		valueOf      func(person) string
		inputPersons []person
		want         map[string][]string
	}{
		// 修正： keyOfとvalueOfに関数を渡す。
		{name: "normal", keyOf: func(p person) string { return p.dept }, valueOf: func(p person) string { return p.name }, inputPersons: []person{{name: "Sato", dept: "Sales"}, {name: "Suzuki", dept: "Dev"}, {name: "Tanaka", dept: "Sales"}}, want: map[string][]string{"Sales": {"Sato", "Tanaka"}, "Dev": {"Suzuki"}}},
		// {name: "invalidKey", inputKey: "age", inputPersons: []person{{name: "Sato", dept: "Sales"}, {name: "Suzuki", dept: "Dev"}, {name: "Tanaka", dept: "Sales"}}, want: map[string][]string{}},
		{name: "duplicateName", keyOf: func(p person) string { return p.name }, valueOf: func(p person) string { return p.dept }, inputPersons: []person{{name: "Sato", dept: "Dev"}, {name: "Saito", dept: "Sales"}, {name: "Sato", dept: "Sales"}}, want: map[string][]string{"Sato": {"Dev", "Sales"}, "Saito": {"Sales"}}},
		{name: "duplicateNameAndDept", keyOf: func(p person) string { return p.dept }, valueOf: func(p person) string { return p.name }, inputPersons: []person{{name: "Sato", dept: "Dev"}, {name: "Tanaka", dept: "Sales"}, {name: "Sato", dept: "Dev"}}, want: map[string][]string{"Dev": {"Sato", "Sato"}, "Sales": {"Tanaka"}}},
		{name: "EmptyInput", keyOf: func(p person) string { return p.dept }, valueOf: func(p person) string { return p.name }, inputPersons: []person{}, want: map[string][]string{}},
		{name: "NilInput", keyOf: func(p person) string { return p.dept }, valueOf: func(p person) string { return p.name }, inputPersons: nil, want: map[string][]string{}},
		// 修正：重複がないテストケースを追加
		{name: "noDuplicate", keyOf: func(p person) string { return p.dept }, valueOf: func(p person) string { return p.name }, inputPersons: []person{{name: "Sato", dept: "Dev"}, {name: "Tanaka", dept: "Account"}, {name: "Suzuki", dept: "Sales"}}, want: map[string][]string{"Dev": {"Sato"}, "Account": {"Tanaka"}, "Sales": {"Suzuki"}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			gotValue := groupBy(c.keyOf, c.valueOf, c.inputPersons)

			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}

		})
	}
}
