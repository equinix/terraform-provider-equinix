package equinix

import (
	"context"
	"fmt"
	"testing"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFabricGatewayCreate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkFabricGatewayDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricGatewayCreateConfig("fg_tf_acc_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_gateway.test", "name", fmt.Sprint("fg_tf_acc_test")),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccFabricGatewayCreateConfig("fg_tf_acc_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_gateway.test", "name", fmt.Sprint("fg_tf_acc_update")),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func checkFabricGatewayDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, testAccProvider.Meta().(*Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_gateway" {
			continue
		}
		err := waitUntilFGDeprovisioned(rs.Primary.ID, testAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}

func testAccFabricGatewayCreateConfig(name string) string {
	return fmt.Sprintf(`resource "equinix_fabric_gateway" "test"{
			type = "XF_GATEWAY"
			name = "%s"
			location{
			  metro_code  = "SV"
			}
			package{
				  code = "PRO"
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
				   project_id = "776847000642406"
			}
			  account {
				account_number = 203612
			}
		}`, name)
}

func TestAccFabricGatewayRead(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricGatewayReadConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_gateway.test", "name", fmt.Sprint("fg_tf_acc_test")),
				),
			},
		},
	})
}

func testAccFabricGatewayReadConfig() string {
	return fmt.Sprint(`data "equinix_fabric_gateway" "test" {
		uuid = "3e91216d-526a-45d2-9029-0c8c8ba48b60"
	}`)
}
