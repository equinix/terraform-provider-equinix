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

var neBGPSchemaNames = map[string]string{
	"UUID":               "uuid",
	"ConnectionUUID":     "connection_uuid",
	"DeviceUUID":         "device_uuid",
	"LocalIPAddress":     "local_ip_address",
	"LocalASN":           "local_asn",
	"RemoteIPAddress":    "remote_ip_address",
	"RemoteASN":          "remote_asn",
	"AuthenticationKey":  "authentication_key",
	"State":              "state",
	"ProvisioningStatus": "provisioning_status",
}

func resourceNeBGP() *schema.Resource {
	return &schema.Resource{
		Create: resourceNeBGPCreate,
		Read:   resourceNeBGPRead,
		Update: resourceNeBGPUpdate,
		Delete: resourceNeBGPDelete,
		Schema: createNeBGPResourceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func createNeBGPResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neBGPSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neBGPSchemaNames["ConnectionUUID"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		neBGPSchemaNames["DeviceUUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neBGPSchemaNames["LocalIPAddress"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsCIDR,
		},
		neBGPSchemaNames["LocalASN"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		neBGPSchemaNames["RemoteIPAddress"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsCIDR,
		},
		neBGPSchemaNames["RemoteASN"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		neBGPSchemaNames["AuthenticationKey"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		neBGPSchemaNames["State"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neBGPSchemaNames["ProvisioningStatus"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceNeBGPCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	bgp := createNEBGPConfiguration(d)
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
		return resourceNeBGPRead(d, m)
	}
	restErr, ok := err.(rest.Error)
	if !ok || restErr.HTTPCode != http.StatusNotFound {
		return fmt.Errorf("failed to fetch BGP configuration for connection '%s': %s", bgp.ConnectionUUID, err)
	}
	uuid, err := conf.ne.CreateBGPConfiguration(bgp)
	if err != nil {
		return err
	}
	d.SetId(uuid)
	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			ne.BGPProvisioningStatusProvisioning,
		},
		Target: []string{
			ne.BGPProvisioningStatusProvisioned,
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Second,
		MinTimeout: 2 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := conf.ne.GetBGPConfiguration(uuid)
			if err != nil {
				return nil, "", err
			}
			return resp, resp.ProvisioningStatus, nil
		},
	}
	if _, err := createStateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for BGP configuration (%s) to be created: %s", d.Id(), err)
	}
	return resourceNeBGPRead(d, m)
}

func resourceNeBGPRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	bgp, err := conf.ne.GetBGPConfiguration(d.Id())
	if err != nil {
		return err
	}
	if err := updateNeBGPResource(bgp, d); err != nil {
		return err
	}
	return nil
}

func resourceNeBGPUpdate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	updateReq := conf.ne.NewBGPConfigurationUpdateRequest(d.Id())
	if v, ok := d.GetOk(neBGPSchemaNames["LocalIPAddress"]); ok && d.HasChange(neBGPSchemaNames["LocalIPAddress"]) {
		updateReq.WithLocalIPAddress(v.(string))
	}
	if v, ok := d.GetOk(neBGPSchemaNames["LocalASN"]); ok && d.HasChange(neBGPSchemaNames["LocalASN"]) {
		updateReq.WithLocalASN(v.(int))
	}
	if v, ok := d.GetOk(neBGPSchemaNames["RemoteIPAddress"]); ok && d.HasChange(neBGPSchemaNames["RemoteIPAddress"]) {
		updateReq.WithRemoteIPAddress(v.(string))
	}
	if v, ok := d.GetOk(neBGPSchemaNames["RemoteASN"]); ok && d.HasChange(neBGPSchemaNames["RemoteASN"]) {
		updateReq.WithRemoteASN(v.(int))
	}
	if v, ok := d.GetOk(neBGPSchemaNames["AuthenticationKey"]); ok && d.HasChange(neBGPSchemaNames["AuthenticationKey"]) {
		updateReq.WithAuthenticationKey(v.(string))
	}
	if err := updateReq.Execute(); err != nil {
		return err
	}
	return resourceNeBGPRead(d, m)
}

func resourceNeBGPDelete(d *schema.ResourceData, m interface{}) error {
	//BGP configuration removal is not possible with NE public APIs
	d.SetId("")
	return nil
}

func createNEBGPConfiguration(d *schema.ResourceData) ne.BGPConfiguration {
	bgp := ne.BGPConfiguration{}
	if v, ok := d.GetOk(neBGPSchemaNames["UUID"]); ok {
		bgp.UUID = v.(string)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["ConnectionUUID"]); ok {
		bgp.ConnectionUUID = v.(string)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["DeviceUUID"]); ok {
		bgp.DeviceUUID = v.(string)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["LocalIPAddress"]); ok {
		bgp.LocalIPAddress = v.(string)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["LocalASN"]); ok {
		bgp.LocalASN = v.(int)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["RemoteIPAddress"]); ok {
		bgp.RemoteIPAddress = v.(string)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["RemoteASN"]); ok {
		bgp.RemoteASN = v.(int)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["AuthenticationKey"]); ok {
		bgp.AuthenticationKey = v.(string)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["State"]); ok {
		bgp.State = v.(string)
	}
	if v, ok := d.GetOk(neBGPSchemaNames["ProvisioningStatus"]); ok {
		bgp.ProvisioningStatus = v.(string)
	}
	return bgp
}

func updateNeBGPResource(bgp *ne.BGPConfiguration, d *schema.ResourceData) error {
	if err := d.Set(neBGPSchemaNames["UUID"], bgp.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["ConnectionUUID"], bgp.ConnectionUUID); err != nil {
		return fmt.Errorf("error reading ConnectionUUID: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["DeviceUUID"], bgp.DeviceUUID); err != nil {
		return fmt.Errorf("error reading DeviceUUID: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["LocalIPAddress"], bgp.LocalIPAddress); err != nil {
		return fmt.Errorf("error reading LocalIPAddress: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["LocalASN"], bgp.LocalASN); err != nil {
		return fmt.Errorf("error reading LocalASN: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["RemoteIPAddress"], bgp.RemoteIPAddress); err != nil {
		return fmt.Errorf("error reading RemoteIPAddress: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["RemoteASN"], bgp.RemoteASN); err != nil {
		return fmt.Errorf("error reading RemoteASN: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["AuthenticationKey"], bgp.AuthenticationKey); err != nil {
		return fmt.Errorf("error reading AuthenticationKey: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["State"], bgp.State); err != nil {
		return fmt.Errorf("error reading State: %s", err)
	}
	if err := d.Set(neBGPSchemaNames["ProvisioningStatus"], bgp.ProvisioningStatus); err != nil {
		return fmt.Errorf("error reading ProvisioningStatus: %s", err)
	}
	return nil
}
