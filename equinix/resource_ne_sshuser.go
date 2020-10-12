package equinix

import (
	"fmt"
	"log"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var neSSHUserSchemaNames = map[string]string{
	"UUID":        "uuid",
	"Username":    "username",
	"Password":    "password",
	"DeviceUUIDs": "devices",
}

func resourceNeSSHUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceNeSSHUserCreate,
		Read:   resourceNeSSHUserRead,
		Update: resourceNeSSHUserUpdate,
		Delete: resourceNeSSHUserDelete,
		Schema: createNeSSHUserResourceSchema(),
	}
}

func createNeSSHUserResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neSSHUserSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neSSHUserSchemaNames["Username"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 32),
		},
		neSSHUserSchemaNames["Password"]: {
			Type:         schema.TypeString,
			Sensitive:    true,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(8, 20),
		},
		neSSHUserSchemaNames["DeviceUUIDs"]: {
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

func resourceNeSSHUserCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	user := createNeSSHUser(d)
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
	return resourceNeSSHUserRead(d, m)
}

func resourceNeSSHUserRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	user, err := conf.ne.GetSSHUser(d.Id())
	if err != nil {
		return err
	}
	if err := updateNeSSHUserResource(user, d); err != nil {
		return err
	}
	return nil
}

func resourceNeSSHUserUpdate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	updateReq := conf.ne.NewSSHUserUpdateRequest(d.Id())
	if v, ok := d.GetOk(neSSHUserSchemaNames["Password"]); ok && d.HasChange(neSSHUserSchemaNames["Password"]) {
		updateReq.WithNewPassword(v.(string))
	}
	if d.HasChange(neSSHUserSchemaNames["DeviceUUIDs"]) {
		a, b := d.GetChange(neSSHUserSchemaNames["DeviceUUIDs"])
		aList := expandSetToStringList(a.(*schema.Set))
		bList := expandSetToStringList(b.(*schema.Set))
		updateReq.WithDeviceChange(aList, bList)
	}
	if err := updateReq.Execute(); err != nil {
		return err
	}
	return resourceNeSSHUserRead(d, m)
}

func resourceNeSSHUserDelete(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	if err := conf.ne.DeleteSSHUser(d.Id()); err != nil {
		return err
	}
	return nil
}

func createNeSSHUser(d *schema.ResourceData) ne.SSHUser {
	user := ne.SSHUser{}
	if v, ok := d.GetOk(neSSHUserSchemaNames["UUID"]); ok {
		user.UUID = v.(string)
	}
	if v, ok := d.GetOk(neSSHUserSchemaNames["Username"]); ok {
		user.Username = v.(string)
	}
	if v, ok := d.GetOk(neSSHUserSchemaNames["Password"]); ok {
		user.Password = v.(string)
	}
	if v, ok := d.GetOk(neSSHUserSchemaNames["DeviceUUIDs"]); ok {
		user.DeviceUUIDs = expandSetToStringList(v.(*schema.Set))
	}
	return user
}

func updateNeSSHUserResource(user *ne.SSHUser, d *schema.ResourceData) error {
	if err := d.Set(neSSHUserSchemaNames["UUID"], user.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(neSSHUserSchemaNames["Username"], user.Username); err != nil {
		return fmt.Errorf("error reading Username: %s", err)
	}
	if user.Password != "" {
		if err := d.Set(neSSHUserSchemaNames["Password"], user.Password); err != nil {
			return fmt.Errorf("error reading Password: %s", err)
		}
	}
	if err := d.Set(neSSHUserSchemaNames["DeviceUUIDs"], user.DeviceUUIDs); err != nil {
		return fmt.Errorf("error reading DeviceUUIDs: %s", err)
	}
	return nil
}
