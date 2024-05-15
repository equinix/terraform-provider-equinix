package equinix_test

import (
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	_ "github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccDataSourceFabricCloudRouter_PFCR(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ConfigCreateCloudRouterResource_PFCR(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "name", "Test_PFCR"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "type", "XF_ROUTER"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "order.0.purchase_order_number", "1-323292"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "location.0.metro_code", "SV"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "package.0.code", "STANDARD"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "project.0.project_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "account.0.account_number"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "href"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "state", "PROVISIONED"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "equinix_asn"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "bgp_ipv4_routes_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "bgp_ipv6_routes_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_router.example", "connections_count", "0"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "change_log.0.updated_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_router.example", "change_log.0.deleted_date_time"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.name", "Test_PFCR"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.type", "XF_ROUTER"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.order.0.purchase_order_number", "1-323292"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.location.0.metro_code", "SV"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.package.0.code", "STANDARD"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.project.0.project_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.account.0.account_number"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.href"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.state", "PROVISIONED"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.equinix_asn"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.bgp_ipv4_routes_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.bgp_ipv6_routes_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_cloud_routers.test", "data.0.connections_count", "0"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.change_log.0.updated_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_cloud_routers.test", "data.0.change_log.0.deleted_date_time"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func ConfigCreateCloudRouterResource_PFCR() string {
	return fmt.Sprintf(`
		resource "equinix_fabric_cloud_router" "example" {
		name = "Test_PFCR"
		type = "XF_ROUTER"
		notifications{
			type="ALL"
			emails= ["test@equinix.com"]
		}
		order {
			purchase_order_number= "1-323292"
		}
		location {
			metro_code= "SV"
		}
		package {
			code="STANDARD"
		}
		project {
			project_id = "291639000636552"
		}
		account {
			account_number = 201257
		}

	}
	data "equinix_fabric_cloud_router" "example"{
		uuid = equinix_fabric_cloud_router.example.id
	}

	data "equinix_fabric_cloud_routers" "test" {
		filter {
			property = "/name"
			operator = "="
			values 	 = [equinix_fabric_cloud_router.example.name]
		}
	}
`)
}
