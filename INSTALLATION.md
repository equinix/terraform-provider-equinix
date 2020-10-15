# Equinix provider manual installation

## Determine Terraform version

Manual installation of third-party plugins process differs depending on
Terraform version.

Run `terraform version` command to determine version of your Terraform installation.

## Terraform 0.13 and newer

1. [Obtain Equinix provider binary](#Obtaining-Equinix-provider-binary)

2. Create directory structure under default Terraform plugin directory

   2.1. **Official Equinix Provider releases**

      * create `developer.equinix.com/terraform/equinix` directories under:
         * `~/.terraform.d/plugins` (Mac and Linux)
         * `%APPDATA%\terraform.d\plugins` (Windows)

      * copy Equinix provider **zip archive** there

        Example *(for Mac*):

        ```sh
        mkdir -p ~/.terraform.d/plugins/developer.equinix.com/terraform/equinix
        cp terraform-provider-equinix_1.0.0_darwin_amd64.zip ~/.terraform.d/plugins/developer.equinix.com/terraform/equinix
        ```

   2.2 **Development Equinix Provider builds**

      * create `developer.equinix.com/terraform/equinix/1.0.0/darwin_amd64` directories
      under:
         * `~/.terraform.d/plugins` (Mac and Linux)
         * `%APPDATA%\terraform.d\plugins` (Windows)

        Note: adjust `darwin_amd64` from above structure to match your *os_arch*

      * copy Equinix provider **binary file** there

        Example *(for Mac*):

        ```sh
        mkdir -p ~/.terraform.d/plugins/developer.equinix.com/terraform/equinix/1.0.0/darwin_amd64
        cp terraform-provider-equinix ~/.terraform.d/plugins/developer.equinix.com/terraform/equinix/1.0.0/darwin_amd64
        ```

3. In every Terraform template directory that uses Equinix provider, ship below
 `terraform.tf` file *(in addition to other Terraform files)*

   ```hcl
   terraform {
     required_providers {
       equinix = {
         source = "developer.equinix.com/terraform/equinix"
       }
     }
   }
   ```

4. **Done!**

   Local Equinix provider plugin will be used after `terraform init`
command execution in Terraform template directory

## Terraform 0.12 and older

1. [Obtain Equinix provider binary](#Obtaining-Equinix-provider-binary)
2. Unpack provider binary archive (if applicable)

   ```sh
   user@host Downloads % unzip terraform-provider-equinix_1.0.0_darwin_amd64.zip
   Archive:  terraform-provider-equinix_1.0.0_darwin_amd64.zip
     inflating: CHANGELOG.md
     inflating: LICENSE
     inflating: README.md
     inflating: terraform-provider-equinix_v1.0.0
   user@host Downloads %
   ```

3. Copy provider binary to Terraform plugin directory

   * On **Linux and Mac**, copy to `~/.terraform.d/plugins`

     ```sh
     cp terraform-provider-equinix_v1.0.0 ~/.terraform.d/plugins
     ```

   * On **Windows**, copy to `%APPDATA%\terraform.d\plugins`

     ```sh
     cp terraform-provider-equinix_v1.0.0.exe $APPDATA/terraform.d/plugins/
     ```

4. **Done!**

   Local Equinix provider plugin will be used after `terraform init`
command execution in selected Terraform template directory

## Obtaining Equinix provider binary

### Released versions

1. Browse to [Equinix Terraform Provider Releases](https://github.com/equinix/terraform-provider-equinix/releases)
Github page
2. Locate desired release and **download archive** for target OS and architecture

   * example for 64bit x86 Windows: `terraform-provider-equinix_1.0.0_windows_amd64.zip`
   * example for 32bit x86 Linux: `terraform-provider-equinix_1.0.0_linux_386.zip`

3. Optionally, download checksum file and verify downloaded archive

### Development version

NOTE: running development version requires you to build plugin from
the source code

**Prerequisites:**

* [Go](https://golang.org/doc/install) 1.14
* (optionally) [GNU Make](https://www.gnu.org/software/make)

**Steps**:

1. Clone Equinix Terraform Provider repository

   ```sh
   git clone https://github.com/equinix/terraform-provider-equinix.git
   ```

2. Enter provider directory and build

   * using GNU Make

     ```sh
     cd terraform-provider-equinix
     make
     ```

   * using Go only

     ```sh
     cd terraform-provider-equinix
     go build
     ```

3. **Done!** Plugin binary will be compiled as file named:

   * `terraform-provider-equinix` (on Mac and Linux)
   * `terraform-provider-equinix.exe` (on Windows)
