# ECX Fabric Layer2 Connection between two own ports

This example shows how create layer 2 connection between two, own ECX Fabric ports.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
obtained same way as above
* `equinix_aside_port_name` - name of ECX Fabric port on a-side
* `equinix_zside_port_name` - name of ECX Fabric port on z-side

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
