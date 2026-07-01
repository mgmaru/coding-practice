package collections

import (
	"reflect"
	"testing"
)

func TestAggregateDateStatus(t *testing.T) {
	cases := []struct {
		name  string
		input []Status
		want  map[string]map[string]int
	}{
		// 構造体StatusのDateとLevelは共にstring型なので、nilは入らない。
		{name: "normal", input: []Status{{Date: "06-01", Level: "INFO"}, {Date: "06-01", Level: "ERROR"}, {Date: "06-01", Level: "INFO"}, {Date: "06-02", Level: "INFO"}}, want: map[string]map[string]int{"06-01": {"INFO": 2, "ERROR": 1}, "06-02": {"INFO": 1}}},
		{name: "EmptyInput", input: []Status{}, want: map[string]map[string]int{}},
		{name: "NilInput", input: nil, want: map[string]map[string]int{}},
		{name: "MissingDate(EmptyString)", input: []Status{{Date: "06-01", Level: "INFO"}, {Date: "", Level: "ERROR"}, {Date: "06-02", Level: "INFO"}}, want: map[string]map[string]int{"06-01": {"INFO": 1}, "06-02": {"INFO": 1}, "": {"ERROR": 1}}},
		{name: "MissingAllDate(EmptyString)", input: []Status{{Date: "", Level: "INFO"}, {Date: "", Level: "ERROR"}, {Date: "", Level: "INFO"}}, want: map[string]map[string]int{"": {"INFO": 2, "ERROR": 1}}},
		{name: "MissingLevel(EmptyString)", input: []Status{{Date: "06-01", Level: ""}, {Date: "06-01", Level: "ERROR"}, {Date: "06-02", Level: "INFO"}}, want: map[string]map[string]int{"06-01": {"ERROR": 1, "": 1}, "06-02": {"INFO": 1}}},
		{name: "MissingAlLevel(EmptyString)", input: []Status{{Date: "06-01", Level: ""}, {Date: "06-01", Level: ""}, {Date: "06-02", Level: ""}}, want: map[string]map[string]int{"06-01": {"": 2}, "06-02": {"": 1}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue := aggregateDateStatus(c.input)
			if !reflect.DeepEqual(gotValue, c.want) {
				t.Errorf("期待する値ではありません。gotValue= %v, want= %v", gotValue, c.want)
			}
		})
	}
}
