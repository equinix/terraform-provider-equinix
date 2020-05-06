package equinix

import (
	ecx "ecx-go-client/v3"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	//recreateDelay in seconds
	recreateDelay = 5
)

var ecxL2ConnectionSchemaNames = map[string]string{
	"UUID":                "uuid",
	"Name":                "name",
	"ProfileUUID":         "profile_uuid",
	"Speed":               "speed",
	"SpeedUnit":           "speed_unit",
	"Status":              "status",
	"Notifications":       "notifications",
	"PurchaseOrderNumber": "purchase_order_number",
	"PortUUID":            "port_uuid",
	"VlanSTag":            "vlan_stag",
	"VlanCTag":            "vlan_ctag",
	"ZSidePortUUID":       "zside_port_uuid",
	"ZSideVlanSTag":       "zside_vlan_stag",
	"ZSideVlanCTag":       "zside_vlan_ctag",
	"SellerRegion":        "seller_region",
	"SellerMetroCode":     "seller_metro_code",
	"AuthorizationKey":    "authorization_key",
	"RedundantUUID":       "redundant_uuid",
	"SecondaryConnection": "secondary_connection",
}

func resourceECXL2Connection() *schema.Resource {
	return &schema.Resource{
		Create: resourceECXL2ConnectionCreate,
		Read:   resourceECXL2ConnectionRead,
		Update: resourceECXL2ConnectionUpdate,
		Delete: resourceECXL2ConnectionDelete,
		Schema: createECXL2ConnectionResourceSchema(),
	}
}

func createECXL2ConnectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxL2ConnectionSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		ecxL2ConnectionSchemaNames["Name"]: {
			Type:     schema.TypeString,
			Required: true,
		},
		ecxL2ConnectionSchemaNames["ProfileUUID"]: {
			Type:     schema.TypeString,
			Required: true,
		},
		ecxL2ConnectionSchemaNames["Speed"]: {
			Type:     schema.TypeInt,
			Required: true,
		},
		ecxL2ConnectionSchemaNames["SpeedUnit"]: {
			Type:     schema.TypeString,
			Required: true,
		},
		ecxL2ConnectionSchemaNames["Status"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		ecxL2ConnectionSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]: {
			Type:     schema.TypeString,
			Optional: true,
		},
		ecxL2ConnectionSchemaNames["PortUUID"]: {
			Type:     schema.TypeString,
			Required: true,
		},
		ecxL2ConnectionSchemaNames["VlanSTag"]: {
			Type:     schema.TypeInt,
			Required: true,
		},
		ecxL2ConnectionSchemaNames["VlanCTag"]: {
			Type:     schema.TypeInt,
			Optional: true,
		},
		ecxL2ConnectionSchemaNames["ZSidePortUUID"]: {
			Type:     schema.TypeString,
			Optional: true,
		},
		ecxL2ConnectionSchemaNames["ZSideVlanSTag"]: {
			Type:     schema.TypeInt,
			Optional: true,
		},
		ecxL2ConnectionSchemaNames["ZSideVlanCTag"]: {
			Type:     schema.TypeInt,
			Optional: true,
		},
		ecxL2ConnectionSchemaNames["SellerRegion"]: {
			Type:     schema.TypeString,
			Required: true,
		},
		ecxL2ConnectionSchemaNames["SellerMetroCode"]: {
			Type:     schema.TypeString,
			Required: true,
		},
		ecxL2ConnectionSchemaNames["AuthorizationKey"]: {
			Type:     schema.TypeString,
			Optional: true,
		},
		ecxL2ConnectionSchemaNames["RedundantUUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		ecxL2ConnectionSchemaNames["SecondaryConnection"]: {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					ecxL2ConnectionSchemaNames["UUID"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					ecxL2ConnectionSchemaNames["Name"]: {
						Type:     schema.TypeString,
						Required: true,
					},
					ecxL2ConnectionSchemaNames["Status"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					ecxL2ConnectionSchemaNames["PortUUID"]: {
						Type:     schema.TypeString,
						Required: true,
					},
					ecxL2ConnectionSchemaNames["VlanSTag"]: {
						Type:     schema.TypeInt,
						Required: true,
					},
					ecxL2ConnectionSchemaNames["VlanCTag"]: {
						Type:     schema.TypeInt,
						Optional: true,
					},
					ecxL2ConnectionSchemaNames["ZSidePortUUID"]: {
						Type:     schema.TypeString,
						Optional: true,
					},
					ecxL2ConnectionSchemaNames["ZSideVlanSTag"]: {
						Type:     schema.TypeInt,
						Optional: true,
					},
					ecxL2ConnectionSchemaNames["ZSideVlanCTag"]: {
						Type:     schema.TypeInt,
						Optional: true,
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
	log.Printf("Primary's redundant UUID %v", primary.RedundantUUID)
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
	resourceECXL2ConnectionDelete(d, m)
	time.Sleep(recreateDelay * time.Second)
	return resourceECXL2ConnectionCreate(d, m)
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

func hasECXErrorCode(errors []ecx.Error, code string) bool {
	for _, err := range errors {
		if err.ErrorCode == code {
			return true
		}
	}
	return false
}

func setToStringSlice(set *schema.Set) []string {
	return listToStringSlice(set.List())
}

func listToStringSlice(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = fmt.Sprint(v)
	}
	return result
}

func createECXL2Connections(d *schema.ResourceData) (*ecx.L2Connection, *ecx.L2Connection) {
	primary := ecx.L2Connection{}
	primary.UUID = d.Get(ecxL2ConnectionSchemaNames["UUID"]).(string)
	primary.Name = d.Get(ecxL2ConnectionSchemaNames["Name"]).(string)
	primary.ProfileUUID = d.Get(ecxL2ConnectionSchemaNames["ProfileUUID"]).(string)
	primary.Speed = d.Get(ecxL2ConnectionSchemaNames["Speed"]).(int)
	primary.SpeedUnit = d.Get(ecxL2ConnectionSchemaNames["SpeedUnit"]).(string)
	if status, ok := d.GetOk(ecxL2ConnectionSchemaNames["Status"]); ok {
		primary.Status = status.(string)
	}
	primary.Notifications = setToStringSlice(d.Get(ecxL2ConnectionSchemaNames["Notifications"]).(*schema.Set))
	if poNum, ok := d.GetOk(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]); ok {
		primary.PurchaseOrderNumber = poNum.(string)
	}
	primary.PortUUID = d.Get(ecxL2ConnectionSchemaNames["PortUUID"]).(string)
	primary.VlanSTag = d.Get(ecxL2ConnectionSchemaNames["VlanSTag"]).(int)
	if cTag, ok := d.GetOk(ecxL2ConnectionSchemaNames["VlanCTag"]); ok {
		primary.VlanCTag = cTag.(int)
	}
	if zPortUUID, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSidePortUUID"]); ok {
		primary.ZSidePortUUID = zPortUUID.(string)
	}
	if zSTag, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSideVlanSTag"]); ok {
		primary.ZSideVlanSTag = zSTag.(int)
	}
	if zCTag, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSideVlanCTag"]); ok {
		primary.ZSideVlanCTag = zCTag.(int)
	}
	primary.SellerRegion = d.Get(ecxL2ConnectionSchemaNames["SellerRegion"]).(string)
	primary.SellerMetroCode = d.Get(ecxL2ConnectionSchemaNames["SellerMetroCode"]).(string)
	if authKey, ok := d.GetOk(ecxL2ConnectionSchemaNames["AuthorizationKey"]); ok {
		primary.AuthorizationKey = authKey.(string)
	}
	if redID, ok := d.GetOk(ecxL2ConnectionSchemaNames["RedundantUUID"]); ok {
		primary.RedundantUUID = redID.(string)
	}
	if secConns, ok := d.GetOk(ecxL2ConnectionSchemaNames["SecondaryConnection"]); ok {
		secConnsList := secConns.(*schema.Set).List()
		if len(secConnsList) > 0 {
			secConn := secConnsList[0].(map[string]interface{})
			secondary := ecx.L2Connection{}
			secondary.UUID = secConn[ecxL2ConnectionSchemaNames["UUID"]].(string)
			secondary.Name = secConn[ecxL2ConnectionSchemaNames["Name"]].(string)
			if status, ok := secConn[ecxL2ConnectionSchemaNames["Status"]]; ok {
				secondary.Status = status.(string)
			}
			secondary.PortUUID = secConn[ecxL2ConnectionSchemaNames["PortUUID"]].(string)
			secondary.VlanSTag = secConn[ecxL2ConnectionSchemaNames["VlanSTag"]].(int)
			if cTag, ok := secConn[ecxL2ConnectionSchemaNames["VlanCTag"]]; ok {
				secondary.VlanCTag = cTag.(int)
			}
			if zPortUUID, ok := secConn[ecxL2ConnectionSchemaNames["ZSidePortUUID"]]; ok {
				secondary.ZSidePortUUID = zPortUUID.(string)
			}
			if zSTag, ok := secConn[ecxL2ConnectionSchemaNames["ZSideVlanSTag"]]; ok {
				secondary.ZSideVlanSTag = zSTag.(int)
			}
			if zCTag, ok := secConn[ecxL2ConnectionSchemaNames["ZSideVlanCTag"]]; ok {
				secondary.ZSideVlanCTag = zCTag.(int)
			}
			return &primary, &secondary
		}
	}
	return &primary, nil
}

func updateECXL2ConnectionResource(primary *ecx.L2Connection, secondary *ecx.L2Connection, d *schema.ResourceData) error {
	if err := d.Set(ecxL2ConnectionSchemaNames["UUID"], primary.UUID); err != nil {
		return fmt.Errorf("Error reading UUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Name"], primary.Name); err != nil {
		return fmt.Errorf("Error reading Name: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ProfileUUID"], primary.ProfileUUID); err != nil {
		return fmt.Errorf("Error reading ProfileUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Speed"], primary.Speed); err != nil {
		return fmt.Errorf("Error reading Speed: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SpeedUnit"], primary.SpeedUnit); err != nil {
		return fmt.Errorf("Error reading SpeedUnit: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Status"], primary.Status); err != nil {
		return fmt.Errorf("Error reading Status: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Notifications"], primary.Notifications); err != nil {
		return fmt.Errorf("Error reading Notifications: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"], primary.PurchaseOrderNumber); err != nil {
		return fmt.Errorf("Error reading PurchaseOrderNumber: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["PortUUID"], primary.PortUUID); err != nil {
		return fmt.Errorf("Error reading PortUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["VlanSTag"], primary.VlanSTag); err != nil {
		return fmt.Errorf("Error reading VlanSTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["VlanCTag"], primary.VlanCTag); err != nil {
		return fmt.Errorf("Error reading VlanCTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSidePortUUID"], primary.ZSidePortUUID); err != nil {
		return fmt.Errorf("Error reading ZSidePortUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSideVlanSTag"], primary.ZSideVlanSTag); err != nil {
		return fmt.Errorf("Error reading ZSideVlanSTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSideVlanCTag"], primary.ZSideVlanCTag); err != nil {
		return fmt.Errorf("Error reading ZSideVlanCTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SellerRegion"], primary.SellerRegion); err != nil {
		return fmt.Errorf("Error reading SellerRegion: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SellerMetroCode"], primary.SellerMetroCode); err != nil {
		return fmt.Errorf("Error reading SellerMetroCode: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["AuthorizationKey"], primary.AuthorizationKey); err != nil {
		return fmt.Errorf("Error reading AuthorizationKey: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["RedundantUUID"], primary.RedundantUUID); err != nil {
		return fmt.Errorf("Error reading RedundantUUID: %s", err)
	}
	if secondary != nil {
		secConn := make(map[string]interface{})
		secConn[ecxL2ConnectionSchemaNames["UUID"]] = secondary.UUID
		secConn[ecxL2ConnectionSchemaNames["Name"]] = secondary.Name
		secConn[ecxL2ConnectionSchemaNames["Status"]] = secondary.Status
		secConn[ecxL2ConnectionSchemaNames["PortUUID"]] = secondary.PortUUID
		secConn[ecxL2ConnectionSchemaNames["VlanSTag"]] = secondary.VlanSTag
		secConn[ecxL2ConnectionSchemaNames["VlanCTag"]] = secondary.VlanCTag
		secConn[ecxL2ConnectionSchemaNames["ZSidePortUUID"]] = secondary.ZSidePortUUID
		secConn[ecxL2ConnectionSchemaNames["ZSideVlanSTag"]] = secondary.ZSideVlanSTag
		secConn[ecxL2ConnectionSchemaNames["ZSideVlanCTag"]] = secondary.ZSideVlanCTag
		secConnList := []map[string]interface{}{secConn}
		if err := d.Set(ecxL2ConnectionSchemaNames["SecondaryConnection"], secConnList); err != nil {
			return fmt.Errorf("Error reading SecondaryConnection: %s", err)
		}
	}
	return nil
}
