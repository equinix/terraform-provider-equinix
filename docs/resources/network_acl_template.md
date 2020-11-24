---
layout: "equinix"
page_title: "Equinix: equinix_network_acl_template"
subcategory: ""
description: |-
 Provides Network Edge ACL template resource
---

# Resource: equinix_network_acl_template

Resource `equinix_network_acl_template` allows creation and management of
Network Edge device Access Control List templates.

Device ACL templates give possibility to define set of rules will allowed inbound
traffic. Templates can be assigned to the network devices.

## Example Usage

```hcl
# Creates ACL template and assigns it to the network device
resource "equinix_network_acl_template" "myacl" {
  name        = "test"
  description = "Test ACL template"
  metro_code  = equinix_network_device.csr1000v.metro_code
  inbound_rule {
    subnets  = ["1.1.1.1/32"]
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
  inbound_rule {
    subnets  = ["172.16.25.0/24"]
    protocol = "UDP"
    src_port = "any"
    dst_port = "53,1045,2041"
  }
  inbound_rule {
    subnets  = ["192.168.0.0/16", "10.0.0.0/8"]
    protocol = "TCP"
    src_port = "any"
    dst_port = "22-23"
  }
}
```

## Argument Reference

* `name` - (Required) ACL template name
* `description` - (Optional) ACL template description
* `metro_code` - (Required) ACL template location metro code
* `inbound_rule` - (Required) One or more rules to specify allowed inbound traffic.
Rules are ordered, matching traffic rule stops processing subsequent ones.
  * `inbound_rule.#.subnets` - (Required) Inbound traffic source IP subnets
  un CIDR format
  * `inbound_rule.#.protocol` - (Required) Inbound traffic protocol.
  One of: `IP`, `TCP`, `UDP`
  * `inbound_rule.#.src_port` - (Required) Inbound traffic source ports.
  Allowed values are:
    * up to 10, comma separated ports (i.e. `20,22,23`)
    * port range (i.e. `1023-1040`
    * word `any`
  * `inbound_rule.#.dst_port` - (Required) Inbound traffic destination ports.
  Allowed values are:
    * up to 10, comma separated ports (i.e. `20,22,23`)
    * port range (i.e. `1023-1040`)
    * word `any`

## Attributes Reference

* `uuid` - Unique universal identifier of ACL template resource
* `device_id` - Identifier of a network device that template was applied on
* `device_acl_status` - Status of ACL template provisioning process on a device,
that template was applied on. One of:
  * PROVISIONING
  * PROVISIONED
