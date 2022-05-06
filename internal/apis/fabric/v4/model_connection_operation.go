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

import (
	"time"
)

// Connection type-specific operational data
type ConnectionOperation struct {
	ProviderStatus *ProviderStatus `json:"providerStatus,omitempty"`
	EquinixStatus  *EquinixStatus  `json:"equinixStatus,omitempty"`
	// Connection operational status
	OperationalStatus string       `json:"operationalStatus,omitempty"`
	Errors            []ModelError `json:"errors,omitempty"`
	// When connection transitioned into current operational status
	OpStatusChangedAt time.Time `json:"opStatusChangedAt,omitempty"`
}
