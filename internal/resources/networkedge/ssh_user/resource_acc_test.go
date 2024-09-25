package ssh_user

import (
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// SSH User acc tests are in Device acc test tests
// reason: SSH User requires device to be provisioned and that is time consuming operation

type testAccConfig struct {
	ctx    map[string]interface{}
	config string
}

func newTestAccConfig(ctx map[string]interface{}) *testAccConfig {
	return &testAccConfig{
		ctx:    ctx,
		config: "",
	}
}

func (t *testAccConfig) build() string {
	return t.config
}

func testAccNeSSHUserExists(resourceName string, user *ne.SSHUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		client := acceptance.TestAccProvider.Meta().(*config.Config).Ne
		resp, err := client.GetSSHUser(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching SSH user '%s': %s", rs.Primary.ID, err)
		}
		*user = *resp
		return nil
	}
}

func testAccNeSSHUserAttributes(user *ne.SSHUser, devices []*ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["username"]; ok && ne.StringValue(user.Username) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ne.StringValue(user.Username), v)
		}
		deviceIDs := make([]string, len(devices))
		for i := range devices {
			deviceIDs[i] = ne.StringValue(devices[i].UUID)
		}
		if !comparisons.SlicesMatch(deviceIDs, user.DeviceUUIDs) {
			return fmt.Errorf("device_ids does not match %v - %v", deviceIDs, user.DeviceUUIDs)
		}
		return nil
	}
}
