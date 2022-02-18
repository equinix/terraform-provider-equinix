package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccMetalReservedIPBlockConfig_global(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	type        = "global_ipv4"
	description = "testdesc"
	quantity    = 1
}`, name)
}

func testAccMetalReservedIPBlockConfig_public(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	facility    = "ewr1"
	type        = "public_ipv4"
	quantity    = 2
	tags        = ["Tag1", "Tag2"]
}`, name)
}

func testAccMetalReservedIPBlockConfig_metro(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	metro       = "sv"
	type        = "public_ipv4"
	quantity    = 2
}`, name)
}

func TestAccMetalReservedIPBlock_global(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_global(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "quantity", "1"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "description", "testdesc"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "type", "global_ipv4"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "netmask", "255.255.255.255"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "public", "true"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "management", "false"),
				),
			},
		},
	})
}

func TestAccMetalReservedIPBlock_public(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_public(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "facility", "ewr1"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "type", "public_ipv4"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "quantity", "2"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "netmask", "255.255.255.254"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "public", "true"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "management", "false"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccMetalReservedIPBlock_metro(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_metro(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "sv"),
				),
			},
		},
	})
}

func TestAccMetalReservedIPBlock_importBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_public(rs),
			},
			{
				ResourceName:      "equinix_metal_reserved_ip_block.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalReservedIPBlockCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_reserved_ip_block" {
			continue
		}
		if _, _, err := client.ProjectIPs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Reserved IP block still exists")
		}
	}

	return nil
}

func testAccMetalReservedIPBlockConfig_facilityToMetro(line string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block_fac_met_test"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	%s
	type        = "public_ipv4"
	quantity    = 2
	tags        = ["Tag1", "Tag2"]
}`, line)
}

func TestAccMetalReservedIPBlock_facilityToMetro(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_facilityToMetro(`   facility = "ny5"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "facility", "ny5"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "ny"),
				),
			},
			{
				Config: testAccMetalReservedIPBlockConfig_facilityToMetro(`   metro = "ny"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "ny"),
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccMetalReservedIP_device(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	facility    = "ewr1"
	type        = "public_ipv4"
	quantity    = 2
}

resource "equinix_metal_device" "test" {
  project_id       = equinix_metal_project.foobar.id
  facilities       = ["ewr1"]
  plan             = "t1.small.x86"
  operating_system = "ubuntu_16_04"
  hostname         = "tfacc-reserved-ip-device"
  billing_cycle    = "hourly"
  ip_address {
	 type = "public_ipv4"
	 cidr = 31
	 reservation_ids = [equinix_metal_reserved_ip_block.test.id]
  }
  ip_address {
	 type = "private_ipv4"
  }
}
`, name)
}

func TestAccMetalReservedIPBlock_device(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIP_device(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_reserved_ip_block.test", "gateway",
						"equinix_metal_device.test", "network.0.gateway",
					),
				),
			},
		},
	})
}
