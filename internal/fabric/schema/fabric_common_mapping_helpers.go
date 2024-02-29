package schema

import (
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func OrderTerraformToGo(orderTerraform []interface{}) *fabricv4.Order {
	if orderTerraform == nil || len(orderTerraform) == 0 {
		return nil
	}
	var order *fabricv4.Order

	orderMap := orderTerraform[0].(map[string]interface{})
	purchaseOrderNumber := orderMap["purchase_order_number"].(*string)
	billingTier := orderMap["billing_tier"].(*string)
	orderId := orderMap["order_id"].(*string)
	orderNumber := orderMap["order_number"].(*string)
	order = &fabricv4.Order{PurchaseOrderNumber: purchaseOrderNumber, BillingTier: billingTier, OrderId: orderId, OrderNumber: orderNumber}

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

func NotificationsTerraformToGo(notificationsTerraform []interface{}) []fabricv4.SimplifiedNotification {
	if notificationsTerraform == nil || len(notificationsTerraform) == 0 {
		return nil
	}
	notifications := make([]fabricv4.SimplifiedNotification, len(notificationsTerraform))
	for index, notification := range notificationsTerraform {
		notificationMap := notification.(map[string]interface{})
		notificationType := fabricv4.SimplifiedNotificationType(notificationMap["type"].(string))
		interval := notificationMap["send_interval"].(*string)
		emailsRaw := notificationMap["emails"].([]interface{})
		emails := converters.IfArrToStringArr(emailsRaw)
		notifications[index] = fabricv4.SimplifiedNotification{
			Type:         notificationType,
			SendInterval: interval,
			Emails:       emails,
		}
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

func LocationToFabric(locationList []interface{}) *fabricv4.SimplifiedLocation {
	if locationList == nil || len(locationList) == 0 {
		return nil
	}

	var location *fabricv4.SimplifiedLocation
	locationListMap := locationList[0].(map[string]interface{})
	metroName := locationListMap["metro_name"].(*string)
	region := locationListMap["region"].(*string)
	mc := locationListMap["metro_code"].(*string)
	ibx := locationListMap["ibx"].(*string)
	location = &fabricv4.SimplifiedLocation{MetroCode: mc, Region: region, Ibx: ibx, MetroName: metroName}
	return location
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

func ProjectTerraformToGo(projectTerraform []interface{}) *fabricv4.Project {
	if projectTerraform == nil || len(projectTerraform) == 0 {
		return nil
	}
	var project *fabricv4.Project
	projectMap := projectTerraform[0].(map[string]interface{})
	projectId := projectMap["project_id"].(string)
	project.ProjectId = projectId

	return project
}

func ProjectToTerra(project *v4.Project) *schema.Set {
	if project == nil {
		return nil
	}
	mappedProject := make(map[string]interface{})
	mappedProject["project_id"] = project.ProjectId
	projectSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: ProjectSch()}),
		[]interface{}{mappedProject})
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
