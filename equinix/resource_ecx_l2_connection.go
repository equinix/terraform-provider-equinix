package equinix

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/equinix/ecx-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var ecxL2ConnectionSchemaNames = map[string]string{
	"UUID":                "uuid",
	"Name":                "name",
	"ProfileUUID":         "profile_uuid",
	"Speed":               "speed",
	"SpeedUnit":           "speed_unit",
	"Status":              "status",
	"ProviderStatus":      "provider_status",
	"Notifications":       "notifications",
	"PurchaseOrderNumber": "purchase_order_number",
	"PortUUID":            "port_uuid",
	"DeviceUUID":          "device_uuid",
	"VlanSTag":            "vlan_stag",
	"VlanCTag":            "vlan_ctag",
	"NamedTag":            "named_tag",
	"AdditionalInfo":      "additional_info",
	"ZSidePortUUID":       "zside_port_uuid",
	"ZSideVlanSTag":       "zside_vlan_stag",
	"ZSideVlanCTag":       "zside_vlan_ctag",
	"SellerRegion":        "seller_region",
	"SellerMetroCode":     "seller_metro_code",
	"AuthorizationKey":    "authorization_key",
	"RedundantUUID":       "redundant_uuid",
	"RedundancyType":      "redundancy_type",
	"SecondaryConnection": "secondary_connection",
}

var ecxL2ConnectionAdditionalInfoSchemaNames = map[string]string{
	"Name":  "name",
	"Value": "value",
}

func resourceECXL2Connection() *schema.Resource {
	return &schema.Resource{
		Create: resourceECXL2ConnectionCreate,
		Read:   resourceECXL2ConnectionRead,
		Update: resourceECXL2ConnectionUpdate,
		Delete: resourceECXL2ConnectionDelete,
		Schema: createECXL2ConnectionResourceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func createECXL2ConnectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxL2ConnectionSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		ecxL2ConnectionSchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 24),
		},
		ecxL2ConnectionSchemaNames["ProfileUUID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{ecxL2ConnectionSchemaNames["ProfileUUID"], ecxL2ConnectionSchemaNames["ZSidePortUUID"]},
			ValidateFunc: validation.StringIsNotEmpty,
		},
		ecxL2ConnectionSchemaNames["Speed"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		ecxL2ConnectionSchemaNames["SpeedUnit"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"MB", "GB"}, false),
		},
		ecxL2ConnectionSchemaNames["Status"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		ecxL2ConnectionSchemaNames["ProviderStatus"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		ecxL2ConnectionSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^ @]+@[^ @]+$"), "Notification list can contain only valid email addresses"),
			},
		},
		ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 30),
		},
		ecxL2ConnectionSchemaNames["PortUUID"]: {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			AtLeastOneOf:  []string{ecxL2ConnectionSchemaNames["PortUUID"], ecxL2ConnectionSchemaNames["DeviceUUID"]},
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"]},
		},
		ecxL2ConnectionSchemaNames["DeviceUUID"]: {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			AtLeastOneOf:  []string{ecxL2ConnectionSchemaNames["PortUUID"], ecxL2ConnectionSchemaNames["DeviceUUID"]},
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["PortUUID"]},
		},
		ecxL2ConnectionSchemaNames["VlanSTag"]: {
			Type:          schema.TypeInt,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ValidateFunc:  validation.IntBetween(2, 4092),
			RequiredWith:  []string{ecxL2ConnectionSchemaNames["PortUUID"]},
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"]},
		},
		ecxL2ConnectionSchemaNames["VlanCTag"]: {
			Type:          schema.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.IntBetween(2, 4092),
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"]},
		},
		ecxL2ConnectionSchemaNames["NamedTag"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"Private", "Public", "Microsoft", "Manual"}, false),
		},
		ecxL2ConnectionSchemaNames["AdditionalInfo"]: {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					ecxL2ConnectionAdditionalInfoSchemaNames["Name"]: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					ecxL2ConnectionAdditionalInfoSchemaNames["Value"]: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		ecxL2ConnectionSchemaNames["ZSidePortUUID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		ecxL2ConnectionSchemaNames["ZSideVlanSTag"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(2, 4092),
		},
		ecxL2ConnectionSchemaNames["ZSideVlanCTag"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(2, 4092),
		},
		ecxL2ConnectionSchemaNames["SellerRegion"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		ecxL2ConnectionSchemaNames["SellerMetroCode"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile("^[A-Z]{2}$"), "MetroCode must consist of two capital letters"),
		},
		ecxL2ConnectionSchemaNames["AuthorizationKey"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		ecxL2ConnectionSchemaNames["RedundantUUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		ecxL2ConnectionSchemaNames["RedundancyType"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		ecxL2ConnectionSchemaNames["SecondaryConnection"]: {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					ecxL2ConnectionSchemaNames["UUID"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					ecxL2ConnectionSchemaNames["Name"]: {
						Type:         schema.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 24),
					},
					ecxL2ConnectionSchemaNames["Status"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					ecxL2ConnectionSchemaNames["ProviderStatus"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					ecxL2ConnectionSchemaNames["PortUUID"]: {
						Type:          schema.TypeString,
						ForceNew:      true,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						AtLeastOneOf:  []string{ecxL2ConnectionSchemaNames["PortUUID"], ecxL2ConnectionSchemaNames["DeviceUUID"]},
						ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"]},
					},
					ecxL2ConnectionSchemaNames["DeviceUUID"]: {
						Type:          schema.TypeString,
						ForceNew:      true,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						AtLeastOneOf:  []string{ecxL2ConnectionSchemaNames["PortUUID"], ecxL2ConnectionSchemaNames["DeviceUUID"]},
						ConflictsWith: []string{ecxL2ConnectionSchemaNames["PortUUID"]},
					},
					ecxL2ConnectionSchemaNames["VlanSTag"]: {
						Type:          schema.TypeInt,
						ForceNew:      true,
						Optional:      true,
						ValidateFunc:  validation.IntBetween(2, 4092),
						RequiredWith:  []string{ecxL2ConnectionSchemaNames["PortUUID"]},
						ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"]},
					},
					ecxL2ConnectionSchemaNames["VlanCTag"]: {
						Type:          schema.TypeInt,
						ForceNew:      true,
						Optional:      true,
						ValidateFunc:  validation.IntBetween(2, 4092),
						ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"]},
					},
					ecxL2ConnectionSchemaNames["ZSidePortUUID"]: {
						Type:         schema.TypeString,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					ecxL2ConnectionSchemaNames["ZSideVlanSTag"]: {
						Type:         schema.TypeInt,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: validation.IntBetween(2, 4092),
					},
					ecxL2ConnectionSchemaNames["ZSideVlanCTag"]: {
						Type:         schema.TypeInt,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: validation.IntBetween(2, 4092),
					},
					ecxL2ConnectionSchemaNames["RedundantUUID"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					ecxL2ConnectionSchemaNames["RedundancyType"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func resourceECXL2ConnectionCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	primary, secondary := createECXL2Connections(d)
	var resp *ecx.L2Connection
	var err error
	if secondary != nil {
		resp, err = conf.ecx.CreateL2RedundantConnection(*primary, *secondary)
	} else {
		resp, err = conf.ecx.CreateL2Connection(*primary)
	}
	if err != nil {
		return err
	}
	d.SetId(resp.UUID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			ecx.ConnectionStatusProvisioning,
			ecx.ConnectionStatusPendingAutoApproval,
		},
		Target: []string{
			ecx.ConnectionStatusProvisioned,
			ecx.ConnectionStatusPendingApproval,
			ecx.ConnectionStatusPendingBGPPeering,
			ecx.ConnectionStatusPendingProviderVlan,
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 2 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := conf.ecx.GetL2Connection(d.Id())
			if err != nil {
				return nil, "", err
			}
			return resp, resp.Status, nil
		},
	}
	if _, err := createStateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for connection (%s) to be created: %s", d.Id(), err)
	}
	return resourceECXL2ConnectionRead(d, m)
}

func resourceECXL2ConnectionRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	var err error
	var primary *ecx.L2Connection
	var secondary *ecx.L2Connection

	primary, err = conf.ecx.GetL2Connection(d.Id())
	if err != nil {
		return fmt.Errorf("cannot fetch primary connection due to %v", err)
	}
	if isStringInSlice(primary.Status, []string{
		ecx.ConnectionStatusPendingDelete,
		ecx.ConnectionStatusDeprovisioning,
		ecx.ConnectionStatusDeprovisioned,
		ecx.ConnectionStatusDeleted,
	}) {
		d.SetId("")
		return nil
	}
	if primary.RedundantUUID != "" {
		secondary, err = conf.ecx.GetL2Connection(primary.RedundantUUID)
		if err != nil {
			return fmt.Errorf("cannot fetch secondary connection due to %v", err)
		}
	}
	if err := updateECXL2ConnectionResource(primary, secondary, d); err != nil {
		return err
	}
	return nil
}

func resourceECXL2ConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	updateReq := conf.ecx.NewL2ConnectionUpdateRequest(d.Id())
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["Name"]); ok && d.HasChange(ecxL2ConnectionSchemaNames["Name"]) {
		updateReq.WithName(v.(string))
	}
	if d.HasChanges(ecxL2ConnectionSchemaNames["Speed"], ecxL2ConnectionSchemaNames["SpeedUnit"]) {
		updateReq.WithBandwidth(d.Get(ecxL2ConnectionSchemaNames["Speed"]).(int),
			d.Get(ecxL2ConnectionSchemaNames["SpeedUnit"]).(string))
	}
	if err := updateReq.Execute(); err != nil {
		return err
	}
	return resourceECXL2ConnectionRead(d, m)
}

func resourceECXL2ConnectionDelete(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	if err := conf.ecx.DeleteL2Connection(d.Id()); err != nil {
		ecxRestErr, ok := err.(ecx.RestError)
		if ok {
			//IC-LAYER2-4021 = Connection already deleted
			if hasECXErrorCode(ecxRestErr.Errors, "IC-LAYER2-4021") {
				return nil
			}
		}
		return err
	}
	//remove secondary connection, don't fail on error as there is no partial state on delete
	if redID, ok := d.GetOk(ecxL2ConnectionSchemaNames["RedundantUUID"]); ok {
		if err := conf.ecx.DeleteL2Connection(redID.(string)); err != nil {
			log.Printf("[WARN] error removing secondary connection with UUID %s, due to %s", redID.(string), err)
		}
	}
	return nil
}

func createECXL2Connections(d *schema.ResourceData) (*ecx.L2Connection, *ecx.L2Connection) {
	primary := ecx.L2Connection{}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["UUID"]); ok {
		primary.UUID = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["Name"]); ok {
		primary.Name = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ProfileUUID"]); ok {
		primary.ProfileUUID = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["Speed"]); ok {
		primary.Speed = v.(int)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SpeedUnit"]); ok {
		primary.SpeedUnit = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["Status"]); ok {
		primary.Status = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ProviderStatus"]); ok {
		primary.ProviderStatus = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["Notifications"]); ok {
		primary.Notifications = expandSetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]); ok {
		primary.PurchaseOrderNumber = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["PortUUID"]); ok {
		primary.PortUUID = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["DeviceUUID"]); ok {
		primary.DeviceUUID = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["VlanSTag"]); ok {
		primary.VlanSTag = v.(int)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["VlanCTag"]); ok {
		primary.VlanCTag = v.(int)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["NamedTag"]); ok {
		primary.NamedTag = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["AdditionalInfo"]); ok {
		primary.AdditionalInfo = expandECXL2ConnectionAdditionalInfo(v.(*schema.Set))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSidePortUUID"]); ok {
		primary.ZSidePortUUID = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSideVlanSTag"]); ok {
		primary.ZSideVlanSTag = v.(int)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSideVlanCTag"]); ok {
		primary.ZSideVlanCTag = v.(int)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SellerRegion"]); ok {
		primary.SellerRegion = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SellerMetroCode"]); ok {
		primary.SellerMetroCode = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["AuthorizationKey"]); ok {
		primary.AuthorizationKey = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["RedundantUUID"]); ok {
		primary.RedundantUUID = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["RedundancyType"]); ok {
		primary.RedundancyType = v.(string)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SecondaryConnection"]); ok {
		secConnSet := v.(*schema.Set)
		if secConnSet.Len() > 0 {
			secConns := expandECXL2ConnectionSecondary(secConnSet)
			return &primary, &secConns[0]
		}
	}
	return &primary, nil
}

func updateECXL2ConnectionResource(primary *ecx.L2Connection, secondary *ecx.L2Connection, d *schema.ResourceData) error {
	if err := d.Set(ecxL2ConnectionSchemaNames["UUID"], primary.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Name"], primary.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ProfileUUID"], primary.ProfileUUID); err != nil {
		return fmt.Errorf("error reading ProfileUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Speed"], primary.Speed); err != nil {
		return fmt.Errorf("error reading Speed: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SpeedUnit"], primary.SpeedUnit); err != nil {
		return fmt.Errorf("error reading SpeedUnit: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Status"], primary.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ProviderStatus"], primary.ProviderStatus); err != nil {
		return fmt.Errorf("error reading ProviderStatus: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Notifications"], primary.Notifications); err != nil {
		return fmt.Errorf("error reading Notifications: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"], primary.PurchaseOrderNumber); err != nil {
		return fmt.Errorf("error reading PurchaseOrderNumber: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["PortUUID"], primary.PortUUID); err != nil {
		return fmt.Errorf("error reading PortUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["DeviceUUID"], primary.DeviceUUID); err != nil {
		return fmt.Errorf("error reading PortUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["VlanSTag"], primary.VlanSTag); err != nil {
		return fmt.Errorf("error reading VlanSTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["VlanCTag"], primary.VlanCTag); err != nil {
		return fmt.Errorf("error reading VlanCTag: %s", err)
	}
	if primary.NamedTag != "" {
		if err := d.Set(ecxL2ConnectionSchemaNames["NamedTag"], primary.NamedTag); err != nil {
			return fmt.Errorf("error reading NamedTag: %s", err)
		}
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["AdditionalInfo"], flattenECXL2ConnectionAdditionalInfo(primary.AdditionalInfo)); err != nil {
		return fmt.Errorf("error reading AdditionalInfo: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSidePortUUID"], primary.ZSidePortUUID); err != nil {
		return fmt.Errorf("error reading ZSidePortUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSideVlanSTag"], primary.ZSideVlanSTag); err != nil {
		return fmt.Errorf("error reading ZSideVlanSTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSideVlanCTag"], primary.ZSideVlanCTag); err != nil {
		return fmt.Errorf("error reading ZSideVlanCTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SellerRegion"], primary.SellerRegion); err != nil {
		return fmt.Errorf("error reading SellerRegion: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SellerMetroCode"], primary.SellerMetroCode); err != nil {
		return fmt.Errorf("error reading SellerMetroCode: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["AuthorizationKey"], primary.AuthorizationKey); err != nil {
		return fmt.Errorf("error reading AuthorizationKey: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["RedundantUUID"], primary.RedundantUUID); err != nil {
		return fmt.Errorf("error reading RedundantUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["RedundancyType"], primary.RedundancyType); err != nil {
		return fmt.Errorf("error reading RedundancyType: %s", err)
	}
	if secondary != nil {
		var prev *ecx.L2Connection
		if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SecondaryConnection"]); ok {
			vSet := v.(*schema.Set)
			if vSet.Len() > 0 {
				prev = &expandECXL2ConnectionSecondary(vSet)[0]
			}
		}
		if err := d.Set(ecxL2ConnectionSchemaNames["SecondaryConnection"], flattenECXL2ConnectionSecondary(prev, secondary)); err != nil {
			return fmt.Errorf("error reading SecondaryConnection: %s", err)
		}
	}
	return nil
}

func flattenECXL2ConnectionSecondary(prev, conn *ecx.L2Connection) interface{} {
	transformed := make(map[string]interface{})
	transformed[ecxL2ConnectionSchemaNames["UUID"]] = conn.UUID
	transformed[ecxL2ConnectionSchemaNames["Name"]] = conn.Name
	transformed[ecxL2ConnectionSchemaNames["Status"]] = conn.Status
	transformed[ecxL2ConnectionSchemaNames["ProviderStatus"]] = conn.ProviderStatus
	transformed[ecxL2ConnectionSchemaNames["PortUUID"]] = conn.PortUUID
	transformed[ecxL2ConnectionSchemaNames["DeviceUUID"]] = conn.DeviceUUID
	if prev == nil || (prev != nil && prev.VlanSTag != 0) {
		transformed[ecxL2ConnectionSchemaNames["VlanSTag"]] = conn.VlanSTag
	}
	transformed[ecxL2ConnectionSchemaNames["VlanCTag"]] = conn.VlanCTag
	if prev == nil || (prev != nil && prev.ZSidePortUUID != "") {
		transformed[ecxL2ConnectionSchemaNames["ZSidePortUUID"]] = conn.ZSidePortUUID
	}
	if prev == nil || (prev != nil && prev.ZSideVlanSTag != 0) {
		transformed[ecxL2ConnectionSchemaNames["ZSideVlanSTag"]] = conn.ZSideVlanSTag
	}
	if prev == nil || (prev != nil && prev.ZSideVlanCTag != 0) {
		transformed[ecxL2ConnectionSchemaNames["ZSideVlanCTag"]] = conn.ZSideVlanCTag
	}
	transformed[ecxL2ConnectionSchemaNames["RedundantUUID"]] = conn.RedundantUUID
	transformed[ecxL2ConnectionSchemaNames["RedundancyType"]] = conn.RedundancyType
	return []map[string]interface{}{transformed}
}

func expandECXL2ConnectionSecondary(connections *schema.Set) []ecx.L2Connection {
	transformed := make([]ecx.L2Connection, 0, connections.Len())
	for _, conn := range connections.List() {
		connMap := conn.(map[string]interface{})
		c := ecx.L2Connection{}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["UUID"]]; ok {
			c.UUID = v.(string)
		}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["Name"]]; ok {
			c.Name = v.(string)
		}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["Status"]]; ok {
			c.Status = v.(string)
		}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["ProviderStatus"]]; ok {
			c.ProviderStatus = v.(string)
		}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["PortUUID"]]; ok {
			c.PortUUID = v.(string)
		}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["DeviceUUID"]]; ok {
			c.DeviceUUID = v.(string)
		}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["VlanSTag"]]; ok {
			c.VlanSTag = v.(int)
		}
		if cTag, ok := connMap[ecxL2ConnectionSchemaNames["VlanCTag"]]; ok {
			c.VlanCTag = cTag.(int)
		}
		if zPortUUID, ok := connMap[ecxL2ConnectionSchemaNames["ZSidePortUUID"]]; ok {
			c.ZSidePortUUID = zPortUUID.(string)
		}
		if zSTag, ok := connMap[ecxL2ConnectionSchemaNames["ZSideVlanSTag"]]; ok {
			c.ZSideVlanSTag = zSTag.(int)
		}
		if zCTag, ok := connMap[ecxL2ConnectionSchemaNames["ZSideVlanCTag"]]; ok {
			c.ZSideVlanCTag = zCTag.(int)
		}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["RedundantUUID"]]; ok {
			c.RedundantUUID = v.(string)
		}
		if v, ok := connMap[ecxL2ConnectionSchemaNames["RedundancyType"]]; ok {
			c.RedundancyType = v.(string)
		}
		transformed = append(transformed, c)
	}
	return transformed
}

func flattenECXL2ConnectionAdditionalInfo(infos []ecx.L2ConnectionAdditionalInfo) interface{} {
	transformed := make([]interface{}, 0, len(infos))
	for _, info := range infos {
		transformed = append(transformed, map[string]interface{}{
			ecxL2ConnectionAdditionalInfoSchemaNames["Name"]:  info.Name,
			ecxL2ConnectionAdditionalInfoSchemaNames["Value"]: info.Value,
		})
	}
	return transformed
}

func expandECXL2ConnectionAdditionalInfo(infos *schema.Set) []ecx.L2ConnectionAdditionalInfo {
	transformed := make([]ecx.L2ConnectionAdditionalInfo, 0, infos.Len())
	for _, info := range infos.List() {
		infoMap := info.(map[string]interface{})
		transformed = append(transformed, ecx.L2ConnectionAdditionalInfo{
			Name:  infoMap[ecxL2ConnectionAdditionalInfoSchemaNames["Name"]].(string),
			Value: infoMap[ecxL2ConnectionAdditionalInfoSchemaNames["Value"]].(string),
		})
	}
	return transformed
}
