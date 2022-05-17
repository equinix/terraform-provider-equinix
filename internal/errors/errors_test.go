package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/equinix/rest-go"
	"github.com/stretchr/testify/assert"
)

func TestIsRestNotFoundError(t *testing.T) {
	// given
	input := []error{
		rest.Error{HTTPCode: http.StatusNotFound, Message: "Not Found"},
		rest.Error{HTTPCode: http.StatusInternalServerError, Message: "Internal Server Error"},
		fmt.Errorf("some bogus error"),
	}
	expected := []bool{
		true,
		false,
		false,
	}
	// when
	result := make([]bool, len(input))
	for i := range input {
		result[i] = IsRestNotFoundError(input[i])
	}
	// then
	assert.Equal(t, expected, result, "Result matches expected output")
}
