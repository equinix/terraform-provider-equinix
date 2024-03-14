package ssh_user

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"

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

var networkSSHUserDescriptions = map[string]string{
	"UUID":        "SSH user unique identifier",
	"Username":    "SSH user login name",
	"Password":    "SSH user password",
	"DeviceUUIDs": "list of device identifiers to which user will have access",
}

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkSSHUserCreate,
		ReadContext:   resourceNetworkSSHUserRead,
		UpdateContext: resourceNetworkSSHUserUpdate,
		DeleteContext: resourceNetworkSSHUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema:      createNetworkSSHUserResourceSchema(),
		Description: "Resource allows creation and management of Equinix Network Edge SSH users",
	}
}

func createNetworkSSHUserResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkSSHUserSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkSSHUserDescriptions["UUID"],
		},
		networkSSHUserSchemaNames["Username"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 32),
			Description:  networkSSHUserDescriptions["Username"],
		},
		networkSSHUserSchemaNames["Password"]: {
			Type:         schema.TypeString,
			Sensitive:    true,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(8, 20),
			Description:  networkSSHUserDescriptions["Password"],
		},
		networkSSHUserSchemaNames["DeviceUUIDs"]: {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			Description: networkSSHUserDescriptions["DeviceUUIDs"],
		},
	}
}

func resourceNetworkSSHUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)

	var diags diag.Diagnostics
	user := createNetworkSSHUser(d)
	if len(user.DeviceUUIDs) < 0 {
		return diag.Errorf("create ssh-user failed: user needs to have at least one device defined")
	}
	uuid, err := client.CreateSSHUser(ne.StringValue(user.Username), ne.StringValue(user.Password), user.DeviceUUIDs[0])
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ne.StringValue(uuid))
	userUpdateReq := client.NewSSHUserUpdateRequest(ne.StringValue(uuid))
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
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	user, err := client.GetSSHUser(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := updateNetworkSSHUserResource(user, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceNetworkSSHUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	updateReq := client.NewSSHUserUpdateRequest(d.Id())
	if v, ok := d.GetOk(networkSSHUserSchemaNames["Password"]); ok && d.HasChange(networkSSHUserSchemaNames["Password"]) {
		updateReq.WithNewPassword(v.(string))
	}
	if d.HasChange(networkSSHUserSchemaNames["DeviceUUIDs"]) {
		a, b := d.GetChange(networkSSHUserSchemaNames["DeviceUUIDs"])
		aList := converters.SetToStringList(a.(*schema.Set))
		bList := converters.SetToStringList(b.(*schema.Set))
		updateReq.WithDeviceChange(aList, bList)
	}
	if err := updateReq.Execute(); err != nil {
		return diag.FromErr(err)
	}
	diags = append(diags, resourceNetworkSSHUserRead(ctx, d, m)...)
	return diags
}

func resourceNetworkSSHUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	if err := client.DeleteSSHUser(d.Id()); err != nil {
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
		user.DeviceUUIDs = converters.SetToStringList(v.(*schema.Set))
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
