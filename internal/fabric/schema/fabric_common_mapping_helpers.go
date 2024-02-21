package schema

import (
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func OrderToFabric(schemaOrder []interface{}) v4.Order {
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

func OrderToTerra(order *v4.Order) *schema.Set {
	if order == nil {
		return nil
	}
	orders := []*v4.Order{order}
	mappedOrders := make([]interface{}, len(orders))
	for _, order := range orders {
		mappedOrder := make(map[string]interface{})
		mappedOrder["purchase_order_number"] = order.PurchaseOrderNumber
		mappedOrder["billing_tier"] = order.BillingTier
		mappedOrder["order_id"] = order.OrderId
		mappedOrder["order_number"] = order.OrderNumber
		mappedOrders = append(mappedOrders, mappedOrder)
	}
	orderSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: OrderSch()}),
		mappedOrders,
	)
	return orderSet
}

func AccountToTerra[Account *v4.SimplifiedAccount | *v4.AllOfServiceProfileAccount](account Account) *schema.Set {
	if account == nil {
		return nil
	}
	var mappedAccount map[string]interface{}
	switch any(account).(type) {
	case *v4.SimplifiedAccount:
		simplifiedAccount := any(account).(*v4.SimplifiedAccount)
		mappedAccount = map[string]interface{}{
			"account_number":           int(simplifiedAccount.AccountNumber),
			"account_name":             simplifiedAccount.AccountName,
			"org_id":                   int(simplifiedAccount.OrgId),
			"organization_name":        simplifiedAccount.OrganizationName,
			"global_org_id":            simplifiedAccount.GlobalOrgId,
			"global_organization_name": simplifiedAccount.GlobalOrganizationName,
			"global_cust_id":           simplifiedAccount.GlobalCustId,
			"ucm_id":                   simplifiedAccount.UcmId,
		}
	case *v4.AllOfServiceProfileAccount:
		allSPAccount := any(account).(*v4.AllOfServiceProfileAccount)
		mappedAccount = map[string]interface{}{
			"account_number":           int(allSPAccount.AccountNumber),
			"account_name":             allSPAccount.AccountName,
			"org_id":                   int(allSPAccount.OrgId),
			"organization_name":        allSPAccount.OrganizationName,
			"global_org_id":            allSPAccount.GlobalOrgId,
			"global_organization_name": allSPAccount.GlobalOrganizationName,
			"global_cust_id":           allSPAccount.GlobalCustId,
			"ucm_id":                   allSPAccount.UcmId,
		}
	default:
		return nil
	}

	accountSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: AccountSch()}),
		[]interface{}{mappedAccount},
	)

	return accountSet
}

func NotificationsToFabric(schemaNotifications []interface{}) []v4.SimplifiedNotification {
	if schemaNotifications == nil {
		return []v4.SimplifiedNotification{}
	}
	var notifications []v4.SimplifiedNotification
	for _, n := range schemaNotifications {
		ntype := n.(map[string]interface{})["type"].(string)
		interval := n.(map[string]interface{})["send_interval"].(string)
		emailsRaw := n.(map[string]interface{})["emails"].([]interface{})
		emails := converters.IfArrToStringArr(emailsRaw)
		notifications = append(notifications, v4.SimplifiedNotification{
			Type_:        ntype,
			SendInterval: interval,
			Emails:       emails,
		})
	}
	return notifications
}

func NotificationsToTerra(notifications []v4.SimplifiedNotification) []map[string]interface{} {
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

func LocationToFabric(locationList []interface{}) v4.SimplifiedLocation {
	sl := v4.SimplifiedLocation{}
	for _, ll := range locationList {
		llMap := ll.(map[string]interface{})
		metroName := llMap["metro_name"]
		var metroNamestr string
		if metroName != nil {
			metroNamestr = metroName.(string)
		}
		region := llMap["region"].(string)
		mc := llMap["metro_code"].(string)
		ibx := llMap["ibx"].(string)
		sl = v4.SimplifiedLocation{MetroCode: mc, Region: region, Ibx: ibx, MetroName: metroNamestr}
	}
	return sl
}

func LocationToTerra(location *v4.SimplifiedLocation) *schema.Set {
	if location == nil {
		return nil
	}
	mappedLocations := make(map[string]interface{})
	mappedLocations["region"] = location.Region
	mappedLocations["metro_name"] = location.MetroName
	mappedLocations["metro_code"] = location.MetroCode
	mappedLocations["ibx"] = location.Ibx

	locationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: LocationSch()}),
		[]interface{}{mappedLocations},
	)
	return locationSet
}

func LocationWithoutIBXToFabric(locationList []interface{}) v4.SimplifiedLocationWithoutIbx {
	sl := v4.SimplifiedLocationWithoutIbx{}
	for _, ll := range locationList {
		llMap := ll.(map[string]interface{})
		mc := llMap["metro_code"].(string)
		sl = v4.SimplifiedLocationWithoutIbx{MetroCode: mc}
	}
	return sl
}

func LocationWithoutIBXToTerra(location *v4.SimplifiedLocationWithoutIbx) *schema.Set {
	locations := []*v4.SimplifiedLocationWithoutIbx{location}
	mappedLocations := make([]interface{}, len(locations))
	for i, location := range locations {
		mappedLocations[i] = map[string]interface{}{
			"region":     location.Region,
			"metro_name": location.MetroName,
			"metro_code": location.MetroCode,
		}
	}
	locationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: LocationSch()}),
		mappedLocations,
	)
	return locationSet
}

func ProjectToFabric(projectRequest []interface{}) v4.Project {
	if len(projectRequest) == 0 {
		return v4.Project{}
	}
	mappedPr := v4.Project{}
	for _, pr := range projectRequest {
		prMap := pr.(map[string]interface{})
		projectId := prMap["project_id"].(string)
		mappedPr = v4.Project{ProjectId: projectId}
	}
	return mappedPr
}

func ProjectToTerra(project *v4.Project) *schema.Set {
	if project == nil {
		return nil
	}
	projects := []*v4.Project{project}
	mappedProjects := make([]interface{}, len(projects))
	for _, project := range projects {
		mappedProject := make(map[string]interface{})
		mappedProject["project_id"] = project.ProjectId
		mappedProjects = append(mappedProjects, mappedProject)
	}
	projectSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: ProjectSch()}),
		mappedProjects)
	return projectSet
}

func ChangeLogToTerra[ChangeLog *v4.Changelog | *v4.AllOfServiceProfileChangeLog](changeLog ChangeLog) *schema.Set {
	if changeLog == nil {
		return nil
	}
	var mappedChangeLog map[string]interface{}
	switch any(changeLog).(type) {
	case *v4.Changelog:
		baseChangeLog := any(changeLog).(*v4.Changelog)
		mappedChangeLog = map[string]interface{}{
			"created_by":           baseChangeLog.CreatedBy,
			"created_by_full_name": baseChangeLog.CreatedByFullName,
			"created_by_email":     baseChangeLog.CreatedByEmail,
			"created_date_time":    baseChangeLog.CreatedDateTime.String(),
			"updated_by":           baseChangeLog.UpdatedBy,
			"updated_by_full_name": baseChangeLog.UpdatedByFullName,
			"updated_date_time":    baseChangeLog.UpdatedDateTime.String(),
			"deleted_by":           baseChangeLog.DeletedBy,
			"deleted_by_full_name": baseChangeLog.DeletedByFullName,
			"deleted_by_email":     baseChangeLog.DeletedByEmail,
			"deleted_date_time":    baseChangeLog.DeletedDateTime.String(),
		}
	case *v4.AllOfServiceProfileChangeLog:
		allOfChangeLog := any(changeLog).(*v4.AllOfServiceProfileChangeLog)
		mappedChangeLog = map[string]interface{}{
			"created_by":           allOfChangeLog.CreatedBy,
			"created_by_full_name": allOfChangeLog.CreatedByFullName,
			"created_by_email":     allOfChangeLog.CreatedByEmail,
			"created_date_time":    allOfChangeLog.CreatedDateTime.String(),
			"updated_by":           allOfChangeLog.UpdatedBy,
			"updated_by_full_name": allOfChangeLog.UpdatedByFullName,
			"updated_date_time":    allOfChangeLog.UpdatedDateTime.String(),
			"deleted_by":           allOfChangeLog.DeletedBy,
			"deleted_by_full_name": allOfChangeLog.DeletedByFullName,
			"deleted_by_email":     allOfChangeLog.DeletedByEmail,
			"deleted_date_time":    allOfChangeLog.DeletedDateTime.String(),
		}
	default:
		return nil
	}
	changeLogSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: ChangeLogSch()}),
		[]interface{}{mappedChangeLog},
	)
	return changeLogSet
}

func ErrorToTerra(errors []v4.ModelError) []interface{} {
	if errors == nil {
		return nil
	}
	mappedErrors := make([]interface{}, len(errors))
	for index, mError := range errors {
		mappedErrors[index] = map[string]interface{}{
			"error_code":      mError.ErrorCode,
			"error_message":   mError.ErrorMessage,
			"correlation_id":  mError.CorrelationId,
			"details":         mError.Details,
			"help":            mError.Help,
			"additional_info": ErrorAdditionalInfoToTerra(mError.AdditionalInfo),
		}
	}
	return mappedErrors
}

func ErrorAdditionalInfoToTerra(additionalInfol []v4.PriceErrorAdditionalInfo) []interface{} {
	if additionalInfol == nil {
		return nil
	}
	mappedAdditionalInfol := make([]interface{}, len(additionalInfol))
	for index, additionalInfo := range additionalInfol {
		mappedAdditionalInfol[index] = map[string]interface{}{
			"property": additionalInfo.Property,
			"reason":   additionalInfo.Reason,
		}
	}
	return mappedAdditionalInfol
}
