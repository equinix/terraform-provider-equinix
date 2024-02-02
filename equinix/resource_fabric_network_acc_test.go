package equinix_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("equinix_fabric_connection_PFCR", &resource.Sweeper{
		Name: "equinix_fabric_network",
		F:    testSweepNetworks,
	})
}

func testSweepNetworks(region string) error {
	return nil
}

func TestAccFabricNetworkCreateOnlyRequiredParameters_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkNetworkDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkCreateOnlyRequiredParameterConfig_PFCR("Ipwan_tf_acc_PFCR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_fabric_network.test", "name", "Ipwan_tf_acc_PFCR"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "state", "INACTIVE"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "connections_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "type", "EVPLAN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "scope", "GLOBAL"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "operation.0.equinix_status"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccNetworkCreateOnlyRequiredParameterConfig_PFCR("Ipwan_update_PFCR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_network.test", "name", "Ipwan_update_PFCR"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
func testAccNetworkCreateOnlyRequiredParameterConfig_PFCR(name string) string {
	return fmt.Sprintf(`resource "equinix_fabric_network" "test"{
			type = "EVPLAN"
			name = "%s"
			scope = "GLOBAL"
			notifications{
				type = "ALL"
				emails = [
					"test@equinix.com",
					"test1@equinix.com"
				]
			}
		}`, name)
}
func checkNetworkDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, acceptance.TestAccProvider.Meta().(*config.Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_network" {
			continue
		}
		err := equinix.WaitUntilFabricNetworkDeprovisioned(rs.Primary.ID, acceptance.TestAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
func TestAccFabricNetworkCreateMixedParameters_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkNetworkDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkCreateMixedParameterConfig_PFCR(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_fabric_network.test", "name", "Tf_Network_PNFV"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "state", "INACTIVE"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "connections_count", "0"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "type", "EPLAN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "scope", "GLOBAL"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "location.0.metro_code", "SV"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "location.0.metro_name", "Silicon Valley"),
					resource.TestCheckResourceAttr("data.equinix_fabric_network.test", "location.0.region", "AMER"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_network.test", "operation.0.equinix_status"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
func testAccNetworkCreateMixedParameterConfig_PFCR() string {
	return fmt.Sprintf(`
	resource "equinix_fabric_network" "test" {
		type = "IPWAN"
		name = "Tf_Network_PNFV"
		scope = "GLOBAL"
		notifications {
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		location {
			region = "AMER"
			metro_code = "SV"
			metro_name = "Silicon Valley"
		}
		project{
			project_id = "291639000636552"
		}
	}
	data "equinix_fabric_cloud_router" "example"{
		uuid = equinix_fabric_cloud_router.example.id
	}
`)
}
