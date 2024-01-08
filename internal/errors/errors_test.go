package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/assert"
)

func TestProvider_HasApplicationErrorCode(t *testing.T) {
	// given
	code := "ERR-505"
	errors := []rest.ApplicationError{
		{
			Code: "ERR-505",
		},
		{
			Code: acctest.RandString(10),
		},
	}
	// when
	result := HasApplicationErrorCode(errors, code)
	// then
	assert.True(t, result, "Error list contains error with given code")
}

func TestProvider_IsRestNotFoundError(t *testing.T) {
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
