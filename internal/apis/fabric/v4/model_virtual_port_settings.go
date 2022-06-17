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

// Virtual Port configuration settings
type VirtualPortSettings struct {
	PortType string `json:"portType"`
	// Product name
	Product string `json:"product,omitempty"`
	Buyout bool `json:"buyout,omitempty"`
	ViewPortPermission bool `json:"viewPortPermission,omitempty"`
	PlaceVcOrderPermission bool `json:"placeVcOrderPermission,omitempty"`
	Layer3Enabled bool `json:"layer3Enabled,omitempty"`
	ProductCode string `json:"productCode,omitempty"`
}
