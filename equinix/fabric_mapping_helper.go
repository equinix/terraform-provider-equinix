package equinix

import (
	"fmt"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

func serviceTokenToFabric(serviceTokenRequest []interface{}) v4.ServiceToken {
	mappedST := v4.ServiceToken{}
	for _, str := range serviceTokenRequest {
		stMap := str.(map[string]interface{})
		stType := stMap["type"].(interface{}).(string)
		uuid := stMap["uuid"].(interface{}).(string)
		description := stMap["description"].(interface{}).(string)
		expirationDateTime := stMap["expiration_date_time"].(interface{}).(string)
		notifications := stMap["notifications"].(interface{}).(*schema.Set).List()
		var mappedNotifications []v4.SimplifiedNotification
		if len(notifications) != 0 {
			mappedNotifications = notificationToFabric(notifications)
		}
		stTypeObj := v4.ServiceTokenType(stType)

		mappedST = v4.ServiceToken{Type_: &stTypeObj, Uuid: uuid, Description: description, ExpirationDateTime: stringToTime("expiration_date_time", "", expirationDateTime),
			Notifications: mappedNotifications}
	}
	return mappedST
}

func additionalInfoToFabric(additionalInfoRequest []interface{}) []v4.ConnectionSideAdditionalInfo {

	var mappedaiArray []v4.ConnectionSideAdditionalInfo
	for _, ai := range additionalInfoRequest {
		i := 0
		aiMap := ai.(map[string]interface{})
		key := aiMap["key"].(interface{}).(string)
		value := aiMap["value"].(interface{}).(string)
		mappedai := v4.ConnectionSideAdditionalInfo{Key: key, Value: value}
		mappedaiArray[i] = mappedai
		i++
	}
	return mappedaiArray
}

func invitationToFabric(invitationRequest []interface{}) v4.Invitation {

	mappedI := v4.Invitation{}
	for _, i := range invitationRequest {
		iMap := i.(map[string]interface{})
		iType := iMap["type"].(interface{}).(string)
		uuid := iMap["uuid"].(interface{}).(string)
		state := iMap["state"].(interface{}).(string)
		email := iMap["email"].(interface{}).(string)
		expiry := iMap["expiry"].(interface{}).(string)
		mappedI = v4.Invitation{Type_: iType, Uuid: uuid, State: state, Email: email, Expiry: expiry}
	}
	return mappedI
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

		//TODO need to figure out, if we need to map it
		//additionalInfoRaw := accessPointMap["additional_info"].(interface{}).(string)

		//TODO I do not see uuid in the contract need to verify as gateway based connection request only need uuid
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
			accessPoint = v4.AccessPoint{Type_: &apt, Port: &p, LinkProtocol: &slp, AuthenticationKey: authenticationKey,
				Profile: &ssp, Location: &sl, ProviderConnectionId: providerConnectionId, SellerRegion: sellerRegion, PeeringType: &peeringType}
		} else {
			accessPoint = v4.AccessPoint{Type_: &apt, Port: &p, LinkProtocol: &slp, AuthenticationKey: authenticationKey,
				Profile: &ssp, Gateway: &mappedGWr, Location: &sl, ProviderConnectionId: providerConnectionId, SellerRegion: sellerRegion}
		}
	}
	return accessPoint
}

func gatewayToFabric(gatewayRequest []interface{}) v4.VirtualGateway {
	gatewayMapped := v4.VirtualGateway{}
	for _, gwr := range gatewayRequest {
		gwrMap := gwr.(map[string]interface{})
		gwtype := gwrMap["type"].(interface{}).(string)
		gwName := gwrMap["name"].(interface{}).(string)
		gwsl := locationNoIbxToFabric(gwrMap["location"].(interface{}).(*schema.Set).List())
		pr := projectToFabric(gwrMap["project"].(interface{}).(*schema.Set).List())
		gatewayMapped = v4.VirtualGateway{Type_: gwtype, Name: gwName, Location: &gwsl, Project: &pr}
	}
	return gatewayMapped
}

func projectToFabric(projectRequest []interface{}) v4.Project {
	mappedPr := v4.Project{}
	for _, pr := range projectRequest {
		prMap := pr.(map[string]interface{})
		projectId := prMap["project_id"].(interface{}).(string)
		mappedPr = v4.Project{ProjectId: projectId}
	}
	return mappedPr

}

func notificationToFabric(schemaNotifications []interface{}) []v4.SimplifiedNotification {
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
	order := v4.Order{}
	for _, o := range schemaOrder {
		orderMap := o.(map[string]interface{})
		purchaseOrderNumber := orderMap["purchase_order_number"]
		order = v4.Order{PurchaseOrderNumber: purchaseOrderNumber.(string)}
	}
	return order
}

func linkProtocolToFabric(linkProtocolList []interface{}) v4.SimplifiedLinkProtocol {
	slp := v4.SimplifiedLinkProtocol{}
	for _, lp := range linkProtocolList {
		lpMap := lp.(map[string]interface{})
		lpType := lpMap["type"].(interface{}).(string)
		lpVlanSTag := lpMap["vlan_s_tag"].(interface{}).(int)
		lpt := v4.LinkProtocolType(lpType)
		slp = v4.SimplifiedLinkProtocol{Type_: &lpt, VlanSTag: int32(lpVlanSTag)}
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

func locationNoIbxToFabric(locationList []interface{}) v4.SimplifiedLocationWithoutIbx {
	sl := v4.SimplifiedLocationWithoutIbx{}
	for _, ll := range locationList {
		llMap := ll.(map[string]interface{})
		href := llMap["href"].(interface{}).(string)
		metroName := llMap["metro_name"].(interface{}).(string)
		region := llMap["region"].(interface{}).(string)
		mc := llMap["metro_code"].(interface{}).(string)
		sl = v4.SimplifiedLocationWithoutIbx{Href: href, MetroCode: mc, Region: region, MetroName: metroName}
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
		mappedAccount["ucm_id"] = account.UcmId
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
		mappedOperation["operational_status"] = operation.OperationalStatus
		// TODO mappedOperation["errors"] = operation.Errors
		mappedOperation["op_status_changed_at"] = operation.OpStatusChangedAt.String()
		mappedOperations = append(mappedOperations, mappedOperation)
	}
	operationSet := schema.NewSet(
		schema.HashResource(createOperationRes),
		mappedOperations,
	)
	return operationSet
}

//Set
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

//Set
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
		mappedRedundancy["group"] = int(redundancy.Group)
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

//Set
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

//Set TODO not fully implemented
func connectionSideToTerra(connectionSide *v4.ConnectionSide) *schema.Set {
	connectionSides := []*v4.ConnectionSide{connectionSide}
	mappedConnectionSides := make([]interface{}, 0)
	for _, connectionSide := range connectionSides {
		mappedConnectionSide := make(map[string]interface{})
		//mappedConnectionSide["serviceToken"] = connectionSide.ServiceToken
		mappedConnectionSide["access_point"] = accessPointToTerra(connectionSide.AccessPoint)
		//mappedConnectionSide["additionalInfo"] = connectionSide.AdditionalInfo
		mappedConnectionSides = append(mappedConnectionSides, mappedConnectionSide)
	}
	connectionSideSet := schema.NewSet(
		schema.HashResource(createFabricConnectionSideRes),
		mappedConnectionSides,
	)
	return connectionSideSet
}

//Set
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

		//mappedAccessPoint["gateway"] = accessPoint.Gateway
		if accessPoint.LinkProtocol != nil {
			mappedAccessPoint["link_protocol"] = linkedProtocolToTerra(accessPoint.LinkProtocol)
		}

		//mappedAccessPoint["virtualDevice"] = accessPoint.VirtualDevice
		//mappedAccessPoint["interface"] = accessPoint.Interface_
		//mappedAccessPoint["sellerRegion"] = accessPoint.SellerRegion
		//mappedAccessPoint["peeringType"] = accessPoint.PeeringType
		mappedAccessPoint["authentication_key"] = accessPoint.AuthenticationKey
		//mappedAccessPoint["routingProtocols"] = accessPoint.RoutingProtocols
		//mappedAccessPoint["additionalInfo"] = accessPoint.AdditionalInfo
		//mappedAccessPoint["providerConnectionId"] = accessPoint.ProviderConnectionId
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
		mappedLinkedProtocol["unit"] = string(linkedProtocol.Unit)
		mappedLinkedProtocol["vni"] = string(linkedProtocol.Vni)
		mappedLinkedProtocol["int_unit"] = string(linkedProtocol.IntUnit)
		mappedLinkedProtocols = append(mappedLinkedProtocols, mappedLinkedProtocol)
	}
	linkedProtocolSet := schema.NewSet(
		schema.HashResource(createAccessPointLinkProtocolSchRes),
		mappedLinkedProtocols)
	return linkedProtocolSet
}

//Set - No full implementation
func simplifiedServiceProfileToTerra(profile *v4.SimplifiedServiceProfile) *schema.Set {
	profiles := []*v4.SimplifiedServiceProfile{profile}
	mappedProfiles := make([]interface{}, 0)
	for _, profile := range profiles {
		mappedProfile := make(map[string]interface{})
		mappedProfile["href"] = profile.Href
		mappedProfile["type"] = string(*profile.Type_)
		mappedProfile["name"] = profile.Name
		mappedProfile["uuid"] = profile.Uuid
		mappedProfile["access_point_type_configs"] = accessPointTypeToTerra(profile.AccessPointTypeConfigs)
		mappedProfiles = append(mappedProfiles, mappedProfile)
	}
	profileSet := schema.NewSet(
		schema.HashResource(createLocationRes),
		mappedProfiles,
	)
	return profileSet
}

func accessPointTypeToTerra(spAccessPointTypes []v4.ServiceProfileAccessPointType) []interface{} {
	mappedSpAccessPointTypes := make([]interface{}, len(spAccessPointTypes))
	for index, spAccessPointType := range spAccessPointTypes {
		mappedSpAccessPointTypes[index] = map[string]interface{}{
			"type": string(*spAccessPointType.Type_),
			"uuid": spAccessPointType.Uuid,
		}
	}
	return mappedSpAccessPointTypes
}

//TODO this needs full schema implementation
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

func stringToTime(attribute string, format string, date string) time.Time {
	defaultFormat := "2006-01-02T10:30:00Z"
	if format != "" {
		defaultFormat = format
	}
	t, err := time.Parse(defaultFormat, date)
	if err != nil {
		fmt.Errorf(" Error while parsing date %s for the format %s , for the attribute %s", format, date, attribute)
	}
	return t
}

func getUpdateRequest(conn v4.Connection, d *schema.ResourceData) (v4.ConnectionChangeOperation, error) {
	changeOps := v4.ConnectionChangeOperation{}
	existingName := conn.Name
	existingBandwidth := int(conn.Bandwidth)
	updateNameVal := d.Get("name").(string)
	updateBandwidthVal := d.Get("bandwidth").(int)

	log.Printf("Update Name Request %s, Update Bandwidth Request %d ", updateNameVal, updateBandwidthVal)
	if existingName != updateNameVal {
		changeOps = v4.ConnectionChangeOperation{Op: "replace", Path: "/name", Value: updateNameVal}
	} else if existingBandwidth != updateBandwidthVal {
		changeOps = v4.ConnectionChangeOperation{Op: "replace", Path: "/bandwidth", Value: updateBandwidthVal}
	} else {
		return changeOps, fmt.Errorf(" Nothing to update for the connection %s ", existingName)
	}
	return changeOps, nil
}
