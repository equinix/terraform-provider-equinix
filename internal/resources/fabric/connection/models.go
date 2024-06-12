package connection

import (
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func additionalInfoContainsAWSSecrets(info []interface{}) ([]interface{}, bool) {
	var awsSecrets []interface{}

	for _, item := range info {
		if value, _ := item.(map[string]interface{})["key"]; value == "accessKey" {
			awsSecrets = append(awsSecrets, item)
		}

		if value, _ := item.(map[string]interface{})["key"]; value == "secretKey" {
			awsSecrets = append(awsSecrets, item)
		}
	}

	return awsSecrets, len(awsSecrets) == 2
}

func setFabricMap(d *schema.ResourceData, conn *fabricv4.Connection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	connection := connectionMap(conn)
	err := equinix_schema.SetMap(d, connection)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func connectionMap(conn *fabricv4.Connection) map[string]interface{} {
	connection := make(map[string]interface{})
	connection["name"] = conn.GetName()
	connection["uuid"] = conn.GetUuid()
	connection["bandwidth"] = conn.GetBandwidth()
	connection["href"] = conn.GetHref()
	connection["is_remote"] = conn.GetIsRemote()
	connection["type"] = string(conn.GetType())
	connection["state"] = string(conn.GetState())
	connection["direction"] = conn.GetDirection()
	if conn.Operation != nil {
		operation := conn.GetOperation()
		connection["operation"] = connectionOperationGoToTerraform(&operation)
	}
	if conn.Order != nil {
		order := conn.GetOrder()
		connection["order"] = equinix_fabric_schema.OrderGoToTerraform(&order)
	}
	if conn.ChangeLog != nil {
		changeLog := conn.GetChangeLog()
		connection["change_log"] = equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog)
	}
	if conn.Redundancy != nil {
		redundancy := conn.GetRedundancy()
		connection["redundancy"] = connectionRedundancyGoToTerraform(&redundancy)
	}
	if conn.Notifications != nil {
		notifications := conn.GetNotifications()
		connection["notifications"] = equinix_fabric_schema.NotificationsGoToTerraform(notifications)
	}
	if conn.Account != nil {
		account := conn.GetAccount()
		connection["account"] = equinix_fabric_schema.AccountGoToTerraform(&account)
	}
	if &conn.ASide != nil {
		aSide := conn.GetASide()
		connection["a_side"] = connectionSideGoToTerraform(&aSide)
	}
	if &conn.ZSide != nil {
		zSide := conn.GetZSide()
		connection["z_side"] = connectionSideGoToTerraform(&zSide)
	}
	if conn.AdditionalInfo != nil {
		additionalInfo := conn.GetAdditionalInfo()
		connection["additional_info"] = additionalInfoGoToTerraform(additionalInfo)
	}
	if conn.Project != nil {
		project := conn.GetProject()
		connection["project"] = equinix_fabric_schema.ProjectGoToTerraform(&project)
	}

	return connection
}

func connectionRedundancyTerraformToGo(redundancyTerraform []interface{}) fabricv4.ConnectionRedundancy {
	if redundancyTerraform == nil || len(redundancyTerraform) == 0 {
		return fabricv4.ConnectionRedundancy{}
	}
	var redundancy fabricv4.ConnectionRedundancy

	redundancyMap := redundancyTerraform[0].(map[string]interface{})
	connectionPriority := redundancyMap["priority"].(string)
	redundancyGroup := redundancyMap["group"].(string)
	redundancy.SetPriority(fabricv4.ConnectionPriority(connectionPriority))
	if redundancyGroup != "" {
		redundancy.SetGroup(redundancyGroup)
	}

	return redundancy
}

func connectionRedundancyGoToTerraform(redundancy *fabricv4.ConnectionRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	mappedRedundancy := make(map[string]interface{})
	mappedRedundancy["group"] = redundancy.GetGroup()
	mappedRedundancy["priority"] = string(redundancy.GetPriority())
	redundancySet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: connectionRedundancySch()}),
		[]interface{}{mappedRedundancy},
	)
	return redundancySet
}

func serviceTokenTerraformToGo(serviceTokenList []interface{}) fabricv4.ServiceToken {
	if serviceTokenList == nil || len(serviceTokenList) == 0 {
		return fabricv4.ServiceToken{}
	}

	var serviceToken fabricv4.ServiceToken
	serviceTokenMap := serviceTokenList[0].(map[string]interface{})
	serviceTokenType := serviceTokenMap["type"].(string)
	uuid := serviceTokenMap["uuid"].(string)
	serviceToken.SetType(fabricv4.ServiceTokenType(serviceTokenType))
	serviceToken.SetUuid(uuid)

	return serviceToken
}

func additionalInfoTerraformToGo(additionalInfoList []interface{}) []fabricv4.ConnectionSideAdditionalInfo {
	if additionalInfoList == nil || len(additionalInfoList) == 0 {
		return nil
	}

	mappedAdditionalInfoList := make([]fabricv4.ConnectionSideAdditionalInfo, len(additionalInfoList))
	for index, additionalInfo := range additionalInfoList {
		additionalInfoMap := additionalInfo.(map[string]interface{})
		key := additionalInfoMap["key"].(string)
		value := additionalInfoMap["value"].(string)

		additionalInfo := fabricv4.ConnectionSideAdditionalInfo{}
		additionalInfo.SetKey(key)
		additionalInfo.SetValue(value)
		mappedAdditionalInfoList[index] = additionalInfo
	}
	return mappedAdditionalInfoList
}

func connectionSideTerraformToGo(connectionSideTerraform []interface{}) fabricv4.ConnectionSide {
	if connectionSideTerraform == nil || len(connectionSideTerraform) == 0 {
		return fabricv4.ConnectionSide{}
	}

	var connectionSide fabricv4.ConnectionSide

	connectionSideMap := connectionSideTerraform[0].(map[string]interface{})
	accessPoint := connectionSideMap["access_point"].(*schema.Set).List()
	serviceTokenRequest := connectionSideMap["service_token"].(*schema.Set).List()
	additionalInfoRequest := connectionSideMap["additional_info"].([]interface{})
	if len(accessPoint) != 0 {
		ap := accessPointTerraformToGo(accessPoint)
		connectionSide.SetAccessPoint(ap)
	}
	if len(serviceTokenRequest) != 0 {
		serviceToken := serviceTokenTerraformToGo(serviceTokenRequest)
		connectionSide.SetServiceToken(serviceToken)
	}
	if len(additionalInfoRequest) != 0 {
		accessPointAdditionalInfo := additionalInfoTerraformToGo(additionalInfoRequest)
		connectionSide.SetAdditionalInfo(accessPointAdditionalInfo)
	}

	return connectionSide
}

func accessPointTerraformToGo(accessPointTerraform []interface{}) fabricv4.AccessPoint {
	if accessPointTerraform == nil || len(accessPointTerraform) == 0 {
		return fabricv4.AccessPoint{}
	}

	var accessPoint fabricv4.AccessPoint
	accessPointMap := accessPointTerraform[0].(map[string]interface{})
	portList := accessPointMap["port"].(*schema.Set).List()
	profileList := accessPointMap["profile"].(*schema.Set).List()
	locationList := accessPointMap["location"].(*schema.Set).List()
	virtualDeviceList := accessPointMap["virtual_device"].(*schema.Set).List()
	interfaceList := accessPointMap["interface"].(*schema.Set).List()
	networkList := accessPointMap["network"].(*schema.Set).List()
	typeVal := accessPointMap["type"].(string)
	authenticationKey := accessPointMap["authentication_key"].(string)
	if authenticationKey != "" {
		accessPoint.SetAuthenticationKey(authenticationKey)
	}
	providerConnectionId := accessPointMap["provider_connection_id"].(string)
	if providerConnectionId != "" {
		accessPoint.SetProviderConnectionId(providerConnectionId)
	}
	sellerRegion := accessPointMap["seller_region"].(string)
	if sellerRegion != "" {
		accessPoint.SetSellerRegion(sellerRegion)
	}
	peeringTypeRaw := accessPointMap["peering_type"].(string)
	if peeringTypeRaw != "" {
		peeringType := fabricv4.PeeringType(peeringTypeRaw)
		accessPoint.SetPeeringType(peeringType)
	}
	cloudRouterRequest := accessPointMap["router"].(*schema.Set).List()
	if len(cloudRouterRequest) == 0 {
		log.Print("[DEBUG] The router attribute was not used, attempting to revert to deprecated gateway attribute")
		cloudRouterRequest = accessPointMap["gateway"].(*schema.Set).List()
	}

	if len(cloudRouterRequest) != 0 {
		cloudRouter := cloudRouterTerraformToGo(cloudRouterRequest)
		if cloudRouter.GetUuid() != "" {
			accessPoint.SetRouter(cloudRouter)
		}
	}
	accessPoint.SetType(fabricv4.AccessPointType(typeVal))
	if len(portList) != 0 {
		port := portTerraformToGo(portList)
		if port.GetUuid() != "" {
			accessPoint.SetPort(port)
		}
	}

	if len(networkList) != 0 {
		network := networkTerraformToGo(networkList)
		if network.GetUuid() != "" {
			accessPoint.SetNetwork(network)
		}
	}
	linkProtocolList := accessPointMap["link_protocol"].(*schema.Set).List()

	if len(linkProtocolList) != 0 {
		linkProtocol := linkProtocolTerraformToGo(linkProtocolList)
		if linkProtocol.GetType().Ptr() != nil {
			accessPoint.SetLinkProtocol(linkProtocol)
		}
	}

	if len(profileList) != 0 {
		serviceProfile := simplifiedServiceProfileTerraformToGo(profileList)
		if serviceProfile.GetUuid() != "" {
			accessPoint.SetProfile(serviceProfile)
		}
	}

	if len(locationList) != 0 {
		location := equinix_fabric_schema.LocationTerraformToGo(locationList)
		accessPoint.SetLocation(location)
	}

	if len(virtualDeviceList) != 0 {
		virtualDevice := virtualDeviceTerraformToGo(virtualDeviceList)
		accessPoint.SetVirtualDevice(virtualDevice)
	}

	if len(interfaceList) != 0 {
		interface_ := interfaceTerraformToGo(interfaceList)
		accessPoint.SetInterface(interface_)
	}

	return accessPoint
}

func cloudRouterTerraformToGo(cloudRouterRequest []interface{}) fabricv4.CloudRouter {
	if cloudRouterRequest == nil || len(cloudRouterRequest) == 0 {
		return fabricv4.CloudRouter{}
	}
	var cloudRouter fabricv4.CloudRouter
	cloudRouterMap := cloudRouterRequest[0].(map[string]interface{})
	uuid := cloudRouterMap["uuid"].(string)
	cloudRouter.SetUuid(uuid)

	return cloudRouter
}

func linkProtocolTerraformToGo(linkProtocolList []interface{}) fabricv4.SimplifiedLinkProtocol {
	if linkProtocolList == nil || len(linkProtocolList) == 0 {
		return fabricv4.SimplifiedLinkProtocol{}
	}

	var linkProtocol fabricv4.SimplifiedLinkProtocol
	lpMap := linkProtocolList[0].(map[string]interface{})
	lpType := lpMap["type"].(string)
	lpVlanSTag := int32(lpMap["vlan_s_tag"].(int))
	lpVlanTag := int32(lpMap["vlan_tag"].(int))
	lpVlanCTag := int32(lpMap["vlan_c_tag"].(int))
	log.Printf("[DEBUG] linkProtocolMap: %v", lpMap)

	linkProtocol.SetType(fabricv4.LinkProtocolType(lpType))
	if lpVlanSTag != 0 {
		linkProtocol.SetVlanSTag(lpVlanSTag)
	}
	if lpVlanTag != 0 {
		linkProtocol.SetVlanTag(lpVlanTag)
	}
	if lpVlanCTag != 0 {
		linkProtocol.SetVlanCTag(lpVlanCTag)
	}

	return linkProtocol
}

func networkTerraformToGo(networkList []interface{}) fabricv4.SimplifiedNetwork {
	if networkList == nil || len(networkList) == 0 {
		return fabricv4.SimplifiedNetwork{}
	}
	var network fabricv4.SimplifiedNetwork
	networkListMap := networkList[0].(map[string]interface{})
	uuid := networkListMap["uuid"].(string)
	network.SetUuid(uuid)
	return network
}

func simplifiedServiceProfileTerraformToGo(profileList []interface{}) fabricv4.SimplifiedServiceProfile {
	if profileList == nil || len(profileList) == 0 {
		return fabricv4.SimplifiedServiceProfile{}
	}

	var serviceProfile fabricv4.SimplifiedServiceProfile
	profileListMap := profileList[0].(map[string]interface{})
	profileType := profileListMap["type"].(string)
	uuid := profileListMap["uuid"].(string)
	serviceProfile.SetType(fabricv4.ServiceProfileTypeEnum(profileType))
	serviceProfile.SetUuid(uuid)
	return serviceProfile
}

func virtualDeviceTerraformToGo(virtualDeviceList []interface{}) fabricv4.VirtualDevice {
	if virtualDeviceList == nil || len(virtualDeviceList) == 0 {
		return fabricv4.VirtualDevice{}
	}

	var virtualDevice fabricv4.VirtualDevice
	virtualDeviceMap := virtualDeviceList[0].(map[string]interface{})
	href := virtualDeviceMap["href"].(string)
	type_ := virtualDeviceMap["type"].(string)
	uuid := virtualDeviceMap["uuid"].(string)
	name := virtualDeviceMap["name"].(string)
	virtualDevice.SetHref(href)
	virtualDevice.SetType(fabricv4.VirtualDeviceType(type_))
	virtualDevice.SetUuid(uuid)
	virtualDevice.SetName(name)

	return virtualDevice
}

func interfaceTerraformToGo(interfaceList []interface{}) fabricv4.Interface {
	if interfaceList == nil || len(interfaceList) == 0 {
		return fabricv4.Interface{}
	}

	var interface_ fabricv4.Interface
	interfaceMap := interfaceList[0].(map[string]interface{})
	uuid := interfaceMap["uuid"].(string)
	type_ := interfaceMap["type"].(string)
	id := interfaceMap["id"].(int)
	interface_.SetUuid(uuid)
	interface_.SetType(fabricv4.InterfaceType(type_))
	interface_.SetId(int32(id))

	return interface_
}

func connectionOperationGoToTerraform(operation *fabricv4.ConnectionOperation) *schema.Set {
	if operation == nil {
		return nil
	}

	mappedOperation := make(map[string]interface{})
	mappedOperation["provider_status"] = string(operation.GetProviderStatus())
	mappedOperation["equinix_status"] = string(operation.GetEquinixStatus())
	if operation.Errors != nil {
		mappedOperation["errors"] = equinix_fabric_schema.ErrorGoToTerraform(operation.GetErrors())
	}
	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: operationSch()}),
		[]interface{}{mappedOperation},
	)
	return operationSet
}

func serviceTokenGoToTerraform(serviceToken *fabricv4.ServiceToken) *schema.Set {
	if serviceToken == nil {
		return nil
	}
	mappedServiceToken := make(map[string]interface{})
	if serviceToken.Type != nil {
		mappedServiceToken["type"] = string(serviceToken.GetType())
	}
	mappedServiceToken["href"] = serviceToken.GetHref()
	mappedServiceToken["uuid"] = serviceToken.GetUuid()

	serviceTokenSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceTokenSch()}),
		[]interface{}{mappedServiceToken},
	)
	return serviceTokenSet
}

func connectionSideGoToTerraform(connectionSide *fabricv4.ConnectionSide) *schema.Set {
	mappedConnectionSide := make(map[string]interface{})
	serviceToken := connectionSide.GetServiceToken()
	serviceTokenSet := serviceTokenGoToTerraform(&serviceToken)
	if serviceTokenSet != nil {
		mappedConnectionSide["service_token"] = serviceTokenSet
	}
	accessPoint := connectionSide.GetAccessPoint()
	mappedConnectionSide["access_point"] = accessPointGoToTerraform(&accessPoint)
	connectionSideSet := schema.NewSet(
		schema.HashResource(connectionSideSch()),
		[]interface{}{mappedConnectionSide},
	)
	return connectionSideSet
}

func additionalInfoGoToTerraform(additionalInfo []fabricv4.ConnectionSideAdditionalInfo) []map[string]interface{} {
	if additionalInfo == nil {
		return nil
	}
	mappedAdditionalInfo := make([]map[string]interface{}, len(additionalInfo))
	for index, additionalInfo := range additionalInfo {
		mappedAdditionalInfo[index] = map[string]interface{}{
			"key":   additionalInfo.GetKey(),
			"value": additionalInfo.GetValue(),
		}
	}
	return mappedAdditionalInfo
}

func cloudRouterGoToTerraform(cloudRouter *fabricv4.CloudRouter) *schema.Set {
	if cloudRouter == nil {
		return nil
	}
	mappedCloudRouter := make(map[string]interface{})
	mappedCloudRouter["uuid"] = cloudRouter.GetUuid()
	mappedCloudRouter["href"] = cloudRouter.GetHref()

	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: equinix_fabric_schema.ProjectSch()}),
		[]interface{}{mappedCloudRouter})
	return linkedProtocolSet
}

func virtualDeviceGoToTerraform(virtualDevice *fabricv4.VirtualDevice) *schema.Set {
	if virtualDevice == nil {
		return nil
	}
	mappedVirtualDevice := make(map[string]interface{})
	mappedVirtualDevice["name"] = virtualDevice.GetName()
	mappedVirtualDevice["href"] = virtualDevice.GetHref()
	mappedVirtualDevice["type"] = string(virtualDevice.GetType())
	mappedVirtualDevice["uuid"] = virtualDevice.GetUuid()

	virtualDeviceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointVirtualDeviceSch()}),
		[]interface{}{mappedVirtualDevice})
	return virtualDeviceSet
}

func interfaceGoToTerraform(mInterface *fabricv4.Interface) *schema.Set {
	if mInterface == nil {
		return nil
	}
	mappedMInterface := make(map[string]interface{})
	mappedMInterface["id"] = int(mInterface.GetId())
	mappedMInterface["type"] = string(mInterface.GetType())
	mappedMInterface["uuid"] = mInterface.GetUuid()

	mInterfaceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointInterface()}),
		[]interface{}{mappedMInterface})
	return mInterfaceSet
}

func networkGoToTerraform(network *fabricv4.SimplifiedNetwork) *schema.Set {
	if network == nil {
		return nil
	}

	mappedNetwork := make(map[string]interface{})
	mappedNetwork["uuid"] = network.GetUuid()

	return schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: networkSch()}),
		[]interface{}{mappedNetwork})
}

func accessPointGoToTerraform(accessPoint *fabricv4.AccessPoint) *schema.Set {
	mappedAccessPoint := make(map[string]interface{})
	if accessPoint.Type != nil {
		mappedAccessPoint["type"] = string(accessPoint.GetType())
	}
	if accessPoint.Account != nil {
		account := accessPoint.GetAccount()
		mappedAccessPoint["account"] = equinix_fabric_schema.AccountGoToTerraform(&account)
	}
	if accessPoint.Location != nil {
		location := accessPoint.GetLocation()
		mappedAccessPoint["location"] = equinix_fabric_schema.LocationGoToTerraform(&location)
	}
	if accessPoint.Port != nil {
		port := accessPoint.GetPort()
		mappedAccessPoint["port"] = portGoToTerraform(&port)
	}
	if accessPoint.Profile != nil {
		profile := accessPoint.GetProfile()
		mappedAccessPoint["profile"] = simplifiedServiceProfileGoToTerraform(&profile)
	}
	if accessPoint.Router != nil {
		router := accessPoint.GetRouter()
		mappedAccessPoint["router"] = cloudRouterGoToTerraform(&router)
		mappedAccessPoint["gateway"] = cloudRouterGoToTerraform(&router)
	}
	if accessPoint.LinkProtocol != nil {
		linkProtocol := accessPoint.GetLinkProtocol()
		mappedAccessPoint["link_protocol"] = linkedProtocolGoToTerraform(&linkProtocol)
	}
	if accessPoint.VirtualDevice != nil {
		virtualDevice := accessPoint.GetVirtualDevice()
		mappedAccessPoint["virtual_device"] = virtualDeviceGoToTerraform(&virtualDevice)
	}
	if accessPoint.Interface != nil {
		interface_ := accessPoint.GetInterface()
		mappedAccessPoint["interface"] = interfaceGoToTerraform(&interface_)
	}
	if accessPoint.Network != nil {
		network := accessPoint.GetNetwork()
		mappedAccessPoint["network"] = networkGoToTerraform(&network)
	}
	mappedAccessPoint["seller_region"] = accessPoint.GetSellerRegion()
	if accessPoint.PeeringType != nil {
		mappedAccessPoint["peering_type"] = string(accessPoint.GetPeeringType())
	}
	mappedAccessPoint["authentication_key"] = accessPoint.GetAuthenticationKey()
	mappedAccessPoint["provider_connection_id"] = accessPoint.GetProviderConnectionId()

	accessPointSet := schema.NewSet(
		schema.HashResource(accessPointSch()),
		[]interface{}{mappedAccessPoint},
	)
	return accessPointSet
}

func linkedProtocolGoToTerraform(linkedProtocol *fabricv4.SimplifiedLinkProtocol) *schema.Set {

	mappedLinkedProtocol := make(map[string]interface{})
	mappedLinkedProtocol["type"] = string(linkedProtocol.GetType())
	mappedLinkedProtocol["vlan_tag"] = int(linkedProtocol.GetVlanTag())
	mappedLinkedProtocol["vlan_s_tag"] = int(linkedProtocol.GetVlanSTag())
	mappedLinkedProtocol["vlan_c_tag"] = int(linkedProtocol.GetVlanCTag())

	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointLinkProtocolSch()}),
		[]interface{}{mappedLinkedProtocol})
	return linkedProtocolSet
}

func simplifiedServiceProfileGoToTerraform(profile *fabricv4.SimplifiedServiceProfile) *schema.Set {

	mappedProfile := make(map[string]interface{})
	mappedProfile["href"] = profile.GetHref()
	mappedProfile["type"] = string(profile.GetType())
	mappedProfile["name"] = profile.GetName()
	mappedProfile["uuid"] = profile.GetName()
	mappedProfile["access_point_type_configs"] = accessPointTypeConfigGoToTerraform(profile.AccessPointTypeConfigs)

	profileSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceProfileSch()}),
		[]interface{}{mappedProfile},
	)
	return profileSet
}

func accessPointTypeConfigGoToTerraform(spAccessPointTypes []fabricv4.ServiceProfileAccessPointType) []interface{} {
	mappedSpAccessPointTypes := make([]interface{}, len(spAccessPointTypes))
	for index, spAccessPointType := range spAccessPointTypes {
		spAccessPointType := spAccessPointType.GetActualInstance().(*fabricv4.ServiceProfileAccessPointTypeCOLO)
		apiConfig := spAccessPointType.GetApiConfig()
		authKey := spAccessPointType.GetAuthenticationKey()
		mappedSpAccessPointTypes[index] = map[string]interface{}{
			"type":                             string(spAccessPointType.Type),
			"uuid":                             spAccessPointType.GetUuid(),
			"allow_remote_connections":         spAccessPointType.GetAllowRemoteConnections(),
			"allow_custom_bandwidth":           spAccessPointType.GetAllowCustomBandwidth(),
			"allow_bandwidth_auto_approval":    spAccessPointType.GetAllowBandwidthAutoApproval(),
			"enable_auto_generate_service_key": spAccessPointType.GetEnableAutoGenerateServiceKey(),
			"connection_redundancy_required":   spAccessPointType.GetConnectionRedundancyRequired(),
			"connection_label":                 spAccessPointType.GetConnectionLabel(),
			"api_config":                       apiConfigGoToTerraform(&apiConfig),
			"authentication_key":               authenticationKeyGoToTerraform(&authKey),
			"supported_bandwidths":             supportedBandwidthsGoToTerraform(spAccessPointType.GetSupportedBandwidths()),
		}
	}

	return mappedSpAccessPointTypes
}

func apiConfigGoToTerraform(apiConfig *fabricv4.ApiConfig) *schema.Set {

	mappedApiConfig := make(map[string]interface{})
	mappedApiConfig["api_available"] = apiConfig.GetApiAvailable()
	mappedApiConfig["equinix_managed_vlan"] = apiConfig.GetEquinixManagedVlan()
	mappedApiConfig["bandwidth_from_api"] = apiConfig.GetBandwidthFromApi()
	mappedApiConfig["integration_id"] = apiConfig.GetIntegrationId()
	mappedApiConfig["equinix_managed_port"] = apiConfig.GetEquinixManagedPort()

	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createApiConfigSch()}),
		[]interface{}{mappedApiConfig})
	return apiConfigSet
}

func authenticationKeyGoToTerraform(authenticationKey *fabricv4.AuthenticationKey) *schema.Set {
	mappedAuthenticationKey := make(map[string]interface{})
	mappedAuthenticationKey["required"] = authenticationKey.GetRequired()
	mappedAuthenticationKey["label"] = authenticationKey.GetLabel()
	mappedAuthenticationKey["description"] = authenticationKey.GetDescription()

	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createAuthenticationKeySch()}),
		[]interface{}{mappedAuthenticationKey})
	return apiConfigSet
}

func supportedBandwidthsGoToTerraform(supportedBandwidths []int32) []interface{} {
	if supportedBandwidths == nil {
		return nil
	}
	mappedSupportedBandwidths := make([]interface{}, len(supportedBandwidths))
	for index, bandwidth := range supportedBandwidths {
		mappedSupportedBandwidths[index] = int(bandwidth)
	}
	return mappedSupportedBandwidths
}

func getUpdateRequests(conn *fabricv4.Connection, d *schema.ResourceData) ([][]fabricv4.ConnectionChangeOperation, error) {
	var changeOps [][]fabricv4.ConnectionChangeOperation
	existingName := conn.GetName()
	existingBandwidth := int(conn.GetBandwidth())
	updateNameVal := d.Get("name").(string)
	updateBandwidthVal := d.Get("bandwidth").(int)
	additionalInfo := d.Get("additional_info").([]interface{})

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfo)

	if existingName != updateNameVal {
		changeOps = append(changeOps, []fabricv4.ConnectionChangeOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: updateNameVal,
			},
		})
	}

	if existingBandwidth != updateBandwidthVal {
		changeOps = append(changeOps, []fabricv4.ConnectionChangeOperation{
			{
				Op:    "replace",
				Path:  "/bandwidth",
				Value: updateBandwidthVal,
			},
		})
	}

	if *conn.Operation.ProviderStatus == fabricv4.PROVIDERSTATUS_PENDING_APPROVAL && hasAWSSecrets {
		changeOps = append(changeOps, []fabricv4.ConnectionChangeOperation{
			{
				Op:    "add",
				Path:  "",
				Value: map[string]interface{}{"additionalInfo": awsSecrets},
			},
		})
	}

	if len(changeOps) == 0 {
		return changeOps, fmt.Errorf("nothing to update for the connection %s", existingName)
	}

	return changeOps, nil
}

func portTerraformToGo(portList []interface{}) fabricv4.SimplifiedPort {
	if portList == nil || len(portList) == 0 {
		return fabricv4.SimplifiedPort{}
	}
	var port fabricv4.SimplifiedPort
	portListMap := portList[0].(map[string]interface{})
	uuid := portListMap["uuid"].(string)
	port.SetUuid(uuid)

	return port
}

func portGoToTerraform(port *fabricv4.SimplifiedPort) *schema.Set {
	mappedPort := make(map[string]interface{})
	mappedPort["href"] = port.GetHref()
	mappedPort["name"] = port.GetName()
	mappedPort["uuid"] = port.GetUuid()
	if port.Redundancy != nil {
		mappedPort["redundancy"] = portRedundancyGoToTerraform(port.Redundancy)
	}
	portSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: portSch()}),
		[]interface{}{mappedPort},
	)
	return portSet
}

func portRedundancyGoToTerraform(redundancy *fabricv4.PortRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	mappedRedundancy := make(map[string]interface{})
	mappedRedundancy["enabled"] = redundancy.GetEnabled()
	mappedRedundancy["group"] = redundancy.GetGroup()
	mappedRedundancy["priority"] = string(redundancy.GetPriority())

	redundancySet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: portRedundancySch()}),
		[]interface{}{mappedRedundancy},
	)
	return redundancySet
}
