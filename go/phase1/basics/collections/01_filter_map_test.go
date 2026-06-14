package collections

import (
	"reflect"
	"testing"
)

func TestFilterSquaredEvens(t *testing.T) {

	input := []int{1, 2, 3, 4, 5, 6}

	got := filterSquaredEvens(input)

	want := []int{4, 26, 36}

	if reflect.DeepEqual(got, want) {
		t.Errorf("計算結果が違います。")
	}

}
