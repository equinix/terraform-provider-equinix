package equinix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/antihax/optional"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func portDeviceSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port name",
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port device redundancy",
			Elem: &schema.Resource{
				Schema: PortRedundancySch(),
			},
		},
	}
}

func portEncapsulationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port encapsulation protocol type",
		},
		"tag_protocol_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port encapsulation Tag Protocol Identifier",
		},
	}
}

func portOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"operational_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port operation status",
		},
		"connection_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total number of current connections",
		},
		"op_status_changed_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Date and time at which port availability changed",
		},
	}
}

func PortRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Access point redundancy",
		},
		"group": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port redundancy group",
		},
		"priority": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Priority type-Primary or Secondary",
		},
	}
}

func FabricPortResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port type",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port URI information",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix-assigned port identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port name",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port description",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port state",
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port specific operational data",
			Elem: &schema.Resource{
				Schema: portOperationSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Customer account information that is associated with this port",
			Elem: &schema.Resource{
				Schema: equinix_schema.AccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures port lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_schema.ChangeLogSch(),
			},
		},
		"service_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port service type",
		},
		"bandwidth": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port bandwidth in Mbps",
		},
		"available_bandwidth": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port available bandwidth in Mbps",
		},
		"used_bandwidth": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port used bandwidth in Mbps",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port location information",
			Elem: &schema.Resource{
				Schema: equinix_schema.LocationSch(),
			},
		},
		"device": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port device",
			Elem: &schema.Resource{
				Schema: portDeviceSch(),
			},
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port redundancy information",
			Elem: &schema.Resource{
				Schema: PortRedundancySch(),
			},
		},
		"encapsulation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port encapsulation protocol",
			Elem: &schema.Resource{
				Schema: portEncapsulationSch(),
			},
		},
		"lag_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Port Lag",
		},
	}
}

func readFabricPortResourceSchemaUpdated() map[string]*schema.Schema {
	sch := FabricPortResourceSchema()
	sch["uuid"].Computed = true
	sch["uuid"].Optional = false
	sch["uuid"].Required = false
	return sch
}

func readFabricPortsResponseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Ports",
			Elem: &schema.Resource{
				Schema: readFabricPortResourceSchemaUpdated(),
			},
		},
		"filters": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "name",
			Elem: &schema.Resource{
				Schema: readGetPortsByNameQueryParamSch(),
			},
		},
	}
}

func readGetPortsByNameQueryParamSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Query Parameter to Get Ports By Name",
		},
	}
}

func portToFabric(portList []interface{}) v4.SimplifiedPort {
	p := v4.SimplifiedPort{}
	for _, pl := range portList {
		plMap := pl.(map[string]interface{})
		uuid := plMap["uuid"].(string)
		p = v4.SimplifiedPort{Uuid: uuid}
	}
	return p
}

func portToTerra(port *v4.SimplifiedPort) *schema.Set {
	ports := []*v4.SimplifiedPort{port}
	mappedPorts := make([]interface{}, len(ports))
	for _, port := range ports {
		mappedPort := make(map[string]interface{})
		mappedPort["href"] = port.Href
		mappedPort["name"] = port.Name
		mappedPort["uuid"] = port.Uuid
		if port.Redundancy != nil {
			mappedPort["redundancy"] = PortRedundancyToTerra(port.Redundancy)
		}
		mappedPorts = append(mappedPorts, mappedPort)
	}
	portSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: portSch()}),
		mappedPorts,
	)
	return portSet
}

func PortRedundancyToTerra(redundancy *v4.PortRedundancy) *schema.Set {
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
		schema.HashResource(&schema.Resource{Schema: PortRedundancySch()}),
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
		schema.HashResource(&schema.Resource{Schema: operationSch()}),
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
		schema.HashResource(&schema.Resource{Schema: PortRedundancySch()}),
		mappedRedundancies,
	)
	return redundancySet
}

func portDeviceToTerra(device *v4.PortDevice) *schema.Set {
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
		schema.HashResource(&schema.Resource{Schema: portDeviceSch()}),
		mappedDevices,
	)
	return deviceSet
}

func portEncapsulationToTerra(portEncapsulation *v4.PortEncapsulation) *schema.Set {
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
		schema.HashResource(&schema.Resource{Schema: portEncapsulationSch()}),
		mappedPortEncapsulations,
	)
	return portEncapsulationSet
}

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
			"redundancy":          PortRedundancyToTerra(port.Redundancy),
			"account":             accountToTerra(port.Account),
			"change_log":          equinix_schema.ChangeLogToTerra(port.Changelog),
			"location":            equinix_schema.LocationToTerra(port.Location),
			"device":              portDeviceToTerra(port.Device),
			"encapsulation":       portEncapsulationToTerra(port.Encapsulation),
			"lag_enabled":         port.LagEnabled,
		}
	}
	return mappedPortsl
}

func resourceFabricPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	port, _, err := client.PortsApi.GetPortByUuid(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Port %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(port.Uuid)
	return setFabricPortMap(d, port)
}

func setFabricPortMap(d *schema.ResourceData, port v4.Port) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
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
		"redundancy":          PortRedundancyToTerra(port.Redundancy),
		"account":             accountToTerra(port.Account),
		"change_log":          equinix_schema.ChangeLogToTerra(port.Changelog),
		"location":            equinix_schema.LocationToTerra(port.Location),
		"device":              portDeviceToTerra(port.Device),
		"encapsulation":       portEncapsulationToTerra(port.Encapsulation),
		"lag_enabled":         port.LagEnabled,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setPortsListMap(d *schema.ResourceData, spl v4.AllPortsResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"data": fabricPortsListToTerra(spl),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceFabricPortGetByPortName(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Panic occurred during GET /fabric/v4/ports: %+v", r)
			log.Printf("[ERROR] Stack Trace from Panic: %s", debug.Stack())
			diags = diag.FromErr(errors.New(`
				there is a schema error in the return value from the GET /fabric/v4/ports endpoint.
				Set the following env variable TF_LOG=DEBUG and rerun the terraform apply.
				Copy the log output and open an issue with it in the Github for the Equinix Terraform Provider.
				https://github.com/equinix/terraform-provider-equinix
				We will review and correct as soon as possible.
				Thank you!
			`))
		}
	}()

	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	portNameParam := d.Get("filters").(*schema.Set).List()
	portName := portNameQueryParamToFabric(portNameParam)
	ports, _, err := client.PortsApi.GetPorts(ctx, &portName)
	if err != nil {
		log.Printf("[WARN] Ports not found , error %s", err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	if len(ports.Data) != 1 {
		error := fmt.Errorf("incorrect # of records are found for the port name parameter criteria - %d , please change the criteria", len(ports.Data))
		d.SetId("")
		return diag.FromErr(error)
	}

	d.SetId(ports.Data[0].Uuid)
	return setPortsListMap(d, ports)
}

func portNameQueryParamToFabric(portNameParam []interface{}) v4.PortsApiGetPortsOpts {
	if len(portNameParam) == 0 {
		return v4.PortsApiGetPortsOpts{}
	}
	mappedPn := v4.PortsApiGetPortsOpts{}
	for _, pn := range portNameParam {
		pnMap := pn.(map[string]interface{})
		portName := pnMap["name"].(string)
		pName := optional.NewString(portName)
		mappedPn = v4.PortsApiGetPortsOpts{Name: pName}
	}
	return mappedPn
}
