package tests

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestPort2ProfilePublicCreateConnection(t *testing.T) {
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/fabric/v4/portConnectivity/port2serviceprofilepublic",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
	output := terraform.Output(t, terraformOptions, "connection_result")
	terraform.OutputAll(t, terraformOptions)
	assert.NotNil(t, output)
}
