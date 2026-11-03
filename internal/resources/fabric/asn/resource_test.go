package asn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func asnTestProjectID(t *testing.T) string {
	t.Helper()
	pfcr, ok := testinghelpers.GetFabricEiaTestData(t)["pfcr"]
	if !ok {
		t.Skip("skipping: no 'pfcr' key in TF_ACC_FABRIC_EIA_TEST_DATA")
	}
	projectID := pfcr["project_id"]
	if projectID == "" {
		t.Skip("skipping: TF_ACC_FABRIC_EIA_TEST_DATA missing project_id")
	}
	return projectID
}

func TestAccFabricAsnPrivateCreate_PFCR(t *testing.T) {
	projectID := asnTestProjectID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkAsnDeprovisioned,
		Steps: []resource.TestStep{
			{
				Config: testAccAsnPrivateConfig(projectID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_asn.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_asn.test", tfjsonpath.New("state"), knownvalue.StringExact("PROVISIONED")),
					statecheck.ExpectKnownValue("equinix_fabric_asn.test", tfjsonpath.New("number"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_asn.test", tfjsonpath.New("range"), knownvalue.StringExact("16_BIT")),
					statecheck.ExpectKnownValue("equinix_fabric_asn.test", tfjsonpath.New("location").AtMapKey("metro_code"), knownvalue.StringExact("SV")),

					statecheck.ExpectKnownValue("data.equinix_fabric_asn.test", tfjsonpath.New("state"), knownvalue.StringExact("PROVISIONED")),
				},
			},
			{
				ResourceName:            "equinix_fabric_asn.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func testAccAsnPrivateConfig(projectID string) string {
	return fmt.Sprintf(`
resource "equinix_fabric_asn" "test" {
  type  = "PRIVATE_ASN"
  range = "16_BIT"
  location = {
    metro_code = "SV"
  }
  project = {
    project_id = %q
  }
}

data "equinix_fabric_asn" "test" {
  asn_id = equinix_fabric_asn.test.uuid
}
`, projectID)
}

func checkAsnDeprovisioned(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_asn" {
			continue
		}
		asn, resp, err := client.ASNsApi.GetAsnById(ctx, rs.Primary.ID).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				continue
			}
			return fmt.Errorf("error checking ASN %s: %w", rs.Primary.ID, err)
		}
		if asn.GetState() != fabricv4.ASNSTATE_DEPROVISIONED {
			return fmt.Errorf("ASN %s still exists in state %s", rs.Primary.ID, asn.GetState())
		}
	}
	return nil
}
