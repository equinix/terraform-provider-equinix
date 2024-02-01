package equinix

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkCreate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkNetworkDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkCreateConfig("Ipwan_tf_acc_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_network.test", "name", "network_acc_test"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccNetworkCreateConfig("Ipwan_tf_acc_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_network.test", "name", "network_name_update"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func checkNetworkDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, testAccProvider.Meta().(*config.Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_network" {
			continue
		}
		err := waitUntilFabricNetworkDeprovisioned(rs.Primary.ID, testAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}

func testAccNetworkCreateConfig(name string) string {
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

func TestAccNetworkRead(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkReadConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_network.test", "name", "network_acc_test"),
				),
			},
		},
	})
}

func testAccNetworkReadConfig() string {
	return `data "equinix_fabric_network" "test" {
		uuid = ""
	}`
}
