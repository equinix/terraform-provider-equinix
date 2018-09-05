package packet

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccPacketOperatingSystem_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: testOperatingSystemConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.packet_operating_system.example", "id", "coreos_alpha"),
				),
			},
		},
	})
}

const testOperatingSystemConfig_Basic = `
	data "packet_operating_system" "example" {
		name    = "Container"
		distro  = "coreos"
		version = "alpha"
	  }`
