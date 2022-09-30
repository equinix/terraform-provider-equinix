package tests

import (
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateServiceProfileGeneric(t *testing.T) {
	// retryable errors in terraform testing.

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/serviceprofile/generic", Logger: logger.TestingT,
	})
	logger.Log(t, "Testing loging ")
	t.Log("Testing t logger")

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.OutputAll(t, terraformOptions)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output)

}
