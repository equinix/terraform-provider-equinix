package maps

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestSchemaSetToMap(t *testing.T) {
	// given
	type item struct {
		id       string
		valueOne int
		valueTwo int
	}
	setFunc := func(v interface{}) int {
		i := v.(item)
		return converters.HashcodeString(i.id)
	}
	items := []interface{}{
		item{"id1", 100, 200},
		item{"id2", 666, 999},
		item{"id3", 0, 100},
	}
	set := schema.NewSet(setFunc, items)
	// when
	list := SchemaSetToMap(set)
	// then
	assert.Equal(t, items[0], list[setFunc(items[0])])
	assert.Equal(t, items[1], list[setFunc(items[1])])
	assert.Equal(t, items[2], list[setFunc(items[2])])
}

func TestResourceDataChangedKeys(t *testing.T) {
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

func TestResourceDataListElementChanges(t *testing.T) {
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

func TestMapChanges(t *testing.T) {
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
	result := GetMapChangedKeys(keys, old, new)
	// then
	assert.Equal(t, expected, result, "Function returns valid key changes")
}
