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
	locations := []*v4.SimplifiedLocation{location}
	mappedLocations := make([]interface{}, len(locations))
	for i, location := range locations {
		mappedLocations[i] = map[string]interface{}{
			"region":     location.Region,
			"metro_name": location.MetroName,
			"metro_code": location.MetroCode,
			"ibx":        location.Ibx,
		}
	}
	locationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: LocationSch()}),
		mappedLocations,
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

func ChangeLogToTerra(changeLog *v4.Changelog) *schema.Set {
	if changeLog == nil {
		return nil
	}
	changeLogs := []*v4.Changelog{changeLog}
	mappedChangeLogs := make([]interface{}, len(changeLogs))
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
		schema.HashResource(&schema.Resource{Schema: ChangeLogSch()}),
		mappedChangeLogs,
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
