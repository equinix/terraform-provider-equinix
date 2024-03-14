package file

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkFile_VSRX(t *testing.T) {
	context := map[string]interface{}{
		"resourceName":   "test",
		"fileName":       fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)) + ".lic",
		"content":        acctest.RandString(50),
		"metroCode":      "SV",
		"deviceTypeCode": "VSRX",
		"processType":    "LICENSE",
		"selfManaged":    true,
		"byol":           true,
	}
	resourceName := fmt.Sprintf("equinix_network_file.%s", context["resourceName"].(string))
	var file ne.File
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkFile(context),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkFileExists(resourceName, &file),
					testAccNetworkFileAttributes(&file, context),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content", "byol", "self_managed"},
			},
		},
	})
}

func testAccNetworkFile(ctx map[string]interface{}) string {
	return nprintf(`
resource "equinix_network_file" "%{resourceName}" {
  file_name        = "%{fileName}"
  content          = "%{content}"
  metro_code       = "%{metroCode}"
  device_type_code = "%{deviceTypeCode}"
  process_type     = "%{processType}"
  self_managed     = "%{selfManaged}"
  byol             = "%{byol}"
}
`, ctx)
}

func testAccNetworkFileExists(resourceName string, file *ne.File) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		client := testAccProvider.Meta().(*config.Config).Ne
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		resp, err := client.GetFile(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching file '%s': %s", rs.Primary.ID, err)
		}
		*file = *resp
		return nil
	}
}

func testAccNetworkFileAttributes(file *ne.File, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["fileName"]; ok && ne.StringValue(file.FileName) != v.(string) {
			return fmt.Errorf("file_name does not match %v - %v", ne.StringValue(file.FileName), v)
		}
		if v, ok := ctx["metroCode"]; ok && ne.StringValue(file.MetroCode) != v.(string) {
			return fmt.Errorf("metro_code does not match %v - %v", ne.StringValue(file.MetroCode), v)
		}
		if v, ok := ctx["deviceTypeCode"]; ok && ne.StringValue(file.DeviceTypeCode) != v.(string) {
			return fmt.Errorf("device_type_code does not match %v - %v", ne.StringValue(file.DeviceTypeCode), v)
		}
		if v, ok := ctx["processType"]; ok && ne.StringValue(file.ProcessType) != v.(string) {
			return fmt.Errorf("process_type does not match %v - %v", ne.StringValue(file.ProcessType), v)
		}
		return nil
	}
}
