package schema

import (
	"log"
	"reflect"
)

// resourceDataProvider proxies interface to schema.ResourceData
// for convenient mocking purposes
type resourceDataProvider interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
	HasChange(key string) bool
	GetChange(key string) (interface{}, interface{})
}

func getMapChangedKeys(keys []string, old, new map[string]interface{}) map[string]interface{} {
	changed := make(map[string]interface{})
	for _, key := range keys {
		if !reflect.DeepEqual(old[key], new[key]) {
			changed[key] = new[key]
		}
	}
	return changed
}

func GetResourceDataChangedKeys(keys []string, d resourceDataProvider) map[string]interface{} {
	changed := make(map[string]interface{})
	for _, key := range keys {
		if v := d.Get(key); v != nil && d.HasChange(key) {
			changed[key] = v
		}
	}
	return changed
}

func GetResourceDataListElementChanges(keys []string, listKeyName string, listIndex int, d resourceDataProvider) map[string]interface{} {
	changed := make(map[string]interface{})
	if !d.HasChange(listKeyName) {
		return changed
	}
	old, new := d.GetChange(listKeyName)
	oldList := old.([]interface{})
	newList := new.([]interface{})
	if len(oldList) < listIndex || len(newList) < listIndex {
		return changed
	}
	if len(oldList) == 0 {
		log.Println(newList[0].(map[string]interface{}))
		return newList[0].(map[string]interface{})
	}
	return getMapChangedKeys(keys, oldList[listIndex].(map[string]interface{}), newList[listIndex].(map[string]interface{}))
}

func IsDataElementAdded(keys []string, listKeyName string, d resourceDataProvider) bool {
	if !d.HasChange(listKeyName) {
		return false
	}
	old, new := d.GetChange(listKeyName)
	oldList := old.([]interface{})
	newList := new.([]interface{})
	return len(oldList) == 0 && len(newList) > 0
}

func IsDataElementRemoved(keys []string, listKeyName string, d resourceDataProvider) bool {
	if !d.HasChange(listKeyName) {
		return false
	}
	old, new := d.GetChange(listKeyName)
	oldList := old.([]interface{})
	newList := new.([]interface{})
	return len(newList) == 0 && len(oldList) > 0
}
