package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckMetalReservedIPBlockConfig_Global(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = metal_project.foobar.id
	type        = "global_ipv4"
	description = "testdesc"
	quantity    = 1
}`, name)
}

func testAccCheckMetalReservedIPBlockConfig_Public(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = metal_project.foobar.id
	facility    = "ewr1"
	type        = "public_ipv4"
	quantity    = 2
	tags        = ["Tag1", "Tag2"]
}`, name)
}

func testAccCheckMetalReservedIPBlockConfig_Metro(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = metal_project.foobar.id
	metro       = "sv"
	type        = "public_ipv4"
	quantity    = 2
}`, name)
}

func TestAccMetalReservedIPBlock_Global(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalReservedIPBlockConfig_Global(rs),
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

func TestAccMetalReservedIPBlock_Public(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalReservedIPBlockConfig_Public(rs),
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

func TestAccMetalReservedIPBlock_Metro(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalReservedIPBlockConfig_Metro(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "sv"),
				),
			},
		},
	})
}

func TestAccMetalReservedIPBlock_ImportBasic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalReservedIPBlockConfig_Public(rs),
			},
			{
				ResourceName:      "equinix_metal_reserved_ip_block.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMetalReservedIPBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_reserved_ip_block" {
			continue
		}
		if _, _, err := client.ProjectIPs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Reserved IP block still exists")
		}
	}

	return nil
}

func testAccCheckMetalReservedIPBlockConfig_FacilityToMetro(line string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block_fac_met_test"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = metal_project.foobar.id
	%s
	type        = "public_ipv4"
	quantity    = 2
	tags        = ["Tag1", "Tag2"]
}`, line)
}

func TestAccMetalReservedIPBlock_FacilityToMetro(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalReservedIPBlockConfig_FacilityToMetro(`   facility = "ny5"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "facility", "ny5"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "ny"),
				),
			},
			{
				Config: testAccCheckMetalReservedIPBlockConfig_FacilityToMetro(`   metro = "ny"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "ny"),
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccMetalReservedIP_Device(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = metal_project.foobar.id
	facility    = "ewr1"
	type        = "public_ipv4"
	quantity    = 2
}

resource "equinix_metal_device" "test" {
  project_id       = metal_project.foobar.id
  facilities       = ["ewr1"]
  plan             = "t1.small.x86"
  operating_system = "ubuntu_16_04"
  hostname         = "tfacc-reserved-ip-device"
  billing_cycle    = "hourly"
  ip_address {
	 type = "public_ipv4"
	 cidr = 31
	 reservation_ids = [metal_reserved_ip_block.test.id]
  }
  ip_address {
	 type = "private_ipv4"
  }
}
`, name)
}

func TestAccMetalReservedIPBlock_Device(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIP_Device(rs),
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
