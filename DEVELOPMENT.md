# Equinix provider development

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 1.0.11+ (to run tests)
* [Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)
* [GNU Make](https://www.gnu.org/software/make) (to build and test easier)

## Building the provider

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules)
making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH).
The instructions that follow assume a directory in your home directory outside of
the standard GOPATH (i.e `$HOME/development/`).

1. Clone Equinix Terraform Provider repository

   ```sh
   mkdir -p $HOME/development; cd $HOME/development 
   git clone https://github.com/equinix/terraform-provider-equinix.git
   ```

2. Enter provider directory and compile it

   ```sh
   cd terraform-provider-equinix
   make build
   ```

## Testing the provider

* to run unit tests, use make `test` target

  ```sh
  make test
  ```

* to run acceptance tests, use make `testacc` target

  ```sh
  export EQUINIX_API_ENDPOINT=https://api.equinix.com
  export EQUINIX_API_CLIENTID=someID
  export EQUINIX_API_CLIENTSECRET=someSecret
  make testacc
  ```
  
  *NOTE*: acceptance tests create resources on real infrastructure, thus may be
subject for costs.

* to manually clean resources that could leak from acceptance tests,
use `sweep` target

  ```sh
  export EQUINIX_API_ENDPOINT=https://api.equinix.com
  export EQUINIX_API_CLIENTID=someID
  export EQUINIX_API_CLIENTSECRET=someSecre
  make sweep
  ```

## Test parametrization

Acceptance tests can be parametrized by setting up various environmental variables.

Rationale behind parametrization is to allow running acceptance tests against
different Equinix test environments.

* `TF_ACC_FABRIC_PRI_PORT_NAME` alters default name of Equinix Fabric port for
primary, Dot1Q encapsulated connections. Reflected by Fabric connection tests.
* `TF_ACC_FABRIC_SEC_PORT_NAME` alters default name of Equinix Fabric port for
secondary,Dot1Q encapsulated connections. Reflected by Fabric connection tests
* `TF_ACC_FABRIC_AWS_L2_SP_NAME` alters default name of Equinix Fabric AWS seller
profile. Reflected by Fabric connection tests
* `TF_ACC_FABRIC_AWS_AUTH_KEY` alters default authentication key of Equinix Fabric
AWS l2 connection. Reflected by Fabric connection tests
* `TF_ACC_FABRIC_L2_AZURE_SP_NAME` alters default name of Equinix Fabric Azure seller
profile. Reflected by Fabric connection tests.
* `TF_ACC_FABRIC_GCP1_L2_SP_NAME` alters default name of Equinix Fabric GCP
Interconnection Zone 1 seller profile. Reflected by Fabric connection tests.
* `TF_ACC_FABRIC_GCP2_L2_SP_NAME` alters default name of Equinix Fabric GCP
Interconnection Zone 2 seller profile. Reflected by Fabric connection tests.
* `TF_ACC_NETWORK_DEVICE_METRO` alters default metro code for Network Edge resources.
Reflected by Network Edge tests.
* `TF_ACC_NETWORK_DEVICE_LICENSE_FILE` alters default path to device license file.
Reflected by Network Edge tests.

## Manual provider installation

*Note:* manual provider installation is needed only for manual testing of custom
built Equinix provider plugin.

Manual installation process differs depending on Terraform version.
Run `terraform version` command to determine version of your Terraform installation.

### Terraform 0.13 and newer

1. Create `developer.equinix.com/terraform/equinix/9.0.0/darwin_amd64` directories
under:

   * `~/.terraform.d/plugins` (Mac and Linux)
   * `%APPDATA%\terraform.d\plugins` (Windows)

   *Note:* adjust `darwin_amd64` from above structure to match your *os_arch*

   ```sh
   mkdir -p ~/.terraform.d/plugins/developer.equinix.com/terraform/equinix/9.0.0/darwin_amd64
   ```

2. Copy Equinix provider **binary file** there.

   ```sh
   cp terraform-provider-equinix ~/.terraform.d/plugins/developer.equinix.com/terraform/equinix/9.0.0/darwin_amd644
   ```

3. In every Terraform template directory that uses Equinix provider, ship below
 `terraform.tf` file *(in addition to other Terraform files)*

   ```hcl
   terraform {
     required_providers {
       equinix = {
         source = "developer.equinix.com/terraform/equinix"
         version = "9.0.0"
       }
     }
   }
   ```

4. **Done!**

   Local Equinix provider plugin will be used after `terraform init`
   command execution in Terraform template directory

### Terraform 0.12 and older

1. Copy provider binary to Terraform plugin directory

   * On **Linux and Mac**, copy to `~/.terraform.d/plugins`

     ```sh
     cp terraform-provider-equinix_v1.0.0 ~/.terraform.d/plugins
     ```

   * On **Windows**, copy to `%APPDATA%\terraform.d\plugins`

     ```sh
     cp terraform-provider-equinix_v1.0.0.exe $APPDATA/terraform.d/plugins/
     ```

2. **Done!**

   Local Equinix provider plugin will be used after `terraform init`
   command execution in selected Terraform template directory
