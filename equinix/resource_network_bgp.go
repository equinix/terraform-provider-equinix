package equinix

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/ne-go"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var networkBGPSchemaNames = map[string]string{
	"UUID":               "uuid",
	"ConnectionUUID":     "connection_id",
	"DeviceUUID":         "device_id",
	"LocalIPAddress":     "local_ip_address",
	"LocalASN":           "local_asn",
	"RemoteIPAddress":    "remote_ip_address",
	"RemoteASN":          "remote_asn",
	"AuthenticationKey":  "authentication_key",
	"State":              "state",
	"ProvisioningStatus": "provisioning_status",
}

var networkBGPDescriptions = map[string]string{
	"UUID":               "BGP peering configuration unique identifier",
	"ConnectionUUID":     "Identifier of a connection established between network device and remote service provider that will be used for peering",
	"DeviceUUID":         "Unique identifier of a network device that is a local peer in a given BGP peering configuration",
	"LocalIPAddress":     "IP address in CIDR format of a local device",
	"LocalASN":           "Local ASN number",
	"RemoteIPAddress":    "IP address of remote peer",
	"RemoteASN":          "Remote ASN number",
	"AuthenticationKey":  "Shared key used for BGP peer authentication",
	"State":              "BGP peer state",
	"ProvisioningStatus": "BGP peering configuration provisioning status",
}

func resourceNetworkBGP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkBGPCreate,
		ReadContext:   resourceNetworkBGPRead,
		UpdateContext: resourceNetworkBGPUpdate,
		DeleteContext: resourceNetworkBGPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: createNetworkBGPResourceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Description: "Resource allows creation and management of Equinix Network Edge BGP peering configurations",
	}
}

func createNetworkBGPResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkBGPSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkBGPDescriptions["UUID"],
		},
		networkBGPSchemaNames["ConnectionUUID"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  networkBGPDescriptions["ConnectionUUID"],
		},
		networkBGPSchemaNames["DeviceUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkBGPDescriptions["DeviceUUID"],
		},
		networkBGPSchemaNames["LocalIPAddress"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsCIDR,
			Description:  networkBGPDescriptions["LocalIPAddress"],
		},
		networkBGPSchemaNames["LocalASN"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  networkBGPDescriptions["LocalASN"],
		},
		networkBGPSchemaNames["RemoteIPAddress"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
			Description:  networkBGPDescriptions["RemoteIPAddress"],
		},
		networkBGPSchemaNames["RemoteASN"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  networkBGPDescriptions["RemoteASN"],
		},
		networkBGPSchemaNames["AuthenticationKey"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(6, 60),
			Description:  networkBGPDescriptions["AuthenticationKey"],
		},
		networkBGPSchemaNames["State"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkBGPDescriptions["State"],
		},
		networkBGPSchemaNames["ProvisioningStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkBGPDescriptions["ProvisioningStatus"],
		},
	}
}

func resourceNetworkBGPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	bgp := createNetworkBGPConfiguration(d)
	existingBGP, err := client.GetBGPConfigurationForConnection(ne.StringValue(bgp.ConnectionUUID))
	if err == nil {
		bgp.UUID = existingBGP.UUID
		if updateErr := createNetworkBGPUpdateRequest(client.NewBGPConfigurationUpdateRequest, &bgp); updateErr != nil {
			return diag.Errorf("failed to update BGP configuration '%s': %s", ne.StringValue(existingBGP.UUID), updateErr)
		}
		d.SetId(ne.StringValue(bgp.UUID))
	} else {
		restErr, ok := err.(rest.Error)
		if !ok || restErr.HTTPCode != http.StatusNotFound {
			return diag.Errorf("failed to fetch BGP configuration for connection '%s': %s", ne.StringValue(bgp.ConnectionUUID), err)
		}
		uuid, err := client.CreateBGPConfiguration(bgp)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(ne.StringValue(uuid))
	}
	if _, err := createBGPConfigStatusProvisioningWaitConfiguration(client.GetBGPConfiguration, d.Id(), 2*time.Second, d.Timeout(schema.TimeoutCreate)).WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for BGP configuration (%s) to be created: %s", d.Id(), err)
	}
	diags = append(diags, resourceNetworkBGPRead(ctx, d, m)...)
	return diags
}

func resourceNetworkBGPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	bgp, err := client.GetBGPConfiguration(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := updateNetworkBGPResource(bgp, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceNetworkBGPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	bgpConfig := createNetworkBGPConfiguration(d)
	if err := createNetworkBGPUpdateRequest(client.NewBGPConfigurationUpdateRequest, &bgpConfig).Execute(); err != nil {
		return diag.FromErr(err)
	}
	diags = append(diags, resourceNetworkBGPRead(ctx, d, m)...)
	return diags
}

func resourceNetworkBGPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// BGP configuration removal is not possible with NE public APIs
	d.SetId("")
	return nil
}

func createNetworkBGPConfiguration(d *schema.ResourceData) ne.BGPConfiguration {
	bgp := ne.BGPConfiguration{}
	if v, ok := d.GetOk(networkBGPSchemaNames["UUID"]); ok {
		bgp.UUID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["ConnectionUUID"]); ok {
		bgp.ConnectionUUID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["LocalIPAddress"]); ok {
		bgp.LocalIPAddress = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["LocalASN"]); ok {
		bgp.LocalASN = ne.Int(v.(int))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["RemoteIPAddress"]); ok {
		bgp.RemoteIPAddress = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["RemoteASN"]); ok {
		bgp.RemoteASN = ne.Int(v.(int))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["AuthenticationKey"]); ok {
		bgp.AuthenticationKey = ne.String(v.(string))
	}
	return bgp
}

func updateNetworkBGPResource(bgp *ne.BGPConfiguration, d *schema.ResourceData) error {
	if err := d.Set(networkBGPSchemaNames["UUID"], bgp.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["ConnectionUUID"], bgp.ConnectionUUID); err != nil {
		return fmt.Errorf("error reading ConnectionUUID: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["DeviceUUID"], bgp.DeviceUUID); err != nil {
		return fmt.Errorf("error reading DeviceUUID: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["LocalIPAddress"], bgp.LocalIPAddress); err != nil {
		return fmt.Errorf("error reading LocalIPAddress: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["LocalASN"], bgp.LocalASN); err != nil {
		return fmt.Errorf("error reading LocalASN: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["RemoteIPAddress"], bgp.RemoteIPAddress); err != nil {
		return fmt.Errorf("error reading RemoteIPAddress: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["RemoteASN"], bgp.RemoteASN); err != nil {
		return fmt.Errorf("error reading RemoteASN: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["AuthenticationKey"], bgp.AuthenticationKey); err != nil {
		return fmt.Errorf("error reading AuthenticationKey: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["State"], bgp.State); err != nil {
		return fmt.Errorf("error reading State: %s", err)
	}
	if err := d.Set(networkBGPSchemaNames["ProvisioningStatus"], bgp.ProvisioningStatus); err != nil {
		return fmt.Errorf("error reading ProvisioningStatus: %s", err)
	}
	return nil
}

type bgpUpdateRequest func(uuid string) ne.BGPUpdateRequest

func createNetworkBGPUpdateRequest(requestFunc bgpUpdateRequest, bgp *ne.BGPConfiguration) ne.BGPUpdateRequest {
	return requestFunc(ne.StringValue(bgp.UUID)).
		WithRemoteIPAddress(ne.StringValue(bgp.RemoteIPAddress)).
		WithRemoteASN(ne.IntValue(bgp.RemoteASN)).
		WithLocalIPAddress(ne.StringValue(bgp.LocalIPAddress)).
		WithLocalASN(ne.IntValue(bgp.LocalASN)).
		WithAuthenticationKey(ne.StringValue(bgp.AuthenticationKey))
}

type getBGPConfig func(uuid string) (*ne.BGPConfiguration, error)

func createBGPConfigStatusProvisioningWaitConfiguration(fetchFunc getBGPConfig, id string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			ne.BGPProvisioningStatusProvisioning,
			ne.BGPProvisioningStatusPendingUpdate,
		},
		Target: []string{
			ne.BGPProvisioningStatusProvisioned,
		},
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: delay,
		Refresh: func() (interface{}, string, error) {
			resp, err := fetchFunc(id)
			if err != nil {
				return nil, "", err
			}
			return resp, ne.StringValue(resp.ProvisioningStatus), nil
		},
	}
}
