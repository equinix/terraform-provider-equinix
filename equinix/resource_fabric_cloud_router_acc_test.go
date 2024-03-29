package equinix_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("equinix_fabric_cloud_router_PFCR", &resource.Sweeper{
		Name: "equinix_fabric_cloud_router",
		F:    testSweepCloudRouters,
	})
}

func testSweepCloudRouters(region string) error {
	return nil
}

func TestAccCloudRouterCreateOnlyRequiredParameters_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkCloudRouterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRouterCreateOnlyRequiredParameterConfig_PFCR("fcr_acctest_PFCR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "name", "fcr_acctest_PFCR"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "type", "XF_ROUTER"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "order.0.purchase_order_number", "1-234567"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "location.0.metro_code", "SV"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "package.0.code", "STANDARD"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "project.0.project_id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "account.0.account_number"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "state", "PROVISIONED"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "equinix_asn"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "bgp_ipv4_routes_count", "0"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "bgp_ipv6_routes_count", "0"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "connections_count", "0"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.updated_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.deleted_date_time"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccCloudRouterCreateOnlyRequiredParameterConfig_PFCR("fcr_update_PFCR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.test", "name", "fcr_update_PFCR"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccCloudRouterCreateOnlyRequiredParameterConfig_PFCR(name string) string {
	return fmt.Sprintf(`resource "equinix_fabric_cloud_router" "test"{
		type = "XF_ROUTER"
		name = "%s"
		location{
			metro_code  = "SV"
		}
		package{
			code = "STANDARD"
		}
		order{
			purchase_order_number = "1-234567"
		}
		notifications{
			type = "ALL"
			emails = [
				"test@equinix.com",
				"test1@equinix.com"
			]
		}
		project{
			project_id = "291639000636552"
		}
		account {
			account_number = 201257
		}
	}`, name)
}

func TestAccCloudRouterCreateMixedParameters_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkCloudRouterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRouterCreateMixedParameterConfig_PFCR(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.example", "name", "fcr_acc_test_PFCR"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "type", "XF_ROUTER"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "order.0.purchase_order_number", "1-234567"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "location.0.metro_code", "SV"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "location.0.region", "AMER"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "location.0.metro_name", "Silicon Valley"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "package.0.code", "STANDARD"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "project.0.project_id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "account.0.account_number"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "state", "PROVISIONED"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "equinix_asn"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "bgp_ipv4_routes_count", "0"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "bgp_ipv6_routes_count", "0"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "connections_count", "0"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.updated_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.deleted_date_time"),
				),
			},
		},
	})
}
func testAccCloudRouterCreateMixedParameterConfig_PFCR() string {
	return fmt.Sprintf(`resource "equinix_fabric_cloud_router" "example"{
		type = "XF_ROUTER"
		name = "fcr_acc_test_PFCR"
		location{
			region      = "AMER"
			metro_code  = "SV"
			metro_name = "Silicon Valley"
		}
		package{
			code = "STANDARD"
		}
		order{
			purchase_order_number = "1-234567"
		}
		notifications{
			type = "ALL"
			emails = [
				"test@equinix.com",
				"test1@equinix.com"
					]
		}
		project{
			project_id = "291639000636552"
		}
		account {
			account_number = 201257
		}
	}`)
}

func checkCloudRouterDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, acceptance.TestAccProvider.Meta().(*config.Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_cloud_router" {
			continue
		}
		err := equinix.WaitUntilCloudRouterDeprovisioned(rs.Primary.ID, acceptance.TestAccProvider.Meta(), ctx, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
