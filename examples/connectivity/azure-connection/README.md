# ECX Fabric Layer2 Connection to Microsoft Azure

**NOTE:** There is an
[Equinix Fabric L2 Connection To Microsoft Azure ExpressRoute Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-azure/equinix/latest)
available with full-fledged examples of connections from Fabric Ports, Network Edge Devices
or Service Tokens.

This example shows how create redundant layer 2 connection between ECX Fabric port
and Microsoft Azure ExpressRoute.

Example covers **provisioning of both sides** of the connection.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
obtained same way as above
* `equinix_pri_port_name` - name of ECX Fabric primary port that you want
to connect to Azure i.e. *EQUINIX_SVC-FR4-CX-PRI-01*
* `equinix_sec_port_name` - name of ECX Fabric secondary port that you want to
connect to Azure i.e. *EQUINIX_SVC-FR4-CX-SEC-01*
* `azure_location` - The name of the Azure resource group location
* `azure_peering_location` - The name of the ExpressRoute peering location
* `azure_metro_code` - Equinix metro location with Azure presence for connection's
destination. Given metro location has to correspond with given Azure peering location

## Azure login

Log in to Azure using CLI and use account that has permission to create
necessary resources.

Refer to [this guide](https://www.terraform.io/docs/providers/azurerm/guides/azure_cli.html)
from Azure provider documentation for help.

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
