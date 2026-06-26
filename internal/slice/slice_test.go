package slice

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("nil slice returns nil", func(t *testing.T) {
		out := Map(nil, func(n int) string { return strconv.Itoa(n + 3) })

		if out != nil {
			t.Errorf(`Output = %q, want nil`, out)
		}
	})

	t.Run("empty slice returns empty slice", func(t *testing.T) {
		in := []int{}
		out := Map(in, func(n int) string { return strconv.Itoa(n + 3) })

		expected := []string{}

		if !reflect.DeepEqual(out, expected) {
			t.Errorf(`Output = %q, want match for %q`, out, expected)
		}
	})

	t.Run("slice returns mapped slice", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := Map(in, func(n int) string { return strconv.Itoa(n + 3) })

		expected := []string{"4", "5", "6"}

		if !reflect.DeepEqual(out, expected) {
			t.Errorf(`Output = %q, want match for %q`, out, expected)
		}
	})
}
