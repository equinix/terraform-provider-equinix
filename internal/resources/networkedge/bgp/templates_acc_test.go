package bgp_test

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
	ctx    map[string]any
	config string
}

func newTestAccConfig(ctx map[string]any) *testAccConfig {
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

func testAccNetworkDeviceUser(ctx map[string]any) string {
	config := nprintf.NPrintf(`
resource "equinix_network_ssh_user" "%{user-resourceName}" {
  username = "%{user-username}"
  password = "%{user-password}"
  device_ids = [
    equinix_network_device.%{device-resourceName}.id`, ctx)
	if _, ok := ctx["device-secondary_name"]; ok {
		config += nprintf.NPrintf(`,
    equinix_network_device.%{device-resourceName}.redundant_id`, ctx)
	}
	config += `
  ]
}`
	return config
}

func testAccNetworkDevice(ctx map[string]any) string {
	var config string
	config += nprintf.NPrintf(`
data "equinix_network_account" "test" {
  metro_code = "%{device-metro_code}"
  status     = "Active"
  project_id = "%{device-project_id}"`, ctx)
	if v, ok := ctx["device-account_name"]; ok && !comparisons.IsEmpty(v) {
		config += nprintf.NPrintf(`
  name = "%{device-account_name}"`, ctx)
	}
	config += nprintf.NPrintf(`
}`, ctx)
	if _, ok := ctx["device-secondary_metro_code"]; ok {
		config += nprintf.NPrintf(`
data "equinix_network_account" "test-secondary" {
  metro_code = "%{device-secondary_metro_code}"
  status     = "Active"`, ctx)
		if v, ok := ctx["device-secondary_account_name"]; ok && !comparisons.IsEmpty(v) {
			config += nprintf.NPrintf(`
  name = "%{device-secondary_account_name}"`, ctx)
		}
		config += nprintf.NPrintf(` 
}`, ctx)
	}
	config += nprintf.NPrintf(`
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
		config += nprintf.NPrintf(`
  purchase_order_number = "%{device-purchase_order_number}"`, ctx)
	}
	if _, ok := ctx["device-purchase_order_number"]; ok {
		config += nprintf.NPrintf(`
  order_reference       = "%{device-order_reference}"`, ctx)
	}
	if _, ok := ctx["device-additional_bandwidth"]; ok {
		config += nprintf.NPrintf(`
  additional_bandwidth  = %{device-additional_bandwidth}`, ctx)
	}
	if _, ok := ctx["device-throughput"]; ok {
		config += nprintf.NPrintf(`
  throughput            = %{device-throughput}
  throughput_unit       = "%{device-throughput_unit}"`, ctx)
	}
	if _, ok := ctx["device-hostname"]; ok {
		config += nprintf.NPrintf(`
  hostname              = "%{device-hostname}"`, ctx)
	}
	if _, ok := ctx["device-interface_count"]; ok {
		config += nprintf.NPrintf(`
  interface_count       = %{device-interface_count}`, ctx)
	}
	if _, ok := ctx["acl-resourceName"]; ok {
		config += nprintf.NPrintf(`
  acl_template_id       = equinix_network_acl_template.%{acl-resourceName}.id`, ctx)
	}
	if _, ok := ctx["mgmtAcl-resourceName"]; ok {
		config += nprintf.NPrintf(`
  mgmt_acl_template_uuid = equinix_network_acl_template.%{mgmtAcl-resourceName}.id`, ctx)
	}
	if _, ok := ctx["sshkey-resourceName"]; ok {
		config += nprintf.NPrintf(`
  ssh_key {
    username = "test"
    key_name = equinix_network_ssh_key.%{sshkey-resourceName}.name
  }`, ctx)
	}
	if _, ok := ctx["device-license_file"]; ok {
		config += nprintf.NPrintf(`
  license_file          = "%{device-license_file}"`, ctx)
	}
	if _, ok := ctx["device-vendorConfig_enabled"]; ok {
		config += nprintf.NPrintf(`
  vendor_configuration  = {`, ctx)
		if _, ok := ctx["device-vendorConfig_siteId"]; ok {
			config += nprintf.NPrintf(`
    siteId          = "%{device-vendorConfig_siteId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_systemIpAddress"]; ok {
			config += nprintf.NPrintf(`
    systemIpAddress = "%{device-vendorConfig_systemIpAddress}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_licenseKey"]; ok {
			config += nprintf.NPrintf(`
    licenseKey = "%{device-vendorConfig_licenseKey}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_licenseSecret"]; ok {
			config += nprintf.NPrintf(`
    licenseSecret = "%{device-vendorConfig_licenseSecret}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_controller1"]; ok {
			config += nprintf.NPrintf(`
    controller1 = "%{device-vendorConfig_controller1}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_controller2"]; ok {
			config += nprintf.NPrintf(`
    controller2 = "%{device-vendorConfig_controller2}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_localId"]; ok {
			config += nprintf.NPrintf(`
    localId = "%{device-vendorConfig_localId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_remoteId"]; ok {
			config += nprintf.NPrintf(`
    remoteId = "%{device-vendorConfig_remoteId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_serialNumber"]; ok {
			config += nprintf.NPrintf(`
    serialNumber = "%{device-vendorConfig_serialNumber}"`, ctx)
		}
		config += nprintf.NPrintf(`
  }`, ctx)
	}
	if _, ok := ctx["device-secondary_name"]; ok {
		config += nprintf.NPrintf(`
  secondary_device {
    name                 = "%{device-secondary_name}"`, ctx)
		if _, ok := ctx["device-secondary_metro_code"]; ok {
			config += nprintf.NPrintf(`
    metro_code           = "%{device-secondary_metro_code}"
    account_number       = data.equinix_network_account.test-secondary.number`, ctx)
		} else {
			config += nprintf.NPrintf(`
    metro_code           = "%{device-metro_code}"
    account_number       = data.equinix_network_account.test.number`, ctx)
		}
		config += nprintf.NPrintf(`
    notifications        = %{device-secondary_notifications}`, ctx)
		if _, ok := ctx["device-secondary_additional_bandwidth"]; ok {
			config += nprintf.NPrintf(`
    additional_bandwidth = %{device-secondary_additional_bandwidth}`, ctx)
		}
		if _, ok := ctx["device-secondary_hostname"]; ok {
			config += nprintf.NPrintf(`
    hostname             = "%{device-secondary_hostname}"`, ctx)
		}
		if _, ok := ctx["acl-secondary_resourceName"]; ok {
			config += nprintf.NPrintf(`
    acl_template_id      = equinix_network_acl_template.%{acl-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["mgmtAcl-secondary_resourceName"]; ok {
			config += nprintf.NPrintf(`
    mgmt_acl_template_uuid = equinix_network_acl_template.%{mgmtAcl-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["sshkey-resourceName"]; ok {
			config += nprintf.NPrintf(`
    ssh_key {
      username = "test"
      key_name = equinix_network_ssh_key.%{sshkey-resourceName}.name
    }`, ctx)
		}
		if _, ok := ctx["device-secondary_license_file"]; ok {
			config += nprintf.NPrintf(`
    license_file         = "%{device-secondary_license_file}"`, ctx)
		}
		if _, ok := ctx["device-secondary_vendorConfig_enabled"]; ok {
			config += nprintf.NPrintf(`
    vendor_configuration  = {`, ctx)
			if _, ok := ctx["device-secondary_vendorConfig_siteId"]; ok {
				config += nprintf.NPrintf(`
      siteId          = "%{device-secondary_vendorConfig_siteId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_systemIpAddress"]; ok {
				config += nprintf.NPrintf(`
      systemIpAddress = "%{device-secondary_vendorConfig_systemIpAddress}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_licenseKey"]; ok {
				config += nprintf.NPrintf(`
      licenseKey = "%{device-secondary_vendorConfig_licenseKey}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_licenseSecret"]; ok {
				config += nprintf.NPrintf(`
      licenseSecret = "%{device-secondary_vendorConfig_licenseSecret}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_controller1"]; ok {
				config += nprintf.NPrintf(`
      controller1 = "%{device-secondary_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_controller2"]; ok {
				config += nprintf.NPrintf(`
      controller2 = "%{device-secondary_vendorConfig_controller2}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_localId"]; ok {
				config += nprintf.NPrintf(`
      localId = "%{device-secondary_vendorConfig_localId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_remoteId"]; ok {
				config += nprintf.NPrintf(`
      remoteId = "%{device-secondary_vendorConfig_remoteId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_serialNumber"]; ok {
				config += nprintf.NPrintf(`
      serialNumber = "%{device-secondary_vendorConfig_serialNumber}"`, ctx)
			}
			config += nprintf.NPrintf(`
    }`, ctx)
		}
		config += `
  }`
	}
	if _, ok := ctx["device-cluster_name"]; ok {
		config += nprintf.NPrintf(`
  cluster_details {
    cluster_name        = "%{device-cluster_name}"`, ctx)
		config += `
    node0 {`
		if _, ok := ctx["device-node0_license_file_id"]; ok {
			config += nprintf.NPrintf(`
      license_file_id   = "%{device-node0_license_file_id}"`, ctx)
		}
		if _, ok := ctx["device-node0_license_token"]; ok {
			config += nprintf.NPrintf(`
      license_token     = "%{device-node0_license_token}"`, ctx)
		}
		if _, ok := ctx["device-node0_vendorConfig_enabled"]; ok {
			config += nprintf.NPrintf(`
      vendor_configuration {`, ctx)
			if _, ok := ctx["device-node0_vendorConfig_hostname"]; ok {
				config += nprintf.NPrintf(`
        hostname        = "%{device-node0_vendorConfig_hostname}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_adminPassword"]; ok {
				config += nprintf.NPrintf(`
        admin_password  = "%{device-node0_vendorConfig_adminPassword}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_controller1"]; ok {
				config += nprintf.NPrintf(`
        controller1     = "%{device-node0_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_activationKey"]; ok {
				config += nprintf.NPrintf(`
        activation_key  = "%{device-node0_vendorConfig_activationKey}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_controllerFqdn"]; ok {
				config += nprintf.NPrintf(`
        controller_fqdn = "%{device-node0_vendorConfig_controllerFqdn}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_rootPassword"]; ok {
				config += nprintf.NPrintf(`
        root_password   = "%{device-node0_vendorConfig_rootPassword}"`, ctx)
			}
			config += nprintf.NPrintf(`
      }`, ctx)
		}
		config += `
    }`
		config += `
    node1 {`
		if _, ok := ctx["device-node1_license_file_id"]; ok {
			config += nprintf.NPrintf(`
      license_file_id   = "%{device-node1_license_file_id}"`, ctx)
		}
		if _, ok := ctx["device-node1_license_token"]; ok {
			config += nprintf.NPrintf(`
      license_token     = "%{device-node1_license_token}"`, ctx)
		}
		if _, ok := ctx["device-node1_vendorConfig_enabled"]; ok {
			config += nprintf.NPrintf(`
      vendor_configuration {`, ctx)
			if _, ok := ctx["device-node1_vendorConfig_hostname"]; ok {
				config += nprintf.NPrintf(`
        hostname        = "%{device-node1_vendorConfig_hostname}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_adminPassword"]; ok {
				config += nprintf.NPrintf(`
        admin_password  = "%{device-node1_vendorConfig_adminPassword}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_controller1"]; ok {
				config += nprintf.NPrintf(`
        controller1     = "%{device-node1_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_activationKey"]; ok {
				config += nprintf.NPrintf(`
        activation_key  = "%{device-node1_vendorConfig_activationKey}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_controllerFqdn"]; ok {
				config += nprintf.NPrintf(`
        controller_fqdn = "%{device-node1_vendorConfig_controllerFqdn}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_rootPassword"]; ok {
				config += nprintf.NPrintf(`
        root_password   = "%{device-node1_vendorConfig_rootPassword}"`, ctx)
			}
			config += nprintf.NPrintf(`
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

func testAccNetworkDeviceACL(ctx map[string]any) string {
	var config string
	if _, ok := ctx["acl-name"]; ok {
		config += nprintf.NPrintf(`
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
		config += nprintf.NPrintf(`
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
		config += nprintf.NPrintf(`
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
		config += nprintf.NPrintf(`
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

func testAccNetworkDeviceSSHKey(ctx map[string]any) string {
	return nprintf.NPrintf(`
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
