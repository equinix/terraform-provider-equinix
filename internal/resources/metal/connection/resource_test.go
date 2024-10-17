package connection_test

import (
	"context"
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
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_connection" {
			continue
		}
		if _, _, err := client.InterconnectionsApi.GetInterconnection(context.Background(), rs.Primary.ID).Execute(); err == nil {
			return fmt.Errorf("Metal Connection %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccMetalConnectionHasID(resourceName string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		*id = rs.Primary.ID

		return nil
	}
}

func testAccMetalConnectionRecreated(resourceName string, prevID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		if rs.Primary.ID == prevID {
			return fmt.Errorf("expected ID not to be %q", prevID)
		}

		return nil
	}

}

func testAccMetalConnectionConfig_SharedVlan(randstr string) string {
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

        resource "equinix_metal_connection" "test" {
            name               = "tfacc-conn-%s"
            project_id         = equinix_metal_project.test.id
            type               = "shared"
            redundancy         = "redundant"
            metro              = "sv"
			speed              = "50Mbps"
			service_token_type = "a_side"
			contact_email      = "tfacc@example.com"
			vlans              = [
				equinix_metal_vlan.test1.vxlan,
				equinix_metal_vlan.test2.vxlan,
			]
        }`,
		randstr, randstr, randstr, randstr)
}

func testAccMetalConnectionConfig_SharedPort(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }
		resource "equinix_metal_vlan" "test1" {
			description = "tfacc-conn-vlan1-%s"
			metro       = "sv"
			project_id  = equinix_metal_project.test.id
		}
		resource "equinix_metal_connection" "test" {
			name               = "tfacc-conn-%s"
			project_id         = equinix_metal_project.test.id
			type               = "shared_port_vlan"
			redundancy         = "primary"
			metro              = "sv"
			speed              = "50Mbps"
			contact_email      = "tfacc@example.com"
			vlans              = [
				equinix_metal_vlan.test1.vxlan,
			]
        }`,
		randstr, randstr, randstr)
}

func testAccMetalConnectionConfig_SharedPrimaryVrf(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        resource "equinix_metal_vrf" "test1" {
            name        = "tfacc-conn-vrf1-%s"
            metro       = "sv"
            local_asn   = "65001"
            ip_ranges   = ["10.99.1.0/24"]
            project_id  = equinix_metal_project.test.id
        }

        resource "equinix_metal_connection" "test" {
            name               = "tfacc-conn-%s"
            project_id         = equinix_metal_project.test.id
            type               = "shared"
            redundancy         = "primary"
            metro              = "sv"
			speed              = "200Mbps"
			service_token_type = "a_side"
			contact_email      = "tfacc@example.com"
			vrfs               = [
				equinix_metal_vrf.test1.id,
			]
        }`,
		randstr, randstr, randstr)
}

func testAccMetalConnectionConfig_SharedRedundantVrf(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        resource "equinix_metal_vrf" "test1" {
            name        = "tfacc-conn-vrf1-%s"
            metro       = "sv"
            local_asn   = "65001"
            ip_ranges   = ["10.99.1.0/24"]
            project_id  = equinix_metal_project.test.id
        }

        resource "equinix_metal_vrf" "test2" {
            name        = "tfacc-conn-vrf2-%s"
            metro       = "sv"
            local_asn   = "65002"
            ip_ranges   = ["10.99.2.0/24",]
            project_id  = equinix_metal_project.test.id
        }

        resource "equinix_metal_connection" "test" {
            name               = "tfacc-conn-%s"
            project_id         = equinix_metal_project.test.id
            type               = "shared"
            redundancy         = "redundant"
            metro              = "sv"
			speed              = "200Mbps"
			service_token_type = "a_side"
			contact_email      = "tfacc@example.com"
			vrfs               = [
				equinix_metal_vrf.test1.id,
				equinix_metal_vrf.test2.id,
			]
        }`,
		randstr, randstr, randstr, randstr)
}

func testAccMetalConnectionConfig_SharedVlan_zside(randstr string) string {
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

        resource "equinix_metal_connection" "test" {
            name               = "tfacc-conn-%s"
            project_id         = equinix_metal_project.test.id
            type               = "shared"
            redundancy         = "redundant"
            metro              = "sv"
			service_token_type = "z_side"
			vlans              = [
				equinix_metal_vlan.test1.vxlan,
				equinix_metal_vlan.test2.vxlan,
			]
        }`,
		randstr, randstr, randstr, randstr)
}

func TestAccMetalConnection_sharedVlan_zside(t *testing.T) {
	rs := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_SharedVlan_zside(rs),
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

func TestAccMetalConnection_sharedVlan(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_SharedVlan(rs),
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
				Config: testAccMetalConnectionConfig_SharedVlan(rs) + testDataSourceMetalConnectionConfig(),
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

func TestAccMetalConnection_sharedPort(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_SharedPort(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "contact_email", "tfacc@example.com"),
					resource.TestCheckResourceAttrSet("equinix_metal_connection.test", "authorization_code"),
					resource.TestCheckResourceAttrSet("equinix_metal_connection.test", "redundancy"),
				),
			},
			{
				ResourceName:            "equinix_metal_connection.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "vlans"},
			},
		},
	})
}

func TestAccMetalConnection_sharedRedundantVrf(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_SharedRedundantVrf(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.type", "a_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_token_type", "a_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.max_allowed_speed", "200Mbps"),
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
				Config: testAccMetalConnectionConfig_SharedRedundantVrf(rs) + testDataSourceMetalConnectionConfig(),
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
						"data.equinix_metal_connection.test", "service_tokens.0.max_allowed_speed", "200Mbps"),
				),
			},
		},
	})
}

func TestAccMetalConnection_sharedVrfUpgradeRedundant(t *testing.T) {
	rs := acctest.RandString(10)

	var connID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_SharedPrimaryVrf(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.type", "a_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_token_type", "a_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.max_allowed_speed", "200Mbps"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "contact_email", "tfacc@example.com"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "vrfs.#", "1",
					),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "redundancy", "primary",
					),
					testAccMetalConnectionHasID("equinix_metal_connection.test", &connID),
				),
			},
			{
				ResourceName:      "equinix_metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// The update to redundancy and VRFs will cause a recreate
				Config: testAccMetalConnectionConfig_SharedRedundantVrf(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.type", "a_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_token_type", "a_side"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "service_tokens.0.max_allowed_speed", "200Mbps"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "contact_email", "tfacc@example.com"),
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "redundancy", "redundant",
					),
					testAccMetalConnectionRecreated("equinix_metal_connection.test", connID),
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
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
// TODO (ocobles): once migrated, this test may be removed
func TestAccMetalConnection_shared_zside_upgradeFromVersion(t *testing.T) {
	rs := acctest.RandString(10)
	cfg := testAccMetalConnectionConfig_SharedVlan_zside(rs)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckMetal(t); acceptance.TestAccPreCheckProviderConfigured(t) },
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
				ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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
