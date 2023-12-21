package equinix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_stringsFound(t *testing.T) {
	// given
	needles := []string{"key1", "key5"}
	hay := []string{"key1", "key2", "Key3", "key4", "key5"}
	// when
	result := stringsFound(needles, hay)
	// then
	assert.True(t, result, "Given strings were found")
}

func TestProvider_atLeastOneStringFound(t *testing.T) {
	// given
	needles := []string{"key4", "key2"}
	hay := []string{"key1", "key2"}
	// when
	result := atLeastOneStringFound(needles, hay)
	// then
	assert.True(t, result, "Given strings were found")
}

func TestProvider_stringsFound_negative(t *testing.T) {
	// given
	needles := []string{"key1", "key6"}
	hay := []string{"key1", "key2", "Key3", "key4", "key5"}
	// when
	result := stringsFound(needles, hay)
	// then
	assert.False(t, result, "Given strings were found")
}

func TestProvider_isEmpty(t *testing.T) {
	// given
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
	// when then
	for i := range input {
		assert.Equal(t, expected[i], isEmpty(input[i]), "Input %v produces expected result %v", input[i], expected[i])
	}
}
