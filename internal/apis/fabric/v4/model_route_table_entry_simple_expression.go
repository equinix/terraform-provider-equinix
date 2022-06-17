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

type RouteTableEntrySimpleExpression struct {
	// Possible field names to use on filters:  * `/type` - Route table entry type  * `/prefix` - Route table entry type  * `/nextHop` - Route table entry type  * `/state` - Route table entry type  * `/_*` - all-category search 
	Property string `json:"property,omitempty"`
	// Possible operators to use on filters:  * `=` - equal  * `!=` - not equal  * `>` - greater than  * `>=` - greater than or equal to  * `<` - less than  * `<=` - less than or equal to  * `[NOT] BETWEEN` - (not) between  * `[NOT] LIKE` - (not) like  * `[NOT] IN` - (not) in  * `~*` - case-insensitive like 
	Operator string `json:"operator,omitempty"`
	Values []string `json:"values,omitempty"`
}
