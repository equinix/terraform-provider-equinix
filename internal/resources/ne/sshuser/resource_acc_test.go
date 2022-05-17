package sshuser

import (
	"context"
	"fmt"
	"log"

	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/tfacc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// SSH User acc tests are in Device acc test tests
// reason: SSH User requires device to be provisioned and that is time consuming operation

func init() {
	resource.AddTestSweepers("equinix_network_ssh_user", &resource.Sweeper{
		Name: "equinix_network_ssh_user",
		F:    testSweepNetworkSSHUser,
	})
}

func testSweepNetworkSSHUser(region string) error {
	config, err := tfacc.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping Network SSH users: %s", err)
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	users, err := config.NEClient.GetSSHUsers()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching NetworkSSHUser list: %s", err)
		return err
	}
	for _, user := range users {
		if !tfacc.IsSweepableTestResource(ne.StringValue(user.Username)) {
			continue
		}
		if err := config.NEClient.DeleteSSHUser(ne.StringValue(user.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkSSHUser resource %s (%s): %s", ne.StringValue(user.UUID), ne.StringValue(user.Username), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkSSHUser resource %s (%s)", ne.StringValue(user.UUID), ne.StringValue(user.Username))
		}
	}
	return nil
}

func testAccNetworkDeviceUser(ctx map[string]interface{}) string {
	config := tfacc.NPrintf(`
resource "equinix_network_ssh_user" "%{user-resourceName}" {
  username = "%{user-username}"
  password = "%{user-password}"
  device_ids = [
    equinix_network_device.%{device-resourceName}.id`, ctx)
	if _, ok := ctx["device-secondary_name"]; ok {
		config += tfacc.NPrintf(`,
    equinix_network_device.%{device-resourceName}.redundant_id`, ctx)
	}
	config += `
  ]
}`
	return config
}

func withSSHUser(t *tfacc.TestAccConfig) {
	t.Use(testAccNetworkDeviceUser)
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
		client := tfacc.AccProvider.Meta().(*config.Config).NEClient
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
