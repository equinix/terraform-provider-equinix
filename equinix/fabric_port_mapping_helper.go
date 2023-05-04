package equinix

import (
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func portApiPortRedundancyToTerr(redundancy *v4.PortRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	redundancies := []*v4.PortRedundancy{redundancy}
	mappedRedundancies := make([]interface{}, 0)
	for _, redundancy := range redundancies {
		mappedRedundancy := make(map[string]interface{})
		mappedRedundancy["enabled"] = redundancy.Enabled
		mappedRedundancy["group"] = redundancy.Group
		mappedRedundancy["priority"] = string(*redundancy.Priority)
		mappedRedundancies = append(mappedRedundancies, mappedRedundancy)
	}
	redundancySet := schema.NewSet(
		schema.HashResource(readPortsRedundancyRes),
		mappedRedundancies,
	)
	return redundancySet
}

func portOperationToTerra(operation *v4.PortOperation) *schema.Set {
	if operation == nil {
		return nil
	}
	operations := []*v4.PortOperation{operation}
	mappedOperations := make([]interface{}, 0)
	for _, operation := range operations {
		mappedOperation := make(map[string]interface{})
		mappedOperation["operational_status"] = operation.OperationalStatus
		mappedOperation["connection_count"] = operation.ConnectionCount
		mappedOperation["op_status_changed_at"] = operation.OpStatusChangedAt.String()
		mappedOperations = append(mappedOperations, mappedOperation)
	}
	operationSet := schema.NewSet(
		schema.HashResource(createOperationRes),
		mappedOperations,
	)
	return operationSet
}

func portDeviceRedundancyToTerra(redundancy *v4.PortDeviceRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	redundancies := []*v4.PortDeviceRedundancy{redundancy}
	mappedRedundancies := make([]interface{}, 0)
	for _, redundancy := range redundancies {
		mappedRedundancy := make(map[string]interface{})
		mappedRedundancy["group"] = redundancy.Group
		mappedRedundancy["priority"] = redundancy.Priority
		mappedRedundancies = append(mappedRedundancies, mappedRedundancy)
	}
	redundancySet := schema.NewSet(
		schema.HashResource(readPortDeviceRedundancyRes),
		mappedRedundancies,
	)
	return redundancySet
}

func deviceToTerra(device *v4.PortDevice) *schema.Set {
	if device == nil {
		return nil
	}
	devices := []*v4.PortDevice{device}
	mappedDevices := make([]interface{}, 0)
	for _, device := range devices {
		mappedDevice := make(map[string]interface{})
		mappedDevice["name"] = device.Name
		if device.Redundancy != nil {
			mappedDevice["redundancy"] = portDeviceRedundancyToTerra(device.Redundancy)
		}
		mappedDevices = append(mappedDevices, mappedDevice)
	}
	deviceSet := schema.NewSet(
		schema.HashResource(readPortDeviceRes),
		mappedDevices,
	)
	return deviceSet
}

func encapsulationToTerra(portEncapsulation *v4.PortEncapsulation) *schema.Set {
	if portEncapsulation == nil {
		return nil
	}
	portEncapsulations := []*v4.PortEncapsulation{portEncapsulation}
	mappedPortEncapsulations := make([]interface{}, 0)
	for _, portEncapsulation := range portEncapsulations {
		mappedPortEncapsulation := make(map[string]interface{})
		mappedPortEncapsulation["type"] = portEncapsulation.Type_
		mappedPortEncapsulation["tag_protocol_id"] = portEncapsulation.TagProtocolId
		mappedPortEncapsulations = append(mappedPortEncapsulations, mappedPortEncapsulation)
	}
	portEncapsulationSet := schema.NewSet(
		schema.HashResource(readPortEncapsulationRes),
		mappedPortEncapsulations,
	)
	return portEncapsulationSet
}

//func lagToTerra(portLag *v4.PortLag) *schema.Set {
//	if portLag == nil {
//		return nil
//	}
//	portLags := []*v4.PortLag{portLag}
//	mappedPortLags := make([]interface{}, 0)
//	for _, portLag := range portLags {
//		mappedPortLag := make(map[string]interface{})
//		mappedPortLag["enabled"] = portLag.Enabled
//		mappedPortLag["id"] = portLag.Id
//		mappedPortLag["name"] = portLag.Name
//		mappedPortLag["member_status"] = portLag.MemberStatus
//		mappedPortLags = append(mappedPortLags, mappedPortLag)
//	}
//	portLagSet := schema.NewSet(
//		schema.HashResource(readPortLagRes),
//		mappedPortLags,
//	)
//	return portLagSet
//}

func fabricPortsListToTerra(ports v4.AllPortsResponse) []map[string]interface{} {
	portsl := ports.Data
	if portsl == nil {
		return nil
	}
	mappedPortsl := make([]map[string]interface{}, len(portsl))
	for index, port := range portsl {
		mappedPortsl[index] = map[string]interface{}{
			"uuid":                port.Uuid,
			"name":                port.Name,
			"bandwidth":           port.Bandwidth,
			"available_bandwidth": port.AvailableBandwidth,
			"used_bandwidth":      port.UsedBandwidth,
			"href":                port.Href,
			"description":         port.Description,
			"type":                port.Type_,
			"state":               port.State,
			"service_type":        port.ServiceType,
			"operation":           portOperationToTerra(port.Operation),
			"redundancy":          portApiPortRedundancyToTerr(port.Redundancy),
			"account":             accountToTerra(port.Account),
			"change_log":          changeLogToTerra(port.Changelog),
			"location":            locationToTerra(port.Location),
			"device":              deviceToTerra(port.Device),
			"encapsulation":       encapsulationToTerra(port.Encapsulation),
			"lag":                 port.LagEnabled,
		}
	}
	return mappedPortsl
}
