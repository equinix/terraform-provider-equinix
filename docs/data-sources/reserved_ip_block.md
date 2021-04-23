---
page_title: "Equinix Metal: reserved_ip_block"
subcategory: ""
description: |-
Look up an IP address block
---

# metal\_reserved\_ip\_block

Use this data source to find IP address blocks in Equinix Metal. You can use IP address or a block ID for lookup.

## Example Usage

Look up an IP address for a domain name, then use the IP to look up the containing IP block and run a device with IP address from the block:

```hcl
data "dns_a_record_set" "www" {
  host = "www.example.com"
}

data "metal_reserved_ip_block" "www" {
  project_id = local.my_project_id
  address = data.dns_a_record_set.www.addrs[0]
}

resource "metal_device" "www" {
  project_id = local.my_project_id
  [...]
  ip_address {
    type = "public_ipv4"
    reservation_ids = [data.metal_reserved_ip_block.www.id]
  }
}
```

## Argument Reference

You should pass either `id`, or both `project_id` and `ip_address`.

* `id` - (Required) UUID of the IP address block to look up
* `project_id` - (Required) UUID of the project where the searched block should be
* `ip_address` - (Required) Block containing this IP address will be returned

## Attributes Reference

This datasource exposes the same attributes as the [metal_reserved_ip_block resource](../resources/reserved_ip_block.md).

