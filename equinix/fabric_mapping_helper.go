package equinix

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func serviceTokenToFabric(serviceTokenRequest []interface{}) v4.ServiceToken {
	mappedST := v4.ServiceToken{}
	for _, str := range serviceTokenRequest {
		stMap := str.(map[string]interface{})
		stType := stMap["type"].(interface{}).(string)
		uuid := stMap["uuid"].(interface{}).(string)
		stTypeObj := v4.ServiceTokenType(stType)
		mappedST = v4.ServiceToken{Type_: &stTypeObj, Uuid: uuid}
	}
	return mappedST
}

func additionalInfoToFabric(additionalInfoRequest []interface{}) []v4.ConnectionSideAdditionalInfo {
	var mappedaiArray []v4.ConnectionSideAdditionalInfo
	for i, ai := range additionalInfoRequest {
		aiMap := ai.(map[string]interface{})
		key := aiMap["key"].(interface{}).(string)
		value := aiMap["value"].(interface{}).(string)
		mappedai := v4.ConnectionSideAdditionalInfo{Key: key, Value: value}
		mappedaiArray[i] = mappedai
	}
	return mappedaiArray
}

func accessPointToFabric(accessPointRequest []interface{}) v4.AccessPoint {
	accessPoint := v4.AccessPoint{}
	for _, ap := range accessPointRequest {
		accessPointMap := ap.(map[string]interface{})
		portList := accessPointMap["port"].(interface{}).(*schema.Set).List()
		profileList := accessPointMap["profile"].(interface{}).(*schema.Set).List()
		locationList := accessPointMap["location"].(interface{}).(*schema.Set).List()
		typeVal := accessPointMap["type"].(interface{}).(string)
		authenticationKey := accessPointMap["authentication_key"].(interface{}).(string)
		providerConnectionId := accessPointMap["provider_connection_id"].(interface{}).(string)
		sellerRegion := accessPointMap["seller_region"].(interface{}).(string)
		peeringTypeRaw := accessPointMap["peering_type"].(interface{}).(string)
		gatewayRequest := accessPointMap["gateway"].(interface{}).(*schema.Set).List()

		mappedGWr := v4.VirtualGateway{}
		if len(gatewayRequest) != 0 {
			mappedGWr = gatewayToFabric(gatewayRequest)
		}

		apt := v4.AccessPointType(typeVal)
		p := v4.Port{}
		if len(portList) != 0 {
			p = portToFabric(portList)
		}
		linkProtocolList := accessPointMap["link_protocol"].(interface{}).(*schema.Set).List()
		slp := v4.SimplifiedLinkProtocol{}
		if len(linkProtocolList) != 0 {
			slp = linkProtocolToFabric(linkProtocolList)
		}
		ssp := v4.SimplifiedServiceProfile{}
		if len(profileList) != 0 {
			ssp = simplifiedServiceProfileToFabric(profileList)
		}

		sl := v4.SimplifiedLocation{}
		if len(locationList) != 0 {
			sl = locationToFabric(locationList)
		}

		if peeringTypeRaw != "" {
			peeringType := v4.PeeringType(peeringTypeRaw)
			accessPoint = v4.AccessPoint{
				Type_: &apt, Port: &p, LinkProtocol: &slp, AuthenticationKey: authenticationKey,
				Profile: &ssp, Location: &sl, ProviderConnectionId: providerConnectionId, SellerRegion: sellerRegion, PeeringType: &peeringType,
			}
		} else {
			accessPoint = v4.AccessPoint{
				Type_: &apt, Port: &p, LinkProtocol: &slp, AuthenticationKey: authenticationKey,
				Profile: &ssp, Gateway: &mappedGWr, Location: &sl, ProviderConnectionId: providerConnectionId, SellerRegion: sellerRegion,
			}
		}
	}
	return accessPoint
}

func gatewayToFabric(gatewayRequest []interface{}) v4.VirtualGateway {
	gatewayMapped := v4.VirtualGateway{}
	for _, gwr := range gatewayRequest {
		gwrMap := gwr.(map[string]interface{})
		gwHref := gwrMap["href"].(interface{}).(string)
		gwuuid := gwrMap["uuid"].(interface{}).(string)
		gatewayMapped = v4.VirtualGateway{Uuid: gwuuid, Href: gwHref}
	}
	return gatewayMapped
}

func projectToFabric(projectRequest []interface{}) v4.Project {
	if projectRequest == nil || len(projectRequest) == 0 {
		return v4.Project{}
	}
	mappedPr := v4.Project{}
	for _, pr := range projectRequest {
		prMap := pr.(map[string]interface{})
		projectId := prMap["project_id"].(interface{}).(string)
		href := prMap["href"].(interface{}).(string)
		mappedPr = v4.Project{ProjectId: projectId, Href: href}
	}
	return mappedPr
}

func notificationToFabric(schemaNotifications []interface{}) []v4.SimplifiedNotification {
	if schemaNotifications == nil {
		return []v4.SimplifiedNotification{}
	}
	var notifications []v4.SimplifiedNotification
	for _, n := range schemaNotifications {
		ntype := n.(map[string]interface{})["type"].(string)
		interval := n.(map[string]interface{})["send_interval"].(string)
		emailsRaw := n.(map[string]interface{})["emails"].([]interface{})
		emails := expandListToStringList(emailsRaw)
		notifications = append(notifications, v4.SimplifiedNotification{
			Type_:        ntype,
			SendInterval: interval,
			Emails:       emails,
		})
	}
	return notifications
}

func redundancyToFabric(schemaRedundancy []interface{}) v4.ConnectionRedundancy {
	if schemaRedundancy == nil {
		return v4.ConnectionRedundancy{}
	}
	red := v4.ConnectionRedundancy{}
	for _, r := range schemaRedundancy {
		redundancyMap := r.(map[string]interface{})
		priority := redundancyMap["priority"]
		priorityCont := v4.ConnectionPriority(priority.(string))
		red = v4.ConnectionRedundancy{
			Priority: &priorityCont,
		}
	}
	return red
}

func orderToFabric(schemaOrder []interface{}) v4.Order {
	if schemaOrder == nil {
		return v4.Order{}
	}
	order := v4.Order{}
	for _, o := range schemaOrder {
		orderMap := o.(map[string]interface{})
		purchaseOrderNumber := orderMap["purchase_order_number"]
		billingTier := orderMap["billing_tier"]
		orderId := orderMap["order_id"]
		orderNumber := orderMap["order_number"]
		order = v4.Order{PurchaseOrderNumber: purchaseOrderNumber.(string), BillingTier: billingTier.(string), OrderId: orderId.(string), OrderNumber: orderNumber.(string)}
	}
	return order
}

func linkProtocolToFabric(linkProtocolList []interface{}) v4.SimplifiedLinkProtocol {
	slp := v4.SimplifiedLinkProtocol{}
	for _, lp := range linkProtocolList {
		lpMap := lp.(map[string]interface{})
		lpType := lpMap["type"].(interface{}).(string)
		lpVlanSTag := lpMap["vlan_s_tag"].(interface{}).(int)
		lpVlanTag := lpMap["vlan_tag"].(interface{}).(int)
		lpVlanCTag := lpMap["vlan_c_tag"].(interface{}).(int)
		lpt := v4.LinkProtocolType(lpType)
		slp = v4.SimplifiedLinkProtocol{Type_: &lpt, VlanSTag: int32(lpVlanSTag), VlanTag: int32(lpVlanTag), VlanCTag: int32(lpVlanCTag)}
	}
	return slp
}

func portToFabric(portList []interface{}) v4.Port {
	p := v4.Port{}
	for _, pl := range portList {
		plMap := pl.(map[string]interface{})
		uuid := plMap["uuid"].(interface{}).(string)
		p = v4.Port{Uuid: uuid}
	}
	return p
}

func simplifiedServiceProfileToFabric(profileList []interface{}) v4.SimplifiedServiceProfile {
	ssp := v4.SimplifiedServiceProfile{}
	for _, pl := range profileList {
		plMap := pl.(map[string]interface{})
		ptype := plMap["type"].(interface{}).(string)
		spte := v4.ServiceProfileTypeEnum(ptype)
		uuid := plMap["uuid"].(interface{}).(string)
		ssp = v4.SimplifiedServiceProfile{Uuid: uuid, Type_: &spte}
	}
	return ssp
}

func locationToFabric(locationList []interface{}) v4.SimplifiedLocation {
	sl := v4.SimplifiedLocation{}
	for _, ll := range locationList {
		llMap := ll.(map[string]interface{})
		metroName := llMap["metro_name"]
		var metroNamestr string
		if metroName != nil {
			metroNamestr = metroName.(interface{}).(string)
		}
		region := llMap["region"].(interface{}).(string)
		mc := llMap["metro_code"].(interface{}).(string)
		ibx := llMap["ibx"].(interface{}).(string)
		sl = v4.SimplifiedLocation{MetroCode: mc, Region: region, Ibx: ibx, MetroName: metroNamestr}
	}
	return sl
}

func accountToTerra(account *v4.SimplifiedAccount) *schema.Set {
	if account == nil {
		return nil
	}
	accounts := []*v4.SimplifiedAccount{account}
	mappedAccounts := make([]interface{}, 0)
	for _, account := range accounts {
		mappedAccount := make(map[string]interface{})
		mappedAccount["account_number"] = int(account.AccountNumber)
		mappedAccount["account_name"] = account.AccountName
		mappedAccount["org_id"] = int(account.OrgId)
		mappedAccount["organization_name"] = account.OrganizationName
		mappedAccount["global_org_id"] = account.GlobalOrgId
		mappedAccount["global_organization_name"] = account.GlobalOrganizationName
		mappedAccount["global_cust_id"] = account.GlobalCustId
		mappedAccounts = append(mappedAccounts, mappedAccount)
	}

	// Setting a Set in a List does not work correctly
	// see https://github.com/hashicorp/terraform/issues/16331 for details
	accountSet := schema.NewSet(
		schema.HashResource(createAccountRes),
		mappedAccounts,
	)
	return accountSet
}

func errorToTerra(errors []v4.ModelError) []map[string]interface{} {
	if errors == nil {
		return nil
	}
	mappedErrors := make([]map[string]interface{}, len(errors))
	for index, mError := range errors {
		mappedErrors[index] = map[string]interface{}{
			"error_code":      mError.ErrorCode,
			"error_message":   mError.ErrorMessage,
			"correlation_id":  mError.CorrelationId,
			"details":         mError.Details,
			"help":            mError.Help,
			"additional_info": errorAdditionalInfoToTerra(mError.AdditionalInfo),
		}
	}
	return mappedErrors
}

func errorAdditionalInfoToTerra(additionalInfol []v4.PriceErrorAdditionalInfo) []map[string]interface{} {
	if additionalInfol == nil {
		return nil
	}
	mappedAdditionalInfol := make([]map[string]interface{}, len(additionalInfol))
	for index, additionalInfo := range additionalInfol {
		mappedAdditionalInfol[index] = map[string]interface{}{
			"property": additionalInfo.Property,
			"reason":   additionalInfo.Reason,
		}
	}
	return mappedAdditionalInfol
}

func operationToTerra(operation *v4.ConnectionOperation) *schema.Set {
	if operation == nil {
		return nil
	}
	operations := []*v4.ConnectionOperation{operation}
	mappedOperations := make([]interface{}, 0)
	for _, operation := range operations {
		mappedOperation := make(map[string]interface{})
		mappedOperation["provider_status"] = string(*operation.ProviderStatus)
		mappedOperation["equinix_status"] = string(*operation.EquinixStatus)
		mappedOperation["errors"] = errorToTerra(operation.Errors)
		mappedOperations = append(mappedOperations, mappedOperation)
	}
	operationSet := schema.NewSet(
		schema.HashResource(createOperationRes),
		mappedOperations,
	)
	return operationSet
}

func orderMappingToTerra(order *v4.Order) *schema.Set {
	if order == nil {
		return nil
	}
	orders := []*v4.Order{order}
	mappedOrders := make([]interface{}, 0)
	for _, order := range orders {
		mappedOrder := make(map[string]interface{})
		mappedOrder["purchase_order_number"] = order.PurchaseOrderNumber
		mappedOrder["billing_tier"] = order.BillingTier
		mappedOrder["order_id"] = order.OrderId
		mappedOrder["order_number"] = order.OrderNumber
		mappedOrders = append(mappedOrders, mappedOrder)
	}
	orderSet := schema.NewSet(
		schema.HashResource(createOrderRes),
		mappedOrders,
	)
	return orderSet
}

func changeLogToTerra(changeLog *v4.Changelog) *schema.Set {
	if changeLog == nil {
		return nil
	}
	changeLogs := []*v4.Changelog{changeLog}
	mappedChangeLogs := make([]interface{}, 0)
	for _, changeLog := range changeLogs {
		mappedChangeLog := make(map[string]interface{})
		mappedChangeLog["created_by"] = changeLog.CreatedBy
		mappedChangeLog["created_by_full_name"] = changeLog.CreatedByFullName
		mappedChangeLog["created_by_email"] = changeLog.CreatedByEmail
		mappedChangeLog["created_date_time"] = changeLog.CreatedDateTime.String()
		mappedChangeLog["updated_by"] = changeLog.UpdatedBy
		mappedChangeLog["updated_by_full_name"] = changeLog.UpdatedByFullName
		mappedChangeLog["updated_date_time"] = changeLog.UpdatedDateTime.String()
		mappedChangeLog["deleted_by"] = changeLog.DeletedBy
		mappedChangeLog["deleted_by_full_name"] = changeLog.DeletedByFullName
		mappedChangeLog["deleted_by_email"] = changeLog.DeletedByEmail
		mappedChangeLog["deleted_date_time"] = changeLog.DeletedDateTime.String()
		mappedChangeLogs = append(mappedChangeLogs, mappedChangeLog)
	}
	changeLogSet := schema.NewSet(
		schema.HashResource(createChangeLogRes),
		mappedChangeLogs,
	)
	return changeLogSet
}

func portRedundancyToTerra(redundancy *v4.PortRedundancy) *schema.Set {
	redundancies := []*v4.PortRedundancy{redundancy}
	mappedRedundancys := make([]interface{}, 0)
	for _, redundancy := range redundancies {
		mappedRedundancy := make(map[string]interface{})
		mappedRedundancy["priority"] = string(*redundancy.Priority)
		mappedRedundancys = append(mappedRedundancys, mappedRedundancy)
	}
	redundancySet := schema.NewSet(
		schema.HashResource(createPortRedundancyRes),
		mappedRedundancys,
	)
	return redundancySet
}

func redundancyToTerra(redundancy *v4.ConnectionRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	redundancies := []*v4.ConnectionRedundancy{redundancy}
	mappedRedundancys := make([]interface{}, 0)
	for _, redundancy := range redundancies {
		mappedRedundancy := make(map[string]interface{})
		mappedRedundancy["group"] = redundancy.Group
		mappedRedundancy["priority"] = string(*redundancy.Priority)
		mappedRedundancys = append(mappedRedundancys, mappedRedundancy)
	}
	redundancySet := schema.NewSet(
		schema.HashResource(createRedundancyRes),
		mappedRedundancys,
	)
	return redundancySet
}

func notificationToTerra(notifications []v4.SimplifiedNotification) []map[string]interface{} {
	if notifications == nil {
		return nil
	}
	mappedNotifications := make([]map[string]interface{}, len(notifications))
	for index, notification := range notifications {
		mappedNotifications[index] = map[string]interface{}{
			"type":          notification.Type_,
			"send_interval": notification.SendInterval,
			"emails":        notification.Emails,
		}
	}
	return mappedNotifications
}

func locationToTerra(location *v4.SimplifiedLocation) *schema.Set {
	locations := []*v4.SimplifiedLocation{location}
	mappedLocations := make([]interface{}, 0)
	for _, location := range locations {
		mappedLocation := make(map[string]interface{})
		mappedLocation["region"] = location.Region
		mappedLocation["metro_name"] = location.MetroName
		mappedLocation["metro_code"] = location.MetroCode
		mappedLocation["ibx"] = location.Ibx
		mappedLocations = append(mappedLocations, mappedLocation)
	}
	locationSet := schema.NewSet(
		schema.HashResource(createLocationRes),
		mappedLocations,
	)
	return locationSet
}

func serviceTokenToTerra(serviceToken *v4.ServiceToken) *schema.Set {
	if serviceToken == nil {
		return nil
	}
	serviceTokens := []*v4.ServiceToken{serviceToken}
	mappedServiceTokens := make([]interface{}, 0)
	for _, serviceToken := range serviceTokens {
		mappedServiceToken := make(map[string]interface{})
		if serviceToken.Type_ != nil {
			mappedServiceToken["type"] = string(*serviceToken.Type_)
		}
		mappedServiceToken["href"] = serviceToken.Href
		mappedServiceToken["uuid"] = serviceToken.Uuid
		mappedServiceToken["description"] = serviceToken.Description
		mappedServiceTokens = append(mappedServiceTokens, mappedServiceToken)
	}
	serviceTokenSet := schema.NewSet(
		schema.HashResource(createServiceTokenRes),
		mappedServiceTokens,
	)
	return serviceTokenSet
}

func connectionSideToTerra(connectionSide *v4.ConnectionSide) *schema.Set {
	connectionSides := []*v4.ConnectionSide{connectionSide}
	mappedConnectionSides := make([]interface{}, 0)
	for _, connectionSide := range connectionSides {
		mappedConnectionSide := make(map[string]interface{})
		serviceTokenSet := serviceTokenToTerra(connectionSide.ServiceToken)
		if serviceTokenSet != nil {
			mappedConnectionSide["service_token"] = serviceTokenSet
		}
		mappedConnectionSide["access_point"] = accessPointToTerra(connectionSide.AccessPoint)
		mappedConnectionSides = append(mappedConnectionSides, mappedConnectionSide)
	}
	connectionSideSet := schema.NewSet(
		schema.HashResource(createFabricConnectionSideRes),
		mappedConnectionSides,
	)
	return connectionSideSet
}

func additionalInfoToTerra(additionalInfol []v4.ConnectionSideAdditionalInfo) []map[string]interface{} {
	if additionalInfol == nil {
		return nil
	}
	mappedadditionalInfol := make([]map[string]interface{}, len(additionalInfol))
	for index, additionalInfo := range additionalInfol {
		mappedadditionalInfol[index] = map[string]interface{}{
			"key":   additionalInfo.Key,
			"value": additionalInfo.Value,
		}
	}
	return mappedadditionalInfol
}

func fabricGatewayToTerra(virtualGateway *v4.VirtualGateway) *schema.Set {
	if virtualGateway == nil {
		return nil
	}
	virtualGateways := []*v4.VirtualGateway{virtualGateway}
	mappedvirtualGateways := make([]interface{}, 0)
	for _, virtualGateway := range virtualGateways {
		mappedvirtualGateway := make(map[string]interface{})
		mappedvirtualGateway["uuid"] = virtualGateway.Uuid
		mappedvirtualGateway["href"] = virtualGateway.Href
		mappedvirtualGateways = append(mappedvirtualGateways, mappedvirtualGateway)
	}
	linkedProtocolSet := schema.NewSet(
		schema.HashResource(createGatewayProjectSchRes),
		mappedvirtualGateways)
	return linkedProtocolSet
}

func projectToTerra(project *v4.Project) *schema.Set {
	if project == nil {
		return nil
	}
	projects := []*v4.Project{project}
	mappedProjects := make([]interface{}, 0)
	for _, project := range projects {
		mappedProject := make(map[string]interface{})
		mappedProject["project_id"] = project.ProjectId
		mappedProject["href"] = project.Href
		mappedProjects = append(mappedProjects, mappedProject)
	}
	projectSet := schema.NewSet(
		schema.HashResource(createGatewayProjectSchRes),
		mappedProjects)
	return projectSet
}

func virtualDeviceToTerra(virtualDevice *v4.VirtualDevice) *schema.Set {
	if virtualDevice == nil {
		return nil
	}
	virtualDevices := []*v4.VirtualDevice{virtualDevice}
	mappedVirtualDevices := make([]interface{}, 0)
	for _, virtualDevice := range virtualDevices {
		mappedVirtualDevice := make(map[string]interface{})
		mappedVirtualDevice["name"] = virtualDevice.Name
		mappedVirtualDevice["href"] = virtualDevice.Href
		mappedVirtualDevice["type"] = virtualDevice.Type_
		mappedVirtualDevice["uuid"] = virtualDevice.Uuid
		mappedVirtualDevices = append(mappedVirtualDevices, mappedVirtualDevice)
	}
	virtualDeviceSet := schema.NewSet(
		schema.HashResource(createAccessPointVirtualDeviceRes),
		mappedVirtualDevices)
	return virtualDeviceSet
}

func interfaceToTerra(mInterface *v4.ModelInterface) *schema.Set {
	if mInterface == nil {
		return nil
	}
	mInterfaces := []*v4.ModelInterface{mInterface}
	mappedMInterfaces := make([]interface{}, 0)
	for _, mInterface := range mInterfaces {
		mappedMInterface := make(map[string]interface{})
		mappedMInterface["id"] = mInterface.Id
		mappedMInterface["type"] = mInterface.Type_
		mappedMInterface["uuid"] = mInterface.Uuid
		mappedMInterfaces = append(mappedMInterfaces, mappedMInterface)
	}
	mInterfaceSet := schema.NewSet(
		schema.HashResource(createAccessPointVirtualDeviceRes),
		mappedMInterfaces)
	return mInterfaceSet
}

func accessPointToTerra(accessPoint *v4.AccessPoint) *schema.Set {
	accessPoints := []*v4.AccessPoint{accessPoint}
	mappedAccessPoints := make([]interface{}, 0)
	for _, accessPoint := range accessPoints {
		mappedAccessPoint := make(map[string]interface{})
		if accessPoint.Type_ != nil {
			mappedAccessPoint["type"] = string(*accessPoint.Type_)
		}
		if accessPoint.Account != nil {
			mappedAccessPoint["account"] = accountToTerra(accessPoint.Account)
		}
		if accessPoint.Location != nil {
			mappedAccessPoint["location"] = locationToTerra(accessPoint.Location)
		}

		if accessPoint.Port != nil {
			mappedAccessPoint["port"] = portToTerra(accessPoint.Port)
		}
		if accessPoint.Profile != nil {
			mappedAccessPoint["profile"] = simplifiedServiceProfileToTerra(accessPoint.Profile)
		}

		if accessPoint.Gateway != nil {
			mappedAccessPoint["gateway"] = fabricGatewayToTerra(accessPoint.Gateway)
		}

		if accessPoint.LinkProtocol != nil {
			mappedAccessPoint["link_protocol"] = linkedProtocolToTerra(accessPoint.LinkProtocol)
		}

		if accessPoint.VirtualDevice != nil {
			mappedAccessPoint["virtual_device"] = virtualDeviceToTerra(accessPoint.VirtualDevice)
		}
		if accessPoint.Interface_ != nil {
			mappedAccessPoint["interface"] = interfaceToTerra(accessPoint.Interface_)
		}
		mappedAccessPoint["seller_region"] = accessPoint.SellerRegion
		if accessPoint.PeeringType != nil {
			mappedAccessPoint["peering_type"] = string(*accessPoint.PeeringType)
		}
		mappedAccessPoint["authentication_key"] = accessPoint.AuthenticationKey
		mappedAccessPoint["provider_connection_id"] = accessPoint.ProviderConnectionId
		mappedAccessPoints = append(mappedAccessPoints, mappedAccessPoint)
	}
	accessPointSet := schema.NewSet(
		schema.HashResource(createConnectionSideAccessPointRes),
		mappedAccessPoints,
	)
	return accessPointSet
}

func linkedProtocolToTerra(linkedProtocol *v4.SimplifiedLinkProtocol) *schema.Set {
	linkedProtocols := []*v4.SimplifiedLinkProtocol{linkedProtocol}
	mappedLinkedProtocols := make([]interface{}, 0)
	for _, linkedProtocol := range linkedProtocols {
		mappedLinkedProtocol := make(map[string]interface{})
		mappedLinkedProtocol["type"] = string(*linkedProtocol.Type_)
		mappedLinkedProtocol["vlan_tag"] = int(linkedProtocol.VlanTag)
		mappedLinkedProtocol["vlan_s_tag"] = int(linkedProtocol.VlanSTag)
		mappedLinkedProtocol["vlan_c_tag"] = int(linkedProtocol.VlanCTag)
		mappedLinkedProtocols = append(mappedLinkedProtocols, mappedLinkedProtocol)
	}
	linkedProtocolSet := schema.NewSet(
		schema.HashResource(createAccessPointLinkProtocolSchRes),
		mappedLinkedProtocols)
	return linkedProtocolSet
}

func simplifiedServiceProfileToTerra(profile *v4.SimplifiedServiceProfile) *schema.Set {
	profiles := []*v4.SimplifiedServiceProfile{profile}
	mappedProfiles := make([]interface{}, 0)
	for _, profile := range profiles {
		mappedProfile := make(map[string]interface{})
		mappedProfile["href"] = profile.Href
		mappedProfile["type"] = string(*profile.Type_)
		mappedProfile["name"] = profile.Name
		mappedProfile["uuid"] = profile.Uuid
		mappedProfile["access_point_type_configs"] = accessPointTypeConfigToTerra(profile.AccessPointTypeConfigs)
		mappedProfiles = append(mappedProfiles, mappedProfile)
	}
	profileSet := schema.NewSet(
		schema.HashResource(createLocationRes),
		mappedProfiles,
	)
	return profileSet
}

func accessPointTypeConfigToTerra(spAccessPointTypes []v4.ServiceProfileAccessPointType) []interface{} {
	mappedSpAccessPointTypes := make([]interface{}, len(spAccessPointTypes))
	for index, spAccessPointType := range spAccessPointTypes {
		mappedSpAccessPointTypes[index] = map[string]interface{}{
			"type": string(*spAccessPointType.Type_),
			"uuid": spAccessPointType.Uuid,
		}
	}
	return mappedSpAccessPointTypes
}

func portToTerra(port *v4.Port) *schema.Set {
	ports := []*v4.Port{port}
	mappedPorts := make([]interface{}, 0)
	for _, port := range ports {
		mappedPort := make(map[string]interface{})
		mappedPort["href"] = port.Href
		mappedPort["name"] = port.Name
		mappedPort["uuid"] = port.Uuid
		if port.Redundancy != nil {
			mappedPort["redundancy"] = portRedundancyToTerra(port.Redundancy)
		}
		mappedPorts = append(mappedPorts, mappedPort)
	}
	portSet := schema.NewSet(
		schema.HashResource(createPortRes),
		mappedPorts,
	)
	return portSet
}

func getUpdateRequest(conn v4.Connection, d *schema.ResourceData) (v4.ConnectionChangeOperation, error) {
	changeOps := v4.ConnectionChangeOperation{}
	existingName := conn.Name
	existingBandwidth := int(conn.Bandwidth)
	updateNameVal := d.Get("name").(string)
	updateBandwidthVal := d.Get("bandwidth").(int)

	log.Printf("existing name %s, existing bandwidth %d, Update Name Request %s, Update Bandwidth Request %d ",
		existingName, existingBandwidth, updateNameVal, updateBandwidthVal)

	if existingName != updateNameVal {
		changeOps = v4.ConnectionChangeOperation{Op: "replace", Path: "/name", Value: updateNameVal}
	} else if existingBandwidth != updateBandwidthVal {
		changeOps = v4.ConnectionChangeOperation{Op: "replace", Path: "/bandwidth", Value: updateBandwidthVal}
	} else {
		return changeOps, fmt.Errorf("Nothing to update for the connection %s", existingName)
	}
	return changeOps, nil
}

const allowed_charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#$&@"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func CorrelationIdWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func CorrelationId(length int) string {
	return CorrelationIdWithCharset(length, allowed_charset)
}
