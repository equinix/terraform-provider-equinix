package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	metalDedicatedConnIDEnvVar = "TF_ACC_METAL_DEDICATED_CONNECTION_ID"
)

func TestSpeedConversion(t *testing.T) {
	speedUint, err := speedStrToUint("50Mbps")
	if err != nil {
		t.Errorf("Error converting speed string to uint64: %s", err)

	}
	if speedUint != 50*mega {
		t.Errorf("Speed string conversion failed. Expected: %d, got: %d", 50*mega, speedUint)
	}

	speedStr, err := speedUintToStr(50 * mega)
	if err != nil {
		t.Errorf("Error converting speed uint to string: %s", err)
	}
	if speedStr != "50Mbps" {
		t.Errorf("Speed uint conversion failed. Expected: %s, got: %s", "50Mbps", speedStr)
	}

	speedUint, err = speedStrToUint("100Gbps")
	if err == nil {
		t.Errorf("Expected error converting invalid speed string to uint, got: %d", speedUint)
	}

	speedStr, err = speedUintToStr(100 * giga)
	if err == nil {
		t.Errorf("Expected error converting invalid speed uint to string, got: %s", speedStr)
	}
}

func testAccMetalConnectionCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

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
        }`,
		randstr, randstr)
}

func TestAccMetalConnection_shared(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalConnectionCheckDestroyed,
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

        resource "equinix_metal_connection" "test" {
            name            = "tfacc-conn-%s"
            metro           = "sv"
			project_id      = equinix_metal_project.test.id
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalConnectionCheckDestroyed,
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_tunnel(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "mode", "tunnel"),
				),
			},
			{
				ResourceName:      "equinix_metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
