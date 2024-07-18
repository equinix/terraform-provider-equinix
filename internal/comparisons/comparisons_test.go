package comparisons_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/stretchr/testify/assert"
)

func TestSubsets(t *testing.T) {
	// given
	needles := []string{"key1", "key5"}
	hay := []string{"key1", "key2", "Key3", "key4", "key5"}
	// when
	result := comparisons.Subsets(needles, hay)
	// then
	assert.True(t, result, "Given strings were found")
}

func TestSubsets_negative(t *testing.T) {
	// given
	needles := []string{"key1", "key6"}
	hay := []string{"key1", "key2", "Key3", "key4", "key5"}
	// when
	result := comparisons.Subsets(hay, needles)
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
		assert.Equal(t, expected[i], comparisons.IsEmpty(input[i]), "Input %v produces expected result %v", input[i], expected[i])
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
		results[i] = comparisons.SlicesMatch(input[i][0], input[i][1])
	}
	// then
	for i := range expected {
		assert.Equal(t, expected[i], results[i])
	}
}
