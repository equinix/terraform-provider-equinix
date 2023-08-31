package tests

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCloudRouterCreate(t *testing.T) {
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/fabric-cloud-router",
	})

	terraform.InitAndApply(t, terraformOptions)
	output := terraform.Output(t, terraformOptions, "fcr_result")
	assert.NotNil(t, output)

	terraformOptions = terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/fabric-cloud-router/fcr",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
	output = terraform.Output(t, terraformOptions, "fcr_result")

	assert.NotNil(t, output)
}
