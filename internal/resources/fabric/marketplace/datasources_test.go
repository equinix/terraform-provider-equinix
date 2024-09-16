package marketplace_test

import (
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	_ "github.com/hashicorp/terraform-plugin-testing/terraform"
	"log"
	"testing"
)

func TestAccFabricDataSourceMarketPlace_PFCR(t *testing.T) {
	susbcriptionId := testing_helpers.GetFabricMarketPlaceSubscriptionId(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configGetMarketplaceSubscriptionResource(susbcriptionId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_market_place_subscription.test", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_market_place_subscription.test", "href"),
					resource.TestCheckResourceAttr("data.equinix_fabric_market_place_subscription.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.equinix_fabric_market_place_subscription.test", "marketplace", "AWS"),
					resource.TestCheckResourceAttr("data.equinix_fabric_market_place_subscription.test", "offer_type", "PUBLIC"),
					resource.TestCheckResourceAttr("data.equinix_fabric_market_place_subscription.test", "is_auto_renew", "false"),
				),
			},
		},
	})

}
func configGetMarketplaceSubscriptionResource(subscription_id string) string {
	log.Printf("!! debugging")
	return fmt.Sprintf(`
	data "equinix_fabric_market_place_subscription" "test"{
		uuid = "%s"
	}
`, subscription_id)
}
