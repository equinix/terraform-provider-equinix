package equinix

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

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

type testAccConfig struct {
	ctx    map[string]interface{}
	config string
}

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"equinix": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_hasApplicationErrorCode(t *testing.T) {
	//given
	code := "ERR-505"
	errors := []rest.ApplicationError{
		{
			Code: "ERR-505",
		},
		{
			Code: randString(10),
		},
	}
	//when
	result := hasApplicationErrorCode(errors, code)
	//then
	assert.True(t, result, "Error list contains error with given code")
}

func TestProvider_stringsFound(t *testing.T) {
	//given
	needles := []string{"key1", "key5"}
	hay := []string{"key1", "key2", "Key3", "key4", "key5"}
	//when
	result := stringsFound(needles, hay)
	//then
	assert.True(t, result, "Given strings were found")
}

func TestProvider_atLeastOneStringFound(t *testing.T) {
	//given
	needles := []string{"key4", "key2"}
	hay := []string{"key1", "key2"}
	//when
	result := atLeastOneStringFound(needles, hay)
	//then
	assert.True(t, result, "Given strings were found")
}

func TestProvider_stringsFound_negative(t *testing.T) {
	//given
	needles := []string{"key1", "key6"}
	hay := []string{"key1", "key2", "Key3", "key4", "key5"}
	//when
	result := stringsFound(needles, hay)
	//then
	assert.False(t, result, "Given strings were found")
}

func TestProvider_resourceDataChangedKeys(t *testing.T) {
	//given
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
	//when
	result := getResourceDataChangedKeys(keys, rd)
	//then
	assert.Equal(t, expected, result, "Function returns valid key changes")
}

func TestProvider_resourceDataListElementChanges(t *testing.T) {
	//given
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
	//when
	result := getResourceDataListElementChanges(keys, listKeyName, 0, rd)
	//then
	assert.Equal(t, expected, result, "Function returns valid key changes")
}

func TestProvider_mapChanges(t *testing.T) {
	//given
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
	//when
	result := getMapChangedKeys(keys, old, new)
	//then
	assert.Equal(t, expected, result, "Function returns valid key changes")
}

func TestProvider_isEmpty(t *testing.T) {
	//given
	input := []interface{}{
		"test",
		"",
		nil,
		123,
		0,
		43.43,
	}
	expected := []bool{
		false,
		true,
		true,
		false,
		true,
		false,
		true,
	}
	//when then
	for i := range input {
		assert.Equal(t, expected[i], isEmpty(input[i]), "Input %v produces expected result %v", input[i], expected[i])
	}
}

func TestProvider_setSchemaValueIfNotEmpty(t *testing.T) {
	//given
	key := "test"
	s := map[string]*schema.Schema{
		key: {
			Type:     schema.TypeString,
			Optional: true,
		}}
	var b *int = nil
	d := schema.TestResourceDataRaw(t, s, make(map[string]interface{}))
	//when
	setSchemaValueIfNotEmpty(key, b, d)
	//then
	_, ok := d.GetOk(key)
	assert.False(t, ok, "Key was not set")

}

func TestProvider_slicesMatch(t *testing.T) {
	//given
	input := [][][]string{
		{
			{"DC", "SV", "FR"},
			{"FR", "SV", "DC"},
		},
		{
			{"SV"},
			{},
		},
		{
			{"DC", "DC", "DC"},
			{"DC", "SV", "DC"},
		},
		{
			{}, {},
		},
	}
	expected := []bool{
		true,
		false,
		false,
		true,
	}
	//when
	results := make([]bool, len(expected))
	for i := range input {
		results[i] = slicesMatch(input[i][0], input[i][1])
	}
	//then
	for i := range expected {
		assert.Equal(t, expected[i], results[i])
	}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Test helper functions
//_______________________________________________________________________

func testAccPreCheck(t *testing.T) {
	if _, err := getFromEnv(endpointEnvVar); err != nil {
		t.Fatalf("%s", err)
	}
	if _, err := getFromEnv(clientIDEnvVar); err != nil {
		t.Fatalf("%s", err)
	}
	if _, err := getFromEnv(clientSecretEnvVar); err != nil {
		t.Fatalf("%s", err)
	}
}

func newTestAccConfig(ctx map[string]interface{}) *testAccConfig {
	return &testAccConfig{
		ctx:    ctx,
		config: "",
	}
}

func (t *testAccConfig) build() string {
	return t.config
}

func nprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		var strVal string
		switch val.(type) {
		case []string:
			r := regexp.MustCompile(`" "`)
			strVal = r.ReplaceAllString(fmt.Sprintf("%q", val), `", "`)
		default:
			strVal = fmt.Sprintf("%v", val)
		}
		format = strings.Replace(format, "%{"+key+"}", strVal, -1)
	}
	return format
}

func randInt(n int) int {
	src := rand.NewSource(time.Now().UnixNano())
	var mu sync.Mutex
	mu.Lock()
	i := rand.New(src).Intn(n)
	mu.Unlock()
	return i
}

func randString(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	result := make([]byte, length)
	set := "abcdefghijklmnopqrstuvwxyz012346789"
	var mu sync.Mutex
	mu.Lock()
	r := rand.New(src)
	for i := 0; i < length; i++ {
		result[i] = set[r.Intn(len(set))]
	}
	mu.Unlock()
	return string(result)
}

func getFromEnv(varName string) (string, error) {
	if v := os.Getenv(varName); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("environmental variable '%s' is not set", varName)
}

func copyMap(source map[string]interface{}) map[string]interface{} {
	target := make(map[string]interface{})
	for k, v := range source {
		target[k] = v
	}
	return target
}

func setSchemaValueIfNotEmpty(key string, value interface{}, d *schema.ResourceData) error {
	if !isEmpty(value) {
		return d.Set(key, value)
	}
	return nil
}
