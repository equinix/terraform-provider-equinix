package comparisons

import (
	"testing"

	"github.com/equinix/rest-go"
	"github.com/stretchr/testify/assert"
)

func TestHasApplicationErrorCode(t *testing.T) {
	// given
	code := "ERR-505"
	errors := []rest.ApplicationError{
		{
			Code: "ERR-505",
		},
		{
			Code: "anything-else",
		},
	}
	// when
	result := HasApplicationErrorCode(errors, code)
	// then
	assert.True(t, result, "Error list contains error with given code")
}

func TestStringsFound(t *testing.T) {
	// given
	needles := []string{"key1", "key5"}
	hay := []string{"key1", "key2", "Key3", "key4", "key5"}
	// when
	result := StringsFound(needles, hay)
	// then
	assert.True(t, result, "Given strings were found")
}

func TestAtLeastOneStringFound(t *testing.T) {
	// given
	needles := []string{"key4", "key2"}
	hay := []string{"key1", "key2"}
	// when
	result := AtLeastOneStringFound(needles, hay)
	// then
	assert.True(t, result, "Given strings were found")
}

func TestStringsFound_negative(t *testing.T) {
	// given
	needles := []string{"key1", "key6"}
	hay := []string{"key1", "key2", "Key3", "key4", "key5"}
	// when
	result := StringsFound(needles, hay)
	// then
	assert.False(t, result, "Given strings were found")
}

func TestIsEmpty(t *testing.T) {
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
		assert.Equal(t, expected[i], IsEmpty(input[i]), "Input %v produces expected result %v", input[i], expected[i])
	}
}

func TestSlicesMatch(t *testing.T) {
	// given
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
	// when
	results := make([]bool, len(expected))
	for i := range input {
		results[i] = SlicesMatch(input[i][0], input[i][1])
	}
	// then
	for i := range expected {
		assert.Equal(t, expected[i], results[i])
	}
}
