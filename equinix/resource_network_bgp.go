package equinix

import (
	"fmt"
	"net/http"
	"time"

	"github.com/equinix/ne-go"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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

func resourceNetworkBGP() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkBGPCreate,
		Read:   resourceNetworkBGPRead,
		Update: resourceNetworkBGPUpdate,
		Delete: resourceNetworkBGPDelete,
		Schema: createNetworkBGPResourceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func createNetworkBGPResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkBGPSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkBGPSchemaNames["ConnectionUUID"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		networkBGPSchemaNames["DeviceUUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkBGPSchemaNames["LocalIPAddress"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsCIDR,
		},
		networkBGPSchemaNames["LocalASN"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		networkBGPSchemaNames["RemoteIPAddress"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},
		networkBGPSchemaNames["RemoteASN"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		networkBGPSchemaNames["AuthenticationKey"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(6, 60),
		},
		networkBGPSchemaNames["State"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkBGPSchemaNames["ProvisioningStatus"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceNetworkBGPCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	bgp := createNetworkBGPConfiguration(d)
	existingBGP, err := conf.ne.GetBGPConfigurationForConnection(bgp.ConnectionUUID)
	//Reuse existing configuration, as there was no possibility to remove it due to API limitations
	if err == nil {
		updateErr := conf.ne.NewBGPConfigurationUpdateRequest(existingBGP.UUID).
			WithRemoteIPAddress(bgp.RemoteIPAddress).
			WithRemoteASN(bgp.RemoteASN).
			WithLocalIPAddress(bgp.LocalIPAddress).
			WithLocalASN(bgp.LocalASN).
			WithAuthenticationKey(bgp.AuthenticationKey).
			Execute()
		if updateErr != nil {
			return fmt.Errorf("failed to update BGP configuration '%s': %s", existingBGP.UUID, updateErr)
		}
		d.SetId(existingBGP.UUID)
	} else {
		restErr, ok := err.(rest.Error)
		if !ok || restErr.HTTPCode != http.StatusNotFound {
			return fmt.Errorf("failed to fetch BGP configuration for connection '%s': %s", bgp.ConnectionUUID, err)
		}
		uuid, err := conf.ne.CreateBGPConfiguration(bgp)
		if err != nil {
			return err
		}
		d.SetId(uuid)
	}
	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			ne.BGPProvisioningStatusProvisioning,
			ne.BGPProvisioningStatusPendingUpdate,
		},
		Target: []string{
			ne.BGPProvisioningStatusProvisioned,
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Second,
		MinTimeout: 2 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := conf.ne.GetBGPConfiguration(d.Id())
			if err != nil {
				return nil, "", err
			}
			return resp, resp.ProvisioningStatus, nil
		},
	}
	if _, err := createStateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for BGP configuration (%s) to be created: %s", d.Id(), err)
	}
	return resourceNetworkBGPRead(d, m)
}

func resourceNetworkBGPRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	bgp, err := conf.ne.GetBGPConfiguration(d.Id())
	if err != nil {
		return err
	}
	if err := updateNetworkBGPResource(bgp, d); err != nil {
		return err
	}
	return nil
}

func resourceNetworkBGPUpdate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	updateReq := conf.ne.NewBGPConfigurationUpdateRequest(d.Id())
	if v, ok := d.GetOk(networkBGPSchemaNames["LocalIPAddress"]); ok {
		updateReq.WithLocalIPAddress(v.(string))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["LocalASN"]); ok {
		updateReq.WithLocalASN(v.(int))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["RemoteIPAddress"]); ok {
		updateReq.WithRemoteIPAddress(v.(string))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["RemoteASN"]); ok {
		updateReq.WithRemoteASN(v.(int))
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["AuthenticationKey"]); ok {
		updateReq.WithAuthenticationKey(v.(string))
	}
	if err := updateReq.Execute(); err != nil {
		return err
	}
	return resourceNetworkBGPRead(d, m)
}

func resourceNetworkBGPDelete(d *schema.ResourceData, m interface{}) error {
	//BGP configuration removal is not possible with NE public APIs
	d.SetId("")
	return nil
}

func createNetworkBGPConfiguration(d *schema.ResourceData) ne.BGPConfiguration {
	bgp := ne.BGPConfiguration{}
	if v, ok := d.GetOk(networkBGPSchemaNames["UUID"]); ok {
		bgp.UUID = v.(string)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["ConnectionUUID"]); ok {
		bgp.ConnectionUUID = v.(string)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["DeviceUUID"]); ok {
		bgp.DeviceUUID = v.(string)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["LocalIPAddress"]); ok {
		bgp.LocalIPAddress = v.(string)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["LocalASN"]); ok {
		bgp.LocalASN = v.(int)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["RemoteIPAddress"]); ok {
		bgp.RemoteIPAddress = v.(string)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["RemoteASN"]); ok {
		bgp.RemoteASN = v.(int)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["AuthenticationKey"]); ok {
		bgp.AuthenticationKey = v.(string)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["State"]); ok {
		bgp.State = v.(string)
	}
	if v, ok := d.GetOk(networkBGPSchemaNames["ProvisioningStatus"]); ok {
		bgp.ProvisioningStatus = v.(string)
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
