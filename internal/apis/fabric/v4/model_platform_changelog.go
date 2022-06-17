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

import (
	"time"
)

// Change log
type PlatformChangelog struct {
	// Created by User Key
	CreatedBy string `json:"createdBy,omitempty"`
	// Created by Date and Time
	CreatedDateTime time.Time `json:"createdDateTime,omitempty"`
	// Updated by User Key
	UpdatedBy string `json:"updatedBy,omitempty"`
	// Updated by Date and Time
	UpdatedDateTime time.Time `json:"updatedDateTime,omitempty"`
	// Deleted by User Key
	DeletedBy string `json:"deletedBy,omitempty"`
	// Deleted by Date and Time
	DeletedDateTime time.Time `json:"deletedDateTime,omitempty"`
}
