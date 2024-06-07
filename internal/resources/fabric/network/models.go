package network

import (
	"fmt"
	"log"

	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func networkMap(nt *fabricv4.Network) map[string]interface{} {
	operation := nt.GetOperation()
	change := nt.GetChange()
	location := nt.GetLocation()
	notifications := nt.GetNotifications()
	project := nt.GetProject()
	changeLog := nt.GetChangeLog()
	network := map[string]interface{}{
		"name":              nt.GetName(),
		"href":              nt.GetHref(),
		"uuid":              nt.GetUuid(),
		"type":              string(nt.GetType()),
		"scope":             string(nt.GetScope()),
		"state":             string(nt.GetState()),
		"operation":         fabricNetworkOperationGoToTerraform(&operation),
		"change":            simplifiedFabricNetworkChangeGoToTerraform(&change),
		"location":          equinix_fabric_schema.LocationGoToTerraform(&location),
		"notifications":     equinix_fabric_schema.NotificationsGoToTerraform(notifications),
		"project":           equinix_fabric_schema.ProjectGoToTerraform(&project),
		"change_log":        equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog),
		"connections_count": nt.GetConnectionsCount(),
	}

	return network
}

func setFabricNetworkMap(d *schema.ResourceData, nt *fabricv4.Network) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, networkMap(nt))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func fabricNetworkOperationGoToTerraform(operation *fabricv4.NetworkOperation) *schema.Set {
	if operation == nil {
		return nil
	}
	mappedOperation := make(map[string]interface{})
	mappedOperation["equinix_status"] = string(*operation.EquinixStatus)

	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricNetworkOperationSch()}),
		[]interface{}{mappedOperation},
	)
	return operationSet
}

func simplifiedFabricNetworkChangeGoToTerraform(networkChange *fabricv4.SimplifiedNetworkChange) *schema.Set {

	mappedChange := make(map[string]interface{})
	mappedChange["href"] = networkChange.GetHref()
	mappedChange["type"] = string(networkChange.GetType())
	mappedChange["uuid"] = networkChange.GetUuid()

	changeSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricNetworkChangeSch()}),
		[]interface{}{mappedChange},
	)
	return changeSet
}

func getFabricNetworkUpdateRequest(network *fabricv4.Network, d *schema.ResourceData) (fabricv4.NetworkChangeOperation, error) {
	changeOps := fabricv4.NetworkChangeOperation{}
	existingName := network.GetName()
	updateNameVal := d.Get("name").(string)

	log.Printf("existing name %s, Update Name Request %s ", existingName, updateNameVal)

	if existingName != updateNameVal {
		changeOps = fabricv4.NetworkChangeOperation{Op: "replace", Path: "/name", Value: updateNameVal}
	} else {
		return changeOps, fmt.Errorf("nothing to update for the Fabric Network: %s", existingName)
	}
	return changeOps, nil
}
