package equinix

import (
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func customFieldFabricSpToTerra(customFieldl []v4.CustomField) []map[string]interface{} {
	if customFieldl == nil {
		return nil
	}
	mappedCustomFieldl := make([]map[string]interface{}, len(customFieldl))
	for index, customField := range customFieldl {
		mappedCustomFieldl[index] = map[string]interface{}{
			"label":       customField.Label,
			"description": customField.Description,
			"required":    customField.Required,
			"data_type":   customField.DataType,
			"options":     customField.Options,
		}
	}
	return mappedCustomFieldl
}

func processStepFabricSpToTerra(processStepl []v4.ProcessStep) []map[string]interface{} {
	if processStepl == nil {
		return nil
	}
	mappedProcessStepl := make([]map[string]interface{}, len(processStepl))
	for index, processStep := range processStepl {
		mappedProcessStepl[index] = map[string]interface{}{
			"title":       processStep.Title,
			"sub_title":   processStep.SubTitle,
			"description": processStep.Description,
		}
	}
	return mappedProcessStepl
}

func marketingInfoMappingToTerra(mkinfo *v4.MarketingInfo) *schema.Set {
	if mkinfo == nil {
		return nil
	}
	mkinfos := []*v4.MarketingInfo{mkinfo}
	mappedMkInfos := make([]interface{}, 0)
	for _, mkinfo := range mkinfos {
		mappedMkInfo := make(map[string]interface{})
		mappedMkInfo["logo"] = mkinfo.Logo
		mappedMkInfo["promotion"] = mkinfo.Promotion
		processStepl := processStepFabricSpToTerra(mkinfo.ProcessSteps)
		if processStepl != nil {
			mappedMkInfo["process_step"] = processStepl
		}
		mappedMkInfos = append(mappedMkInfos, mappedMkInfo)
	}
	marketingInfoSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: equinix_schema.OrderSch()}),
		mappedMkInfos,
	)
	return marketingInfoSet
}

func accessPointColoFabricSpToTerra(accessPointColol []v4.ServiceProfileAccessPointColo) []map[string]interface{} {
	if accessPointColol == nil {
		return nil
	}
	mappedAccessPointColol := make([]map[string]interface{}, len(accessPointColol))
	for index, accessPointColo := range accessPointColol {
		mappedAccessPointColol[index] = map[string]interface{}{
			"type":                      accessPointColo.Type_,
			"uuid":                      accessPointColo.Uuid,
			"location":                  locationToTerra(accessPointColo.Location),
			"seller_region":             accessPointColo.SellerRegion,
			"seller_region_description": accessPointColo.SellerRegionDescription,
			"cross_connect_id":          accessPointColo.CrossConnectId,
		}
	}
	return mappedAccessPointColol
}

func serviceMetroFabricSpToTerra(serviceMetrol []v4.ServiceMetro) []map[string]interface{} {
	if serviceMetrol == nil {
		return nil
	}
	mappedServiceMetrol := make([]map[string]interface{}, len(serviceMetrol))
	for index, serviceMetro := range serviceMetrol {
		mappedServiceMetrol[index] = map[string]interface{}{
			"code":           serviceMetro.Code,
			"name":           serviceMetro.Name,
			"ibxs":           serviceMetro.Ibxs,
			"in_trail":       serviceMetro.InTrail,
			"display_name":   serviceMetro.DisplayName,
			"seller_regions": serviceMetro.SellerRegions,
		}
	}
	return mappedServiceMetrol
}

func serviceProfileAccountFabricSpToTerra(account *v4.AllOfServiceProfileAccount) *schema.Set {
	if account == nil {
		return nil
	}
	accounts := []*v4.AllOfServiceProfileAccount{account}
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
		mappedAccount["ucm_id"] = account.UcmId
		mappedAccounts = append(mappedAccounts, mappedAccount)
	}

	accountSet := schema.NewSet(
		schema.HashResource(readSpAccountRes),
		mappedAccounts,
	)
	return accountSet
}

func tagsFabricSpToTerra(tagsl *[]string) []interface{} {
	if tagsl == nil {
		return nil
	}
	mappedTags := make([]interface{}, 0)
	for _, tag := range *tagsl {
		mappedTags = append(mappedTags, tag)
	}
	return mappedTags
}

func allowedEmailsFabricSpToTerra(allowedemails []string) []interface{} {
	if allowedemails == nil {
		return nil
	}
	mappedEmails := make([]interface{}, len(allowedemails))
	for index, email := range allowedemails {
		mappedEmails[index] = email
	}
	return mappedEmails
}

func allOfServiceProfileChangeLogToTerra(changeLog *v4.AllOfServiceProfileChangeLog) *schema.Set {
	if changeLog == nil {
		return nil
	}
	changeLogs := []*v4.AllOfServiceProfileChangeLog{changeLog}
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
		schema.HashResource(readAllOfServiceProfileChangeLogRes),
		mappedChangeLogs,
	)
	return changeLogSet
}

func accessPointTypeConfigsToFabric(schemaAccessPointTypeConfigs []interface{}) []v4.ServiceProfileAccessPointType {
	if schemaAccessPointTypeConfigs == nil {
		return []v4.ServiceProfileAccessPointType{}
	}
	var accessPointTypeConfigs []v4.ServiceProfileAccessPointType
	for _, accessPoint := range schemaAccessPointTypeConfigs {
		spType := v4.ServiceProfileAccessPointTypeEnum(accessPoint.(map[string]interface{})["type"].(string))
		spConnectionRedundancyRequired := accessPoint.(map[string]interface{})["connection_redundancy_required"].(bool)
		spAllowBandwidthAutoApproval := accessPoint.(map[string]interface{})["allow_bandwidth_auto_approval"].(bool)
		spAllowRemoteConnections := accessPoint.(map[string]interface{})["allow_remote_connections"].(bool)
		spConnectionLabel := accessPoint.(map[string]interface{})["connection_label"].(string)
		spEnableAutoGenerateServiceKey := accessPoint.(map[string]interface{})["enable_auto_generate_service_key"].(bool)
		spBandwidthAlertThreshold := accessPoint.(map[string]interface{})["bandwidth_alert_threshold"].(float64)
		spAllowCustomBandwidth := accessPoint.(map[string]interface{})["allow_custom_bandwidth"].(bool)

		var spApiConfig v4.ApiConfig
		if accessPoint.(map[string]interface{})["api_config"] != nil {
			apiConfig := accessPoint.(map[string]interface{})["api_config"].(interface{}).(*schema.Set).List()
			spApiConfig = apiConfigToFabric(apiConfig)
		}

		var spAuthenticationKey v4.AuthenticationKey
		if accessPoint.(map[string]interface{})["authentication_key"] != nil {
			authenticationKey := accessPoint.(map[string]interface{})["authentication_key"].(interface{}).(*schema.Set).List()
			spAuthenticationKey = authenticationKeyToFabric(authenticationKey)
		}

		supportedBandwidthsRaw := accessPoint.(map[string]interface{})["supported_bandwidths"].([]interface{})
		spSupportedBandwidths := converters.ListToInt32List(supportedBandwidthsRaw)

		accessPointTypeConfigs = append(accessPointTypeConfigs, v4.ServiceProfileAccessPointType{
			Type_:                        &spType,
			ConnectionRedundancyRequired: spConnectionRedundancyRequired,
			AllowBandwidthAutoApproval:   spAllowBandwidthAutoApproval,
			AllowRemoteConnections:       spAllowRemoteConnections,
			ConnectionLabel:              spConnectionLabel,
			EnableAutoGenerateServiceKey: spEnableAutoGenerateServiceKey,
			BandwidthAlertThreshold:      spBandwidthAlertThreshold,
			AllowCustomBandwidth:         spAllowCustomBandwidth,
			ApiConfig:                    &spApiConfig,
			AuthenticationKey:            &spAuthenticationKey,
			SupportedBandwidths:          &spSupportedBandwidths,
		})
	}
	return accessPointTypeConfigs
}

func apiConfigToFabric(apiConfigl []interface{}) v4.ApiConfig {
	if apiConfigl == nil {
		return v4.ApiConfig{}
	}
	var apiConfigs v4.ApiConfig
	for _, apiCongig := range apiConfigl {
		psApiAvailable := apiCongig.(map[string]interface{})["api_available"].(interface{}).(bool)
		psEquinixManagedVlan := apiCongig.(map[string]interface{})["equinix_managed_vlan"].(interface{}).(bool)
		psBandwidthFromApi := apiCongig.(map[string]interface{})["bandwidth_from_api"].(interface{}).(bool)
		psIntegrationId := apiCongig.(map[string]interface{})["integration_id"].(interface{}).(string)
		psEquinixManagedPort := apiCongig.(map[string]interface{})["equinix_managed_port"].(interface{}).(bool)
		apiConfigs = v4.ApiConfig{
			ApiAvailable:       psApiAvailable,
			EquinixManagedVlan: psEquinixManagedVlan,
			BandwidthFromApi:   psBandwidthFromApi,
			IntegrationId:      psIntegrationId,
			EquinixManagedPort: psEquinixManagedPort,
		}
	}
	return apiConfigs
}

func authenticationKeyToFabric(authenticationKeyl []interface{}) v4.AuthenticationKey {
	if authenticationKeyl == nil {
		return v4.AuthenticationKey{}
	}
	var authenticationKeys v4.AuthenticationKey
	for _, authKey := range authenticationKeyl {
		psRequired := authKey.(map[string]interface{})["required"].(interface{}).(bool)
		psLabel := authKey.(map[string]interface{})["label"].(interface{}).(string)
		psDescription := authKey.(map[string]interface{})["description"].(interface{}).(string)
		authenticationKeys = v4.AuthenticationKey{
			Required:    psRequired,
			Label:       psLabel,
			Description: psDescription,
		}
	}
	return authenticationKeys
}

func customFieldsToFabric(schemaCustomField []interface{}) []v4.CustomField {
	if schemaCustomField == nil {
		return []v4.CustomField{}
	}
	var customFields []v4.CustomField
	for _, customField := range schemaCustomField {
		cfLabel := customField.(map[string]interface{})["label"].(string)
		cfDescription := customField.(map[string]interface{})["description"].(string)
		cfRequired := customField.(map[string]interface{})["required"].(bool)
		cfDataType := customField.(map[string]interface{})["data_type"].(string)
		optionsRaw := customField.(map[string]interface{})["options"].([]interface{})
		cfOptions := converters.IfArrToStringArr(optionsRaw)
		cfCaptureInEmail := customField.(map[string]interface{})["capture_in_email"].(bool)
		customFields = append(customFields, v4.CustomField{
			Label:          cfLabel,
			Description:    cfDescription,
			Required:       cfRequired,
			DataType:       cfDataType,
			Options:        cfOptions,
			CaptureInEmail: cfCaptureInEmail,
		})
	}
	return customFields
}

func marketingInfoToFabric(schemaMarketingInfo []interface{}) v4.MarketingInfo {
	if schemaMarketingInfo == nil {
		return v4.MarketingInfo{}
	}
	marketingInfos := v4.MarketingInfo{}
	for _, marketingInfo := range schemaMarketingInfo {
		miLogo := marketingInfo.(map[string]interface{})["logo"].(string)
		miPromotion := marketingInfo.(map[string]interface{})["promotion"].(bool)

		var miProcessSteps []v4.ProcessStep
		if marketingInfo.(map[string]interface{})["process_steps"] != nil {
			processStepsList := marketingInfo.(map[string]interface{})["process_steps"].([]interface{})
			miProcessSteps = processStepToFabric(processStepsList)
		}

		marketingInfos = v4.MarketingInfo{
			Logo:         miLogo,
			Promotion:    miPromotion,
			ProcessSteps: miProcessSteps,
		}
	}
	return marketingInfos
}

func processStepToFabric(processStepl []interface{}) []v4.ProcessStep {
	if processStepl == nil {
		return []v4.ProcessStep{}
	}
	var processSteps []v4.ProcessStep
	for _, processStep := range processStepl {
		psTitle := processStep.(map[string]interface{})["title"].(interface{}).(string)
		psSubTitle := processStep.(map[string]interface{})["sub_title"].(interface{}).(string)
		psDescription := processStep.(map[string]interface{})["description"].(interface{}).(string)
		processSteps = append(processSteps, v4.ProcessStep{
			Title:       psTitle,
			SubTitle:    psSubTitle,
			Description: psDescription,
		})
	}
	return processSteps
}

func portsToFabric(schemaPorts []interface{}) []v4.ServiceProfileAccessPointColo {
	if schemaPorts == nil {
		return []v4.ServiceProfileAccessPointColo{}
	}
	var ports []v4.ServiceProfileAccessPointColo
	for _, port := range schemaPorts {
		pType := port.(map[string]interface{})["type"].(string)
		pUuid := port.(map[string]interface{})["uuid"].(string)
		locationList := port.(map[string]interface{})["location"].(interface{}).(*schema.Set).List()
		pLocation := v4.SimplifiedLocation{}
		if len(locationList) != 0 {
			pLocation = locationToFabric(locationList)
		}
		pSellerRegion := port.(map[string]interface{})["seller_region"].(string)
		pSellerRegionDescription := port.(map[string]interface{})["seller_region_description"].(string)
		pCrossConnectId := port.(map[string]interface{})["cross_connect_id"].(string)
		ports = append(ports, v4.ServiceProfileAccessPointColo{
			Type_:                   pType,
			Uuid:                    pUuid,
			Location:                &pLocation,
			SellerRegion:            pSellerRegion,
			SellerRegionDescription: pSellerRegionDescription,
			CrossConnectId:          pCrossConnectId,
		})
	}
	return ports
}

func virtualDevicesToFabric(schemaVirtualDevices []interface{}) []v4.ServiceProfileAccessPointVd {
	if schemaVirtualDevices == nil {
		return []v4.ServiceProfileAccessPointVd{}
	}
	var virtualDevices []v4.ServiceProfileAccessPointVd
	for _, virtualDevice := range schemaVirtualDevices {
		vType := virtualDevice.(map[string]interface{})["type"].(string)
		vUuid := virtualDevice.(map[string]interface{})["uuid"].(string)
		locationList := virtualDevice.(map[string]interface{})["location"].(interface{}).(*schema.Set).List()
		vLocation := v4.SimplifiedLocation{}
		if len(locationList) != 0 {
			vLocation = locationToFabric(locationList)
		}
		pInterfaceUuid := virtualDevice.(map[string]interface{})["interface_uuid"].(string)
		virtualDevices = append(virtualDevices, v4.ServiceProfileAccessPointVd{
			Type_:         vType,
			Uuid:          vUuid,
			Location:      &vLocation,
			InterfaceUuid: pInterfaceUuid,
		})
	}
	return virtualDevices
}

func metrosToFabric(schemaMetros []interface{}) []v4.ServiceMetro {
	if schemaMetros == nil {
		return []v4.ServiceMetro{}
	}
	var metros []v4.ServiceMetro
	for _, metro := range schemaMetros {
		mCode := metro.(map[string]interface{})["code"].(string)
		mName := metro.(map[string]interface{})["name"].(string)
		ibxsRaw := metro.(map[string]interface{})["ibxs"].([]interface{})
		mIbxs := converters.IfArrToStringArr(ibxsRaw)
		mInTrail := metro.(map[string]interface{})["in_trail"].(bool)
		mDisplayName := metro.(map[string]interface{})["display_name"].(string)
		mSellerRegions := metro.(map[string]interface{})["seller_regions"].(map[string]string)
		metros = append(metros, v4.ServiceMetro{
			Code:          mCode,
			Name:          mName,
			Ibxs:          mIbxs,
			InTrail:       mInTrail,
			DisplayName:   mDisplayName,
			SellerRegions: mSellerRegions,
		})
	}
	return metros
}

func fabricServiceProfilesListToTerra(serviceProfiles v4.ServiceProfiles) []map[string]interface{} {
	serviceProfile := serviceProfiles.Data
	if serviceProfile == nil {
		return nil
	}
	mappedServiceProfile := make([]map[string]interface{}, len(serviceProfile))
	for index, serviceProfile := range serviceProfile {
		mappedServiceProfile[index] = map[string]interface{}{
			"href":                      serviceProfile.Href,
			"type":                      serviceProfile.Type_,
			"name":                      serviceProfile.Name,
			"uuid":                      serviceProfile.Uuid,
			"description":               serviceProfile.Description,
			"notifications":             notificationToTerra(serviceProfile.Notifications),
			"tags":                      tagsFabricSpToTerra(serviceProfile.Tags),
			"visibility":                serviceProfile.Visibility,
			"access_point_type_configs": accessPointTypeConfigToTerra(serviceProfile.AccessPointTypeConfigs),
			"custom_fields":             customFieldFabricSpToTerra(serviceProfile.CustomFields),
			"marketing_info":            marketingInfoMappingToTerra(serviceProfile.MarketingInfo),
			"ports":                     accessPointColoFabricSpToTerra(serviceProfile.Ports),
			"allowed_emails":            allowedEmailsFabricSpToTerra(serviceProfile.AllowedEmails),
			"metros":                    serviceMetroFabricSpToTerra(serviceProfile.Metros),
			"self_profile":              serviceProfile.SelfProfile,
			"state":                     serviceProfile.State,
			"account":                   serviceProfileAccountFabricSpToTerra(serviceProfile.Account),
			"project":                   projectToTerra(serviceProfile.Project),
			"change_log":                allOfServiceProfileChangeLogToTerra(serviceProfile.ChangeLog),
		}
	}
	return mappedServiceProfile
}

func serviceProfilesSearchFilterRequestToFabric(schemaServiceProfileFilterRequest []interface{}) v4.ServiceProfileSimpleExpression {
	if schemaServiceProfileFilterRequest == nil {
		return v4.ServiceProfileSimpleExpression{}
	}
	mappedFilter := v4.ServiceProfileSimpleExpression{}
	for _, s := range schemaServiceProfileFilterRequest {
		sProperty := s.(map[string]interface{})["property"].(string)
		operator := s.(map[string]interface{})["operator"].(string)
		valuesRaw := s.(map[string]interface{})["values"].([]interface{})
		values := converters.IfArrToStringArr(valuesRaw)
		mappedFilter = v4.ServiceProfileSimpleExpression{Property: sProperty, Operator: operator, Values: values}
	}
	return mappedFilter
}

func serviceProfilesSearchPaginationRequestToFabric(schemaServiceProfilePaginationRequest []interface{}) v4.PaginationRequest {
	if schemaServiceProfilePaginationRequest == nil {
		return v4.PaginationRequest{}
	}
	mappedPaginationReq := v4.PaginationRequest{}

	for _, s := range schemaServiceProfilePaginationRequest {
		sMap := s.(map[string]interface{})
		sOffset := sMap["offset"].(interface{}).(int)
		limit := sMap["limit"].(interface{}).(int)
		mappedPaginationReq = v4.PaginationRequest{Offset: int32(sOffset), Limit: int32(limit)}
	}
	return mappedPaginationReq
}

func serviceProfilesSearchSortRequestToFabric(schemaServiceProfilesSearchSortRequest []interface{}) []v4.ServiceProfileSortCriteria {
	if schemaServiceProfilesSearchSortRequest == nil {
		return []v4.ServiceProfileSortCriteria{}
	}
	var spSortCriteria []v4.ServiceProfileSortCriteria
	for _, sp := range schemaServiceProfilesSearchSortRequest {
		serviceProfileSortCriteriaMap := sp.(map[string]interface{})
		direction := serviceProfileSortCriteriaMap["direction"]
		directionCont := v4.ServiceProfileSortDirection(direction.(string))
		property := serviceProfileSortCriteriaMap["property"]
		propertyCont := v4.ServiceProfileSortBy(property.(string))
		spSortCriteria = append(spSortCriteria, v4.ServiceProfileSortCriteria{
			Direction: &directionCont,
			Property:  &propertyCont,
		})

	}
	return spSortCriteria
}
