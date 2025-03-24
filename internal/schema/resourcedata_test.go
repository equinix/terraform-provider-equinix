package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedResourceDataProvider struct {
	actual map[string]interface{}
	old    map[string]interface{}
}

func (r mockedResourceDataProvider) Get(key string) interface{} {
	return r.actual[key]
}

func (r mockedResourceDataProvider) GetOk(key string) (interface{}, bool) {
	v, ok := r.actual[key]
	return v, ok
}

func (r mockedResourceDataProvider) HasChange(key string) bool {
	return !reflect.DeepEqual(r.old[key], r.actual[key])
}

func (r mockedResourceDataProvider) GetChange(key string) (interface{}, interface{}) {
	return r.old[key], r.actual[key]
}

func TestProvider_resourceDataChangedKeys(t *testing.T) {
	// given
	keys := []string{"key", "keyTwo", "keyThree"}
	rd := mockedResourceDataProvider{
		actual: map[string]interface{}{
			"key":    "value",
			"keyTwo": "newValueTwo",
		},
		old: map[string]interface{}{
			"key":    "value",
			"keyTwo": "valueTwo",
		},
	}
	expected := map[string]interface{}{
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
		old: map[string]interface{}{
			listKeyName: []interface{}{
				map[string]interface{}{
					"key":      "value",
					"keyTwo":   "valueTwo",
					"keyThree": 50,
				},
			},
		},
		actual: map[string]interface{}{
			listKeyName: []interface{}{
				map[string]interface{}{
					"key":      "value",
					"keyTwo":   "newValueTwo",
					"keyThree": 100,
				},
			},
		},
	}
	expected := map[string]interface{}{
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
	old := map[string]interface{}{
		"key":    "value",
		"keyTwo": "valueTwo",
	}
	new := map[string]interface{}{
		"key":    "newValue",
		"keyTwo": "valueTwo",
	}
	expected := map[string]interface{}{
		"key": "newValue",
	}
	// when
	result := getMapChangedKeys(keys, old, new)
	// then
	assert.Equal(t, expected, result, "Function returns valid key changes")
}

func TestProvider_IsDataElementAdded(t *testing.T) {
	// given
	keys := []string{"key", "keyTwo", "keyThree"}
	listKeyName := "myList"
	rd := mockedResourceDataProvider{
		old: map[string]interface{}{
			listKeyName: []interface{}{
				map[string]interface{}{
					"key":      "value",
					"keyTwo":   "valueTwo",
					"keyThree": 50,
				},
			},
		},
		actual: map[string]interface{}{
			listKeyName: []interface{}{
				map[string]interface{}{
					"key":      "value",
					"keyTwo":   "newValueTwo",
					"keyThree": 100,
				},
			},
		},
	}
	// when
	result := IsDataElementAdded(keys, listKeyName, rd)
	// then
	assert.Equal(t, false, result, "Function returns valid key changes")
}

func TestProvider_IsDataElementRemoved(t *testing.T) {
	// given
	keys := []string{"key", "keyTwo", "keyThree"}
	listKeyName := "myList"
	rd := mockedResourceDataProvider{
		old: map[string]interface{}{
			listKeyName: []interface{}{
				map[string]interface{}{
					"key":      "value",
					"keyTwo":   "valueTwo",
					"keyThree": 50,
				},
			},
		},
		actual: map[string]interface{}{
			listKeyName: []interface{}{
				map[string]interface{}{
					"key":      "value",
					"keyTwo":   "newValueTwo",
					"keyThree": 100,
				},
			},
		},
	}
	// when
	result := IsDataElementRemoved(keys, listKeyName, rd)
	// then
	assert.Equal(t, false, result, "Function returns valid key changes")
}
