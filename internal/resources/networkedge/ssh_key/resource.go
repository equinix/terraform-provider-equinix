package ssh_key

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/ne-go"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var networkSSHKeySchemaNames = map[string]string{
	"UUID":      "uuid",
	"Name":      "name",
	"Value":     "public_key",
	"Type":      "type",
	"ProjectID": "project_id",
}

var networkSSHKeyDescriptions = map[string]string{
	"UUID":      "The unique identifier of the key",
	"Name":      "The name of SSH key used for identification",
	"Value":     "The SSH public key. If this is a file, it can be read using the file interpolation function",
	"Type":      "The type of SSH key: RSA (default) or DSA",
	"ProjectID": "The unique identifier of Project Resource to which ssh key is scoped to",
}

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkSSHKeyCreate,
		ReadContext:   resourceNetworkSSHKeyRead,
		DeleteContext: resourceNetworkSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: createNetworkSSHKeyResourceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Description: "Resource allows creation and management of Equinix Network Edge SSH keys",
	}
}

func createNetworkSSHKeyResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkSSHKeySchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkSSHKeyDescriptions["UUID"],
		},
		networkSSHKeySchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  networkSSHKeyDescriptions["Name"],
		},
		networkSSHKeySchemaNames["Value"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  networkSSHKeyDescriptions["Value"],
		},
		networkSSHKeySchemaNames["Type"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "RSA",
			ValidateFunc: validation.StringInSlice([]string{"RSA", "DSA"}, false),
			Description:  networkSSHKeyDescriptions["Type"],
		},
		networkSSHKeySchemaNames["ProjectID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.IsUUID,
			Description:  networkSSHKeyDescriptions["ProjectID"],
		},
	}
}

func resourceNetworkSSHKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	key := createNetworkSSHKey(d)
	uuid, err := client.CreateSSHPublicKey(key)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ne.StringValue(uuid))
	diags = append(diags, resourceNetworkSSHKeyRead(ctx, d, m)...)
	return diags
}

func resourceNetworkSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	key, err := client.GetSSHPublicKey(d.Id())
	if err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if restErr.HTTPCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}
	if err := updateNetworkSSHKeyResource(key, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceNetworkSSHKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	if err := client.DeleteSSHPublicKey(d.Id()); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			for _, detailedErr := range restErr.ApplicationErrors {
				if detailedErr.Code == ne.ErrorCodeSSHPublicKeyInvalid {
					return nil
				}
			}
		}
		return diag.FromErr(err)
	}
	return diags
}

func createNetworkSSHKey(d *schema.ResourceData) ne.SSHPublicKey {
	key := ne.SSHPublicKey{}
	if v, ok := d.GetOk(networkSSHKeySchemaNames["Name"]); ok {
		key.Name = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkSSHKeySchemaNames["Value"]); ok {
		key.Value = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkSSHKeySchemaNames["Type"]); ok {
		key.Type = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkSSHKeySchemaNames["ProjectID"]); ok {
		key.ProjectID = ne.String(v.(string))
	}
	return key
}

func updateNetworkSSHKeyResource(key *ne.SSHPublicKey, d *schema.ResourceData) error {
	if err := d.Set(networkSSHKeySchemaNames["UUID"], key.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(networkSSHKeySchemaNames["Name"], key.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(networkSSHKeySchemaNames["Value"], key.Value); err != nil {
		return fmt.Errorf("error reading Value: %s", err)
	}
	if err := d.Set(networkSSHKeySchemaNames["Type"], key.Type); err != nil {
		return fmt.Errorf("error reading Type: %s", err)
	}
	if err := d.Set(networkSSHKeySchemaNames["ProjectID"], key.ProjectID); err != nil {
		return fmt.Errorf("error reading Type: %s", err)
	}
	return nil
}
