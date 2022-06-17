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

type NexusConnectionRoutingProtocol struct {
	// Routing protocol type
	Type_ string `json:"type,omitempty"`
	// Routing protocol identifier
	Uuid string `json:"uuid,omitempty"`
	// Customer asn
	CustomerAsn int32 `json:"customerAsn,omitempty"`
	// Peer asn
	PeerAsn int32 `json:"peerAsn,omitempty"`
	// BGP authorization key
	BgpAuthKey string          `json:"bgpAuthKey,omitempty"`
	Ipv4       *ConnectionIpv4 `json:"ipv4,omitempty"`
	// Route filters values
	RouteFilters []string `json:"routeFilters,omitempty"`
}
