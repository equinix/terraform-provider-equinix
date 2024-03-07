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

func OrderGoToTerraform(order *fabricv4.Order) *schema.Set {
	if order == nil {
		return nil
	}
	mappedOrder := make(map[string]interface{})
	mappedOrder["purchase_order_number"] = order.GetPurchaseOrderNumber()
	mappedOrder["billing_tier"] = order.GetBillingTier()
	mappedOrder["order_id"] = order.GetOrderId()
	mappedOrder["order_number"] = order.GetOrderNumber()
	orderSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: OrderSch()}),
		[]interface{}{mappedOrder},
	)
	return orderSet
}

func AccountGoToTerraform[Account *fabricv4.SimplifiedAccount | *v4.AllOfServiceProfileAccount](account Account) *schema.Set {
	if account == nil {
		return nil
	}
	var mappedAccount map[string]interface{}
	switch any(account).(type) {
	case *fabricv4.SimplifiedAccount:
		simplifiedAccount := any(account).(*fabricv4.SimplifiedAccount)
		mappedAccount = map[string]interface{}{
			"account_number":           int(simplifiedAccount.GetAccountNumber()),
			"account_name":             simplifiedAccount.GetAccountName(),
			"org_id":                   int(simplifiedAccount.GetOrgId()),
			"organization_name":        simplifiedAccount.GetOrganizationName(),
			"global_org_id":            simplifiedAccount.GetGlobalOrgId(),
			"global_organization_name": simplifiedAccount.GetGlobalOrganizationName(),
			"global_cust_id":           simplifiedAccount.GetGlobalCustId(),
			"ucm_id":                   simplifiedAccount.GetUcmId(),
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
		notificationType, _ := fabricv4.NewSimplifiedNotificationTypeFromValue(notificationMap["type"].(string))
		interval := notificationMap["send_interval"].(*string)
		emailsRaw := notificationMap["emails"].([]interface{})
		emails := converters.IfArrToStringArr(emailsRaw)
		notifications[index] = fabricv4.SimplifiedNotification{
			Type:         *notificationType,
			SendInterval: interval,
			Emails:       emails,
		}
	}
	return notifications
}

func NotificationsGoToTerraform(notifications []fabricv4.SimplifiedNotification) []map[string]interface{} {
	if notifications == nil {
		return nil
	}
	mappedNotifications := make([]map[string]interface{}, len(notifications))
	for index, notification := range notifications {
		mappedNotifications[index] = map[string]interface{}{
			"type":          notification.GetType(),
			"send_interval": notification.GetSendInterval(),
			"emails":        notification.GetEmails(),
		}
	}
	return mappedNotifications
}

func LocationTerraformToGo(locationList []interface{}) *fabricv4.SimplifiedLocation {
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

func LocationGoToTerraform(location *fabricv4.SimplifiedLocation) *schema.Set {
	if location == nil {
		return nil
	}
	mappedLocations := make(map[string]interface{})
	mappedLocations["region"] = location.GetRegion()
	mappedLocations["metro_name"] = location.GetMetroName()
	mappedLocations["metro_code"] = location.GetMetroCode()
	mappedLocations["ibx"] = location.GetIbx()

	locationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: LocationSch()}),
		[]interface{}{mappedLocations},
	)
	return locationSet
}

func LocationWithoutIBXTerraformToGo(locationList []interface{}) v4.SimplifiedLocationWithoutIbx {
	sl := v4.SimplifiedLocationWithoutIbx{}
	for _, ll := range locationList {
		llMap := ll.(map[string]interface{})
		mc := llMap["metro_code"].(string)
		sl = v4.SimplifiedLocationWithoutIbx{MetroCode: mc}
	}
	return sl
}

func LocationWithoutIBXGoToTerraform(location *fabricv4.SimplifiedLocationWithoutIBX) *schema.Set {
	mappedLocation := map[string]interface{}{
		"region":     location.GetRegion(),
		"metro_name": location.GetMetroName(),
		"metro_code": location.GetMetroCode(),
	}

	locationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: LocationSch()}),
		[]interface{}{mappedLocation},
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

func ProjectGoToTerraform(project *fabricv4.Project) *schema.Set {
	if project == nil {
		return nil
	}
	mappedProject := make(map[string]interface{})
	mappedProject["project_id"] = project.GetProjectId()
	projectSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: ProjectSch()}),
		[]interface{}{mappedProject})
	return projectSet
}

func ChangeLogGoToTerraform[ChangeLog *fabricv4.Changelog | *v4.AllOfServiceProfileChangeLog](changeLog ChangeLog) *schema.Set {
	if changeLog == nil {
		return nil
	}
	var mappedChangeLog map[string]interface{}
	switch any(changeLog).(type) {
	case *fabricv4.Changelog:
		baseChangeLog := any(changeLog).(*fabricv4.Changelog)
		mappedChangeLog = map[string]interface{}{
			"created_by":           baseChangeLog.GetCreatedBy(),
			"created_by_full_name": baseChangeLog.GetCreatedByFullName(),
			"created_by_email":     baseChangeLog.GetCreatedByEmail(),
			"created_date_time":    baseChangeLog.GetCreatedDateTime().String(),
			"updated_by":           baseChangeLog.GetUpdatedBy(),
			"updated_by_full_name": baseChangeLog.GetUpdatedByFullName(),
			"updated_date_time":    baseChangeLog.GetUpdatedDateTime().String(),
			"deleted_by":           baseChangeLog.GetDeletedBy(),
			"deleted_by_full_name": baseChangeLog.GetDeletedByFullName(),
			"deleted_by_email":     baseChangeLog.GetDeletedByEmail(),
			"deleted_date_time":    baseChangeLog.GetDeletedDateTime().String(),
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

func ErrorGoToTerraform(errors []fabricv4.Error) []interface{} {
	if errors == nil || len(errors) == 0 {
		return nil
	}
	mappedErrors := make([]interface{}, len(errors))
	for index, mError := range errors {
		mappedErrors[index] = map[string]interface{}{
			"error_code":      mError.GetErrorCode(),
			"error_message":   mError.GetErrorMessage(),
			"correlation_id":  mError.GetCorrelationId(),
			"details":         mError.GetDetails(),
			"help":            mError.GetHelp(),
			"additional_info": ErrorAdditionalInfoGoToTerraform(mError.GetAdditionalInfo()),
		}
	}
	return mappedErrors
}

func ErrorAdditionalInfoGoToTerraform(additionalInfol []fabricv4.PriceErrorAdditionalInfo) []interface{} {
	if additionalInfol == nil {
		return nil
	}
	mappedAdditionalInfol := make([]interface{}, len(additionalInfol))
	for index, additionalInfo := range additionalInfol {
		mappedAdditionalInfol[index] = map[string]interface{}{
			"property": additionalInfo.GetProperty(),
			"reason":   additionalInfo.GetReason(),
		}
	}
	return mappedAdditionalInfol
}
