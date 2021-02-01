package equinix

import (
	"context"
	"fmt"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var networkSSHUserSchemaNames = map[string]string{
	"UUID":        "uuid",
	"Username":    "username",
	"Password":    "password",
	"DeviceUUIDs": "device_ids",
}

func resourceNetworkSSHUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkSSHUserCreate,
		ReadContext:   resourceNetworkSSHUserRead,
		UpdateContext: resourceNetworkSSHUserUpdate,
		DeleteContext: resourceNetworkSSHUserDelete,
		Schema:        createNetworkSSHUserResourceSchema(),
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

func resourceNetworkSSHUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	user := createNetworkSSHUser(d)
	if len(user.DeviceUUIDs) < 0 {
		return diag.Errorf("create ssh-user failed: user needs to have at least one device defined")
	}
	uuid, err := conf.ne.CreateSSHUser(ne.StringValue(user.Username), ne.StringValue(user.Password), user.DeviceUUIDs[0])
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ne.StringValue(uuid))
	userUpdateReq := conf.ne.NewSSHUserUpdateRequest(ne.StringValue(uuid))
	userUpdateReq.WithDeviceChange([]string{}, user.DeviceUUIDs[1:len(user.DeviceUUIDs)])
	if err := userUpdateReq.Execute(); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       "Failed to assign all devices to newly created user",
			Detail:        err.Error(),
			AttributePath: cty.GetAttrPath(networkSSHUserSchemaNames["DeviceUUIDs"]),
		})
	}
	diags = append(diags, resourceNetworkSSHUserRead(ctx, d, m)...)
	return diags
}

func resourceNetworkSSHUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	user, err := conf.ne.GetSSHUser(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := updateNetworkSSHUserResource(user, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceNetworkSSHUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
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
		return diag.FromErr(err)
	}
	diags = append(diags, resourceNetworkSSHUserRead(ctx, d, m)...)
	return diags
}

func resourceNetworkSSHUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	if err := conf.ne.DeleteSSHUser(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createNetworkSSHUser(d *schema.ResourceData) ne.SSHUser {
	user := ne.SSHUser{}
	if v, ok := d.GetOk(networkSSHUserSchemaNames["UUID"]); ok {
		user.UUID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkSSHUserSchemaNames["Username"]); ok {
		user.Username = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkSSHUserSchemaNames["Password"]); ok {
		user.Password = ne.String(v.(string))
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
	if ne.StringValue(user.Password) != "" {
		if err := d.Set(networkSSHUserSchemaNames["Password"], user.Password); err != nil {
			return fmt.Errorf("error reading Password: %s", err)
		}
	}
	if err := d.Set(networkSSHUserSchemaNames["DeviceUUIDs"], user.DeviceUUIDs); err != nil {
		return fmt.Errorf("error reading DeviceUUIDs: %s", err)
	}
	return nil
}
