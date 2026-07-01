package collections

import (
	"reflect"
	"testing"
)

func TestAggregateDateStatus(t *testing.T) {
	cases := []struct {
		name  string
		input []logEntry
		want  map[string]map[string]int
	}{
		// 構造体logEntryのDateとLevelは共にstring型なので、nilは入らない。
		{name: "normal", input: []logEntry{{date: "06-01", level: "INFO"}, {date: "06-01", level: "ERROR"}, {date: "06-01", level: "INFO"}, {date: "06-02", level: "INFO"}}, want: map[string]map[string]int{"06-01": {"INFO": 2, "ERROR": 1}, "06-02": {"INFO": 1}}},
		{name: "EmptyInput", input: []logEntry{}, want: map[string]map[string]int{}},
		{name: "NilInput", input: nil, want: map[string]map[string]int{}},
		{name: "MissingDate(EmptyString)", input: []logEntry{{date: "06-01", level: "INFO"}, {date: "", level: "ERROR"}, {date: "06-02", level: "INFO"}}, want: map[string]map[string]int{"06-01": {"INFO": 1}, "06-02": {"INFO": 1}, "": {"ERROR": 1}}},
		{name: "MissingAllDate(EmptyString)", input: []logEntry{{date: "", level: "INFO"}, {date: "", level: "ERROR"}, {date: "", level: "INFO"}}, want: map[string]map[string]int{"": {"INFO": 2, "ERROR": 1}}},
		{name: "MissingLevel(EmptyString)", input: []logEntry{{date: "06-01", level: ""}, {date: "06-01", level: "ERROR"}, {date: "06-02", level: "INFO"}}, want: map[string]map[string]int{"06-01": {"ERROR": 1, "": 1}, "06-02": {"INFO": 1}}},
		{name: "MissingAllLevel(EmptyString)", input: []logEntry{{date: "06-01", level: ""}, {date: "06-01", level: ""}, {date: "06-02", level: ""}}, want: map[string]map[string]int{"06-01": {"": 2}, "06-02": {"": 1}}},
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
