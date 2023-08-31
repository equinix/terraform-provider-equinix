# ECX Fabric Layer2 Connection to AWS

**NOTE:** There is an
[Equinix Fabric L2 Connection To AWS Direct Connect Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-aws/equinix/latest)
available with full-fledged examples of connections from Fabric Ports, Network Edge Devices
or Service Tokens.

This example shows how create layer 2 connection between ECX Fabric port
and AWS Direct Connect, including creation of Direct Connect private
virtual interface.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
obtained same way as above
* `equinix_port_name` - name of ECX Fabric port that you want to connect
to AWS i.e. *EQUINIX_SVC-FR4-CX-PRI-01*
* `aws_account_id` - AWS account identifier
* `aws_access_key` - AWS access key
* `aws_secret_key` - AWS secret key
* `aws_region` - AWS region
* `aws_metro_code` - Equinix metro location with AWS presence for connection's destination.
Given metro location has to correspond with given AWS region

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
