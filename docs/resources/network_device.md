---
subcategory: "Network Edge"
---

# equinix_network_device (Resource)

Resource `equinix_network_device` allows creation and management of Equinix Network Edge virtual network devices.

Network Edge virtual network devices can be created in two modes:

* **managed** - (default) Where Equinix manages connectivity and services in the device and customer gets limited access to the device.
* **self-configured** - Where customer provisions and manages own services in the device with less restricted access. Some device types are offered only in this mode.

In addition to management modes, there are two software license modes available:

* **subscription** - Where Equinix provides software license, including end-to-end support, and bills for the service respectively.
* **BYOL** - [bring your own license] Where customer brings his own, already procured device software license. There are no charges associated with such license. It is the only licensing mode for `self-configured` devices.

## Example Usage

```terraform
# Create pair of redundant, managed CSR1000V routers with license subscription
# in two different metro locations

data "equinix_network_account" "dc" {
  metro_code = "DC"
}

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "csr1000v-ha" {
  name            = "tf-csr1000v-p"
  throughput      = 500
  throughput_unit = "Mbps"
  metro_code      = data.equinix_network_account.dc.metro_code
  type_code       = "CSR1000V"
  self_managed    = false
  connectivity    = "INTERNET-ACCESS"
  byol            = false
  package_code    = "SEC"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  hostname        = "csr1000v-p"
  term_length     = 12
  account_number  = data.equinix_network_account.dc.number
  version         = "16.09.05"
  core_count      = 2
  secondary_device {
    name            = "tf-csr1000v-s"
    metro_code      = data.equinix_network_account.sv.metro_code
    hostname        = "csr1000v-s"
    notifications   = ["john@equinix.com", "marry@equinix.com"]
    account_number  = data.equinix_network_account.sv.number
  }
}
```

```terraform
# Create self configured PANW cluster with BYOL license

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "panw-cluster" {
  name            = "tf-panw"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "PA-VM"
  self_managed    = true
  byol            = true
  package_code    = "VM100"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "10.1.3"
  interface_count = 10
  core_count      = 2
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
  acl_template_id = "0bff6e05-f0e7-44cd-804a-25b92b835f8b"
  cluster_details {
    cluster_name    = "tf-panw-cluster"
    node0 {
      vendor_configuration {
        hostname = "panw-node0"
      }
      license_token = "licenseToken"
    }
    node1 {
      vendor_configuration {
        hostname = "panw-node1"
      }
      license_token = "licenseToken"
    }
  }
}
```

```terraform
# Create self configured single Aviatrix device with cloud init file

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

variable "filepath" { default = "cloudInitFileFolder/TF-AVX-cloud-init-file.txt" }

resource "equinix_network_file" "aviatrix-cloudinit-file" {
  file_name = "TF-AVX-cloud-init-file.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = data.equinix_network_account.sv.metro_code
  device_type_code = "AVIATRIX_EDGE"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}

resource "equinix_network_device" "aviatrix-single" {
  name            = "tf-aviatrix"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "AVIATRIX_EDGE"
  self_managed    = true
  byol            = true
  package_code    = "STD"
  notifications   = ["john@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "6.9"
  core_count      = 2
  cloud_init_file_id = equinix_network_file.aviatrix-cloudinit-file.uuid
  acl_template_id = "c06150ea-b604-4ad1-832a-d63936e9b938"
}
```

```terraform
# Create self configured single Catalyst 8000V (Autonomous Mode) router with license token

data "equinix_network_account" "sv" {
  name = "account-name"
  metro_code = "SV"
}

resource "equinix_network_device" "c8kv-single" {
  name            = "tf-c8kv"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "C8000V"
  self_managed    = true
  byol            = true
  package_code    = "network-essentials"
  notifications   = ["test@equinix.com"]
  hostname        = "C8KV"
  account_number  = data.equinix_network_account.sv.number
  version         = "17.06.01a"
  core_count      = 2
  term_length     = 12
  license_token = "valid-license-token"
  additional_bandwidth = 5
  ssh_key {
    username = "test-username"
    key_name = "valid-key-name"
  }
  acl_template_id = "3e548c02-9164-4197-aa23-05b1f644883c"
}
```

```terraform
# Create self configured single VSRX device with BYOL License

data "equinix_network_account" "sv" {
  name = "account-name"
  metro_code = "SV"
}

resource "equinix_network_device" "vsrx-single" {
  name            = "tf-c8kv-sdwan"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "VSRX"
  self_managed    = true
  byol            = true
  package_code    = "STD"
  notifications   = ["test@equinix.com"]
  hostname        = "VSRX"
  account_number  = data.equinix_network_account.sv.number
  version         = "23.2R1.13"
  core_count      = 2
  term_length     = 12
  additional_bandwidth = 5
  project_id      = "a86d7112-d740-4758-9c9c-31e66373746b"
  diverse_device_id = "ed7891bd-15b4-4f72-ac56-d96cfdacddcc"
  ssh_key {
    username = "test-username"
    key_name = "valid-key-name"
  }
  acl_template_id = "3e548c02-9164-4197-aa23-05b1f644883c"
}
```

```terraform
# Create self configured redundant Arista router with DSA key

data "equinix_network_account" "sv" {
  name = "account-name"
  metro_code = "SV"
}

resource "equinix_network_ssh_key" "test-public-key" {
  name = "key-name"
  public_key = "ssh-dss key-value"
  type = "DSA"
}

resource "equinix_network_device" "arista-ha" {
  name            = "tf-arista-p"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "ARISTA-ROUTER"
  self_managed    = true
  connectivity    = "PRIVATE"
  byol            = true
  package_code    = "CloudEOS"
  notifications   = ["test@equinix.com"]
  hostname        = "arista-p"
  account_number  = data.equinix_network_account.sv.number
  version         = "4.29.0"
  core_count      = 4
  term_length     = 12
  additional_bandwidth = 5
  ssh_key {
    username = "test-username"
    key_name = equinix_network_ssh_key.test-public-key.name
  }
  acl_template_id = "c637a17b-7a6a-4486-924b-30e6c36904b0"
  secondary_device {
    name            = "tf-arista-s"
    metro_code      = data.equinix_network_account.sv.metro_code
    hostname        = "arista-s"
    notifications   = ["test@eq.com"]
    account_number  = data.equinix_network_account.sv.number
    acl_template_id = "fee5e2c0-6198-4ce6-9cbd-bbe6c1dbe138"
  }
}
```

```terraform
# Create self configured redundant BlueCat DNS and DHCP Server
data "equinix_network_account" "sv" {
  name       = "account-name"
  metro_code = "SV"
}

resource "equinix_network_ssh_key" "test-public-key" {
  name       = "key-name"
  public_key = "ssh-dss key-value"
  type       = "DSA"
}

resource "equinix_network_device" "bluecat-bdds-ha" {
  name                 = "tf-bluecat-bdds-p"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "BLUECAT"
  self_managed         = true
  connectivity         = "PRIVATE"
  byol                 = true
  package_code         = "STD"
  notifications        = ["test@equinix.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "9.6.0"
  core_count           = 2
  term_length          = 12
  vendor_configuration = {
    "hostname" = "test"
    "privateAddress" : "x.x.x.x"
    "privateCidrMask" : "24"
    "privateGateway" : "x.x.x.x"
    "licenseKey" : "xxxxx-xxxxx-xxxxx-xxxxx-xxxxx"
    "licenseId" : "xxxxxxxxxxxxxxx"
  }
  ssh_key {
    username = "test-username"
    key_name = equinix_network_ssh_key.test-public-key.name
  }
  secondary_device {
    name                 = "tf-bluecat-bdds-s"
    metro_code           = data.equinix_network_account.sv.metro_code
    notifications        = ["test@eq.com"]
    account_number       = data.equinix_network_account.sv.number
    vendor_configuration = {
      "hostname" = "test"
      "privateAddress" : "x.x.x.x"
      "privateCidrMask" : "24"
      "privateGateway" : "x.x.x.x"
      "licenseKey" : "xxxxx-xxxxx-xxxxx-xxxxx-xxxxx"
      "licenseId" : "xxxxxxxxxxxxxxx"
    }
  }
}
```

```terraform
# Create self configured redundant BlueCat Edge Service Point
data "equinix_network_account" "sv" {
  name = "account-name"
  metro_code = "SV"
}

resource "equinix_network_file" "bluecat-edge-service-point-cloudinit-primary-file" {
  file_name = "TF-BLUECAT-ESP-cloud-init-file.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = data.equinix_network_account.sv.metro_code
  device_type_code = "BLUECAT-EDGE-SERVICE-POINT"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}

resource "equinix_network_file" "bluecat-edge-service-point-cloudinit-secondary-file" {
  file_name = "TF-BLUECAT-ESP-cloud-init-file.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = data.equinix_network_account.sv.metro_code
  device_type_code = "BLUECAT-EDGE-SERVICE-POINT"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}

resource "equinix_network_device" "bluecat-edge-service-point-ha" {
  name            = "tf-bluecat-edge-service-point-p"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "BLUECAT-EDGE-SERVICE-POINT"
  self_managed    = true
  connectivity    = "PRIVATE"
  byol            = true
  package_code    = "STD"
  notifications   = ["test@equinix.com"]
  account_number  = data.equinix_network_account.sv.number
  cloud_init_file_id = equinix_network_file.bluecat-edge-service-point-cloudinit-primary-file.uuid
  version         = "4.6.3"
  core_count      = 4
  term_length     = 12
  secondary_device {
    name            = "tf-bluecat-edge-service-point-s"
    metro_code      = data.equinix_network_account.sv.metro_code
    notifications   = ["test@eq.com"]
    account_number  = data.equinix_network_account.sv.number
    cloud_init_file_id = equinix_network_file.bluecat-edge-service-point-cloudinit-secondary-file.uuid
  }
}
```

```terraform
# Create PA-VM firewall cluster with Panorama Server Integration
# with Panorama Server IP and Panorama Auth Key in vendor Configuration

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "panw-cluster" {
  name            = "tf-panw"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "PA-VM"
  self_managed    = true
  byol            = true
  package_code    = "VM100"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "11.1.3"
  interface_count = 10
  core_count      = 2
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
  acl_template_id = "0bff6e05-f0e7-44cd-804a-25b92b835f8b"
  cluster_details {
    cluster_name = "tf-panw-cluster"
    node0 {
      vendor_configuration {
        hostname            = "panw-node0"
        panorama_ip_address = "x.x.x.x"
        panorama_auth_key   = "xxxxxxxxxxx"
      }
      license_token = "licenseToken"
    }
    node1 {
      vendor_configuration {
        hostname            = "panw-node1"
        panorama_ip_address = "x.x.x.x"
        panorama_auth_key   = "xxxxxxxxxxx"
      }
      license_token = "licenseToken"
    }
  }
}
```

```terraform
# Create C8000V BYOL device with bandwidth tier information

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "c8000v-byol-tier" {
  name            = "tf-c8000v-byol"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "C8000V"
  self_managed    = true
  byol            = true
  package_code    = "VM100"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "17.11.01a"
  interface_count = 10
  core_count      = 2
  tier            = 1
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
  acl_template_id = "0bff6e05-f0e7-44cd-804a-25b92b835f8b"
}
```

```terraform
# Create C8000V BYOL device with numeric bandwidth throughput information

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "c8000v-byol-throughput" {
  name            = "tf-c8000v-byol"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "C8000V"
  self_managed    = true
  byol            = true
  package_code    = "VM100"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "17.11.01a"
  interface_count = 10
  core_count      = 2
  throughput 	  = "100"
  throughput_unit = "Mbps"
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
  acl_template_id = "0bff6e05-f0e7-44cd-804a-25b92b835f8b"
}
```

```terraform
# Create self configured single Aviatrix Transit Edge device with cloud init file

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

variable "filepath" { default = "cloudInitFileFolder/TF-AVX-cloud-init-file.txt" }

resource "equinix_network_file" "aviatrix-cloudinit-file" {
  file_name = "TF-AVX-cloud-init-file.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = data.equinix_network_account.sv.metro_code
  device_type_code = "AVIATRIX_TRANSIT_EDGE"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}

resource "equinix_network_device" "aviatrix-transit-edge-single" {
  name            = "tf-aviatrix"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "AVIATRIX_TRANSIT_EDGE"
  self_managed    = true
  byol            = true
  package_code    = "STD"
  notifications   = ["john@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "7.2.a"
  core_count      = 2
  cloud_init_file_id = equinix_network_file.aviatrix-cloudinit-file.uuid
  acl_template_id = "c06150ea-b604-4ad1-832a-d63936e9b938"
}
```

```terraform
# Create ZSCALER APPC device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "zscaler-appc-single" {
  name            = "tf-zscaler-appc"
  project_id      = "XXXXXX"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "ZSCALER-APPC"
  self_managed    = true
  byol            = true
  connectivity    = "PRIVATE"
  package_code    = "STD"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "23.395.1"
  interface_count = 1
  core_count      = 4
  vendor_configuration   = {"provisioningKey" = "XXXXXXXXXX", "hostname" = "XXXX"}
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
}
```

```terraform
# Create ZSCALER APPC device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "zscaler-pse-single" {
  name            = "tf-zscaler-pse"
  project_id      = "XXXXXX"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "ZSCALER-PSE"
  self_managed    = true
  byol            = true
  connectivity    = "PRIVATE"
  package_code    = "STD"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "23.395.1"
  interface_count = 1
  core_count      = 4
  vendor_configuration   = {"provisioningKey" = "XXXXXXXXXX", "hostname" = "XXXX"}
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Device name.
* `type_code` - (Required) Device type code.
* `metro_code` - (Required) Device location metro code.
* `hostname` - (Optional) Device hostname prefix.
* `package_code` - (Required) Device software package code.
* `version` - (Required) Device software software version.
* `core_count` - (Required) Number of CPU cores used by device. (**NOTE: Use this field to resize your device. When resizing your HA devices, primary device will be upgraded first. If the upgrade failed, device will be automatically rolled back to the previous state with original core number.**)
* `tier` - (Optional, conflicts with `throughput`,`throughput_unit` ) Select bandwidth tier for your own license, i.e., `0` or `1` or `2` or `3`. Tiers applicable only for C8000V Autonomous or C8000V SDWAN (controller) device types. If not provided, tier is defaulted to '2'.
* `term_length` - (Required) Device term length.
* `self_managed` - (Optional) Boolean value that determines device management mode, i.e., `self-managed` or `Equinix-managed` (default).
* `byol` - (Optional) Boolean value that determines device licensing mode, i.e., `bring your own license` or `subscription` (default).
* `license_token` - (Optional, conflicts with `license_file`) License Token applicable for some device types in BYOL licensing mode.
* `license_file` - (Optional) Path to the license file that will be uploaded and applied on a device. Applicable for some device types in BYOL licensing mode.
* `license_file_id` - (Optional, conflicts with `license_file`) Identifier of a license file that will be applied on the device.
* `cloud_init_file_id` - (Optional) Identifier of a cloud init file that will be applied on the device.
* `throughput` - (Optional) Device license throughput.
* `throughput_unit` - (Optional) License throughput unit. One of `Mbps` or `Gbps`.
* `account_number` - (Required) Billing account number for a device.
* `notifications` - (Required) List of email addresses that will receive device status notifications.
* `purchase_order_number` - (Optional) Purchase order number associated with a device order.
* `order_reference` - (Optional) Name/number used to identify device order on the invoice.
* `acl_template_id` - (Optional) Identifier of a WAN interface ACL template that will be applied on the device.
* `mgmt_acl_template_uuid` - (Optional) Identifier of an MGMT interface ACL template that will be applied on the device.
* `additional_bandwidth` - (Optional) Additional Internet bandwidth, in Mbps, that will be allocated to the device (in addition to default 15Mbps).
* `interface_count` - (Optional) Number of network interfaces on a device. If not specified, default number for a given device type will be used.
* `wan_interafce_id` - (Optional) Specify the WAN/SSH interface id. If not specified, default WAN/SSH interface for a given device type will be used.
* `vendor_configuration` - (Optional) Map of vendor specific configuration parameters for a device (controller1, activationKey, managementType, siteId, systemIpAddress, privateAddress, privateCidrMask, privateGateway, licenseKey, licenseId, panoramaAuthKey, panoramaIpAddress, provisioningKey)
* `ssh-key` - (Optional) Definition of SSH key that will be provisioned on a device (max one key). See [SSH Key](#ssh-key) below for more details.
* `secondary_device` - (Optional) Definition of secondary device for redundant device configurations. See [Secondary Device](#secondary-device) below for more details.
* `cluster_details` - (Optional) An object that has the cluster details. See [Cluster Details](#cluster-details) below for more details.
* `connectivity` - (Optional) Device accessibility (INTERNET-ACCESS or PRIVATE or INTERNET-ACCESS-WITH-PRVT-MGMT). If not specified, default will be INTERNET-ACCESS
* `project_id` - (Optional) Unique Identifier for the project resource where the device is scoped to.If you leave it out, the device will be created under the default project id of your organization.
* `diverse_device_id` - (Optional) Unique ID of an existing device. Use this field to let Equinix know if you want your new device to be in a different location from any existing virtual device. This field is only meaningful for single devices.

### Secondary Device

-> **NOTE:** Network Edge provides different High Availability (HA) options. By defining a `secondary_device` block, terraform will deploy [Redundant Devices](https://docs.equinix.com/en-us/Content/Interconnection/NE/deploy-guide/Reference%20Architecture/NE-High-Availability-Options.htm#:~:text=Redundant%20Devices%20(Active/Active)), useful for customers that require two actively forwarding data planes (Active/Active) on separate hardware stacks. See [Architecting for Resiliency](https://docs.equinix.com/en-us/Content/Interconnection/NE/deploy-guide/NE-architecting-resiliency.htm) documentation to know more about the fault-tolerant solutions that you can achieve.

The `secondary_device` block supports the following arguments:

* `name` - (Required) Secondary device name.
* `metro_code` - (Required) Metro location of a secondary device.
* `hostname` - (Optional) Secondary device hostname.
* `license_token` - (Optional, conflicts with `license_file`) License Token can be provided for some device types o the device.
* `license_file` - (Optional) Path to the license file that will be uploaded and applied on a secondary device. Applicable for some device types in BYOL licensing mode.
* `license_file_id` - (Optional, conflicts with `license_file`) Identifier of a license file that will be applied on a secondary device.
* `cloud_init_file_id` - (Optional) Identifier of a cloud init file that will be applied on a secondary device.
* `account_number` - (Required) Billing account number for secondary device.
* `notifications` - (Required) List of email addresses that will receive notifications about secondary device.
* `additional_bandwidth` - (Optional) Additional Internet bandwidth, in Mbps, for a secondary device.
* `vendor_configuration` - (Optional) Key/Value pairs of vendor specific configuration parameters for a secondary device. Key values are `controller1`, `activationKey`, `managementType`, `siteId`, `systemIpAddress`, `privateAddress`, `privateCidrMask`, `privateGateway`, `licenseKey`, `licenseId`, `panoramaAuthKey`, `panoramaIpAddress`.
* `acl_template_id` - (Optional) Identifier of a WAN interface ACL template that will be applied on a secondary device.
* `mgmt_acl_template_uuid` - (Optional) Identifier of an MGMT interface ACL template that will be applied on a secondary device.
* `ssh-key` - (Optional) Up to one definition of SSH key that will be provisioned on a secondary device.

### SSH Key

The `ssh_key` block supports the following arguments:

* `username` - (Required) username associated with given key.
* `name` - (Optional) reference by name to previously provisioned public SSH key.

### Cluster Details

-> **NOTE:** Network Edge provides different High Availability (HA) options. By defining a `cluster_details` block, terraform will deploy a `Device Clustering`. This option, based on vendor-specific features, allows customers to deploy more advanced resilient configurations than `secondary_device`. See [Network Edge HA Options](https://docs.equinix.com/en-us/Content/Interconnection/NE/deploy-guide/Reference%20Architecture/NE-High-Availability-Options.htm) documentation to know which vendors support clustered devices. See [Architecting for Resiliency](https://docs.equinix.com/en-us/Content/Interconnection/NE/deploy-guide/NE-architecting-resiliency.htm) documentation to know more about the fault-tolerant solutions that you can achieve.

The `cluster_details` block supports the following arguments:

* `cluster_name` - (Required) The name of the cluster device
* `node0` - (Required) An object that has `node0` configuration. See [Cluster Details - Nodes](#cluster-details---nodes) below for more details.
* `node1` - (Required) An object that has `node1` configuration. See [Cluster Details - Nodes](#cluster-details---nodes) below for more details.

### Cluster Details - Nodes

The `node0` and `node1` blocks supports the following arguments:

* `vendor_configuration` - (Optional) An object that has fields relevant to the vendor of the cluster device. See [Cluster Details - Nodes - Vendor Configuration](#cluster-details---nodes---vendor-configuration) below for more details.
* `license_file_id` - (Optional) License file id. This is necessary for Fortinet and Juniper clusters.
* `license_token` - (Optional) License token. This is necessary for Palo Alto clusters.

### Cluster Details - Nodes - Vendor Configuration

The `vendor_configuration` block supports the following arguments:

* `hostname` - (Optional) Hostname. This is necessary for Palo Alto, Juniper, and Fortinet clusters.
* `admin_password` - (Optional) The administrative password of the device. You can use it to log in to the console. This field is not available for all device types.
* `controller1` - (Optional) System IP Address. Mandatory for the Fortinet SDWAN cluster device.
* `activation_key` - (Optional) Activation key. This is required for Velocloud clusters.
* `controller_fqdn` - (Optional) Controller fqdn. This is required for Velocloud clusters.
* `root_password` - (Optional) The CLI password of the device. This field is relevant only for the Velocloud SDWAN cluster.
* `panorama_ip_address` - (Optional) Panorama Server IP Address. This field is relevant only for the PA-VM firewall devices to have integration with Panorama Server.
* `panorama_auth_key` - (Optional) Panorama Server Auth Key. This field is relevant only for the PA-VM firewall devices to have integration with Panorama Server.
* `provisioning_key` - (Optional) Provisioning Key. This field is relevant only for the ZSCALER APPC and ZSCALER PSE devices.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Device unique identifier.
* `status` - Device provisioning status. Possible values are `INITIALIZING`, `PROVISIONING`, `WAITING_FOR_PRIMARY`, `WAITING_FOR_SECONDARY`, `WAITING_FOR_REPLICA_CLUSTER_NODES`, `CLUSTER_SETUP_IN_PROGRESS`, `FAILED`, `PROVISIONED`, `DEPROVISIONING`, `DEPROVISIONED`, `RESOURCE_UPGRADE_IN_PROGRESS`, `RESOURCE_UPGRADE_FAILED`.
* `license_status` - Device license registration status. Possible values are `APPLYING_LICENSE`, `REGISTERED`, `APPLIED`, `WAITING_FOR_CLUSTER_SETUP`, `REGISTRATION_FAILED`.
* `license_file_id` - Unique identifier of applied license file.
* `ibx` - Device location Equinix Business Exchange name.
* `region` - Device location region.
* `acl_template_id` - Unique identifier of applied ACL template.
* `ssh_ip_address` - IP address of SSH enabled interface on the device.
* `ssh_ip_fqdn` - FQDN of SSH enabled interface on the device.
* `redundancy_type` - Device redundancy type applicable for HA devices, either primary or secondary.
* `redundant_id` - Unique identifier for a redundant device applicable for HA devices.
* `interface` - List of device interfaces. See [Interface Attribute](#interface-attribute) below for more details.
* `asn` - (Autonomous System Number) Unique identifier for a network on the internet.
* `zone_code` - Device location zone code.
* `cluster_id` - The ID of the cluster.
* `num_of_nodes` - The number of nodes in the cluster.
* `diverse_device_name` - Name of the device with diverse device UUID. This field is returned in device details if the device is created by passing diverse_device_id.

### Interface Attribute

Each interface attribute has below fields:

* `id` - interface identifier.
* `name` - interface name.
* `status` - interface status. One of `AVAILABLE`, `RESERVED`, `ASSIGNED`.
* `operational_status` - interface operational status. One of `up`, `down`.
* `mac_address` - interface MAC address.
* `ip_address` - interface IP address.
* `assigned_type` - interface management type (Equinix Managed or empty).
* `type` - interface type.

## Timeouts

This resource provides the following [Timeouts configuration](https://www.terraform.io/language/resources/syntax#operation-timeouts) options:

* create - Default is 90 minutes
* update - Default is 90 minutes
* delete - Default is 30 minutes

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_device.example {existing_id}
```

The `license_token`, `mgmt_acl_template_uuid` and `cloud_init_file_id` fields can not be imported.
