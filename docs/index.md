---
layout: ""
page_title: "Provider: Equinix Metal"
description: |-
  The Equinix Metal provider is used to interact with the Equinix Metal Host API.
---

# Equinix Metal Provider

[Packet is now Equinix Metal!](https://blog.equinix.com/blog/2020/10/06/equinix-metal-metal-and-more/)

The Equinix Metal (`metal`) provider is used to interact with the resources supported by [Equinix Metal](https://metal.equinix.com/).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Equinix Metal Provider.
provider "metal" {
  auth_token = var.auth_token
}

# Declare your project ID
#
# You can find ID of your project form the URL in the Equinix Metal web app.
# For example, if you see your devices listed at
# https://console.equinix.com/projects/352000fb2-ee46-4673-93a8-de2c2bdba33b
# .. then 352000fb2-ee46-4673-93a8-de2c2bdba33b is your project ID.
locals {
  project_id = "<UUID_of_your_project>"
}

# If you want to create a fresh project, you can create one with metal_project
#
# resource "metal_project" "cool_project" {
#   name           = "My First Terraform Project"
# }

# Create a device and add it to tf_project_1
resource "metal_device" "web1" {
  hostname         = "web1"
  plan             = "c1.small.x86"
  facilities       = ["ewr1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id

  # if you have created project with metal_project resource, refer to its ID
  # project_id       = metal_project.cool_project.id
}
```

## Argument Reference

The following arguments are supported:

* `auth_token` - (Required) This is your Equinix Metal API Auth token. This can also be specified
  with the `PACKET_AUTH_TOKEN` shell environment variable.
