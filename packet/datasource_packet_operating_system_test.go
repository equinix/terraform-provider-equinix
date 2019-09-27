package packet

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

var matchErrOSNotFound = regexp.MustCompile(".*There are no operating systems*")

func TestAccPacketOperatingSystem_NotFound(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: testOperatingSystemConfig_NotFound,
				ExpectError: matchErrOSNotFound,
			},
		},
	})
}

const testOperatingSystemConfig_NotFound = `
	data "packet_operating_system" "example" {
		name    = "Container"
		distro  = "NOTEXISTS"
		version = "alpha"
	  }`

var matchErrOSAmbiguous = regexp.MustCompile(".*There is more than one operating system.*")

func TestAccPacketOperatingSystem_Ambiguous(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: testOperatingSystemConfig_Ambiguous,
				ExpectError: matchErrOSAmbiguous,
			},
		},
	})
}

const testOperatingSystemConfig_Ambiguous = `
	data "packet_operating_system" "example" {
		distro  = "ubuntu"
	  }`
