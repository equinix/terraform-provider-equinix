package equinix_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("equinix_fabric_network_PFCR", &resource.Sweeper{
		Name: "equinix_fabric_network",
		F:    testSweepNetworks,
	})
}

func testSweepNetworks(region string) error {
	return nil
}

func TestAccFabricNetworkCreateOnlyRequiredParameters_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkNetworkDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkCreateOnlyRequiredParameterConfig_PFCR("Ipwan_tf_acc_PFCR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_fabric_network.test", "name", "Ipwan_tf_acc_PFCR"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.test", "href"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.test", "state"),
					resource.TestCheckResourceAttr("equinix_fabric_network.test", "connections_count", "0"),
					resource.TestCheckResourceAttr("equinix_fabric_network.test", "type", "EVPLAN"),
					resource.TestCheckResourceAttr("equinix_fabric_network.test", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("equinix_fabric_network.test", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("equinix_fabric_network.test", "scope", "GLOBAL"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.test", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.test", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.test", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.test", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.test", "operation.0.equinix_status"),
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
				project {
					project_id = "291639000636552"
				}
			}`, name)
}
func checkNetworkDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_network" {
			continue
		}
		err := equinix.WaitUntilFabricNetworkDeprovisioned(rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, ctx, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}

func TestAccFabricNetworkCreateMixedParameters_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkNetworkDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkCreateMixedParameterConfig_PFCR(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_fabric_network.example2", "name", "Tf_Network_PFCR"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.example2", "state"),
					resource.TestCheckResourceAttr("equinix_fabric_network.example2", "connections_count", "0"),
					resource.TestCheckResourceAttr("equinix_fabric_network.example2", "type", "IPWAN"),
					resource.TestCheckResourceAttr("equinix_fabric_network.example2", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("equinix_fabric_network.example2", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("equinix_fabric_network.example2", "scope", "REGIONAL"),
					resource.TestCheckResourceAttr("equinix_fabric_network.example2", "location.0.region", "AMER"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.example2", "href"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.example2", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.example2", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.example2", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.example2", "change_log.0.created_by_email"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.example2", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_network.example2", "operation.0.equinix_status"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
func testAccNetworkCreateMixedParameterConfig_PFCR() string {
	return fmt.Sprintf(`
	resource "equinix_fabric_network" "example2" {
		type = "IPWAN"
		name = "Tf_Network_PFCR"
		scope = "REGIONAL"
		notifications {
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		location {
			region = "AMER"
		}
		project{
			project_id = "291639000636552"
		}
	}
`)
}
