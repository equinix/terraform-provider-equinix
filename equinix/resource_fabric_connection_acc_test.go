package equinix

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccFabricConnection_basic(t *testing.T) {
	//var conn v4.Connection
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricConnectionConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					//testAccMetalVRFExists("equinix_metal_vrf.test", &conn),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", fmt.Sprintf("tfacc-terra-cond-%d", rInt)),
					resource.TestCheckResourceAttrSet(
						"equinix_fabric_connection.test", "local_asn"),
				),
			},
			{
				ResourceName:      "equinix_fabric_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFabricConnectionConfig_basic(r int) string {

	return fmt.Sprintln(`
resource "equinix_fabric_connection" "test" {
    name = "tfacc-terra-cond-%d"
}
`, r)
}

/*func testAccMetalVRFCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_vrf" {
			continue
		}
		if _, _, err := client.VRFs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal VRF still exists")
		}
	}

	return nil
}*/

/*func testAccMetalVRFExists(n string, vrf *packngo.VRF) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*Config).metal

		foundResource, _, err := client.VRFs.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundResource.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundResource)
		}

		*vrf = *foundResource

		return nil
	}
}*/
