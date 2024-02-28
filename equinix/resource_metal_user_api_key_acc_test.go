package equinix

// "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

// Commented out because it lists existing user API keys in debug log

/*

func TestAccMetalUserAPIKey_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ExternalProviders:    testExternalProviders,
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalUserAPIKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalUserAPIKeyConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"equinix_metal_user_api_key.test", "token"),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_user_api_key.test", "user_id"),
				),
			},
		},
	})
}

func testAccMetalUserAPIKeyCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*config.Config).Metal
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_user_api_key" {
			continue
		}
		if _, err := client.APIKeys.UserGet(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal UserAPI key still exists")
		}
	}
	return nil
}

func testAccMetalUserAPIKeyConfig_basic() string {
	return fmt.Sprintf(`
resource "equinix_metal_user_api_key" "test" {
    description = "tfacc-user-key"
    read_only   = true
}`)
}
*/
