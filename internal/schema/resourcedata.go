// Package schema for maintaining resource data changes
package schema

import (
	"reflect"
)

// resourceDataProvider proxies interface to schema.ResourceData
// for convenient mocking purposes
type resourceDataProvider interface {
	Get(key string) any
	GetOk(key string) (any, bool)
	HasChange(key string) bool
	GetChange(key string) (any, any)
}

func getMapChangedKeys(keys []string, old, newMap map[string]any) map[string]any {
	changed := make(map[string]any)
	for _, key := range keys {
		if !reflect.DeepEqual(old[key], newMap[key]) {
			changed[key] = newMap[key]
		}
	}
	return changed
}

// GetResourceDataChangedKeys returns changed keys
func GetResourceDataChangedKeys(keys []string, d resourceDataProvider) map[string]any {
	changed := make(map[string]any)
	for _, key := range keys {
		if v := d.Get(key); v != nil && d.HasChange(key) {
			changed[key] = v
		}
	}
	return changed
}

// GetResourceDataListElementChanges returns list element changes
func GetResourceDataListElementChanges(keys []string, listKeyName string, listIndex int, d resourceDataProvider) map[string]any {
	changed := make(map[string]any)
	if !d.HasChange(listKeyName) {
		return changed
	}
	old, newName := d.GetChange(listKeyName)
	oldList := old.([]any)
	newList := newName.([]any)
	if len(oldList) < listIndex || len(newList) < listIndex {
		return changed
	}
	if len(oldList) == 0 {
		return newList[0].(map[string]any)
	}
	return getMapChangedKeys(keys, oldList[listIndex].(map[string]any), newList[listIndex].(map[string]any))
}

// IsDataElementAdded - checks if a data element added
func IsDataElementAdded(listKeyName string, d resourceDataProvider) bool {
	if !d.HasChange(listKeyName) {
		return false
	}
	old, newName := d.GetChange(listKeyName)
	oldList := old.([]any)
	newList := newName.([]any)
	return len(oldList) == 0 && len(newList) > 0
}

// IsDataElementRemoved - checks if a data element removed
func IsDataElementRemoved(listKeyName string, d resourceDataProvider) bool {
	if !d.HasChange(listKeyName) {
		return false
	}
	old, newName := d.GetChange(listKeyName)
	oldList := old.([]any)
	newList := newName.([]any)
	return len(newList) == 0 && len(oldList) > 0
}
