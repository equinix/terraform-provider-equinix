package asn_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFabricAsnsSearchDataSource_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricAsnDataSourcesConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_asns.search",
						tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_asns.search",
						tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_asns.search",
						tfjsonpath.New("pagination").AtMapKey("total"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccFabricAsnGetByUUIDDataSource_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricAsnDataSourcesConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_asn.by_id",
						tfjsonpath.New("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_asn.by_id",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_asn.by_id",
						tfjsonpath.New("project").AtMapKey("project_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

// The search matches any ASN that is not in a terminal state, so the test does not
// depend on a particular provisioning state of the ASNs present in the environment.
func testAccFabricAsnDataSourcesConfig() string {
	return `
data "equinix_fabric_asns" "search" {
  filter = {
    property = "/state"
    operator = "IN"
    values   = ["REQUESTED", "PENDING", "PROVISIONING", "PROVISIONED", "ATTACHED"]
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

data "equinix_fabric_asn" "by_id" {
  asn_id = data.equinix_fabric_asns.search.data[0].uuid
}
`
}
