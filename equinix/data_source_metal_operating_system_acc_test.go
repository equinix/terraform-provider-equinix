package equinix_test

import (
	"regexp"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalOperatingSystem_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalOperatingSystemConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.equinix_metal_operating_system.example", "slug", "ubuntu_20_04"),
				),
			},
		},
	})
}

const testAccDataSourceMetalOperatingSystemConfig_basic = `
	data "equinix_metal_operating_system" "example" {
		distro  = "ubuntu"
		version = "20.04"
	  }`

var matchErrOSNotFound = regexp.MustCompile(".*There are no operating systems*")

func TestAccDataSourceMetalOperatingSystem_notFound(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceMetalOperatingSystemConfig_notFound,
				ExpectError: matchErrOSNotFound,
			},
		},
	})
}

const testAccDataSourceMetalOperatingSystemConfig_notFound = `
	data "equinix_metal_operating_system" "example" {
		distro  = "NOTEXISTS"
		version = "alpha"
	  }`

var matchErrOSAmbiguous = regexp.MustCompile(".*There is more than one operating system.*")

func TestAccDataSourceMetalOperatingSystem_ambiguous(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceMetalOperatingSystemConfig_ambiguous,
				ExpectError: matchErrOSAmbiguous,
			},
		},
	})
}

const testAccDataSourceMetalOperatingSystemConfig_ambiguous = `
	data "equinix_metal_operating_system" "example" {
		distro  = "ubuntu"
	  }`
