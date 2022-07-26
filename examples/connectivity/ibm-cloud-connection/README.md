# ECX Fabric Layer2 Connection to IBM Cloud

**NOTE:** There is an
[Equinix Fabric L2 Connection To IBM Cloud Direct Link Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-ibm/equinix/latest)
available with full-fledged examples of connections from Fabric Ports, Network Edge Devices
or Service Tokens.

This example shows how create layer 2 connection between ECX Fabric port
and IBM Direct Link.

Example covers **provisioning of Equinix side** of the connection.
Please refer to [IBM's Direct Link documentation page](https://cloud.ibm.com/docs/terraform?topic=terraform-dl-gateway-resource)
for information about setting up IBM side.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
obtained same way as above
* `equinix_port_name` - name of ECX Fabric port that you want to connect
to IBM i.e. *EQUINIX_SVC-FR4-CX-PRI-01*
* `ibm_account_id` - IBM Cloud Account ID
* `ibm_region` - IBM Cloud region for a connection
* `ibm_metro_code` - Equinix metro location with IBM presence for connection's
destination. Given metro location code has to correspond with given IBM Cloud region

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
