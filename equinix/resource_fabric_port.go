package equinix

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/antihax/optional"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFabricPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	port, _, err := client.PortsApi.GetPortByUuid(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Port %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	d.SetId(port.Uuid)
	return setFabricPortMap(d, port)
}

func setFabricPortMap(d *schema.ResourceData, port v4.Port) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := setMap(d, map[string]interface{}{
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
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setPortsListMap(d *schema.ResourceData, spl v4.AllPortsResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := setMap(d, map[string]interface{}{
		"data": fabricPortsListToTerra(spl),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceFabricPortGetByPortName(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	portNameParam := d.Get("filters").(*schema.Set).List()
	portName := portNameQueryParamToFabric(portNameParam)
	ports, _, err := client.PortsApi.GetPorts(ctx, &portName)
	if err != nil {
		log.Printf("[WARN] Ports not found , error %s", err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
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
