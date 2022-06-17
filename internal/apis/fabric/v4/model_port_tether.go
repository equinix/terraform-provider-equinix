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

// Port physical connection
type PortTether struct {
	// Port cross connect identifier
	CrossConnectId string `json:"crossConnectId,omitempty"`
	// Port cabinet number
	CabinetNumber string `json:"cabinetNumber,omitempty"`
	// Port system name
	SystemName string `json:"systemName,omitempty"`
	// Port patch panel
	PatchPanel string `json:"patchPanel,omitempty"`
	// Port patch panel port A
	PatchPanelPortA string `json:"patchPanelPortA,omitempty"`
	// Port patch panel port B
	PatchPanelPortB string `json:"patchPanelPortB,omitempty"`
}
