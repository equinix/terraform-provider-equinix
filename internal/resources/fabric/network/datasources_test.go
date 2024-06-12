package network_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceFabricNetwork_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configCreateNetworkResource_PFCR(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "name", "Test_Network_PFCR"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "state"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "connections_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "type", "EVPLAN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.example", "scope", "GLOBAL"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.example", "operation.0.equinix_status"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_networks.example", "data.0.href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_networks.example", "data.0.uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_networks.example", "data.0.name", "Test_Network_PFCR"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_networks.example", "data.0.state"),
					resource.TestCheckResourceAttr("data.equinix_fabric_networks.example", "data.0.connections_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_networks.example", "data.0.type", "EVPLAN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_networks.example", "data.0.notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_networks.example", "data.0.notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("data.equinix_fabric_networks.example", "data.0.scope", "GLOBAL"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_networks.example", "data.0.change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_networks.example", "data.0.change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_networks.example", "data.0.change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_networks.example", "data.0.change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_networks.example", "data.0.operation.0.equinix_status"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func configCreateNetworkResource_PFCR() string {
	return fmt.Sprintf(`
	resource "equinix_fabric_network" "example" {
		type = "EVPLAN"
		name = "Test_Network_PFCR"
		scope = "GLOBAL"
		notifications {
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		project{
			project_id = "291639000636552"
		}
	}
	data "equinix_fabric_network" "example"{
		uuid = equinix_fabric_network.example.id
	}
	data "equinix_fabric_networks" "example" {
		outer_operator = "AND"
		filter {
			property = "/type"
			operator = "="
			values 	 = ["EVPLAN"]
		}
		filter {
			property = "/name"
			operator = "="
			values   = ["Test_Network_PFCR"]
		}
		filter {
			property = "/uuid"
			operator = "="
			values   = [equinix_fabric_network.example.id]
		}
	}
`)
}
