package collections

import (
	"reflect"
	"testing"
)

func TestGroupBy(t *testing.T) {
	cases := []struct {
		name         string
		inputKey     string
		inputPersons []person
		want         map[string][]string
	}{
		{name: "normal", inputKey: "dept", inputPersons: []person{{name: "Sato", dept: "Sales"}, {name: "Suzuki", dept: "Dev"}, {name: "Tanaka", dept: "Sales"}}, want: map[string][]string{"Sales": {"Sato", "Tanaka"}, "Dev": {"Suzuki"}}},
		{name: "invalidKey", inputKey: "age", inputPersons: []person{{name: "Sato", dept: "Sales"}, {name: "Suzuki", dept: "Dev"}, {name: "Tanaka", dept: "Sales"}}, want: map[string][]string{}},
		{name: "duplicateName", inputKey: "name", inputPersons: []person{{name: "Sato", dept: "Dev"}, {name: "Saito", dept: "Sales"}, {name: "Sato", dept: "Sales"}}, want: map[string][]string{"Sato": {"Dev", "Sales"}, "Saito": {"Sales"}}},
		{name: "duplicateNameAndDept", inputKey: "dept", inputPersons: []person{{name: "Sato", dept: "Dev"}, {name: "Tanaka", dept: "Sales"}, {name: "Sato", dept: "Dev"}}, want: map[string][]string{"Dev": {"Sato", "Sato"}, "Sales": {"Tanaka"}}},
		{name: "EmptyInput", inputKey: "dept", inputPersons: []person{}, want: map[string][]string{}},
		{name: "NilInput", inputKey: "dept", inputPersons: nil, want: map[string][]string{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			gotValue := groupBy(c.inputKey, c.inputPersons)

			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値と違います。gotValue= %v, want= %v", gotValue, c.want)
			}

		})
	}
}
