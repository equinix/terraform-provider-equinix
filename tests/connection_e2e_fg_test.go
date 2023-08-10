package tests

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestFabricGatewayCreate(t *testing.T) {
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/fabric-gateway",
	})

	terraform.InitAndApply(t, terraformOptions)
	output := terraform.Output(t, terraformOptions, "fg_result")
	assert.NotNil(t, output)

	terraformOptions = terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/fg2port",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
	output = terraform.Output(t, terraformOptions, "fg2port_result")

	assert.NotNil(t, output)
}
