package internetaccess_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFabricInternetAccessServicesSearchDataSource_PFCR(t *testing.T) {
	testData := testinghelpers.GetFabricInternetAccessTestData(t)
	if testData["pfcr"]["project_id"] == "" {
		t.Skipf("TF_ACC_FABRIC_INTERNET_ACCESS_TEST_DATA not set or missing pfcr.project_id")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricInternetAccessDataSourcesConfig(testData["pfcr"]["project_id"]),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_services.search",
						tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("uuid"),
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
	testData := testinghelpers.GetFabricInternetAccessTestData(t)
	if testData["pfcr"]["project_id"] == "" {
		t.Skipf("TF_ACC_FABRIC_INTERNET_ACCESS_TEST_DATA not set or missing pfcr.project_id")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricInternetAccessDataSourcesConfig(testData["pfcr"]["project_id"]),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_service.by_id",
						tfjsonpath.New("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_internet_access_service.by_id",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccFabricInternetAccessDataSourcesConfig(projectID string) string {
	return fmt.Sprintf(`
data "equinix_fabric_internet_access_services" "search" {
  filter = [
    {
      property = "/project/projectId"
      operator = "="
      values   = [%q]
    },
    {
      property = "/state"
      operator = "IN"
      values   = ["ACTIVE", "PROVISIONED", "PROVISIONING"]
    }
  ]
  pagination = {
    limit  = 1
    offset = 0
  }
}

data "equinix_fabric_internet_access_service" "by_id" {
  internet_access_service_id = data.equinix_fabric_internet_access_services.search.data[0].uuid
}
`, projectID)
}
