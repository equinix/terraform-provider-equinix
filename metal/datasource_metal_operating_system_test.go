package metal

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMetalOperatingSystem_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: testOperatingSystemConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.metal_operating_system.example", "slug", "alpine_3"),
				),
			},
		},
	})
}

const testOperatingSystemConfig_Basic = `
	data "metal_operating_system" "example" {
		name    = "Alpine 3"
		distro  = "alpine"
		version = "3"
	  }`

var matchErrOSNotFound = regexp.MustCompile(".*There are no operating systems*")

func TestAccMetalOperatingSystem_NotFound(t *testing.T) {

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
	data "metal_operating_system" "example" {
		name    = "Container"
		distro  = "NOTEXISTS"
		version = "alpha"
	  }`

var matchErrOSAmbiguous = regexp.MustCompile(".*There is more than one operating system.*")

func TestAccMetalOperatingSystem_Ambiguous(t *testing.T) {

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
	data "metal_operating_system" "example" {
		distro  = "ubuntu"
	  }`
