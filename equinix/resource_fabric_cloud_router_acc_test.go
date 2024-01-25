package equinix_test

import (
	"context"
	"fmt"
	"testing"

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
				Config: testAccCloudRouterCreateOnlyRequiredParameterConfig("fcr_tf_acc_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.test", "name", "fcr_tf_acc_test"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccCloudRouterCreateOnlyRequiredParameterConfig("fcr_tf_acc_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.test", "name", "fcr_tf_acc_update"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccCloudRouterCreateOnlyRequiredParameterConfig(name string) string {
	return fmt.Sprintf(`resource "equinix_fabric_cloud_router" "test"{
		type = "XF_ROUTER"
		name = "%s"
		location{
			metro_code  = "SV"
		}
		package{
			code = "LAB"
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
				Config: testAccCloudRouterCreateMixedParameterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.example", "name", "fcr_acc_test"),
				),
			},
		},
	})
}
func testAccCloudRouterCreateMixedParameterConfig() string {
	return fmt.Sprintf(`resource "equinix_fabric_cloud_router" "example"{
		type = "XF_ROUTER"
		name = "fcr_acc_test"
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
		err := equinix.WaitUntilCloudRouterDeprovisioned(rs.Primary.ID, acceptance.TestAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
