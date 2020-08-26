# ECX Fabric Layer2 Connection to Microsoft Azure

This example shows how create redundant layer 2 connection between ECX Fabric port and Microsoft Azure ExpressRoute.
Example covers **provisioning of both sides** of the connection.

## Adjust variables
At minimum, you must set below variables in `terrafrom.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after registring app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret), obtained same way as above
* `equinix_pri_port_name` - name of ECX Fabric primary port that you want to connect to Azure i.e. *EQUINIX_SVC-FR4-CX-PRI-01*
* `equinix_sec_port_name` - name of ECX Fabric secondary port that you want to connect to Azure i.e. *EQUINIX_SVC-FR4-CX-SEC-01*

##  Azure login
Log in to Azure using CLI and use account that has permission to create the necessary resources.

Refer to [this guide](https://www.terraform.io/docs/providers/azurerm/guides/azure_cli.html) from Azure provider documentation for help.

## Initialize
Change directory to example directory and initialize Terraform plugins by running `terrafrom init`.

## Deploy template
Apply changes by running `terrafrom apply`, then **inspect proposed plan** and approve it.
