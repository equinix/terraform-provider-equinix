package equinix_test

import (
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func testAccFabricReadServiceProfileConfig(uuid string) string {
	return fmt.Sprintf(`data "equinix_fabric_service_profile" "test" {
	uuid = "%s"
	}
`, uuid)
}

func testAccFabricReadServiceProfilesListConfig(name string) string {
	return fmt.Sprintf(`data "equinix_fabric_service_profiles" "test" {
		filter {
			property = "/name"
			operator = "="
			values = ["%s"]
		}
	}
`, name)
}

func TestAccFabricReadServiceProfileByUuid_SP_FCR_RPAA(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadServiceProfileConfig("bfb74121-7e2c-4f74-99b3-69cdafb03b41"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "name", fmt.Sprint("Azure ExpressRoute")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "uuid", fmt.Sprint("bfb74121-7e2c-4f74-99b3-69cdafb03b41")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "state", fmt.Sprint("ACTIVE")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "visibility", fmt.Sprint("PUBLIC")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "access_point_type_configs.#", fmt.Sprint(1)),
					resource.TestMatchResourceAttr(
						"data.equinix_fabric_service_profile.test", "account.0.organization_name", regexp.MustCompile("^azure")),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_remote_connections"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_custom_bandwidth"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.enable_auto_generate_service_key"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.connection_redundancy_required"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "metros.0.code"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "metros.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "metros.0.display_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "self_profile"),
				),
			},
		},
	})
}

func TestAccFabricSearchServiceProfilesByName_SP_FCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadServiceProfilesListConfig("Azure ExpressRoute"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.#", fmt.Sprint(1)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.name", fmt.Sprint("Azure ExpressRoute")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.uuid", fmt.Sprint("bfb74121-7e2c-4f74-99b3-69cdafb03b41")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.type", fmt.Sprint("L2_PROFILE")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.state", fmt.Sprint("ACTIVE")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.visibility", fmt.Sprint("PUBLIC")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.account.#", fmt.Sprint(1)),
					resource.TestMatchResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.account.0.organization_name", regexp.MustCompile("^azure")),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.allow_remote_connections"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.allow_custom_bandwidth"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.enable_auto_generate_service_key"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.connection_redundancy_required"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.metros.0.code"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.metros.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.metros.0.display_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.self_profile"),
				),
			},
		},
	})
}
