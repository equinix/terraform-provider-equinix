package testinghelpers

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

const (
	FabricDedicatedPortEnvVar          = "TF_ACC_FABRIC_DEDICATED_PORTS"
	FabricConnectionsTestDataEnvVar    = "TF_ACC_FABRIC_CONNECTIONS_TEST_DATA"
	FabricSubscriptionEnvVar           = "TF_ACC_FABRIC_MARKET_PLACE_SUBSCRIPTION_ID"
	FabricStreamEnvVar                 = "TF_ACC_FABRIC_STREAM_TEST_DATA"
	FabricIpBlockEnvVar                = "TF_ACC_FABRIC_IP_BLOCK_TEST_DATA"
	FabricInternetAccessEnvVar         = "TF_ACC_FABRIC_INTERNET_ACCESS_TEST_DATA"
)

type EnvPorts map[string]map[string][]fabricv4.Port

func GetFabricEnvPorts(t *testing.T) EnvPorts {
	var ports EnvPorts
	portJSON := os.Getenv(FabricDedicatedPortEnvVar)
	if err := json.Unmarshal([]byte(portJSON), &ports); portJSON != "" && err != nil {
		t.Fatalf("Failed reading port data from environment: %v, %s", err, portJSON)
	}
	return ports
}

func GetFabricEnvConnectionTestData(t *testing.T) map[string]map[string]string {
	var connectionTestData map[string]map[string]string
	connectionTestDataJSON := os.Getenv(FabricConnectionsTestDataEnvVar)
	if err := json.Unmarshal([]byte(connectionTestDataJSON), &connectionTestData); connectionTestDataJSON != "" && err != nil {
		t.Fatalf("Failed reading connection data from environment: %v, %s", err, connectionTestDataJSON)
	}
	return connectionTestData
}

func GetFabricMarketPlaceSubscriptionID(_ *testing.T) string {
	subscriptionID := os.Getenv(FabricSubscriptionEnvVar)
	return subscriptionID
}

func GetFabricStreamTestData(t *testing.T) map[string]map[string]string {
	var streamTestData map[string]map[string]string
	streamJSON := os.Getenv(FabricStreamEnvVar)
	if err := json.Unmarshal([]byte(streamJSON), &streamTestData); streamJSON != "" && err != nil {
		t.Fatalf("failed reading stream data from environment: %v, %s", err, streamJSON)
	}
	return streamTestData
}

// GetFabricIpBlockTestData reads IP block test data from the TF_ACC_FABRIC_IP_BLOCK_TEST_DATA env var.
// Expected JSON format: {"pfcr": {"uuid": "<ip-block-uuid>", "prefix": "<prefix>", "ownership": "<ownership>"}}
func GetFabricIpBlockTestData(t *testing.T) map[string]map[string]string {
	var ipBlockTestData map[string]map[string]string
	ipBlockJSON := os.Getenv(FabricIpBlockEnvVar)
	if err := json.Unmarshal([]byte(ipBlockJSON), &ipBlockTestData); ipBlockJSON != "" && err != nil {
		t.Fatalf("failed reading ip block data from environment: %v, %s", err, ipBlockJSON)
	}
	return ipBlockTestData
}

// GetFabricInternetAccessTestData reads Internet Access test data from the TF_ACC_FABRIC_INTERNET_ACCESS_TEST_DATA env var.
// Expected JSON format: {"pfcr": {"project_id": "<project-uuid>"}}
func GetFabricInternetAccessTestData(t *testing.T) map[string]map[string]string {
	var testData map[string]map[string]string
	raw := os.Getenv(FabricInternetAccessEnvVar)
	if err := json.Unmarshal([]byte(raw), &testData); raw != "" && err != nil {
		t.Fatalf("failed reading internet access test data from environment: %v, %s", err, raw)
	}
	return testData
}
