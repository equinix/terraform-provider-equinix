---
layout: "packet"
page_title: "Packet: packet_project"
sidebar_current: "docs-packet-resource-project"
description: |-
  Provides a Packet Project resource.
---

# packet\_project

Provides a Packet Project resource to allow you manage devices
in your projects.

## Example Usage

```hcl
# Create a new Project
resource "packet_project" "tf_project_1" {
  name           = "Terraform Fun"
}
```

Example with BGP config
```hcl
# Create a new Project
resource "packet_project" "tf_project_1" {
  name           = "tftest"
  bgp_config {
    deployment_type = "local"
    md5 = "C179c28c41a85b"
    asn = 65000
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Project on Packet.net
* `payment_method_id` - The UUID of payment method for this project. If you keep it empty, Packet API will pick your default Payment Method.
* `organization_id` - The UUID of Organization under which you want to create the project. If you leave it out, the project will be create under your the default Organization of your account.
* bgp_config - Optional BGP settings. Refer to [Packet guide for BGP](https://support.packet.com/kb/articles/bgp).

The `bgp_config` block supports:
* `asn` - Autonomous System Numer for local BGP deployment
* `md5` - (Optional) MD5 sum of password for BGP session
* `deployment_type` - `private` or `public`, the `private` is likely to be usable immediately, the `public` will need to be review by Packet engineers

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the project
* `payment_method_id` - The UUID of payment method for this project.
* `organization_id` - The UUID of this project's parent organization.
* `created` - The timestamp for when the Project was created
* `updated` - The timestamp for the last time the Project was updated

The `bgp_config` block additionally exports: 
* `status` - status of BGP configuration in the project
* `max_prefix` - `private` or `public`, the `private` is likely to be usable immediately, the `public` will need to be review by Packet engineers
* `route_object`
* `ranges`
* `max_prefix`
