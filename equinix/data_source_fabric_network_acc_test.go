package equinix_test

import (
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	_ "github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccDataSourceFabricNetwork_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ConfigCreateNetworkResource_PFCR(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "name", "Test_Network_PNFV"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "state", "INACTIVE"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "connections_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "type", "EVPLAN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "scope", "GLOBAL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "location.0.metro_code", "SV"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "operation.0.equinix_status"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func ConfigCreateNetworkResource_PFCR() string {
	return fmt.Sprintf(`
	resource "equinix_fabric_network" "example" {
		type = "EVPLAN"
		name = "Test_Network_PNFV"
		scope = "GLOBAL"
		notifications {
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		location {
			metro_code = "SV"
		}
	}
	data "equinix_fabric_cloud_router" "example"{
		uuid = equinix_fabric_cloud_router.example.id
	}
`)
}
