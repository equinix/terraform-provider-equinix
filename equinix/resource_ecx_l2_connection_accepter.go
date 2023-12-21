package equinix

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ecxL2ConnectionAccepterSchemaNames = map[string]string{
	"ConnectionId":    "connection_id",
	"AccessKey":       "access_key",
	"SecretKey":       "secret_key",
	"Profile":         "aws_profile",
	"AWSConnectionID": "aws_connection_id",
}

var ecxL2ConnectionAccepterDescriptions = map[string]string{
	"ConnectionId":    "Identifier of layer 2 connection that will be accepted",
	"AccessKey":       "Access Key used to accept connection on provider side",
	"SecretKey":       "Secret Key used to accept connection on provider side",
	"Profile":         "AWS Profile Name for retrieving credentials from shared credentials file",
	"AWSConnectionID": "Identifier of a hosted Direct Connect connection on AWS side, applicable for accepter resource with connections to AWS only",
}

func ResourceECXL2ConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated and will be removed in a future version of the provider. Please use the equivalent resource `aws_dx_connection_confirmation` in the 'AWS' provider instead.",
		CreateContext:      resourceECXL2ConnectionAccepterCreate,
		ReadContext:        resourceECXL2ConnectionAccepterRead,
		DeleteContext:      resourceECXL2ConnectionAccepterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema:      createECXL2ConnectionAccepterResourceSchema(),
		Description: "Resource is used to accept Equinix Fabric layer 2 connection on provider side",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
			Description:  ecxL2ConnectionAccepterDescriptions["ConnectionId"],
		},
		ecxL2ConnectionAccepterSchemaNames["AccessKey"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionAccepterDescriptions["AccessKey"],
		},
		ecxL2ConnectionAccepterSchemaNames["SecretKey"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionAccepterDescriptions["SecretKey"],
		},
		ecxL2ConnectionAccepterSchemaNames["Profile"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionAccepterDescriptions["Profile"],
		},
		ecxL2ConnectionAccepterSchemaNames["AWSConnectionID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionAccepterDescriptions["AWSConnectionID"],
		},
	}
}

func resourceECXL2ConnectionAccepterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Errorf("This resource is deprecated and will be removed in a future version of the provider. Please use the equivalent resource `aws_dx_connection_confirmation` in the 'AWS' provider instead.")
}

func resourceECXL2ConnectionAccepterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Errorf("This resource is deprecated and will be removed in a future version of the provider. Please use the equivalent resource `aws_dx_connection_confirmation` in the 'AWS' provider instead.")
}

func resourceECXL2ConnectionAccepterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Errorf("This resource is deprecated and will be removed in a future version of the provider. Please use the equivalent resource `aws_dx_connection_confirmation` in the 'AWS' provider instead.")
}
