/*
 * Equinix Fabric API v4
 *
 * Equinix Fabric is an advanced software-defined interconnection solution that enables you to directly, securely and dynamically connect to distributed infrastructure and digital ecosystems on platform Equinix via a single port, Customers can use Fabric to connect to: </br> 1. Cloud Service Providers - Clouds, network and other service providers.  </br> 2. Enterprises - Other Equinix customers, vendors and partners.  </br> 3. Myself - Another customer instance deployed at Equinix. </br>
 *
 * API version: 4.3
 * Contact: api-support@equinix.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger
// ProductType : Product type
type ProductType string

// List of ProductType
const (
	VIRTUAL_CONNECTION_PRODUCT_ProductType ProductType = "VIRTUAL_CONNECTION_PRODUCT"
	IP_BLOCK_PRODUCT_ProductType ProductType = "IP_BLOCK_PRODUCT"
	VIRTUAL_PORT_PRODUCT_ProductType ProductType = "VIRTUAL_PORT_PRODUCT"
	FABRIC_GATEWAY_PRODUCT_ProductType ProductType = "FABRIC_GATEWAY_PRODUCT"
)
