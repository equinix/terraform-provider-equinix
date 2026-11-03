package testinghelpers

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

const (
	FabricDedicatedPortEnvVar       = "TF_ACC_FABRIC_DEDICATED_PORTS"
	FabricConnectionsTestDataEnvVar = "TF_ACC_FABRIC_CONNECTIONS_TEST_DATA"
	FabricSubscriptionEnvVar        = "TF_ACC_FABRIC_MARKET_PLACE_SUBSCRIPTION_ID"
	FabricStreamEnvVar              = "TF_ACC_FABRIC_STREAM_TEST_DATA"
	FabricEiaEnvVar                 = "TF_ACC_FABRIC_EIA_TEST_DATA"
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

// GetFabricEiaTestData reads the shared test data of the IP block, ASN and Internet Access (EIA)
// tests from the TF_ACC_FABRIC_EIA_TEST_DATA env var.
// Expected JSON format:
//
//	{"pfcr": {"project_id": "<uuid>", "account_number": "<number>", "port_uuid": "<physical port uuid>",
//	          "service_profile_uuid": "<uuid of the Equinix Internet Access service profile>", "metro_code": "<metro of the port>"}}
func GetFabricEiaTestData(t *testing.T) map[string]map[string]string {
	var testData map[string]map[string]string
	raw := os.Getenv(FabricEiaEnvVar)
	if err := json.Unmarshal([]byte(raw), &testData); raw != "" && err != nil {
		t.Fatalf("failed reading EIA test data from environment: %v, %s", err, raw)
	}
	return testData
}
