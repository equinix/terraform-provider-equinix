package connection_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccMetalConnectionCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_connection" {
			continue
		}
		if _, _, err := client.Connections.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Connection still exists")
		}
	}

	return nil
}

func testAccMetalConnectionConfig_Shared(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        resource "equinix_metal_connection" "test" {
            name               = "tfacc-conn-%s"
            project_id         = equinix_metal_project.test.id
            type               = "shared"
            redundancy         = "redundant"
            metro              = "sv"
			speed              = "50Mbps"
			service_token_type = "a_side"
			contact_email      = "tfacc@example.com"
        }`,
		randstr, randstr)
}

func testAccMetalConnectionConfig_Shared_zside(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        resource "equinix_metal_connection" "test" {
            name               = "tfacc-conn-%s"
            project_id         = equinix_metal_project.test.id
            type               = "shared"
            redundancy         = "redundant"
            metro              = "sv"
			service_token_type = "z_side"
        }`,
		randstr, randstr)
}

func TestAccMetalConnection_shared_zside(t *testing.T) {
	rs := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_Shared_zside(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "speed", "10Gbps"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.type", "z_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_token_type", "z_side"),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_connection.test", "contact_email"),
				),
			},
		},
	})
}

func TestAccMetalConnection_shared(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_Shared(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.type", "a_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_token_type", "a_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.max_allowed_speed", "50Mbps"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "contact_email", "tfacc@example.com"),
				),
			},
			{
				ResourceName:      "equinix_metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalConnectionConfig_Shared(rs) + testDataSourceMetalConnectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "mode", "standard"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "type", "shared"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "redundancy", "redundant"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_connection.test", "service_tokens.0.type", "a_side"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_connection.test", "service_token_type", "a_side"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_connection.test", "service_tokens.0.max_allowed_speed", "50Mbps"),
				),
			},
		},
	})
}

func testAccMetalConnectionConfig_dedicated(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

		// We use the project resource to get organization_id
        resource "equinix_metal_connection" "test" {
            name            = "tfacc-conn-%s"
            metro           = "sv"
			organization_id = equinix_metal_project.test.organization_id
            type            = "dedicated"
            redundancy      = "redundant"
			tags            = ["tfacc"]
			speed           = "50Mbps"
			mode            = "standard"
        }`,
		randstr, randstr)
}

func testDataSourceMetalConnectionConfig() string {
	return `
		data "equinix_metal_connection" "test" {
            connection_id = equinix_metal_connection.test.id
        }`
}

func TestAccMetalConnection_dedicated(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_dedicated(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "mode", "standard"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "type", "dedicated"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "redundancy", "redundant"),
				),
			},
			{
				ResourceName:      "equinix_metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalConnectionConfig_dedicated(rs) + testDataSourceMetalConnectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "mode", "standard"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "type", "dedicated"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "redundancy", "redundant"),
				),
			},
		},
	})
}

func testAccMetalConnectionConfig_tunnel(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

		// We use the project resource to get organization_id internally
		resource "equinix_metal_connection" "test" {
			name            = "tfacc-conn-%s"
			project_id      = equinix_metal_project.test.id
			metro           = "sv"
			redundancy      = "redundant"
			type            = "dedicated"
			mode            = "tunnel"
			speed           = "50Mbps"
        }`,
		randstr, randstr)
}

func TestAccMetalConnection_tunnel(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_tunnel(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "mode", "tunnel"),
				),
			},
			{
				ResourceName:            "equinix_metal_connection.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func testAccMetalConnectionConfig_sharedVlans(randstr string, vlans string) string {
	return fmt.Sprintf(`
		resource "equinix_metal_project" "test" {
			name = "tfacc-conn-pro-%s"
		}

		resource "equinix_metal_vlan" "test1" {
			description = "tfacc-conn-vlan1-%s"
			metro       = "sv"
			project_id  = equinix_metal_project.test.id
		}

		resource "equinix_metal_vlan" "test2" {
			description = "tfacc-conn-vlan2-%s"
			metro       = "sv"
			project_id  = equinix_metal_project.test.id
		}

		resource "equinix_metal_vlan" "test3" {
			description = "tfacc-conn-vlan3-%s"
			metro       = "sv"
			project_id  = equinix_metal_project.test.id
		}

		resource "equinix_metal_connection" "test" {
			name               = "tfacc-conn-%s"
			project_id         = equinix_metal_project.test.id
			type               = "shared"
			redundancy         = "redundant"
			metro              = "sv"
			speed              = "50Mbps"
			service_token_type = "a_side"
			vlans = [
				%s
			]
		}`,
		randstr, randstr, randstr, randstr, randstr, vlans)
}

func TestAccMetalConnection_sharedVlans(t *testing.T) {
	rs := acctest.RandString(10)

	// In the first test step, we will assign 2 VLANs
	step1Vlans := "equinix_metal_vlan.test1.vxlan, equinix_metal_vlan.test2.vxlan,"
	// In the second test step, we will change the primary VLAN
	step2Vlans := "equinix_metal_vlan.test3.vxlan, equinix_metal_vlan.test2.vxlan,"
	// In the third test step, we will remove the secondary VLAN
	step3Vlans := "equinix_metal_vlan.test3.vxlan,"
	// In the fourth test step, we will add a new secondary VLAN
	step4Vlans := "equinix_metal_vlan.test3.vxlan, equinix_metal_vlan.test1.vxlan,"
	// In the fifth test step, we will remove both VLANs
	step5Vlans := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_sharedVlans(rs, step1Vlans),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test1", "vxlan",
						"equinix_metal_connection.test", "vlans.0"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test2", "vxlan",
						"equinix_metal_connection.test", "vlans.1"),
				),
			},
			{
				Config: testAccMetalConnectionConfig_sharedVlans(rs, step2Vlans),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test3", "vxlan",
						"equinix_metal_connection.test", "vlans.0"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test2", "vxlan",
						"equinix_metal_connection.test", "vlans.1"),
				),
			},
			{
				Config: testAccMetalConnectionConfig_sharedVlans(rs, step3Vlans),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test3", "vxlan",
						"equinix_metal_connection.test", "vlans.0"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "vlans.#", "1"),
				),
			},
			{
				Config: testAccMetalConnectionConfig_sharedVlans(rs, step4Vlans),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test3", "vxlan",
						"equinix_metal_connection.test", "vlans.0"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.test1", "vxlan",
						"equinix_metal_connection.test", "vlans.1"),
				),
			},
			{
				Config: testAccMetalConnectionConfig_sharedVlans(rs, step5Vlans),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "vlans.#", "0"),
				),
			},
		},
	})
}

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
// TODO (ocobles): once migrated, this test may be removed
func TestAccMetalConnection_shared_zside_upgradeFromVersion(t *testing.T) {
	rs := acctest.RandString(10)
	cfg := testAccMetalConnectionConfig_Shared_zside(rs)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckMetal(t) },
		CheckDestroy: testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.29.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
				},
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "speed", "10Gbps"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.type", "z_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_token_type", "z_side"),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_connection.test", "contact_email"),
				),
			},
			{
				ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
				Config:                   cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
