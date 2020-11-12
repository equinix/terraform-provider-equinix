package equinix

import (
	"fmt"
	"log"
	"time"

	"github.com/equinix/ecx-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var ecxL2ConnectionAccepterSchemaNames = map[string]string{
	"ConnectionId":    "connection_id",
	"AccessKey":       "access_key",
	"SecretKey":       "secret_key",
	"AWSConnectionID": "aws_connection_id",
}

func resourceECXL2ConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourceECXL2ConnectionAccepterCreate,
		Read:   resourceECXL2ConnectionAccepterRead,
		Delete: resourceECXL2ConnectionAccepterDelete,
		Schema: createECXL2ConnectionAccepterResourceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func createECXL2ConnectionAccepterResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxL2ConnectionAccepterSchemaNames["ConnectionId"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		ecxL2ConnectionAccepterSchemaNames["AccessKey"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		ecxL2ConnectionAccepterSchemaNames["SecretKey"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		ecxL2ConnectionAccepterSchemaNames["AWSConnectionID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceECXL2ConnectionAccepterCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	req := ecx.L2ConnectionToConfirm{}
	connID := d.Get(ecxL2ConnectionAccepterSchemaNames["ConnectionId"]).(string)
	req.AccessKey = d.Get(ecxL2ConnectionAccepterSchemaNames["AccessKey"]).(string)
	req.SecretKey = d.Get(ecxL2ConnectionAccepterSchemaNames["SecretKey"]).(string)

	if _, err := conf.ecx.ConfirmL2Connection(connID, req); err != nil {
		return err
	}
	d.SetId(connID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			ecx.ConnectionStatusPendingApproval,
		},
		Target: []string{
			ecx.ConnectionStatusProvisioned,
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := conf.ecx.GetL2Connection(connID)
			if err != nil {
				return nil, "", err
			}
			return resp, resp.Status, nil
		},
	}
	if _, err := createStateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for connection %q to be approved: %s", connID, err)
	}
	return resourceECXL2ConnectionAccepterRead(d, m)
}

func resourceECXL2ConnectionAccepterRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	conn, err := conf.ecx.GetL2Connection(d.Id())
	if err != nil {
		return err
	}
	if conn == nil || isStringInSlice(conn.Status, []string{
		ecx.ConnectionStatusPendingDelete,
		ecx.ConnectionStatusDeprovisioning,
		ecx.ConnectionStatusDeprovisioned,
		ecx.ConnectionStatusDeleted,
	}) {
		d.SetId("")
		return nil
	}
	if err := updateECXL2ConnectionAccepterResource(conn, d); err != nil {
		return err
	}
	return nil
}

func resourceECXL2ConnectionAccepterDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[WARN] [equinix_ecx_l2_connection_accepter] Will not delete ECX L2 connection (%s)"+
		"Terraform will remove this resource from the state file, however resources may remain.", d.Id())
	return nil
}

func updateECXL2ConnectionAccepterResource(conn *ecx.L2Connection, d *schema.ResourceData) error {
	if err := d.Set(ecxL2ConnectionAccepterSchemaNames["ConnectionId"], conn.UUID); err != nil {
		return fmt.Errorf("error reading connection UUID: %s", err)
	}
	var awsConnectionID string
	for _, action := range conn.Actions {
		if action.OperationID != "CONFIRM_CONNECTION" {
			continue
		}
		for _, actionData := range action.RequiredData {
			if actionData.Key != "awsConnectionId" {
				continue
			}
			awsConnectionID = actionData.Value
		}
	}
	if err := d.Set(ecxL2ConnectionAccepterSchemaNames["AWSConnectionID"], awsConnectionID); err != nil {
		return fmt.Errorf("error reading connection AWSConnectionID: %s", err)
	}
	return nil
}
