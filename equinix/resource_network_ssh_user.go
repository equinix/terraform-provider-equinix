package equinix

import (
	"fmt"
	"log"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var networkSSHUserSchemaNames = map[string]string{
	"UUID":        "uuid",
	"Username":    "username",
	"Password":    "password",
	"DeviceUUIDs": "device_ids",
}

func resourceNetworkSSHUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkSSHUserCreate,
		Read:   resourceNetworkSSHUserRead,
		Update: resourceNetworkSSHUserUpdate,
		Delete: resourceNetworkSSHUserDelete,
		Schema: createNetworkSSHUserResourceSchema(),
	}
}

func createNetworkSSHUserResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkSSHUserSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkSSHUserSchemaNames["Username"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 32),
		},
		networkSSHUserSchemaNames["Password"]: {
			Type:         schema.TypeString,
			Sensitive:    true,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(8, 20),
		},
		networkSSHUserSchemaNames["DeviceUUIDs"]: {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceNetworkSSHUserCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	user := createNetworkSSHUser(d)
	if len(user.DeviceUUIDs) < 0 {
		return fmt.Errorf("create ssh-user failed: user needs to have at least one device defined")
	}
	uuid, err := conf.ne.CreateSSHUser(user.Username, user.Password, user.DeviceUUIDs[0])
	if err != nil {
		return err
	}
	d.SetId(uuid)
	userUpdateReq := conf.ne.NewSSHUserUpdateRequest(uuid)
	userUpdateReq.WithDeviceChange([]string{}, user.DeviceUUIDs[1:len(user.DeviceUUIDs)])
	if err := userUpdateReq.Execute(); err != nil {
		log.Printf("[WARN] failed to assign devices to newly created user")
	}
	return resourceNetworkSSHUserRead(d, m)
}

func resourceNetworkSSHUserRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	user, err := conf.ne.GetSSHUser(d.Id())
	if err != nil {
		return err
	}
	if err := updateNetworkSSHUserResource(user, d); err != nil {
		return err
	}
	return nil
}

func resourceNetworkSSHUserUpdate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	updateReq := conf.ne.NewSSHUserUpdateRequest(d.Id())
	if v, ok := d.GetOk(networkSSHUserSchemaNames["Password"]); ok && d.HasChange(networkSSHUserSchemaNames["Password"]) {
		updateReq.WithNewPassword(v.(string))
	}
	if d.HasChange(networkSSHUserSchemaNames["DeviceUUIDs"]) {
		a, b := d.GetChange(networkSSHUserSchemaNames["DeviceUUIDs"])
		aList := expandSetToStringList(a.(*schema.Set))
		bList := expandSetToStringList(b.(*schema.Set))
		updateReq.WithDeviceChange(aList, bList)
	}
	if err := updateReq.Execute(); err != nil {
		return err
	}
	return resourceNetworkSSHUserRead(d, m)
}

func resourceNetworkSSHUserDelete(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	if err := conf.ne.DeleteSSHUser(d.Id()); err != nil {
		return err
	}
	return nil
}

func createNetworkSSHUser(d *schema.ResourceData) ne.SSHUser {
	user := ne.SSHUser{}
	if v, ok := d.GetOk(networkSSHUserSchemaNames["UUID"]); ok {
		user.UUID = v.(string)
	}
	if v, ok := d.GetOk(networkSSHUserSchemaNames["Username"]); ok {
		user.Username = v.(string)
	}
	if v, ok := d.GetOk(networkSSHUserSchemaNames["Password"]); ok {
		user.Password = v.(string)
	}
	if v, ok := d.GetOk(networkSSHUserSchemaNames["DeviceUUIDs"]); ok {
		user.DeviceUUIDs = expandSetToStringList(v.(*schema.Set))
	}
	return user
}

func updateNetworkSSHUserResource(user *ne.SSHUser, d *schema.ResourceData) error {
	if err := d.Set(networkSSHUserSchemaNames["UUID"], user.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(networkSSHUserSchemaNames["Username"], user.Username); err != nil {
		return fmt.Errorf("error reading Username: %s", err)
	}
	if user.Password != "" {
		if err := d.Set(networkSSHUserSchemaNames["Password"], user.Password); err != nil {
			return fmt.Errorf("error reading Password: %s", err)
		}
	}
	if err := d.Set(networkSSHUserSchemaNames["DeviceUUIDs"], user.DeviceUUIDs); err != nil {
		return fmt.Errorf("error reading DeviceUUIDs: %s", err)
	}
	return nil
}
