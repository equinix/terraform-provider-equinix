---
layout: "equinix"
page_title: "Equinix: equinix_network_ssh_key"
subcategory: ""
description: |-
 Provides Network Edge SSH key resource
---

# Resource: equinix_network_ssh_key

Resource `equinix_network_ssh_key` allows creation and management of Network Edge
SSH keys.

## Example Usage

```hcl
resource "equinix_network_ssh_key" "john" {
  name       = "johnKent"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDpXGdxljAyPp9vH97436U171cX
  2gRkfPnpL8ebrk7ZBeeIpdjtd8mYpXf6fOI0o91TQXZTYtjABzeRgg6/m9hsMOnTHjzWpFyuj/hiPu
  iie1WtT4NffSH1ALQFX//zouBLmdNiYFMLfEVPZleergAqsYOHGCiQuR6Qh5j0yc5Wx+LKxiRZyjsS
  qo+EB8V6xBXi2i5PDJXK+dYG8YU9vdNeQdB84HvTWcGEnLR5w7pgC74pBVwzs3oWLy+3jWS0TKKtfl
  mryeFRufXq87gEkC1MOWX88uQgjyCsemuhPdN++2WS57gu7vcqCMwMDZa7dukRS3JANBtbs7qQhp9N
  w2PB4q6tohqUnSDxNjCqcoGeMNg/0kHeZcoVuznsjOrIDt0HgUApflkbtw1DP7Epfc2MJ0anf5GizM
  8UjMYiXEvv2U/qu8Vb7d5bxAshXM5nh67NSrgst9YzSSodjUCnFQkniz6KLrTkX6c2y2gJ5c9tWhg5
  SPkAc8OqLrmIwf5jGoHGh6eUJy7AtMcwE3iUpbrLw8EEoZDoDXkzh+RbOtSNKXWV4EAXsIhjQusCOW
  WQnuAHCy9N4Td0Sntzu/xhCZ8xN0oO67Cqlsk98xSRLXeg21PuuhOYJw0DLF6L68zU2OO0RzqoNq/F
  jIsltSUJPAIfYKL0yEefeNWOXSrasI1ezw== John.Kent@company.com"
}
```

## Argument Reference

* `name` - (Required) The name of SSH key used for identification
* `public_key` - (Required) The public key. If this is a file, it can be read
using the file interpolation function

## Attributes Reference

* `uuid` - The universally unique identifier of the key
