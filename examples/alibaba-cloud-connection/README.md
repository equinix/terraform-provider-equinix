# ECX Fabric Layer2 Connection to Alibaba Cloud

This example shows how create layer 2 connection between ECX Fabric port
and Alibaba Express Connect.

Example covers **provisioning of Equinix side** of the connection.
Please refer to [Alibaba Cloud Express Connect documentation page](https://www.alibabacloud.com/products/express-connect)
for information about setting up Alibaba side.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret), obtained
same way as above
* `equinix_port_name` - name of ECX Fabric port that you want to connect to
Alibaba i.e. *EQUINIX_SVC-FR4-CX-PRI-01*
* `alibaba_account_id` - Alibaba Cloud Account ID

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
