package equinix

import (
	"context"
	"errors"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"log"
	"runtime/debug"
	"strings"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

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
				Schema: equinix_fabric_schema.AccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures port lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
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
				Schema: equinix_fabric_schema.LocationSch(),
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
			MaxItems:    1,
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

func portTerraformToGo(portList []interface{}) fabricv4.SimplifiedPort {
	if portList == nil || len(portList) == 0 {
		return fabricv4.SimplifiedPort{}
	}
	var port fabricv4.SimplifiedPort
	portListMap := portList[0].(map[string]interface{})
	uuid := portListMap["uuid"].(string)
	port.SetUuid(uuid)

	return port
}

func portGoToTerraform(port *fabricv4.SimplifiedPort) *schema.Set {
	mappedPort := make(map[string]interface{})
	mappedPort["href"] = port.GetHref()
	mappedPort["name"] = port.GetName()
	mappedPort["uuid"] = port.GetUuid()
	if port.Redundancy != nil {
		mappedPort["redundancy"] = portRedundancyGoToTerraform(port.Redundancy)
	}
	portSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: portSch()}),
		[]interface{}{mappedPort},
	)
	return portSet
}

func portRedundancyGoToTerraform(redundancy *fabricv4.PortRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	mappedRedundancy := make(map[string]interface{})
	mappedRedundancy["enabled"] = redundancy.GetEnabled()
	mappedRedundancy["group"] = redundancy.GetGroup()
	mappedRedundancy["priority"] = string(redundancy.GetPriority())

	redundancySet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: PortRedundancySch()}),
		[]interface{}{mappedRedundancy},
	)
	return redundancySet
}

func portOperationGoToTerraform(operation *fabricv4.PortOperation) *schema.Set {
	if operation == nil {
		return nil
	}

	mappedOperation := make(map[string]interface{})
	mappedOperation["operational_status"] = operation.GetOperationalStatus()
	mappedOperation["connection_count"] = operation.GetConnectionCount()
	mappedOperation["op_status_changed_at"] = operation.GetOpStatusChangedAt().String()

	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: operationSch()}),
		[]interface{}{mappedOperation},
	)
	return operationSet
}

func portDeviceRedundancyGoToTerraform(redundancy *fabricv4.PortDeviceRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}

	mappedRedundancy := make(map[string]interface{})
	mappedRedundancy["group"] = redundancy.GetGroup()
	mappedRedundancy["priority"] = string(redundancy.GetPriority())

	redundancySet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: PortRedundancySch()}),
		[]interface{}{mappedRedundancy},
	)
	return redundancySet
}

func portDeviceGoToTerraform(device *fabricv4.PortDevice) *schema.Set {
	if device == nil {
		return nil
	}

	mappedDevice := make(map[string]interface{})
	mappedDevice["name"] = device.GetName()
	redundancy := device.GetRedundancy()
	if &redundancy != nil {
		mappedDevice["redundancy"] = portDeviceRedundancyGoToTerraform(&redundancy)
	}

	deviceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: portDeviceSch()}),
		[]interface{}{mappedDevice},
	)
	return deviceSet
}

func portEncapsulationGoToTerraform(portEncapsulation *fabricv4.PortEncapsulation) *schema.Set {
	if portEncapsulation == nil {
		return nil
	}

	mappedPortEncapsulation := make(map[string]interface{})
	mappedPortEncapsulation["type"] = string(portEncapsulation.GetType())
	mappedPortEncapsulation["tag_protocol_id"] = portEncapsulation.GetTagProtocolId()

	portEncapsulationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: portEncapsulationSch()}),
		[]interface{}{mappedPortEncapsulation},
	)
	return portEncapsulationSet
}

func resourceFabricPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	port, _, err := client.PortsApi.GetPortByUuid(ctx, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Port %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(port.GetUuid())
	return setFabricPortMap(d, port)
}

func fabricPortMap(port *fabricv4.PortResponse) map[string]interface{} {
	operation := port.GetOperation()
	redundancy := port.GetRedundancy()
	account := port.GetAccount()
	changelog := port.GetChangelog()
	location := port.GetLocation()
	device := port.GetDevice()
	encapsulation := port.GetEncapsulation()
	return map[string]interface{}{
		"uuid":                port.GetUuid(),
		"name":                port.GetName(),
		"bandwidth":           port.GetBandwidth(),
		"available_bandwidth": port.GetAvailableBandwidth(),
		"used_bandwidth":      port.GetUsedBandwidth(),
		"href":                port.GetHref(),
		"description":         port.GetDescription(),
		"type":                string(port.GetType()),
		"state":               string(port.GetState()),
		"service_type":        string(port.GetServiceType()),
		"operation":           portOperationGoToTerraform(&operation),
		"redundancy":          portRedundancyGoToTerraform(&redundancy),
		"account":             equinix_fabric_schema.AccountGoToTerraform(&account),
		"change_log":          equinix_fabric_schema.ChangeLogGoToTerraform(&changelog),
		"location":            equinix_fabric_schema.LocationGoToTerraform(&location),
		"device":              portDeviceGoToTerraform(&device),
		"encapsulation":       portEncapsulationGoToTerraform(&encapsulation),
		"lag_enabled":         port.GetLagEnabled(),
	}
}

func setFabricPortMap(d *schema.ResourceData, port *fabricv4.PortResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, fabricPortMap(port))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setPortsListMap(d *schema.ResourceData, portResponse *fabricv4.AllPortsResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	ports := portResponse.Data
	if ports == nil {
		return nil
	}
	mappedPorts := make([]map[string]interface{}, len(ports))
	for index, port := range ports {
		mappedPorts[index] = fabricPortMap(&port)
	}

	err := equinix_schema.SetMap(d, map[string]interface{}{
		"data": mappedPorts,
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

	client := meta.(*config.Config).NewFabricClientForSDK(d)
	portNameParam := d.Get("filters").(*schema.Set).List()
	portName := portName(portNameParam)
	ports, _, err := client.PortsApi.GetPorts(ctx).Name(portName).Execute()
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

	d.SetId(ports.Data[0].GetUuid())
	return setPortsListMap(d, ports)
}

func portName(portNameParam []interface{}) string {
	if len(portNameParam) == 0 {
		return ""
	}

	pnMap := portNameParam[0].(map[string]interface{})
	portName := pnMap["name"].(string)
	return portName
}
