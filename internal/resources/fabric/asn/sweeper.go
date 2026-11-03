package asn

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

const sweepMinAge = 30 * time.Minute

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_asn", &resource.Sweeper{
		Name:         "equinix_fabric_asn",
		Dependencies: []string{},
		F:            testSweepAsns,
	})
}

func isTestAsn(asn *fabricv4.Asn) bool {
	asnType, _ := asn.AdditionalProperties[typeKey].(string)
	return asnType == "PRIVATE_ASN"
}

func testSweepAsns(_ string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping ASNs")

	ctx := context.Background()
	meta, err := testinghelpers.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping ASNs: %s", err)
	}
	if configLoadErr := meta.Load(ctx); configLoadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping ASNs: %s", configLoadErr)
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	projectID := sweepProjectID()
	if projectID == "" {
		log.Printf("[WARN] Skipping ASN sweep: no project_id found. Set TF_ACC_FABRIC_IP_BLOCK_TEST_DATA with {'pfcr': {'project_id': '<uuid>'}}")
		return nil
	}

	property, operator := "/state", fabricv4.SearchExpressionOperator("IN")
	projectProperty, equal := "/project/projectId", fabricv4.SearchExpressionOperator("=")
	searchReq := fabricv4.AsnSearchRequestBody{
		Filter: &fabricv4.SearchExpression{
			And: []fabricv4.SearchExpression{
				{
					Property: &property,
					Operator: &operator,
					Values: []string{
						string(fabricv4.ASNSTATE_REQUESTED),
						string(fabricv4.ASNSTATE_PROVISIONING),
						string(fabricv4.ASNSTATE_PROVISIONED),
						string(fabricv4.ASNSTATE_ATTACHED),
						string(fabricv4.ASNSTATE_PROVISIONING_FAILED),
					},
				},
				{
					Property: &projectProperty,
					Operator: &equal,
					Values:   []string{projectID},
				},
			},
		},
	}

	result, _, err := fabric.ASNsApi.SearchAsns(ctx).AsnSearchRequestBody(searchReq).Execute()
	if err != nil {
		return fmt.Errorf("error searching ASNs for sweep: %s", err)
	}

	for _, asn := range result.GetData() {
		if !isTestAsn(&asn) {
			log.Printf("[DEBUG] Skipping ASN %s (number %d is not a test ASN)", asn.GetUuid(), asn.GetNumber())
			continue
		}
		if age := time.Since(asn.ChangeLog.GetCreatedDateTime()); age < sweepMinAge {
			log.Printf("[DEBUG] Skipping ASN %s: created %s ago, newer than %s", asn.GetUuid(), age.Round(time.Second), sweepMinAge)
			continue
		}
		log.Printf("[DEBUG] Deleting ASN: %s (number: %d)", asn.GetUuid(), asn.GetNumber())
		_, resp, delErr := fabric.ASNsApi.DeleteAsnById(ctx, asn.GetUuid()).Execute()
		if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, delErr) != nil {
			errs = append(errs, fmt.Errorf("error deleting ASN %s: %s", asn.GetUuid(), delErr))
		}
	}

	return errors.Join(errs...)
}

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
		if v, ok := data[env]; ok && v["project_id"] != "" {
			return v["project_id"]
		}
	}
	return ""
}
