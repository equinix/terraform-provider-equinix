package equinix

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	networkDeviceMetroEnvVar = "TF_ACC_NETWORK_DEVICE_METRO"
)

func init() {
	resource.AddTestSweepers("NetworkDevice", &resource.Sweeper{
		Name: "NetworkDevice",
		F:    testSweepNetworkDevice,
	})
}

func testSweepNetworkDevice(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	devices, err := config.ne.GetDevices([]string{
		ne.DeviceStateInitializing,
		ne.DeviceStateProvisioned,
		ne.DeviceStateProvisioning,
		ne.DeviceStateWaitingSecondary,
		ne.DeviceStateFailed})
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching NetworkDevice list: %s", err)
		return err
	}
	for _, device := range devices {
		if !isSweepableTestResource(device.Name) {
			continue
		}
		if device.RedundancyType != "PRIMARY" {
			continue
		}
		if err := config.ne.DeleteDevice(device.UUID); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkDevice resource %s (%s): %s", device.UUID, device.Name, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkDevice resource %s (%s)", device.UUID, device.Name)
		}
	}
	return nil
}

func TestAccNetworkDeviceAndUser(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "DC")()
	context := map[string]interface{}{
		"resourceName":            "tst-csr1000v",
		"name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"throughput":              500,
		"throughput_unit":         "Mbps",
		"metro_code":              metro.(string),
		"type_code":               "CSR1000V",
		"package_code":            "SEC",
		"notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"hostname":                fmt.Sprintf("tf-%s", randString(6)),
		"acls":                    []string{"10.0.0.0/24", "1.1.1.1/32"},
		"term_length":             1,
		"version":                 "16.09.05",
		"core_count":              2,
		"secondary-name":          fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"secondary-hostname":      randString(6),
		"secondary-notifications": []string{"secondary@equinix.com"},
		"secondary-acls":          []string{"2.2.2.2/32"},
		"userResourceName":        "tst-user",
		"username":                fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"password":                randString(10),
	}
	resourceName := fmt.Sprintf("equinix_network_device.%s", context["resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", context["userResourceName"].(string))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkDeviceAndUser(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "license_status"),
					resource.TestCheckResourceAttrSet(resourceName, "ibx"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "interface_count"),
					resource.TestCheckResourceAttrSet(resourceName, "redundant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "redundancy_type"),
					resource.TestCheckResourceAttrSet(userResourceName, "uuid"),
				),
			},
		},
	})
}

func testAccNetworkDeviceAndUser(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_account" "test" {
  metro_code = "%{metro_code}"
  status     = "Active"
}

resource "equinix_network_device" "%{resourceName}" {
	name            = "%{name}"
	throughput      = %{throughput}
	throughput_unit = "%{throughput_unit}"
	metro_code      = data.equinix_network_account.test.metro_code
	type_code       = "%{type_code}"
	package_code    = "%{package_code}"
	notifications   = %{notifications}
	hostname        = "%{hostname}"
	acls            = %{acls}
	term_length     = %{term_length}
	account_number  = data.equinix_network_account.test.number
	version         = "%{version}"
	core_count      = %{core_count}
	secondary_device {
		name           = "%{secondary-name}"
		metro_code     = data.equinix_network_account.test.metro_code
		hostname       = "%{secondary-hostname}"
		acls           = %{secondary-acls} 
		notifications  = %{secondary-notifications}
		account_number = data.equinix_network_account.test.number
	  }
}

resource "equinix_network_ssh_user" "%{userResourceName}" {
	username = "%{username}"
	password = "%{password}"
	device_ids = [
	  equinix_network_device.%{resourceName}.uuid,
	  equinix_network_device.%{resourceName}.redundant_uuid
	]
  }
`, ctx)
}
