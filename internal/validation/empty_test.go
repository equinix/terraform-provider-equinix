package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation_IsEmpty(t *testing.T) {
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
