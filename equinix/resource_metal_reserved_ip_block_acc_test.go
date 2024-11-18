package equinix

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"text/template"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var matchIpBlockErrTimeout = regexp.MustCompile(".* timeout while waiting for state to become 'created'.*")

func testAccMetalReservedIPBlockConfig_global(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	type        = "global_ipv4"
	description = "tfacc-reserved_ip_block-%s"
	quantity    = 1
	custom_data = jsonencode({
		"foo": "bar"
	})
}`, name, name)
}

func testAccMetalReservedIPBlockConfig_public(name, createTimeout string) string {
	if createTimeout == "" {
		createTimeout = "20m"
	}

	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	facility    = "ny5"
	type        = "public_ipv4"
	description = "tfacc-reserved_ip_block-%s"
	quantity    = 2
	tags        = ["Tag1", "Tag2"]

    timeouts {
      create = "%s"
    }
}`, name, name, createTimeout)
}

// testAccMetalReservedIPBlockConfig_metro generates a config for a metro IP
// block with optional tag
func testAccMetalReservedIPBlockConfig_metro(name, tag string) string {
	var b bytes.Buffer
	t, _ := template.New("").Parse(`
	resource "equinix_metal_project" "foobar" {
		name = "tfacc-reserved_ip_block-{{.name}}"
	}

	resource "equinix_metal_reserved_ip_block" "test" {
		project_id  = equinix_metal_project.foobar.id
		metro       = "sv"
		type        = "public_ipv4"
		description = "tfacc-reserved_ip_block-{{.tag}}"
		quantity    = 2
		{{if .tag}}
		tags        = [{{.tag | printf "%q"}}]
		{{end}}
	}`)
	err := t.Execute(&b, map[string]string{"name": name, "tag": tag})

	if err != nil {
		panic(err)
	}

	return b.String()
}

func TestAccMetalReservedIPBlock_global(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_global(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "quantity", "1"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "description", "tfacc-reserved_ip_block-"+rs),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "type", "global_ipv4"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "netmask", "255.255.255.255"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "public", "true"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "management", "false"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "custom_data", `{"foo":"bar"}`),
				),
			},
		},
	})
}

func TestAccMetalReservedIPBlock_public(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_public(rs, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "facility", "ny5"),
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
					resource.TestCheckResourceAttrSet("equinix_metal_reserved_ip_block.test", "network"),
					resource.TestCheckResourceAttrSet("equinix_metal_reserved_ip_block.test", "cidr"),
				),
			},
		},
	})
}

func TestAccMetalReservedIPBlock_metro(t *testing.T) {
	rs := acctest.RandString(10)
	tag := "tag"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_metro(rs, tag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "sv"),
					resource.TestCheckResourceAttr("equinix_metal_reserved_ip_block.test", "tags.0", tag),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "description", "tfacc-reserved_ip_block-"+tag),
				),
			},
			{
				ResourceName:            "equinix_metal_reserved_ip_block.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_state"},
			},
			{
				Config: testAccMetalReservedIPBlockConfig_metro(rs, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "sv"),
					resource.TestCheckResourceAttr("equinix_metal_reserved_ip_block.test", "tags.#", "0"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "description", "tfacc-reserved_ip_block-"),
				),
			},
		},
	})
}

func TestAccMetalReservedIPBlock_importBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_public(rs, ""),
			},
			{
				ResourceName:            "equinix_metal_reserved_ip_block.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_state"},
			},
		},
	})
}

func testAccMetalReservedIPBlockCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*config.Config).Metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_reserved_ip_block" {
			continue
		}
		if _, _, err := client.ProjectIPs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Reserved IP block still exists")
		}
	}

	return nil
}

func testAccMetalReservedIPBlockConfig_facilityToMetro(line string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block_fac_met_test"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	%s
	type        = "public_ipv4"
	quantity    = 2
	tags        = ["Tag1", "Tag2"]
}`, line)
}

func TestAccMetalReservedIPBlock_facilityToMetro(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIPBlockConfig_facilityToMetro(`   facility = "ny5"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "facility", "ny5"),
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "ny"),
				),
			},
			{
				Config: testAccMetalReservedIPBlockConfig_facilityToMetro(`   metro = "ny"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_reserved_ip_block.test", "metro", "ny"),
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccMetalReservedIP_device(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	facility    = local.facility
	type        = "public_ipv4"
	quantity    = 2
}

resource "equinix_metal_device" "test" {
  project_id       = equinix_metal_project.foobar.id
  plan             = local.plan
  facilities       = [local.facility]
  operating_system = local.os
  hostname         = "tfacc-reserved-ip-device"
  billing_cycle    = "hourly"
  ip_address {
	 type = "public_ipv4"
	 cidr = 31
	 reservation_ids = [equinix_metal_reserved_ip_block.test.id]
  }
  ip_address {
	 type = "private_ipv4"
  }
  termination_time = "%s"
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func TestAccMetalReservedIPBlock_device(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalReservedIP_device(rs),
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

func TestAccMetalReservedIPBlockCreate_public_timeout(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testAccMetalReservedIPBlockConfig_public(rs, "2s"),
				ExpectError: matchIpBlockErrTimeout,
			},
		},
	})
}
