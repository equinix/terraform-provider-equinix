---
subcategory: "Metal"
---

# equinix_metal_precreated_ip_block (Data Source)

Use this data source to get CIDR expression for precreated (management) IPv6 and IPv4 blocks in Equinix Metal. You can then use the cidrsubnet TF builtin function to derive subnets.

~> For backward compatibility, this data source will also return reserved (elastic) IP blocks.

-> Precreated (management) IP blocks for a metro will not be available until first device is created in that metro.

-> Public IPv4 blocks auto-assigned (management) to a device cannot be retrieved. If you need that information, consider using the [equinix_metal_device](equinix_metal_device.md) data source instead.

## Example Usage

```terraform
# Create device in your project and then assign /64 subnet from precreated block
# to the new device

# Declare your project ID
locals {
  project_id = "<UUID_of_your_project>"
}

resource "equinix_metal_device" "web1" {
  hostname         = "web1"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id

}

data "equinix_metal_precreated_ip_block" "test" {
  metro          = "sv"
  project_id     = local.project_id
  address_family = 6
  public         = true
}

# The precreated IPv6 blocks are /56, so to get /64, we specify 8 more bits for network.
# The cirdsubnet interpolation will pick second /64 subnet from the precreated block.

resource "equinix_metal_ip_attachment" "from_ipv6_block" {
  device_id     = equinix_metal_device.web1.id
  cidr_notation = cidrsubnet(data.equinix_metal_precreated_ip_block.test.cidr_notation, 8, 2)
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) ID of the project where the searched block should be.
* `address_family` - (Required) 4 or 6, depending on which block you are looking for.
* `public` - (Required) Whether to look for public or private block.
* `global` - (Optional) Whether to look for global block. Default is false for backward compatibility.
* `facility` - (**Deprecated**) Facility of the searched block. (for non-global blocks). Use metro instead; read the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices)
* `metro` - (Optional) Metro of the searched block (for non-global blocks).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cidr_notation` - CIDR notation of the looked up block.
