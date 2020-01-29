---
layout: "packet"
page_title: "Provider: Packet"
sidebar_current: "docs-packet-index"
description: |-
  The Packet provider is used to interact with the resources supported by Packet. The provider needs to be configured with the proper credentials before it can be used.
---

# Packet Provider

The Packet provider is used to interact with the resources supported by Packet.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

Be cautious when using the `packet_project` resource. Packet is invoicing per project, so creating man new projects will cause your Packet bill to fragment. If you want to keep your Packet bill simple, please re-use your existing projects.

## Example Usage

```hcl
# Configure the Packet Provider. 
provider "packet" {
  auth_token = var.auth_token
}

# Declare your project ID
#
# You can find ID of your project form the URL in the Packet web app.
# For example, if you see your devices listed at
# https://app.packet.net/projects/352000fb2-ee46-4673-93a8-de2c2bdba33b
# .. then 352000fb2-ee46-4673-93a8-de2c2bdba33b is your project ID.
locals {
  project_id = "<UUID_of_your_project>"
}

# If you want to create a fresh project, you can create one with packet_project
# 
# resource "packet_project" "cool_project" {
#   name           = "My First Terraform Project"
# }

# Create a device and add it to tf_project_1
resource "packet_device" "web1" {
  hostname         = "web1"
  plan             = "c1.small.x86"
  facilities         = ["ewr1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id

  # if you have created project with packet_project resource, refer to its ID
  # project_id       = packet_project.cool_project.id
}
```

## Argument Reference

The following arguments are supported:

* `auth_token` - (Required) This is your Packet API Auth token. This can also be specified
  with the `PACKET_AUTH_TOKEN` shell environment variable.
