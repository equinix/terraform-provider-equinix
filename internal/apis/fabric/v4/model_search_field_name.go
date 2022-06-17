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

// SearchFieldName : Possible field names to use on filters
type SearchFieldName string

// List of SearchFieldName
const (
	IS_REMOTE_SearchFieldName                                 SearchFieldName = "/isRemote"
	NAME_SearchFieldName                                      SearchFieldName = "/name"
	UUID_SearchFieldName                                      SearchFieldName = "/uuid"
	ACCOUNTORG_ID_SearchFieldName                             SearchFieldName = "/account/orgId"
	A_SIDEACCESS_POINTACCOUNTACCOUNT_NAME_SearchFieldName     SearchFieldName = "/aSide/accessPoint/account/accountName"
	A_SIDEACCESS_POINTACCOUNTACCOUNT_NUMBER_SearchFieldName   SearchFieldName = "/aSide/accessPoint/account/accountNumber"
	A_SIDEACCESS_POINTGATEWAYUUID_SearchFieldName             SearchFieldName = "/aSide/accessPoint/gateway/uuid"
	A_SIDEACCESS_POINTLINK_PROTOCOLVLAN_C_TAG_SearchFieldName SearchFieldName = "/aSide/accessPoint/linkProtocol/vlanCTag"
	A_SIDEACCESS_POINTLINK_PROTOCOLVLAN_S_TAG_SearchFieldName SearchFieldName = "/aSide/accessPoint/linkProtocol/vlanSTag"
	A_SIDEACCESS_POINTLOCATIONMETRO_CODE_SearchFieldName      SearchFieldName = "/aSide/accessPoint/location/metroCode"
	A_SIDEACCESS_POINTLOCATIONMETRO_NAME_SearchFieldName      SearchFieldName = "/aSide/accessPoint/location/metroName"
	A_SIDEACCESS_POINTNAME_SearchFieldName                    SearchFieldName = "/aSide/accessPoint/name"
	A_SIDEACCESS_POINTPORTUUID_SearchFieldName                SearchFieldName = "/aSide/accessPoint/port/uuid"
	A_SIDEACCESS_POINTPORTNAME_SearchFieldName                SearchFieldName = "/aSide/accessPoint/port/name"
	A_SIDEACCESS_POINTTYPE_SearchFieldName                    SearchFieldName = "/aSide/accessPoint/type"
	A_SIDEACCESS_POINTVIRTUAL_DEVICENAME_SearchFieldName      SearchFieldName = "/aSide/accessPoint/virtualDevice/name"
	A_SIDEACCESS_POINTVIRTUAL_DEVICEUUID_SearchFieldName      SearchFieldName = "/aSide/accessPoint/virtualDevice/uuid"
	A_SIDESERVICE_TOKENUUID_SearchFieldName                   SearchFieldName = "/aSide/serviceToken/uuid"
	CHANGESTATUS_SearchFieldName                              SearchFieldName = "/change/status"
	OPERATIONEQUINIX_STATUS_SearchFieldName                   SearchFieldName = "/operation/equinixStatus"
	OPERATIONPROVIDER_STATUS_SearchFieldName                  SearchFieldName = "/operation/providerStatus"
	PROJECTPROJECT_ID_SearchFieldName                         SearchFieldName = "/project/projectId"
	REDUNDANCYGROUP_SearchFieldName                           SearchFieldName = "/redundancy/group"
	REDUNDANCYPRIORITY_SearchFieldName                        SearchFieldName = "/redundancy/priority"
	Z_SIDEACCESS_POINTACCOUNTACCOUNT_NAME_SearchFieldName     SearchFieldName = "/zSide/accessPoint/account/accountName"
	Z_SIDEACCESS_POINTAUTHENTICATION_KEY_SearchFieldName      SearchFieldName = "/zSide/accessPoint/authenticationKey"
	Z_SIDEACCESS_POINTLINK_PROTOCOLVLAN_C_TAG_SearchFieldName SearchFieldName = "/zSide/accessPoint/linkProtocol/vlanCTag"
	Z_SIDEACCESS_POINTLINK_PROTOCOLVLAN_S_TAG_SearchFieldName SearchFieldName = "/zSide/accessPoint/linkProtocol/vlanSTag"
	Z_SIDEACCESS_POINTLOCATIONMETRO_CODE_SearchFieldName      SearchFieldName = "/zSide/accessPoint/location/metroCode"
	Z_SIDEACCESS_POINTLOCATIONMETRO_NAME_SearchFieldName      SearchFieldName = "/zSide/accessPoint/location/metroName"
	Z_SIDEACCESS_POINTNAME_SearchFieldName                    SearchFieldName = "/zSide/accessPoint/name"
	Z_SIDEACCESS_POINTPORTUUID_SearchFieldName                SearchFieldName = "/zSide/accessPoint/port/uuid"
	Z_SIDEACCESS_POINTPORTNAME_SearchFieldName                SearchFieldName = "/zSide/accessPoint/port/name"
	Z_SIDEACCESS_POINTPROFILEUUID_SearchFieldName             SearchFieldName = "/zSide/accessPoint/profile/uuid"
	Z_SIDEACCESS_POINTTYPE_SearchFieldName                    SearchFieldName = "/zSide/accessPoint/type"
	Z_SIDEACCESS_POINTVIRTUAL_DEVICENAME_SearchFieldName      SearchFieldName = "/zSide/accessPoint/virtualDevice/name"
	Z_SIDEACCESS_POINTVIRTUAL_DEVICEUUID_SearchFieldName      SearchFieldName = "/zSide/accessPoint/virtualDevice/uuid"
	Z_SIDESERVICE_TOKENUUID_SearchFieldName                   SearchFieldName = "/zSide/serviceToken/uuid"
	STAR_SearchFieldName                                      SearchFieldName = "*"
)
