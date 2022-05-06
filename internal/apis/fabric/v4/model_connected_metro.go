/*
 * Equinix Fabric API v4
 *
 * Equinix Fabric is an advanced software-defined interconnection solution that enables you to directly, securely and dynamically connect to distributed infrastructure and digital ecosystems on platform Equinix via a single port, Customers can use Fabric to connect to: </br> 1. Cloud Service Providers - Clouds, network and other service providers.  </br> 2. Enterprises - Other Equinix customers, vendors and partners.  </br> 3. Myself - Another customer instance deployed at Equinix. </br>
 *
 * API version: 4.2
 * Contact: api-support@equinix.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package v4

// Arrays of objects containing latency data for the specified metros
type ConnectedMetro struct {
	// The Canonical URL at which the resource resides.
	Href string `json:"href,omitempty"`
	// Code assigned to an Equinix International Business Exchange (IBX) data center in a specified metropolitan area.
	Code string `json:"code,omitempty"`
	// Average latency (in milliseconds[ms]) between two specified metros.
	AvgLatency float64 `json:"avgLatency,omitempty"`
}
