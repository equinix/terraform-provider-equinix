---
layout: "equinix"
page_title: "Equinix: equinix_network_acl_template"
subcategory: ""
description: |-
 Provides Equinix Network Edge ACL template resource
---

# Resource: equinix_network_acl_template

Resource `equinix_network_acl_template` allows creation and management of
Equinix Network Edge device Access Control List templates.

Device ACL templates give possibility to define set of rules will allowed inbound
traffic. Templates can be assigned to the network devices.

## Example Usage

```hcl
# Creates ACL template and assigns it to the network device
resource "equinix_network_acl_template" "myacl" {
  name        = "test"
  description = "Test ACL template"
  inbound_rule {
    subnet  = "1.1.1.1/32"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
  inbound_rule {
    subnet  = "172.16.25.0/24"
    protocol = "UDP"
    src_port = "any"
    dst_port = "53,1045,2041"
  }
}
```

## Argument Reference

* `name` - (Required) ACL template name
* `description` - (Optional) ACL template description
* `metro_code` - (Deprecated) ACL template location metro code
* `inbound_rule` - (Required) One or more rules to specify allowed inbound traffic.
Rules are ordered, matching traffic rule stops processing subsequent ones.
  * `inbound_rule.#.subnets` - (Deprecated) Inbound traffic source IP subnets
  in CIDR format
  * `inbound_rule.#.subnet` - (Required) Inbound traffic source IP subnet
    in CIDR format
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

* `uuid` - Unique identifier of ACL template resource
* `device_id` - (Deprecated) Identifier of a network device where template was applied
* `device_acl_status` - Status of ACL template provisioning process,
  where template was applied. One of:
  * PROVISIONING
  * PROVISIONED
* `device_details` - List of the devices where the ACL template is applied,
  * `uuid` - Device uuid
  * `name` - Device name
  * `acl_status` - Device acl provisioning status
    where template was applied. One of:
    * PROVISIONING
    * PROVISIONED


## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_acl_template.example {existing_id}
```
