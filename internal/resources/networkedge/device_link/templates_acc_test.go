package device_link

import (
	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"
)

const (
	tstResourcePrefix = "tfacc"

	networkDeviceProjectId                  = "TF_ACC_NETWORK_DEVICE_PROJECT_ID"
	networkDeviceAccountNameEnvVar          = "TF_ACC_NETWORK_DEVICE_BILLING_ACCOUNT_NAME"
	networkDeviceSecondaryAccountNameEnvVar = "TF_ACC_NETWORK_DEVICE_SECONDARY_BILLING_ACCOUNT_NAME"
	networkDeviceMetroEnvVar                = "TF_ACC_NETWORK_DEVICE_METRO"
	networkDeviceSecondaryMetroEnvVar       = "TF_ACC_NETWORK_DEVICE_SECONDARY_METRO"
	networkDeviceCSRSDWANLicenseFileEnvVar  = "TF_ACC_NETWORK_DEVICE_CSRSDWAN_LICENSE_FILE"
	networkDeviceVSRXLicenseFileEnvVar      = "TF_ACC_NETWORK_DEVICE_VSRX_LICENSE_FILE"
	networkDeviceVersaController1EnvVar     = "TF_ACC_NETWORK_DEVICE_VERSA_CONTROLLER1"
	networkDeviceVersaController2EnvVar     = "TF_ACC_NETWORK_DEVICE_VERSA_CONTROLLER2"
	networkDeviceVersaLocalIDEnvVar         = "TF_ACC_NETWORK_DEVICE_VERSA_LOCALID"
	networkDeviceVersaRemoteIDEnvVar        = "TF_ACC_NETWORK_DEVICE_VERSA_REMOTEID"
	networkDeviceVersaSerialNumberEnvVar    = "TF_ACC_NETWORK_DEVICE_VERSA_SERIAL"
	networkDeviceCGENIXLicenseKeyEnvVar     = "TF_ACC_NETWORK_DEVICE_CGENIX_LICENSE_KEY"
	networkDeviceCGENIXLicenseSecretEnvVar  = "TF_ACC_NETWORK_DEVICE_CGENIX_LICENSE_SECRET"
	networkDevicePANWLicenseTokenEnvVar     = "TF_ACC_NETWORK_DEVICE_PANW_LICENSE_TOKEN"
)

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

func (t *testAccConfig) withDevice() *testAccConfig {
	t.config += testAccNetworkDevice(t.ctx)
	return t
}

func copyMap(source map[string]interface{}) map[string]interface{} {
	target := make(map[string]interface{})
	for k, v := range source {
		target[k] = v
	}
	return target
}

func testAccNetworkDeviceUser(ctx map[string]interface{}) string {
	config := nprintf.Nprintf(`
resource "equinix_network_ssh_user" "%{user-resourceName}" {
  username = "%{user-username}"
  password = "%{user-password}"
  device_ids = [
    equinix_network_device.%{device-resourceName}.id`, ctx)
	if _, ok := ctx["device-secondary_name"]; ok {
		config += nprintf.Nprintf(`,
    equinix_network_device.%{device-resourceName}.redundant_id`, ctx)
	}
	config += `
  ]
}`
	return config
}

func testAccNetworkDevice(ctx map[string]interface{}) string {
	var config string
	config += nprintf.Nprintf(`
data "equinix_network_account" "test" {
  metro_code = "%{device-metro_code}"
  status     = "Active"
  project_id = "%{device-project_id}"`, ctx)
	if v, ok := ctx["device-account_name"]; ok && !comparisons.IsEmpty(v) {
		config += nprintf.Nprintf(`
  name = "%{device-account_name}"`, ctx)
	}
	config += nprintf.Nprintf(`
}`, ctx)
	if _, ok := ctx["device-secondary_metro_code"]; ok {
		config += nprintf.Nprintf(`
data "equinix_network_account" "test-secondary" {
  metro_code = "%{device-secondary_metro_code}"
  status     = "Active"`, ctx)
		if v, ok := ctx["device-secondary_account_name"]; ok && !comparisons.IsEmpty(v) {
			config += nprintf.Nprintf(`
  name = "%{device-secondary_account_name}"`, ctx)
		}
		config += nprintf.Nprintf(` 
}`, ctx)
	}
	config += nprintf.Nprintf(`
resource "equinix_network_device" "%{device-resourceName}" {
  self_managed          = %{device-self_managed}
  byol                  = %{device-byol}
  name                  = "%{device-name}"
  metro_code            = "%{device-metro_code}"
  type_code             = "%{device-type_code}"
  project_id            = "%{device-project_id}"
  package_code          = "%{device-package_code}"
  notifications         = %{device-notifications}
  term_length           = %{device-term_length}
  account_number        = data.equinix_network_account.test.number
  version               = "%{device-version}"
  core_count            = %{device-core_count}`, ctx)
	if _, ok := ctx["device-purchase_order_number"]; ok {
		config += nprintf.Nprintf(`
  purchase_order_number = "%{device-purchase_order_number}"`, ctx)
	}
	if _, ok := ctx["device-purchase_order_number"]; ok {
		config += nprintf.Nprintf(`
  order_reference       = "%{device-order_reference}"`, ctx)
	}
	if _, ok := ctx["device-additional_bandwidth"]; ok {
		config += nprintf.Nprintf(`
  additional_bandwidth  = %{device-additional_bandwidth}`, ctx)
	}
	if _, ok := ctx["device-throughput"]; ok {
		config += nprintf.Nprintf(`
  throughput            = %{device-throughput}
  throughput_unit       = "%{device-throughput_unit}"`, ctx)
	}
	if _, ok := ctx["device-hostname"]; ok {
		config += nprintf.Nprintf(`
  hostname              = "%{device-hostname}"`, ctx)
	}
	if _, ok := ctx["device-interface_count"]; ok {
		config += nprintf.Nprintf(`
  interface_count       = %{device-interface_count}`, ctx)
	}
	if _, ok := ctx["acl-resourceName"]; ok {
		config += nprintf.Nprintf(`
  acl_template_id       = equinix_network_acl_template.%{acl-resourceName}.id`, ctx)
	}
	if _, ok := ctx["mgmtAcl-resourceName"]; ok {
		config += nprintf.Nprintf(`
  mgmt_acl_template_uuid = equinix_network_acl_template.%{mgmtAcl-resourceName}.id`, ctx)
	}
	if _, ok := ctx["sshkey-resourceName"]; ok {
		config += nprintf.Nprintf(`
  ssh_key {
    username = "test"
    key_name = equinix_network_ssh_key.%{sshkey-resourceName}.name
  }`, ctx)
	}
	if _, ok := ctx["device-license_file"]; ok {
		config += nprintf.Nprintf(`
  license_file          = "%{device-license_file}"`, ctx)
	}
	if _, ok := ctx["device-vendorConfig_enabled"]; ok {
		config += nprintf.Nprintf(`
  vendor_configuration  = {`, ctx)
		if _, ok := ctx["device-vendorConfig_siteId"]; ok {
			config += nprintf.Nprintf(`
    siteId          = "%{device-vendorConfig_siteId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_systemIpAddress"]; ok {
			config += nprintf.Nprintf(`
    systemIpAddress = "%{device-vendorConfig_systemIpAddress}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_licenseKey"]; ok {
			config += nprintf.Nprintf(`
    licenseKey = "%{device-vendorConfig_licenseKey}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_licenseSecret"]; ok {
			config += nprintf.Nprintf(`
    licenseSecret = "%{device-vendorConfig_licenseSecret}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_controller1"]; ok {
			config += nprintf.Nprintf(`
    controller1 = "%{device-vendorConfig_controller1}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_controller2"]; ok {
			config += nprintf.Nprintf(`
    controller2 = "%{device-vendorConfig_controller2}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_localId"]; ok {
			config += nprintf.Nprintf(`
    localId = "%{device-vendorConfig_localId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_remoteId"]; ok {
			config += nprintf.Nprintf(`
    remoteId = "%{device-vendorConfig_remoteId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_serialNumber"]; ok {
			config += nprintf.Nprintf(`
    serialNumber = "%{device-vendorConfig_serialNumber}"`, ctx)
		}
		config += nprintf.Nprintf(`
  }`, ctx)
	}
	if _, ok := ctx["device-secondary_name"]; ok {
		config += nprintf.Nprintf(`
  secondary_device {
    name                 = "%{device-secondary_name}"`, ctx)
		if _, ok := ctx["device-secondary_metro_code"]; ok {
			config += nprintf.Nprintf(`
    metro_code           = "%{device-secondary_metro_code}"
    account_number       = data.equinix_network_account.test-secondary.number`, ctx)
		} else {
			config += nprintf.Nprintf(`
    metro_code           = "%{device-metro_code}"
    account_number       = data.equinix_network_account.test.number`, ctx)
		}
		config += nprintf.Nprintf(`
    notifications        = %{device-secondary_notifications}`, ctx)
		if _, ok := ctx["device-secondary_additional_bandwidth"]; ok {
			config += nprintf.Nprintf(`
    additional_bandwidth = %{device-secondary_additional_bandwidth}`, ctx)
		}
		if _, ok := ctx["device-secondary_hostname"]; ok {
			config += nprintf.Nprintf(`
    hostname             = "%{device-secondary_hostname}"`, ctx)
		}
		if _, ok := ctx["acl-secondary_resourceName"]; ok {
			config += nprintf.Nprintf(`
    acl_template_id      = equinix_network_acl_template.%{acl-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["mgmtAcl-secondary_resourceName"]; ok {
			config += nprintf.Nprintf(`
    mgmt_acl_template_uuid = equinix_network_acl_template.%{mgmtAcl-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["sshkey-resourceName"]; ok {
			config += nprintf.Nprintf(`
    ssh_key {
      username = "test"
      key_name = equinix_network_ssh_key.%{sshkey-resourceName}.name
    }`, ctx)
		}
		if _, ok := ctx["device-secondary_license_file"]; ok {
			config += nprintf.Nprintf(`
    license_file         = "%{device-secondary_license_file}"`, ctx)
		}
		if _, ok := ctx["device-secondary_vendorConfig_enabled"]; ok {
			config += nprintf.Nprintf(`
    vendor_configuration  = {`, ctx)
			if _, ok := ctx["device-secondary_vendorConfig_siteId"]; ok {
				config += nprintf.Nprintf(`
      siteId          = "%{device-secondary_vendorConfig_siteId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_systemIpAddress"]; ok {
				config += nprintf.Nprintf(`
      systemIpAddress = "%{device-secondary_vendorConfig_systemIpAddress}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_licenseKey"]; ok {
				config += nprintf.Nprintf(`
      licenseKey = "%{device-secondary_vendorConfig_licenseKey}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_licenseSecret"]; ok {
				config += nprintf.Nprintf(`
      licenseSecret = "%{device-secondary_vendorConfig_licenseSecret}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_controller1"]; ok {
				config += nprintf.Nprintf(`
      controller1 = "%{device-secondary_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_controller2"]; ok {
				config += nprintf.Nprintf(`
      controller2 = "%{device-secondary_vendorConfig_controller2}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_localId"]; ok {
				config += nprintf.Nprintf(`
      localId = "%{device-secondary_vendorConfig_localId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_remoteId"]; ok {
				config += nprintf.Nprintf(`
      remoteId = "%{device-secondary_vendorConfig_remoteId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_serialNumber"]; ok {
				config += nprintf.Nprintf(`
      serialNumber = "%{device-secondary_vendorConfig_serialNumber}"`, ctx)
			}
			config += nprintf.Nprintf(`
    }`, ctx)
		}
		config += `
  }`
	}
	if _, ok := ctx["device-cluster_name"]; ok {
		config += nprintf.Nprintf(`
  cluster_details {
    cluster_name        = "%{device-cluster_name}"`, ctx)
		config += `
    node0 {`
		if _, ok := ctx["device-node0_license_file_id"]; ok {
			config += nprintf.Nprintf(`
      license_file_id   = "%{device-node0_license_file_id}"`, ctx)
		}
		if _, ok := ctx["device-node0_license_token"]; ok {
			config += nprintf.Nprintf(`
      license_token     = "%{device-node0_license_token}"`, ctx)
		}
		if _, ok := ctx["device-node0_vendorConfig_enabled"]; ok {
			config += nprintf.Nprintf(`
      vendor_configuration {`, ctx)
			if _, ok := ctx["device-node0_vendorConfig_hostname"]; ok {
				config += nprintf.Nprintf(`
        hostname        = "%{device-node0_vendorConfig_hostname}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_adminPassword"]; ok {
				config += nprintf.Nprintf(`
        admin_password  = "%{device-node0_vendorConfig_adminPassword}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_controller1"]; ok {
				config += nprintf.Nprintf(`
        controller1     = "%{device-node0_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_activationKey"]; ok {
				config += nprintf.Nprintf(`
        activation_key  = "%{device-node0_vendorConfig_activationKey}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_controllerFqdn"]; ok {
				config += nprintf.Nprintf(`
        controller_fqdn = "%{device-node0_vendorConfig_controllerFqdn}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_rootPassword"]; ok {
				config += nprintf.Nprintf(`
        root_password   = "%{device-node0_vendorConfig_rootPassword}"`, ctx)
			}
			config += nprintf.Nprintf(`
      }`, ctx)
		}
		config += `
    }`
		config += `
    node1 {`
		if _, ok := ctx["device-node1_license_file_id"]; ok {
			config += nprintf.Nprintf(`
      license_file_id   = "%{device-node1_license_file_id}"`, ctx)
		}
		if _, ok := ctx["device-node1_license_token"]; ok {
			config += nprintf.Nprintf(`
      license_token     = "%{device-node1_license_token}"`, ctx)
		}
		if _, ok := ctx["device-node1_vendorConfig_enabled"]; ok {
			config += nprintf.Nprintf(`
      vendor_configuration {`, ctx)
			if _, ok := ctx["device-node1_vendorConfig_hostname"]; ok {
				config += nprintf.Nprintf(`
        hostname        = "%{device-node1_vendorConfig_hostname}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_adminPassword"]; ok {
				config += nprintf.Nprintf(`
        admin_password  = "%{device-node1_vendorConfig_adminPassword}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_controller1"]; ok {
				config += nprintf.Nprintf(`
        controller1     = "%{device-node1_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_activationKey"]; ok {
				config += nprintf.Nprintf(`
        activation_key  = "%{device-node1_vendorConfig_activationKey}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_controllerFqdn"]; ok {
				config += nprintf.Nprintf(`
        controller_fqdn = "%{device-node1_vendorConfig_controllerFqdn}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_rootPassword"]; ok {
				config += nprintf.Nprintf(`
        root_password   = "%{device-node1_vendorConfig_rootPassword}"`, ctx)
			}
			config += nprintf.Nprintf(`
      }`, ctx)
		}
		config += `
    }`
		config += `
  }`
	}
	config += `
}`
	return config
}

func testAccNetworkDeviceACL(ctx map[string]interface{}) string {
	var config string
	if _, ok := ctx["acl-name"]; ok {
		config += nprintf.Nprintf(`
resource "equinix_network_acl_template" "%{acl-resourceName}" {
  name          = "%{acl-name}"
  description   = "%{acl-description}"
  inbound_rule {
    subnet   = "10.0.0.0/24"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	}
	if _, ok := ctx["mgmtAcl-name"]; ok {
		config += nprintf.Nprintf(`
resource "equinix_network_acl_template" "%{mgmtAcl-resourceName}" {
  name          = "%{mgmtAcl-name}"
  description   = "%{mgmtAcl-description}"
  inbound_rule {
    subnet   = "11.0.0.0/24"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	}
	if _, ok := ctx["acl-secondary_name"]; ok {
		config += nprintf.Nprintf(`
resource "equinix_network_acl_template" "%{acl-secondary_resourceName}" {
  name          = "%{acl-secondary_name}"
  description   = "%{acl-secondary_description}"
  inbound_rule {
    subnet   = "192.0.0.0/24"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	}
	if _, ok := ctx["mgmtAcl-secondary_name"]; ok {
		config += nprintf.Nprintf(`
resource "equinix_network_acl_template" "%{mgmtAcl-secondary_resourceName}" {
  name          = "%{mgmtAcl-secondary_name}"
  description   = "%{mgmtAcl-secondary_description}"
  inbound_rule {
    subnet   = "193.0.0.0/24"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	}
	return config
}

func testAccNetworkDeviceSSHKey(ctx map[string]interface{}) string {
	return nprintf.Nprintf(`
resource "equinix_network_ssh_key" "%{sshkey-resourceName}" {
  name       = "%{sshkey-name}"
  public_key = "%{sshkey-public_key}"
}
`, ctx)
}

func (t *testAccConfig) withACL() *testAccConfig {
	t.config += testAccNetworkDeviceACL(t.ctx)
	return t
}

func (t *testAccConfig) withSSHKey() *testAccConfig {
	t.config += testAccNetworkDeviceSSHKey(t.ctx)
	return t
}

func (t *testAccConfig) withSSHUser() *testAccConfig {
	t.config += testAccNetworkDeviceUser(t.ctx)
	return t
}
