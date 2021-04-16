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
terraform {
  required_providers {
    metal = {
      source = "equinix/metal"
      # version = "1.0.0"
    }
  }
}

# Configure the Equinix Metal Provider.
provider "metal" {
  auth_token = var.auth_token
}

data "metal_project" "project" {
  name = "My Project"
}

# If you want to create a fresh project, you can create one with metal_project
#
# resource "metal_project" "cool_project" {
#   name           = "My First Terraform Project"
# }

# Create a device and add it to tf_project_1
resource "metal_device" "web1" {
  hostname         = "web1"
  plan             = "c3.medium.x86"
  facilities       = ["ny5"]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = data.metal_project.project.id

  # if you created a project with the metal_project resource, refer to its ID
  # project_id       = metal_project.cool_project.id

  # You can find the ID of your project in the URL of the Equinix Metal console.
  # For example, if you see your devices listed at
  # https://console.equinix.com/projects/352000fb2-ee46-4673-93a8-de2c2bdba33b
  # .. then 352000fb2-ee46-4673-93a8-de2c2bdba33b is your project ID.
}
```

## Argument Reference

The following arguments are supported:

* `auth_token` - (Required) This is your Equinix Metal API Auth token. This can
  also be specified with the `METAL_AUTH_TOKEN` environment variable.

  Use of the legacy `PACKET_AUTH_TOKEN` environment variable is deprecated.
* `max_retries` - Maximum number of retries in case of network failure.
* `max_retry_wait_seconds` - Maximum time to wait in case of network failure.
