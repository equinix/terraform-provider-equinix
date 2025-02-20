package schema

import (
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func OrderTerraformToGo(orderTerraform []interface{}) fabricv4.Order {
	if len(orderTerraform) == 0 {
		return fabricv4.Order{}
	}
	var order fabricv4.Order

	orderMap := orderTerraform[0].(map[string]interface{})
	purchaseOrderNumber := orderMap["purchase_order_number"].(string)
	billingTier := orderMap["billing_tier"].(string)
	orderId := orderMap["order_id"].(string)
	orderNumber := orderMap["order_number"].(string)
	termLength := orderMap["term_length"].(int)
	if purchaseOrderNumber != "" {
		order.SetPurchaseOrderNumber(purchaseOrderNumber)
	}
	if billingTier != "" {
		order.SetBillingTier(billingTier)
	}
	if orderId != "" {
		order.SetOrderId(orderId)
	}
	if orderNumber != "" {
		order.SetOrderNumber(orderNumber)
	}
	if termLength >= 1 {
		order.SetTermLength(int32(termLength))
	}

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
	mappedOrder["term_length"] = int(order.GetTermLength())
	orderSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: OrderSch()}),
		[]interface{}{mappedOrder},
	)
	return orderSet
}

func AccountGoToTerraform(account *fabricv4.SimplifiedAccount) *schema.Set {
	if account == nil {
		return nil
	}
	mappedAccount := map[string]interface{}{
		"account_number":           int(account.GetAccountNumber()),
		"account_name":             account.GetAccountName(),
		"org_id":                   int(account.GetOrgId()),
		"organization_name":        account.GetOrganizationName(),
		"global_org_id":            account.GetGlobalOrgId(),
		"global_organization_name": account.GetGlobalOrganizationName(),
		"global_cust_id":           account.GetGlobalCustId(),
		"ucm_id":                   account.GetUcmId(),
	}

	accountSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: AccountSch()}),
		[]interface{}{mappedAccount},
	)

	return accountSet
}

func NotificationsTerraformToGo(notificationsTerraform []interface{}) []fabricv4.SimplifiedNotification {
	if len(notificationsTerraform) == 0 {
		return nil
	}
	notifications := make([]fabricv4.SimplifiedNotification, len(notificationsTerraform))
	for index, notification := range notificationsTerraform {
		notificationMap := notification.(map[string]interface{})
		notificationType := fabricv4.SimplifiedNotificationType(notificationMap["type"].(string))
		sendInterval := notificationMap["send_interval"].(string)
		emailsRaw := notificationMap["emails"].([]interface{})
		emails := converters.IfArrToStringArr(emailsRaw)
		simplifiedNotification := fabricv4.SimplifiedNotification{}
		simplifiedNotification.SetType(notificationType)
		if sendInterval != "" {
			simplifiedNotification.SetSendInterval(sendInterval)
		}
		simplifiedNotification.SetEmails(emails)
		notifications[index] = simplifiedNotification
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
			"type":          string(notification.GetType()),
			"send_interval": notification.GetSendInterval(),
			"emails":        notification.GetEmails(),
		}
	}
	return mappedNotifications
}

func LocationTerraformToGo(locationList []interface{}) fabricv4.SimplifiedLocation {
	if len(locationList) == 0 {
		return fabricv4.SimplifiedLocation{}
	}

	var location fabricv4.SimplifiedLocation
	locationListMap := locationList[0].(map[string]interface{})
	metroName := locationListMap["metro_name"].(string)
	region := locationListMap["region"].(string)
	metroCode := locationListMap["metro_code"].(string)
	ibx := locationListMap["ibx"].(string)
	if metroName != "" {
		location.SetMetroName(metroName)
	}
	if region != "" {
		location.SetRegion(region)
	}
	if metroCode != "" {
		location.SetMetroCode(metroCode)
	}
	if ibx != "" {
		location.SetIbx(ibx)
	}

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

func LocationWithoutIBXTerraformToGo(locationList []interface{}) fabricv4.SimplifiedLocationWithoutIBX {
	if len(locationList) == 0 {
		return fabricv4.SimplifiedLocationWithoutIBX{}
	}

	var locationWithoutIbx fabricv4.SimplifiedLocationWithoutIBX
	locationMap := locationList[0].(map[string]interface{})
	metro_code := locationMap["metro_code"].(string)
	locationWithoutIbx.SetMetroCode(metro_code)
	return locationWithoutIbx
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

func ProjectTerraformToGo(projectTerraform []interface{}) fabricv4.Project {
	if len(projectTerraform) == 0 {
		return fabricv4.Project{}
	}
	var project fabricv4.Project
	projectMap := projectTerraform[0].(map[string]interface{})
	projectId := projectMap["project_id"].(string)
	if projectId != "" {
		project.SetProjectId(projectId)
	}

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

func ChangeLogGoToTerraform(changeLog *fabricv4.Changelog) *schema.Set {
	if changeLog == nil {
		return nil
	}

	mappedChangeLog := map[string]interface{}{
		"created_by":           changeLog.GetCreatedBy(),
		"created_by_full_name": changeLog.GetCreatedByFullName(),
		"created_by_email":     changeLog.GetCreatedByEmail(),
		"created_date_time":    changeLog.GetCreatedDateTime().String(),
		"updated_by":           changeLog.GetUpdatedBy(),
		"updated_by_full_name": changeLog.GetUpdatedByFullName(),
		"updated_date_time":    changeLog.GetUpdatedDateTime().String(),
		"deleted_by":           changeLog.GetDeletedBy(),
		"deleted_by_full_name": changeLog.GetDeletedByFullName(),
		"deleted_by_email":     changeLog.GetDeletedByEmail(),
		"deleted_date_time":    changeLog.GetDeletedDateTime().String(),
	}

	changeLogSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: ChangeLogSch()}),
		[]interface{}{mappedChangeLog},
	)
	return changeLogSet
}

func ErrorGoToTerraform(errors []fabricv4.Error) []interface{} {
	if len(errors) == 0 {
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
