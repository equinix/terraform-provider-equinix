# ECX Fabric Layer2 private service profile

This example shows how create layer2 private service profile.

Example service profile will be available for use only for selected users.
It will also leverage Equinix API integration and allow to derive speed
from connection API calls.
Refer to [ECXF layer 2 service
profile resource documentation](../../docs/resources/ecx_l2_serviceprofile.md) for
more details about possible service profile options..

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret), obtained
same way as above
* `equinix_pri_port_name` - name of ECX Fabric primary port that you want to use
for connections that will use your profile
* `equinix_sec_port_name` - name of ECX Fabric secondary port that you want to use
for connections that will use your profile

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
