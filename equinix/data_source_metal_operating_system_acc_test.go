package equinix

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetalOperatingSystem_Basic(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: testOperatingSystemConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.equinix_metal_operating_system.example", "slug", "ubuntu_20_04"),
				),
			},
		},
	})
}

const testOperatingSystemConfig_Basic = `
	data "equinix_metal_operating_system" "example" {
		distro  = "ubuntu"
		version = "20.04"
	  }`

var matchErrOSNotFound = regexp.MustCompile(".*There are no operating systems*")

func TestAccMetalOperatingSystem_NotFound(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
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
	data "equinix_metal_operating_system" "example" {
		distro  = "NOTEXISTS"
		version = "alpha"
	  }`

var matchErrOSAmbiguous = regexp.MustCompile(".*There is more than one operating system.*")

func TestAccMetalOperatingSystem_Ambiguous(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
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
	data "equinix_metal_operating_system" "example" {
		distro  = "ubuntu"
	  }`
