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

// Create connection post request
type NexusConnectionPostRequest struct {
	Type_ *NexusConnectionType `json:"type"`
	// Customer-provided connection name
	Name string `json:"name,omitempty"`
	// Connection bandwidth in Mbps
	Bandwidth  int32                      `json:"bandwidth"`
	Redundancy *NexusConnectionRedundancy `json:"redundancy,omitempty"`
	ASide      *NexusConnectionSide       `json:"aSide"`
	ZSide      *NexusConnectionSide       `json:"zSide"`
	// connection routing protocols configuration
	RoutingProtocols []RoutingProtocol `json:"routingProtocols,omitempty"`
}
