package testing_helpers

import (
	"encoding/json"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"os"
	"testing"
)

const (
	FabricDedicatedPortEnvVar = "TF_ACC_FABRIC_DEDICATED_PORTS"
)

const (
	FabricConnectionsTestDataEnvVar = "TF_ACC_FABRIC_CONNECTIONS_TEST_DATA"
)

type EnvPorts map[string]map[string][]fabricv4.PortResponse

func GetFabricEnvPorts(t *testing.T) EnvPorts {
	var ports EnvPorts
	portJson := os.Getenv(FabricDedicatedPortEnvVar)
	if err := json.Unmarshal([]byte(portJson), &ports); portJson != "" && err != nil {
		t.Fatalf("Failed reading port data from environment: %v, %s", err, portJson)
	}
	return ports
}

func GetFabricEnvConnectionTestData(t *testing.T) map[string]map[string]string {
	var connectionTestData map[string]map[string]string
	connectionTestDataJson := os.Getenv(FabricConnectionsTestDataEnvVar)
	if err := json.Unmarshal([]byte(connectionTestDataJson), &connectionTestData); connectionTestDataJson != "" && err != nil {
		t.Fatalf("Failed reading connection data from environment: %v, %s", err, connectionTestDataJson)
	}
	return connectionTestData
}
