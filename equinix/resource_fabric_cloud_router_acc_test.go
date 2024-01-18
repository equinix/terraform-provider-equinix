package equinix

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudRouterCreate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkCloudRouterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRouterCreateConfig("fg_tf_acc_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.test", "name", "fg_tf_acc_test"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccCloudRouterCreateConfig("fg_tf_acc_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.test", "name", "fg_tf_acc_update"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func checkCloudRouterDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, testAccProvider.Meta().(*config.Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_cloud_router" {
			continue
		}
		err := waitUntilCloudRouterDeprovisioned(rs.Primary.ID, testAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}

func testAccCloudRouterCreateConfig(name string) string {
	return fmt.Sprintf(`resource "equinix_fabric_cloud_router" "test"{
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

func TestAccCloudRouterRead(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRouterReadConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.test", "name", "fcr_tf_acc_test"),
				),
			},
		},
	})
}

func testAccCloudRouterReadConfig() string {
	return `data "equinix_fabric_cloud_router" "test" {
		uuid = "3e91216d-526a-45d2-9029-0c8c8ba48b60"
	}`
}
