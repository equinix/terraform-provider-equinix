package ipblock_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFabricIpBlocksSearchDataSource_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricIpBlockDataSourcesConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_ip_blocks.search",
						tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_ip_blocks.search",
						tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("ownership"),
						knownvalue.StringExact("EQUINIX"),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_ip_blocks.search",
						tfjsonpath.New("pagination").AtMapKey("total"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccFabricIpBlockGetByUUIDDataSource_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricIpBlockDataSourcesConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_ip_block.by_id",
						tfjsonpath.New("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_ip_block.by_id",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccFabricIpBlockDataSourcesConfig() string {
	return `
data "equinix_fabric_ip_blocks" "search" {
  filter = {
    property = "/ownership"
    operator = "="
    values   = ["EQUINIX"]
  }
  pagination = {
    limit  = 1
    offset = 0
  }
  sort = {
    property  = "/changeLog/createdDateTime"
    direction = "DESC"
  }
}

data "equinix_fabric_ip_block" "by_id" {
  ip_block_id = data.equinix_fabric_ip_blocks.search.data[0].uuid
}
`
}
