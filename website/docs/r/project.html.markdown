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

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Project on Packet.net
* `payment_method_id` - The UUID of payment method for this project. If you keep it empty, Packet API will pick your default Payment Method.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the project
* `payment_method_id` - The UUID of payment method for this project.
* `created` - The timestamp for when the Project was created
* `updated` - The timestamp for the last time the Project was updated
