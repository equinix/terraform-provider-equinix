// Package comparisons for comparison functions
package comparisons

import (
	"cmp"
	"slices"
	"strings"
)

func IsEmpty(v interface{}) bool {
	switch v := v.(type) {
	case int:
		return v == 0
	case *int:
		return v == nil || *v == 0
	case string:
		return v == ""
	case *string:
		return v == nil || *v == ""
	case nil:
		return true
	default:
		return false
	}
}

// Subsets returns true if the first slice is a subset of the second slice
func Subsets[T cmp.Ordered](s1, s2 []T) bool {
	// Iterate over the first slice
	for _, e := range s1 {
		// If the element is not in the second slice, return false
		if !slices.Contains(s2, e) {
			return false
		}
	}

	return true
}

// SlicesMatch comparisons.SlicesMatch returns true if the two slices contain the same elements, regardless of order
func SlicesMatch[T cmp.Ordered](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	// Create copies of the slices to avoid mutating the input slices
	s1Copy := append([]T(nil), s1...)
	s2Copy := append([]T(nil), s2...)

	// Sort the slices
	slices.Sort(s1Copy)
	slices.Sort(s2Copy)

	return slices.Equal(s1Copy, s2Copy)
}

// caseInsensitiveLess is a comparison function for sorting strings case-insensitively
func _(s1, s2 string) int {
	switch {
	case strings.EqualFold(strings.ToLower(s1), strings.ToLower(s2)):
		return 0
	case strings.ToLower(s1) < strings.ToLower(s2):
		return -1
	default:
		return 1
	}
}

// CompareMaps is a comparison function for comparing vendor config maps
func CompareMaps(old, newMap map[string]interface{}) bool {
	var oldTemp = copyMap(old)
	var newTemp = copyMap(newMap)
	if len(oldTemp) != len(newTemp) {
		if len(oldTemp) == 0 || len(newTemp) == 0 {
			return true
		}
		return false

	}

	for key, value := range oldTemp {
		if val, ok := newTemp[key]; !ok || val != value {
			return false
		}
	}

	return true
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{}, len(m))
	for k, v := range m {
		nm[k] = v
	}
	return nm
}
