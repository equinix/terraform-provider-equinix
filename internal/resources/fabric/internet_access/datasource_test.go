package internetaccess_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFabricInternetAccessServicesSearchDataSource_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricInternetAccessDataSourcesConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_services.search",
						tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_services.search",
						tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("state"),
						knownvalue.StringExact("PROVISIONED"),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_services.search",
						tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("routing_protocol").AtMapKey("type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_services.search",
						tfjsonpath.New("pagination").AtMapKey("total"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccFabricInternetAccessServiceGetByUUIDDataSource_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricInternetAccessDataSourcesConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_service.by_id",
						tfjsonpath.New("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_service.by_id",
						tfjsonpath.New("state"),
						knownvalue.StringExact("PROVISIONED"),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_service.by_id",
						tfjsonpath.New("name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_service.by_id",
						tfjsonpath.New("routing_protocol").AtMapKey("connections"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccFabricInternetAccessDataSourcesConfig() string {
	return `
data "equinix_fabric_internet_access_services" "search" {
  filter = {
    property = "/state"
    operator = "="
    values   = ["PROVISIONED"]
  }
  pagination = {
    limit  = 1
    offset = 0
  }
  sort = {
    property  = "/name"
    direction = "DESC"
  }
}

data "equinix_fabric_internet_access_service" "by_id" {
  internet_access_service_id = data.equinix_fabric_internet_access_services.search.data[0].uuid
}
`
}
