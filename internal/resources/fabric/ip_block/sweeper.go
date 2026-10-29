package ipblock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// sweepMinAge keeps the sweeper away from blocks owned by tests that are still running.
const sweepMinAge = 30 * time.Minute

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_ip_block", &resource.Sweeper{
		Name:         "equinix_fabric_ip_block",
		Dependencies: []string{},
		F:            testSweepIpBlocks,
	})
}

// testSweepIpBlocks deletes ACTIVE IP blocks created during acceptance tests.
// IP blocks have no name field and the API cannot filter by prefix, so we search by state=ACTIVE
// within the test project and keep only blocks whose prefix is in testinghelpers.IpBlockTestPool
// and that are older than sweepMinAge.
// Requires TF_ACC_FABRIC_IP_BLOCK_TEST_DATA JSON:
//
//	{"pfcr": {"project_id": "<uuid>", "metro_code": "SV"}}
func testSweepIpBlocks(_ string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping IP Blocks")

	ctx := context.Background()
	meta, err := testinghelpers.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping IP Blocks: %s", err)
	}
	if configLoadErr := meta.Load(ctx); configLoadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping IP Blocks: %s", configLoadErr)
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	projectID := sweepProjectID()
	if projectID == "" {
		log.Printf("[WARN] Skipping IP block sweep: no project_id found. Set TF_ACC_FABRIC_IP_BLOCK_TEST_DATA with {'pfcr': {'project_id': '<uuid>'}}")
		return nil
	}

	searchReq := fabricv4.IpBlocksSearchRequestBody{
		Filter: &fabricv4.IpBlockFilter{
			And: []fabricv4.IpBlockAndQuery{
				{
					Property: "/state",
					Operator: fabricv4.IPBLOCKANDQUERYOPERATOR_EQUAL,
					Values:   []string{string(fabricv4.IPBLOCKSTATE_ACTIVE)},
				},
				{
					Property: "/project/projectId",
					Operator: fabricv4.IPBLOCKANDQUERYOPERATOR_EQUAL,
					Values:   []string{projectID},
				},
			},
		},
	}

	result, _, err := fabric.IPBlocksApi.SearchIpBlocks(ctx).IpBlocksSearchRequestBody(searchReq).Execute()
	if err != nil {
		return fmt.Errorf("error searching IP blocks for sweep: %s", err)
	}

	for _, block := range result.GetData() {
		if !testinghelpers.IsIpBlockTestPrefix(block.GetPrefix()) {
			log.Printf("[DEBUG] Skipping IP block %s (prefix %q is outside the test pool %s)", block.GetUuid(), block.GetPrefix(), testinghelpers.IpBlockTestPool)
			continue
		}
		if age := time.Since(block.ChangeLog.GetCreatedDateTime()); age < sweepMinAge {
			log.Printf("[DEBUG] Skipping IP block %s: created %s ago, newer than %s", block.GetUuid(), age.Round(time.Second), sweepMinAge)
			continue
		}
		log.Printf("[DEBUG] Deleting IP block: %s (prefix: %s)", block.GetUuid(), block.GetPrefix())
		_, resp, delErr := fabric.IPBlocksApi.DeleteIpBlockById(ctx, block.GetUuid()).Execute()
		if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, delErr) != nil {
			errs = append(errs, fmt.Errorf("error deleting IP block %s: %s", block.GetUuid(), delErr))
		}
	}

	return errors.Join(errs...)
}

// sweepProjectID reads the test project ID from TF_ACC_FABRIC_IP_BLOCK_TEST_DATA.
// Expected JSON: {"pfcr": {"project_id": "<uuid>", ...}, ...}
func sweepProjectID() string {
	raw := os.Getenv(testinghelpers.FabricIpBlockEnvVar)
	if raw == "" {
		return ""
	}

	var data map[string]map[string]string
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		log.Printf("[WARN] sweeper: could not parse %s: %s", testinghelpers.FabricIpBlockEnvVar, err)
		return ""
	}

	for _, env := range []string{"pfcr", "pnfv", "ppds"} {
		if v, ok := data[env]; ok {
			if id, ok := v["project_id"]; ok && id != "" {
				return id
			}
		}
	}
	return ""
}
