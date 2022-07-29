---
subcategory: "Network Edge"
---

# equinix_network_acl_template (Resource)

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
    description = "inbound rule description"
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

The following arguments are supported:

* `name` - (Required) ACL template name.
* `description` - (Optional) ACL template description, up to 200 characters.
* `metro_code` - (Deprecated) ACL template location metro code.
* `inbound_rule` - (Required) One or more rules to specify allowed inbound traffic.
Rules are ordered, matching traffic rule stops processing subsequent ones.

The `inbound_rule` block has below fields:

* `subnets` - (Deprecated) Inbound traffic source IP subnets in CIDR format.
* `subnet` - (Required) Inbound traffic source IP subnet in CIDR format.
* `protocol` - (Required) Inbound traffic protocol. One of `IP`, `TCP`, `UDP`.
* `src_port` - (Required) Inbound traffic source ports. Allowed values are a comma separated list
of ports, e.g., `20,22,23`, port range, e.g., `1023-1040` or word `any`.
* `dst_port` - (Required) Inbound traffic destination ports. Allowed values are a comma separated
list of ports, e.g., `20,22,23`, port range, e.g., `1023-1040` or word `any`.
* `description` - (Optional) Inbound rule description, up to 200 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Unique identifier of ACL template resource.
* `device_id` - (Deprecated) Identifier of a network device where template was applied.
* `device_acl_status` - Status of ACL template provisioning process, where template was applied.
One of `PROVISIONING`, `PROVISIONED`.
* `device_details` - List of the devices where the ACL template is applied.

The `device_details` block has below fields:

* `uuid` - Device uuid.
* `name` - Device name.
* `acl_status` - Device ACL provisioning status where template was applied. One of `PROVISIONING`,
`PROVISIONED`.

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_acl_template.example {existing_id}
```
