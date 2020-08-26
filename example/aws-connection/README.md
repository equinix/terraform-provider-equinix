# ECX Fabric Layer2 Connection to AWS

This example shows how create layer 2 connection between ECX Fabric port and AWS Direct Connect.

## Adjust variables
At minimum, you must set below variables in `terrafrom.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after registring app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret), obtained same way as above
* `equinix_port_name` - name of ECX Fabric port that you want to connect to AWS i.e. *EQUINIX_SVC-FR4-CX-PRI-01*
* `aws_account_id` - AWS account identifier
* `aws_access_key` - AWS access key *(used for accepting connection on AWS side)*
* `aws_secret_key` - AWS secret key *(used as above)*

## Initialize
Change directory to example directory and initialize Terraform plugins by running `terrafrom init`.

## Deploy template
Apply changes by running `terrafrom apply`, then **inspect proposed plan** and approve it.