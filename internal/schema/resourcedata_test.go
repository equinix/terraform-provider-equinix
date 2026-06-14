package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedResourceDataProvider struct {
	actual map[string]any
	old    map[string]any
}

func (r mockedResourceDataProvider) Get(key string) any {
	return r.actual[key]
}

func (r mockedResourceDataProvider) GetOk(key string) (any, bool) {
	v, ok := r.actual[key]
	return v, ok
}

func (r mockedResourceDataProvider) HasChange(key string) bool {
	return !reflect.DeepEqual(r.old[key], r.actual[key])
}

func (r mockedResourceDataProvider) GetChange(key string) (any, any) {
	return r.old[key], r.actual[key]
}

func TestProvider_resourceDataChangedKeys(t *testing.T) {
	// given
	keys := []string{"key", "keyTwo", "keyThree"}
	rd := mockedResourceDataProvider{
		actual: map[string]any{
			"key":    "value",
			"keyTwo": "newValueTwo",
		},
		old: map[string]any{
			"key":    "value",
			"keyTwo": "valueTwo",
		},
	}
	expected := map[string]any{
		"keyTwo": "newValueTwo",
	}
	// when
	result := GetResourceDataChangedKeys(keys, rd)
	// then
	assert.Equal(t, expected, result, "Function returns valid key changes")
}

func TestProvider_resourceDataListElementChanges(t *testing.T) {
	// given
	keys := []string{"key", "keyTwo", "keyThree"}
	listKeyName := "myList"
	rd := mockedResourceDataProvider{
		old: map[string]any{
			listKeyName: []any{
				map[string]any{
					"key":      "value",
					"keyTwo":   "valueTwo",
					"keyThree": 50,
				},
			},
		},
		actual: map[string]any{
			listKeyName: []any{
				map[string]any{
					"key":      "value",
					"keyTwo":   "newValueTwo",
					"keyThree": 100,
				},
			},
		},
	}
	expected := map[string]any{
		"keyTwo":   "newValueTwo",
		"keyThree": 100,
	}
	// when
	result := GetResourceDataListElementChanges(keys, listKeyName, 0, rd)
	// then
	assert.Equal(t, expected, result, "Function returns valid key changes")
}

func TestProvider_mapChanges(t *testing.T) {
	// given
	keys := []string{"key", "keyTwo", "keyThree"}
	old := map[string]any{
		"key":    "value",
		"keyTwo": "valueTwo",
	}
	newMap := map[string]any{
		"key":    "newValue",
		"keyTwo": "valueTwo",
	}
	expected := map[string]any{
		"key": "newValue",
	}
	// when
	result := getMapChangedKeys(keys, old, newMap)
	// then
	assert.Equal(t, expected, result, "Function returns valid key changes")
}

func TestProvider_IsDataElementAdded(t *testing.T) {
	// given
	listKeyName := "myList"
	rd := mockedResourceDataProvider{
		old: map[string]any{
			listKeyName: []any{
				map[string]any{
					"key":      "value",
					"keyTwo":   "valueTwo",
					"keyThree": 50,
				},
			},
		},
		actual: map[string]any{
			listKeyName: []any{
				map[string]any{
					"key":      "value",
					"keyTwo":   "newValueTwo",
					"keyThree": 100,
				},
			},
		},
	}
	// when
	result := IsDataElementAdded(listKeyName, rd)
	// then
	assert.Equal(t, false, result, "Function returns valid key changes")
}

func TestProvider_IsDataElementRemoved(t *testing.T) {
	// given
	listKeyName := "myList"
	rd := mockedResourceDataProvider{
		old: map[string]any{
			listKeyName: []any{
				map[string]any{
					"key":      "value",
					"keyTwo":   "valueTwo",
					"keyThree": 50,
				},
			},
		},
		actual: map[string]any{
			listKeyName: []any{
				map[string]any{
					"key":      "value",
					"keyTwo":   "newValueTwo",
					"keyThree": 100,
				},
			},
		},
	}
	// when
	result := IsDataElementRemoved(listKeyName, rd)
	// then
	assert.Equal(t, false, result, "Function returns valid key changes")
}
