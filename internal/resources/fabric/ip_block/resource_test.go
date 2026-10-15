package ipblock_test

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

func TestAccFabricIpBlockCreate_PFCR(t *testing.T) {
	testData := testinghelpers.GetFabricIpBlockTestData(t)
	pfcr, ok := testData["pfcr"]
	if !ok {
		t.Skip("skipping: no 'pfcr' key in TF_ACC_FABRIC_IP_BLOCK_TEST_DATA")
	}
	projectID := pfcr["project_id"]
	prefix := pfcr["prefix"]
	if projectID == "" || prefix == "" {
		t.Skip("skipping: TF_ACC_FABRIC_IP_BLOCK_TEST_DATA missing project_id or prefix")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkIpBlockDeleted,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricIpBlockResourceConfig(projectID, prefix),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test",
						tfjsonpath.New("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("ACTIVE"),
					),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test",
						tfjsonpath.New("href"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("IPV4_IP_BLOCK"),
					),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test",
						tfjsonpath.New("prefix"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test",
						tfjsonpath.New("project"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"project_id": knownvalue.StringExact(projectID),
						}),
					),
				},
			},
		},
	})
}

func testAccFabricIpBlockResourceConfig(projectID, prefix string) string {
	return fmt.Sprintf(`
resource "equinix_fabric_ip_block" "test" {
  type    = "IPV4_IP_BLOCK"
  prefix  = %q
  project = {
    project_id = %q
  }
}
`, prefix, projectID)
}

func checkIpBlockDeleted(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_ip_block" {
			continue
		}
		ipBlock, _, err := client.IPBlocksApi.GetIpBlock(ctx, rs.Primary.ID).Execute()
		if err == nil && ipBlock.GetState() == fabricv4.IPBLOCKSTATE_ACTIVE {
			return fmt.Errorf("IP block %s still exists and is ACTIVE", rs.Primary.ID)
		}
	}
	return nil
}
