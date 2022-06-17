/*
 * Equinix Fabric API v4
 *
 * Equinix Fabric is an advanced software-defined interconnection solution that enables you to directly, securely and dynamically connect to distributed infrastructure and digital ecosystems on platform Equinix via a single port, Customers can use Fabric to connect to: </br> 1. Cloud Service Providers - Clouds, network and other service providers.  </br> 2. Enterprises - Other Equinix customers, vendors and partners.  </br> 3. Myself - Another customer instance deployed at Equinix. </br>
 *
 * API version: 4.2.25
 * Contact: api-support@equinix.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package v4

// Virtual Device AccessPoint Information
type VirtualDevice struct {
	// Virtual Device URI
	Href string `json:"href,omitempty"`
	// Equinix-assigned Virtual Device identifier
	Uuid string `json:"uuid,omitempty"`
	// Customer-assigned Virtual Device name
	Name string `json:"name,omitempty"`
	// Virtual Device type
	Type_ string `json:"type,omitempty"`
}
