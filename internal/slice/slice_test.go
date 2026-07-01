package slice

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	in := []int{1, 2, 3}
	out := Map(in, func(n int) string { return strconv.Itoa(n + 3) })

	expected := []string{"4", "5", "6"}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf(`Output = %q, want match for %q`, out, expected)
	}
}
