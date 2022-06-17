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

// NexusConnectionType : Connection type
type NexusConnectionType string

// List of NexusConnectionType
const (
	EVPL_VC_NexusConnectionType       NexusConnectionType = "EVPL_VC"
	EPL_VC_NexusConnectionType        NexusConnectionType = "EPL_VC"
	EC_VC_NexusConnectionType         NexusConnectionType = "EC_VC"
	GW_VC_NexusConnectionType         NexusConnectionType = "GW_VC"
	ACCESS_EPL_VC_NexusConnectionType NexusConnectionType = "ACCESS_EPL_VC"
	VD_CHAIN_VC_NexusConnectionType   NexusConnectionType = "VD_CHAIN_VC"
)
