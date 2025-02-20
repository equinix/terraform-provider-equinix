package service_token

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildCreateRequest(d *schema.ResourceData) fabricv4.ServiceToken {
	serviceTokenRequest := fabricv4.ServiceToken{}

	typeConfig := d.Get("type").(string)
	serviceTokenRequest.SetType(fabricv4.ServiceTokenType(typeConfig))

	expirationDateTimeConfig := d.Get("expiration_date_time").(string)
	const TimeFormat = "2006-01-02T15:04:05.000Z"
	expirationTime, err := time.Parse(TimeFormat, expirationDateTimeConfig)
	if err != nil {
		fmt.Print("Error Parsing expiration date time: ", err)
	}
	serviceTokenRequest.SetExpirationDateTime(expirationTime)

	descriptionConfig := d.Get("description").(string)
	if descriptionConfig != "" {
		serviceTokenRequest.SetDescription(descriptionConfig)
	}

	connectionConfig := d.Get("service_token_connection").(*schema.Set).List()
	connection := connectionTerraformToGo(connectionConfig)
	serviceTokenRequest.SetConnection(connection)

	notificationsConfig := d.Get("notifications").(*schema.Set).List()
	notifications := equinix_fabric_schema.NotificationsTerraformToGo(notificationsConfig)
	serviceTokenRequest.SetNotifications(notifications)

	return serviceTokenRequest

}

func buildUpdateRequest(d *schema.ResourceData) ([][]fabricv4.ServiceTokenChangeOperation, error) {
	patches := make([][]fabricv4.ServiceTokenChangeOperation, 0)
	oldName, newName := d.GetChange("name")
	if oldName.(string) != newName.(string) {
		patches = append(patches, []fabricv4.ServiceTokenChangeOperation{{
			Op:    "replace",
			Path:  "/name",
			Value: newName.(string),
		}})
	}

	oldDescription, newDescription := d.GetChange("description")
	if oldDescription.(string) != newDescription.(string) {
		patches = append(patches, []fabricv4.ServiceTokenChangeOperation{{
			Op:    "replace",
			Path:  "/description",
			Value: newDescription.(string),
		}})
	}

	oldExpirationDate, newExpirationDate := d.GetChange("expiration_date_time")
	if oldExpirationDate.(string) != newExpirationDate.(string) {
		patches = append(patches, []fabricv4.ServiceTokenChangeOperation{{
			Op:    "replace",
			Path:  "/expirationDateTime",
			Value: newExpirationDate.(string),
		}})
	}

	oldNotifications, newNotifications := d.GetChange("notifications")

	var oldNotificationEmails, newNotificationEmails []string

	if oldNotifications != nil {
		for _, notification := range oldNotifications.(*schema.Set).List() {
			notificationMap := notification.(map[string]interface{})

			if emails, ok := notificationMap["emails"]; ok {
				oldEmailInterface := emails.([]interface{})
				if len(oldEmailInterface) > 0 {
					oldNotificationEmails = converters.IfArrToStringArr(oldEmailInterface)
				}
			}
		}
	}
	if newNotifications != nil {
		for _, notification := range newNotifications.(*schema.Set).List() {
			notificationMap := notification.(map[string]interface{})

			if emails, ok := notificationMap["emails"]; ok {
				newEmailInterface := emails.([]interface{})
				if len(newEmailInterface) > 0 {
					newNotificationEmails = converters.IfArrToStringArr(newEmailInterface)
				}
			}
		}
	}

	if !reflect.DeepEqual(oldNotificationEmails, newNotificationEmails) {
		patches = append(patches, []fabricv4.ServiceTokenChangeOperation{{
			Op:    "replace",
			Path:  "/notifications/emails",
			Value: newNotificationEmails,
		}})
	}

	oldServiceTokenConnection, newServiceTokenConnection := d.GetChange("service_token_connection")

	var oldAsideBandwidthLimit, newAsideBandwidthLimit int

	if oldServiceTokenConnection != nil {
		for _, connection := range oldServiceTokenConnection.(*schema.Set).List() {
			oldBandwidthLimitMap := connection.(map[string]interface{})

			if bandwidth, ok := oldBandwidthLimitMap["bandwidth_limit"]; ok {
				if bandwidthLimitValue, ok := bandwidth.(int); ok {
					oldAsideBandwidthLimit = bandwidthLimitValue
				}
			}
		}
	}

	if newServiceTokenConnection != nil {
		for _, connection := range newServiceTokenConnection.(*schema.Set).List() {
			newBandwidthLimitMap := connection.(map[string]interface{})

			if bandwidth, ok := newBandwidthLimitMap["bandwidth_limit"]; ok {
				if bandwidthLimitValue, ok := bandwidth.(int); ok {
					newAsideBandwidthLimit = bandwidthLimitValue
				}
			}
		}
	}

	if oldAsideBandwidthLimit != newAsideBandwidthLimit {
		patches = append(patches, []fabricv4.ServiceTokenChangeOperation{{
			Op:    "replace",
			Path:  "/connection/bandwidthLimit",
			Value: newAsideBandwidthLimit,
		}})
	}

	var oldZsideBandwidth, newZsideBandwidth []int

	if oldServiceTokenConnection != nil {
		for _, connection := range oldServiceTokenConnection.(*schema.Set).List() {
			olSupportedBandwidthMap := connection.(map[string]interface{})

			if bandwidth, ok := olSupportedBandwidthMap["supported_bandwidths"]; ok {
				oldSupportedBandwidth := bandwidth.([]interface{})
				if len(oldSupportedBandwidth) > 0 {
					oldZsideBandwidth = converters.IfArrToIntArr(oldSupportedBandwidth)
				}
			}
		}
	}

	if newServiceTokenConnection != nil {
		for _, connection := range newServiceTokenConnection.(*schema.Set).List() {
			newSupportedBandwidthMap := connection.(map[string]interface{})

			if bandwidth, ok := newSupportedBandwidthMap["supported_bandwidths"]; ok {
				newSupportedBandwidth := bandwidth.([]interface{})
				if len(newSupportedBandwidth) > 0 {
					newZsideBandwidth = converters.IfArrToIntArr(newSupportedBandwidth)

				}
			}
		}
	}

	if !areSlicesEqual(oldZsideBandwidth, newZsideBandwidth) {
		patches = append(patches, []fabricv4.ServiceTokenChangeOperation{{
			Op:    "replace",
			Path:  "/connection/supportedBandwidths",
			Value: newZsideBandwidth,
		}})
	}

	return patches, nil
}

func areSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Ints(a)
	sort.Ints(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func buildSearchRequest(d *schema.ResourceData) fabricv4.ServiceTokenSearchRequest {
	searchRequest := fabricv4.ServiceTokenSearchRequest{}

	schemaFilters := d.Get("filter").([]interface{})
	filter := filtersTerraformToGo(schemaFilters)
	searchRequest.SetFilter(filter)

	if schemaPagination, ok := d.GetOk("pagination"); ok {
		pagination := paginationTerraformToGo(schemaPagination.(*schema.Set).List())
		searchRequest.SetPagination(pagination)
	}

	return searchRequest
}
func setServiceTokenMap(d *schema.ResourceData, serviceToken *fabricv4.ServiceToken) diag.Diagnostics {
	diags := diag.Diagnostics{}
	serviceTokenMap := serviceTokenResponseMap(serviceToken)
	err := equinix_schema.SetMap(d, serviceTokenMap)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setServiceTokensData(d *schema.ResourceData, routeFilters *fabricv4.ServiceTokens) diag.Diagnostics {
	diags := diag.Diagnostics{}
	mappedRouteFilters := make([]map[string]interface{}, len(routeFilters.Data))
	pagination := routeFilters.GetPagination()
	if routeFilters.Data != nil {
		for index, routeFilter := range routeFilters.Data {
			mappedRouteFilters[index] = serviceTokenResponseMap(&routeFilter)
		}
	} else {
		mappedRouteFilters = nil
	}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"data":       mappedRouteFilters,
		"pagination": paginationGoToTerraform(&pagination),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func serviceTokenResponseMap(token *fabricv4.ServiceToken) map[string]interface{} {
	serviceToken := make(map[string]interface{})
	serviceToken["type"] = string(token.GetType())
	expirationDateTime := token.GetExpirationDateTime()
	const TimeFormat = "2006-01-02T15:04:05.000Z"
	serviceToken["expiration_date_time"] = expirationDateTime.Format(TimeFormat)
	if token.Href != nil {
		serviceToken["href"] = token.GetHref()
	}
	if token.Uuid != nil {
		serviceToken["uuid"] = token.GetUuid()
	}
	if token.State != nil {
		serviceToken["state"] = token.GetState()
	}
	if token.IssuerSide != nil {
		serviceToken["issuer_side"] = token.GetIssuerSide()
	}
	if token.Name != nil {
		serviceToken["name"] = token.GetName()
	}
	if token.Description != nil {
		serviceToken["description"] = token.GetDescription()
	}
	if token.Connection != nil {
		connection := token.GetConnection()
		serviceToken["service_token_connection"] = connectionGoToTerraform(&connection)
	}
	if token.Notifications != nil {
		notifications := token.GetNotifications()
		serviceToken["notifications"] = equinix_fabric_schema.NotificationsGoToTerraform(notifications)
	}
	if token.Account != nil {
		account := token.GetAccount()
		serviceToken["account"] = equinix_fabric_schema.AccountGoToTerraform(&account)
	}
	if token.Changelog != nil {
		changelog := token.GetChangelog()
		serviceToken["change_log"] = equinix_fabric_schema.ChangeLogGoToTerraform(&changelog)
	}
	if token.Project != nil {
		project := token.GetProject()
		serviceToken["project"] = equinix_fabric_schema.ProjectGoToTerraform(&project)
	}

	return serviceToken
}

func connectionTerraformToGo(connectionTerraform []interface{}) fabricv4.ServiceTokenConnection {
	if len(connectionTerraform) == 0 {
		return fabricv4.ServiceTokenConnection{}
	}

	var connection fabricv4.ServiceTokenConnection

	connectionMap := connectionTerraform[0].(map[string]interface{})

	typeVal := connectionMap["type"].(string)
	if typeVal != "" {
		connection.SetType(fabricv4.ServiceTokenConnectionType(typeVal))
	}
	uuid := connectionMap["uuid"].(string)
	if uuid != "" {
		connection.SetUuid(uuid)
	}

	allowRemoteConnection := connectionMap["allow_remote_connection"].(bool)
	connection.SetAllowRemoteConnection(allowRemoteConnection)

	allowCustomBandwidth := connectionMap["allow_custom_bandwidth"].(bool)
	connection.SetAllowCustomBandwidth(allowCustomBandwidth)

	bandwidthLimit := connectionMap["bandwidth_limit"].(int)
	if bandwidthLimit > 0 {
		connection.SetBandwidthLimit(int32(bandwidthLimit))
	}

	supportedBandwidths := connectionMap["supported_bandwidths"].([]interface{})
	if len(supportedBandwidths) > 0 {
		int32Bandwidths := make([]int32, len(supportedBandwidths))
		for i, v := range supportedBandwidths {
			int32Bandwidths[i] = int32(v.(int))
		}
		sort.Slice(int32Bandwidths, func(i, j int) bool {
			return int32Bandwidths[i] < int32Bandwidths[j]
		})
		connection.SetSupportedBandwidths(int32Bandwidths)
	}

	asideRequest := connectionMap["a_side"].(*schema.Set).List()
	zsideRequest := connectionMap["z_side"].(*schema.Set).List()
	if len(asideRequest) != 0 {
		aside := accessPointTerraformToGo(asideRequest)
		connection.SetASide(aside)
	}
	if len(zsideRequest) != 0 {
		zside := accessPointTerraformToGo(zsideRequest)
		connection.SetZSide(zside)
	}
	return connection
}

func accessPointTerraformToGo(accessPoint []interface{}) fabricv4.ServiceTokenSide {
	if len(accessPoint) == 0 {
		return fabricv4.ServiceTokenSide{}
	}

	var apSide fabricv4.ServiceTokenSide

	accessPointMap := accessPoint[0].(map[string]interface{})
	accessPointSelectors := accessPointMap["access_point_selectors"].([]interface{})
	if len(accessPointSelectors) != 0 {
		aps := accessPointSelectorsTerraformToGo(accessPointSelectors)
		apSide.SetAccessPointSelectors(aps)
	}
	return apSide
}

func accessPointSelectorsTerraformToGo(accessPointSelectors []interface{}) []fabricv4.AccessPointSelector {
	if len(accessPointSelectors) == 0 {
		return []fabricv4.AccessPointSelector{}
	}

	var apSelectors fabricv4.AccessPointSelector

	apSelectorsMap := accessPointSelectors[0].(map[string]interface{})
	typeVal := apSelectorsMap["type"].(string)
	apSelectors.SetType(fabricv4.AccessPointSelectorType(typeVal))
	portList := apSelectorsMap["port"].(*schema.Set).List()
	linkProtocolList := apSelectorsMap["link_protocol"].(*schema.Set).List()
	virtualDeviceList := apSelectorsMap["virtual_device"].(*schema.Set).List()
	interfaceList := apSelectorsMap["interface"].(*schema.Set).List()
	networkList := apSelectorsMap["network"].(*schema.Set).List()

	if len(portList) != 0 {
		port := portTerraformToGo(portList)
		apSelectors.SetPort(port)
	}

	if len(linkProtocolList) != 0 {
		linkProtocol := linkProtocolTerraformToGo(linkProtocolList)
		apSelectors.SetLinkProtocol(linkProtocol)
	}

	if len(virtualDeviceList) != 0 {
		virtualDevice := virtualDeviceTerraformToGo(virtualDeviceList)
		apSelectors.SetVirtualDevice(virtualDevice)
	}

	if len(interfaceList) != 0 {
		interfaceInfo := interfaceTerraformToGo(interfaceList)
		apSelectors.SetInterface(interfaceInfo)
	}

	if len(networkList) != 0 {
		network := networkTerraformToGo(networkList)
		apSelectors.SetNetwork(network)
	}

	return []fabricv4.AccessPointSelector{apSelectors}
}

func portTerraformToGo(portList []interface{}) fabricv4.SimplifiedMetadataEntity {
	if len(portList) == 0 {
		return fabricv4.SimplifiedMetadataEntity{}
	}
	var port fabricv4.SimplifiedMetadataEntity
	portListMap := portList[0].(map[string]interface{})
	uuid := portListMap["uuid"].(string)
	href := portListMap["href"].(string)
	portType := portListMap["type"].(string)
	cvpID := portListMap["cvp_id"].(int)
	bandwidth := portListMap["bandwidth"].(int)
	portName := portListMap["port_name"].(string)
	encapsulationProtocolType := portListMap["encapsulation_protocol_type"].(string)
	accountName := portListMap["account_name"].(string)
	priority := portListMap["priority"].(string)
	locationList := portListMap["location"].(*schema.Set).List()

	port.SetUuid(uuid)
	if href != "" {
		port.SetHref(href)
	}
	if portType != "" {
		port.SetType(portType)
	}
	if cvpID != 0 {
		port.SetCvpId(int32(cvpID))
	}
	if bandwidth != 0 {
		port.SetBandwidth(float32(bandwidth))
	}
	if portName != "" {
		port.SetPortName(portName)
	}
	if encapsulationProtocolType != "" {
		port.SetEncapsulationProtocolType(encapsulationProtocolType)
	}
	if accountName != "" {
		port.SetAccountName(accountName)
	}
	if priority != "" {
		port.SetPriority(priority)
	}
	if len(locationList) != 0 {
		location := equinix_fabric_schema.LocationTerraformToGo(locationList)
		port.SetLocation(location)
	}

	return port
}

func linkProtocolTerraformToGo(linkProtocolList []interface{}) fabricv4.SimplifiedLinkProtocol {
	if len(linkProtocolList) == 0 {
		return fabricv4.SimplifiedLinkProtocol{}
	}
	var linkProtocol fabricv4.SimplifiedLinkProtocol
	lpMap := linkProtocolList[0].(map[string]interface{})
	lpType := lpMap["type"].(string)
	lpVlanSTag := int32(lpMap["vlan_s_tag"].(int))
	lpVlanTag := int32(lpMap["vlan_tag"].(int))
	lpVlanCTag := int32(lpMap["vlan_c_tag"].(int))

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

func virtualDeviceTerraformToGo(virtualDeviceList []interface{}) fabricv4.SimplifiedVirtualDevice {
	if len(virtualDeviceList) == 0 {
		return fabricv4.SimplifiedVirtualDevice{}
	}

	var virtualDevice fabricv4.SimplifiedVirtualDevice
	virtualDeviceMap := virtualDeviceList[0].(map[string]interface{})
	href := virtualDeviceMap["href"].(string)
	virtualDeviceType := virtualDeviceMap["type"].(string)
	uuid := virtualDeviceMap["uuid"].(string)
	name := virtualDeviceMap["name"].(string)
	cluster := virtualDeviceMap["cluster"].(string)

	if href != "" {
		virtualDevice.SetHref(href)
	}
	if virtualDeviceType != "" {
		virtualDevice.SetType(fabricv4.SimplifiedVirtualDeviceType(virtualDeviceType))
	}
	virtualDevice.SetUuid(uuid)
	if name != "" {
		virtualDevice.SetName(name)
	}
	if cluster != "" {
		virtualDevice.SetCluster(cluster)
	}

	return virtualDevice
}

func interfaceTerraformToGo(interfaceList []interface{}) fabricv4.VirtualDeviceInterface {
	if len(interfaceList) == 0 {
		return fabricv4.VirtualDeviceInterface{}
	}

	var interfaceInfo fabricv4.VirtualDeviceInterface
	interfaceMap := interfaceList[0].(map[string]interface{})
	uuid := interfaceMap["uuid"].(string)
	interfaceType := interfaceMap["type"].(string)
	id := interfaceMap["id"].(int)

	if uuid != "" {
		interfaceInfo.SetUuid(uuid)
	}
	interfaceInfo.SetType(fabricv4.VirtualDeviceInterfaceType(interfaceType))
	if id >= 0 {
		interfaceInfo.SetId(int32(id))
	}

	return interfaceInfo
}

func networkTerraformToGo(networkList []interface{}) fabricv4.SimplifiedTokenNetwork {
	if len(networkList) == 0 {
		return fabricv4.SimplifiedTokenNetwork{}
	}
	var network fabricv4.SimplifiedTokenNetwork
	networkListMap := networkList[0].(map[string]interface{})
	uuid := networkListMap["uuid"].(string)
	href := networkListMap["href"].(string)
	networkType := networkListMap["type"].(string)
	name := networkListMap["name"].(string)
	scope := networkListMap["scope"].(string)
	locationList := networkListMap["location"].(*schema.Set).List()

	network.SetUuid(uuid)
	if href != "" {
		network.SetHref(href)
	}
	if networkType != "" {
		network.SetType(fabricv4.SimplifiedTokenNetworkType(networkType))
	}
	if name != "" {
		network.SetName(name)
	}
	if scope != "" {
		network.SetScope(fabricv4.SimplifiedTokenNetworkScope(scope))
	}
	if len(locationList) != 0 {
		location := equinix_fabric_schema.LocationTerraformToGo(locationList)
		network.SetLocation(location)
	}
	return network
}

func filtersTerraformToGo(tokens []interface{}) fabricv4.ServiceTokenSearchExpression {
	if tokens == nil {
		return fabricv4.ServiceTokenSearchExpression{}
	}

	searchTokensList := make([]fabricv4.ServiceTokenSearchExpression, 0)

	for _, filter := range tokens {
		filterMap := filter.(map[string]interface{})
		filterItem := fabricv4.ServiceTokenSearchExpression{}
		if property, ok := filterMap["property"]; ok {
			filterItem.SetProperty(fabricv4.ServiceTokenSearchFieldName(property.(string)))
		}
		if operator, ok := filterMap["operator"]; ok {
			filterItem.SetOperator(fabricv4.ServiceTokenSearchExpressionOperator(operator.(string)))
		}
		if values, ok := filterMap["values"]; ok {
			stringValues := converters.IfArrToStringArr(values.([]interface{}))
			filterItem.SetValues(stringValues)
		}
		searchTokensList = append(searchTokensList, filterItem)
	}

	searchTokens := fabricv4.ServiceTokenSearchExpression{}
	searchTokens.SetAnd(searchTokensList)

	return searchTokens
}

func paginationTerraformToGo(pagination []interface{}) fabricv4.PaginationRequest {
	if pagination == nil {
		return fabricv4.PaginationRequest{}
	}
	paginationRequest := fabricv4.PaginationRequest{}
	for _, page := range pagination {
		pageMap := page.(map[string]interface{})
		if offset, ok := pageMap["offset"]; ok {
			paginationRequest.SetOffset(int32(offset.(int)))
		}
		if limit, ok := pageMap["limit"]; ok {
			paginationRequest.SetLimit(int32(limit.(int)))
		}
	}

	return paginationRequest
}

func connectionGoToTerraform(connection *fabricv4.ServiceTokenConnection) *schema.Set {
	mappedConnection := make(map[string]interface{})
	if connection.Type != nil {
		mappedConnection["type"] = string(connection.GetType())
	}
	mappedConnection["allow_remote_connection"] = connection.GetAllowRemoteConnection()
	mappedConnection["allow_custom_bandwidth"] = connection.GetAllowCustomBandwidth()
	if connection.SupportedBandwidths != nil {
		supportedBandwidths := connection.GetSupportedBandwidths()
		interfaceBandwidths := make([]interface{}, len(supportedBandwidths))

		for i, v := range supportedBandwidths {
			interfaceBandwidths[i] = int(v) // Convert each int32 to interface{}
		}

		sort.Slice(interfaceBandwidths, func(i, j int) bool {
			return interfaceBandwidths[i].(int) < interfaceBandwidths[j].(int)
		})

		mappedConnection["supported_bandwidths"] = interfaceBandwidths
	}
	if connection.BandwidthLimit != nil {
		mappedConnection["bandwidth_limit"] = int(connection.GetBandwidthLimit())
	}
	if connection.ASide != nil {
		accessPoint := connection.GetASide()
		mappedConnection["a_side"] = accessPointGoToTerraform(&accessPoint)
	}
	if connection.ZSide != nil {
		accessPoint := connection.GetZSide()
		mappedConnection["z_side"] = accessPointGoToTerraform(&accessPoint)
	}
	connectionSet := schema.NewSet(
		schema.HashResource(serviceTokenConnectionSch()),
		[]interface{}{mappedConnection},
	)
	return connectionSet
}

func accessPointGoToTerraform(accessPoint *fabricv4.ServiceTokenSide) *schema.Set {
	return schema.NewSet(
		schema.HashResource(serviceTokenAccessPointSch()),
		[]interface{}{map[string]interface{}{
			"access_point_selectors": accessPointSelectorsGoToTerraform(accessPoint.GetAccessPointSelectors()),
		}},
	)
}

func accessPointSelectorsGoToTerraform(apSelectors []fabricv4.AccessPointSelector) []interface{} {
	mappedSelectors := make([]interface{}, len(apSelectors))
	for index, selector := range apSelectors {
		mappedAccessPointSelector := make(map[string]interface{})
		if selector.Type != nil {
			mappedAccessPointSelector["type"] = string(selector.GetType())
		}
		if selector.Port != nil {
			port := selector.GetPort()
			mappedAccessPointSelector["port"] = portGoToTerraform(&port)
		}
		if selector.LinkProtocol != nil {
			linkProtocol := selector.GetLinkProtocol()
			mappedAccessPointSelector["link_protocol"] = linkedProtocolGoToTerraform(&linkProtocol)
		}
		if selector.VirtualDevice != nil {
			virtualDevice := selector.GetVirtualDevice()
			mappedAccessPointSelector["virtual_device"] = virtualDeviceGoToTerraform(&virtualDevice)
		}
		if selector.Interface != nil {
			interfaceInfo := selector.GetInterface()
			mappedAccessPointSelector["interface"] = interfaceGoToTerraform(&interfaceInfo)
		}
		if selector.Network != nil {
			network := selector.GetNetwork()
			mappedAccessPointSelector["network"] = networkGoToTerraform(&network)
		}
		mappedSelectors[index] = mappedAccessPointSelector
	}

	return mappedSelectors
}

func portGoToTerraform(port *fabricv4.SimplifiedMetadataEntity) *schema.Set {
	if port == nil {
		return nil
	}
	mappedPort := make(map[string]interface{})
	if href := port.GetHref(); href != "" {
		mappedPort["href"] = href
	}
	if uuid := port.GetUuid(); uuid != "" {
		mappedPort["uuid"] = uuid
	}
	if port.GetType() != "" {
		mappedPort["type"] = port.GetType()
	}
	if cvpID := port.GetCvpId(); cvpID != 0 {
		mappedPort["cvp_id"] = port.GetCvpId()
	}
	if bandwidth := port.GetBandwidth(); bandwidth != 0 {
		mappedPort["bandwidth"] = port.GetBandwidth()
	}
	if portName := port.GetPortName(); portName != "" {
		mappedPort["port_name"] = port.GetPortName()
	}
	if encapsulationProtocolType := port.GetEncapsulationProtocolType(); encapsulationProtocolType != "" {
		mappedPort["encapsulation_protocol_type"] = port.GetEncapsulationProtocolType()
	}
	if accountName := port.GetAccountName(); accountName != "" {
		mappedPort["account_name"] = port.GetAccountName()
	}
	if priority := port.GetPriority(); priority != "" {
		mappedPort["priority"] = port.GetPriority()
	}
	if port.Location != nil {
		location := port.GetLocation()
		mappedPort["location"] = equinix_fabric_schema.LocationGoToTerraform(&location)
	}

	portSet := schema.NewSet(
		schema.HashResource(portSch()),
		[]interface{}{mappedPort},
	)
	return portSet
}

func linkedProtocolGoToTerraform(linkedProtocol *fabricv4.SimplifiedLinkProtocol) *schema.Set {

	mappedLinkedProtocol := make(map[string]interface{})
	mappedLinkedProtocol["type"] = string(linkedProtocol.GetType())
	mappedLinkedProtocol["vlan_tag"] = int(linkedProtocol.GetVlanTag())
	mappedLinkedProtocol["vlan_s_tag"] = int(linkedProtocol.GetVlanSTag())
	mappedLinkedProtocol["vlan_c_tag"] = int(linkedProtocol.GetVlanCTag())

	linkedProtocolSet := schema.NewSet(
		schema.HashResource(linkProtocolSch()),
		[]interface{}{mappedLinkedProtocol},
	)
	return linkedProtocolSet
}

func virtualDeviceGoToTerraform(virtualDevice *fabricv4.SimplifiedVirtualDevice) *schema.Set {
	if virtualDevice == nil {
		return nil
	}
	mappedVirtualDevice := make(map[string]interface{})
	if name := virtualDevice.GetName(); name != "" {
		mappedVirtualDevice["name"] = name
	}
	if href := virtualDevice.GetHref(); href != "" {
		mappedVirtualDevice["href"] = href
	}
	if virtualDevice.GetType() != "" {
		mappedVirtualDevice["type"] = string(virtualDevice.GetType())
	}
	if uuid := virtualDevice.GetUuid(); uuid != "" {
		mappedVirtualDevice["uuid"] = uuid
	}
	if virtualDevice.Cluster != nil && virtualDevice.GetCluster() != "" {
		mappedVirtualDevice["cluster"] = virtualDevice.GetCluster()
	}
	virtualDeviceSet := schema.NewSet(
		schema.HashResource(virtualDeviceSch()),
		[]interface{}{mappedVirtualDevice},
	)
	return virtualDeviceSet
}

func interfaceGoToTerraform(mInterface *fabricv4.VirtualDeviceInterface) *schema.Set {
	if mInterface == nil {
		return nil
	}
	mappedMInterface := make(map[string]interface{})
	mappedMInterface["id"] = int(mInterface.GetId())
	mappedMInterface["type"] = string(mInterface.GetType())
	mappedMInterface["uuid"] = mInterface.GetUuid()

	mInterfaceSet := schema.NewSet(
		schema.HashResource(interfaceSch()),
		[]interface{}{mappedMInterface},
	)
	return mInterfaceSet
}

func networkGoToTerraform(network *fabricv4.SimplifiedTokenNetwork) *schema.Set {
	if network == nil {
		return nil
	}

	mappedNetwork := make(map[string]interface{})
	if uuid := network.GetUuid(); uuid != "" {
		mappedNetwork["uuid"] = uuid
	}
	if href := network.GetHref(); href != "" {
		mappedNetwork["href"] = href
	}
	if networkType := network.GetType(); networkType != "" {
		mappedNetwork["type"] = string(networkType)
	}
	if name := network.GetName(); name != "" {
		mappedNetwork["name"] = name
	}
	if scope := network.GetName(); scope != "" {
		mappedNetwork["scope"] = string(network.GetScope())
	}
	if network.Location != nil {
		location := network.GetLocation()
		mappedNetwork["location"] = equinix_fabric_schema.LocationGoToTerraform(&location)
	}
	return schema.NewSet(
		schema.HashResource(networkSch()),
		[]interface{}{mappedNetwork},
	)
}

func paginationGoToTerraform(pagination *fabricv4.Pagination) *schema.Set {
	if pagination == nil {
		return nil
	}
	mappedPagination := make(map[string]interface{})
	mappedPagination["offset"] = int(pagination.GetOffset())
	mappedPagination["limit"] = int(pagination.GetLimit())
	mappedPagination["total"] = int(pagination.GetTotal())
	mappedPagination["next"] = pagination.GetNext()
	mappedPagination["previous"] = pagination.GetPrevious()

	return schema.NewSet(
		schema.HashResource(paginationSchema()),
		[]interface{}{mappedPagination},
	)
}
