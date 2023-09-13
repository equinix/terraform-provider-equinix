package tests

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCloudRouter2AzureTwoConnectionsCreateConnection(t *testing.T) {
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/fabric/v4/cloudRouterConnectivity/cloudRouter2azure/two_connections",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
	primaryConnection := terraform.Output(t, terraformOptions, "primary_connection_result")
	secondaryConnection := terraform.Output(t, terraformOptions, "secondary_connection_result")
	assert.NotNil(t, primaryConnection)
	assert.NotNil(t, secondaryConnection)
}
