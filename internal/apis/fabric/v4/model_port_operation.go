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
import (
	"time"
)

// Operational specifications for ports.
type PortOperation struct {
	// Availability of a given physical port.
	OperationalStatus string `json:"operationalStatus,omitempty"`
	// Total number of connections.
	ConnectionCount int32 `json:"connectionCount,omitempty"`
	// Date and time at which port availability changed.
	OpStatusChangedAt time.Time `json:"opStatusChangedAt,omitempty"`
}
