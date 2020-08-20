---
layout: "packet"
page_title: "Packet: packet_project"
sidebar_current: "docs-packet-resource-project"
description: |-
  Provides a Packet Project resource.
---

# packet\_project

Provides a Packet project resource to allow you manage devices
in your projects.

## Example Usage

```hcl
# Create a new project
resource "packet_project" "tf_project_1" {
  name = "Terraform Fun"
}
```

Example with BGP config

```hcl
# Create a new Project
resource "packet_project" "tf_project_1" {
  name = "tftest"
  bgp_config {
    deployment_type = "local"
    md5             = "C179c28c41a85b"
    asn             = 65000
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project
* `organization_id` - The UUID of organization under which you want to create the project. If you leave it out, the project will be create under your the default organization of your account.
* `payment_method_id` - The UUID of payment method for this project. The payment method and the project need to belong to the same organization (passed with `organization_id`, or default).
* `backend_transfer` - Enable or disable [Backend Transfer](https://www.packet.com/developers/docs/network/basic/backend-transfer/), default is false
* `bgp_config` - Optional BGP settings. Refer to [Packet guide for BGP](https://www.packet.com/developers/docs/network/advanced/local-and-global-bgp/).

Once you set the BGP config in a project, it can't be removed (due to a limitation in the Packet API). It can be updated.

The `bgp_config` block supports:

* `asn` - Autonomous System Number for local BGP deployment
* `md5` - (Optional) Password for BGP session in plaintext (not a checksum)
* `deployment_type` - `private` or `public`, the `private` is likely to be usable immediately, the `public` will need to be review by Packet engineers

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the project
* `payment_method_id` - The UUID of payment method for this project. 
* `organization_id` - The UUID of this project's parent organization.
* `backend_transfer` - Whether Backend Transfer is enabled for this project.
* `created` - The timestamp for when the project was created
* `updated` - The timestamp for the last time the project was updated

The `bgp_config` block additionally exports: 

* `status` - status of BGP configuration in the project
* `max_prefix` - The maximum number of route filters allowed per server
