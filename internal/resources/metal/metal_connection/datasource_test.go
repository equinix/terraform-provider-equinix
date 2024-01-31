package metal_connection_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalConnection_withoutVlans(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalConnectionConfig_withoutVlans(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_connection.test", "id",
						"data.equinix_metal_connection.test", "id"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_connection.test", "vlans.#",
						"data.equinix_metal_connection.test", "vlans.#"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_connection.test", "vlans.#", "0"),
				),
			},
		},
	})
}

func testDataSourceMetalConnectionConfig_withoutVlans(r int) string {
	return fmt.Sprintf(`
		resource "equinix_metal_project" "test" {
			name = "tfacc-conn-project-%d"
		}

		resource "equinix_metal_connection" "test" {
			name               = "tfacc-conn-%d"
			project_id         = equinix_metal_project.test.id
			type               = "shared"
			redundancy         = "redundant"
			metro              = "sv"
			speed              = "50Mbps"
			service_token_type = "a_side"
		}

		data "equinix_metal_connection" "test" {
			connection_id = equinix_metal_connection.test.id
		}`,
		r, r)
}

func TestAccDataSourceMetalConnection_withVlans(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalConnectionConfig_withVlans(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_connection.test", "id",
						"data.equinix_metal_connection.test", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_connection.test", "vlans.#", "2"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test1", "vxlan",
						"data.equinix_metal_connection.test", "vlans.0"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test2", "vxlan",
						"data.equinix_metal_connection.test", "vlans.1"),
				),
			},
		},
	})
}

func testDataSourceMetalConnectionConfig_withVlans(r int) string {
	return fmt.Sprintf(`
		resource "equinix_metal_project" "test" {
			name = "tfacc-conn-pro-%d"
		}

		resource "equinix_metal_vlan" "test1" {
			description = "tfacc-conn-vlan1-%d"
			metro       = "sv"
			project_id  = equinix_metal_project.test.id
		}

		resource "equinix_metal_vlan" "test2" {
			description = "tfacc-conn-vlan2-%d"
			metro       = "sv"
			project_id  = equinix_metal_project.test.id
		}

		resource "equinix_metal_connection" "test" {
			name               = "tfacc-conn-%d"
			project_id         = equinix_metal_project.test.id
			type               = "shared"
			redundancy         = "redundant"
			metro              = "sv"
			speed              = "50Mbps"
			service_token_type = "a_side"
			vlans = [
				equinix_metal_vlan.test1.vxlan,
				equinix_metal_vlan.test2.vxlan
			]
		}

		data "equinix_metal_connection" "test" {
    		connection_id = equinix_metal_connection.test.id
		}`,
		r, r, r, r)
}
