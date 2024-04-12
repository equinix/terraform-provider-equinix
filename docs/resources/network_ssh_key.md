---
subcategory: "Network Edge"
---

# equinix_network_ssh_key (Resource)

Resource `equinix_network_ssh_key` allows creation and management of Equinix Network Edge SSH keys.

## Example Usage

```terraform
resource "equinix_network_ssh_key" "john" {
  name       = "johnKent"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDpXGdxljAyPp9vH97436U171cX
  2gRkfPnpL8ebrk7ZBeeIpdjtd8mYpXf6fOI0o91TQXZTYtjABzeRgg6/m9hsMOnTHjzWpFyuj/hiPu
  iie1WtT4NffSH1ALQFX/azouBLmdNiYFMLfEVPZleergAqsYOHGCiQuR6Qh5j0yc5Wx+LKxiRZyjsS
  qo+EB8V6xBXi2i5PDJXK+dYG8YU9vdNeQdB84HvTWcGEnLR5w7pgC74pBVwzs3oWLy+3jWS0TKKtfl
  mryeFRufXq87gEkC1MOWX88uQgjyCsemuhPdN++2WS57gu7vcqCMwMDZa7dukRS3JANBtbs7qQhp9N
  w2PB4q6tohqUnSDxNjCqcoGeMNg/0kHeZcoVuznsjOrIDt0HgUApflkbtw1DP7Epfc2MJ0anf5GizM
  8UjMYiXEvv2U/qu8Vb7d5bxAshXM5nh67NSrgst9YzSSodjUCnFQkniz6KLrTkX6c2y2gJ5c9tWhg5
  SPkAc8OqLrmIwf5jGoHGh6eUJy7AtMcwE3iUpbrLw8EEoZDoDXkzh+RbOtSNKXWV4EAXsIhjQusCOW
  WQnuAHCy9N4Td0Sntzu/xhCZ8xN0oO67Cqlsk98xSRLXeg21PuuhOYJw0DLF6L68zU2OO0RzqoNq/F
  jIsltSUJPAIfYKL0yEefeNWOXSrasI1ezw== John.Kent@company.com"
  type       = "RSA"
  project_id = "a86d7112-d740-4758-9c9c-31e66373746b"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of SSH key used for identification.
* `public_key` - (Required) The SSH public key. If this is a file, it can be read using the file interpolation function.
* `type` - (Optional) The type of SSH key: `RSA` (default) or `DSA`.
* `project_id` - (Optional) Unique Identifier for the project resource where the SSH key is scoped to.If you leave it out, the ssh key will be created under the default project id of your organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - The unique identifier of the key

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_ssh_key.example {existing_id}
```
