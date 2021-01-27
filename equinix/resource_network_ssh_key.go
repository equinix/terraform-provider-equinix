package equinix

import (
	"fmt"
	"net/http"
	"time"

	"github.com/equinix/ne-go"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var networkSSHKeySchemaNames = map[string]string{
	"UUID":  "uuid",
	"Name":  "name",
	"Value": "public_key",
}

func resourceNetworkSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkSSHKeyCreate,
		Read:   resourceNetworkSSHKeyRead,
		Delete: resourceNetworkSSHKeyDelete,
		Schema: createNetworkSSHKeyResourceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func createNetworkSSHKeyResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkSSHKeySchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkSSHKeySchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		networkSSHKeySchemaNames["Value"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func resourceNetworkSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	key := createNetworkSSHKey(d)
	uuid, err := conf.ne.CreateSSHPublicKey(key)
	if err != nil {
		return err
	}
	d.SetId(ne.StringValue(uuid))
	return resourceNetworkSSHKeyRead(d, m)
}

func resourceNetworkSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	key, err := conf.ne.GetSSHPublicKey(d.Id())
	if err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if restErr.HTTPCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}
		return err
	}
	if err := updateNetworkSSHKeyResource(key, d); err != nil {
		return err
	}
	return nil
}

func resourceNetworkSSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	if err := conf.ne.DeleteSSHPublicKey(d.Id()); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			for _, detailedErr := range restErr.ApplicationErrors {
				if detailedErr.Code == ne.ErrorCodeSSHPublicKeyInvalid {
					return nil
				}
			}
		}
		return err
	}
	return nil
}

func createNetworkSSHKey(d *schema.ResourceData) ne.SSHPublicKey {
	key := ne.SSHPublicKey{}
	if v, ok := d.GetOk(networkSSHKeySchemaNames["Name"]); ok {
		key.Name = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkSSHKeySchemaNames["Value"]); ok {
		key.Value = ne.String(v.(string))
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
	return nil
}
