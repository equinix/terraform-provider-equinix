package framework

import (
	"context"
	"fmt"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func GetResourceMeta(
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) *config.Config {
	meta, ok := req.ProviderData.(*config.Config)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *http.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return nil
	}

	return meta
}

// NewBaseResource returns a new instance of the BaseResource
// struct for cleaner initialization.
func NewBaseResource(cfg BaseResourceConfig) BaseResource {
	return BaseResource{
		Config: cfg,
	}
}

// BaseResourceConfig contains all configurable base resource fields.
type BaseResourceConfig struct {
	Name   string
	IDAttr string

	// Optional
	Schema *schema.Schema
}

// BaseResource contains various re-usable fields and methods
// intended for use in resource implementations by composition.
type BaseResource struct {
	Config BaseResourceConfig
	Meta   *config.Config
}

func (r *BaseResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.Meta = GetResourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BaseResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = r.Config.Name
}

func (r *BaseResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	if r.Config.Schema == nil {
		resp.Diagnostics.AddError(
			"Missing Schema",
			"Base resource was not provided a schema. "+
				"Please provide a Schema config attribute or implement, the Schema(...) function.",
		)
		return
	}

	resp.Schema = *r.Config.Schema
}

// ImportState should be overridden for resources with
// complex read logic (e.g. parent ID).
func (r *BaseResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Enforce defaults
	idAttr := r.Config.IDAttr
	if idAttr == "" {
		idAttr = "id"
	}

	attrPath := path.Root(idAttr)

	if attrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState path must be set to a valid attribute path.",
		)
		return
	}

	// Handle type conversion
	var err error
	var idValue any

	idValue = req.ID

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to convert ID attribute",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, idValue)...)
}

// WithTimeouts is intended to be embedded in resources which use the special "timeouts" nested block.
// See https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts.
type WithTimeouts struct {
	defaultCreateTimeout, defaultReadTimeout, defaultUpdateTimeout, defaultDeleteTimeout time.Duration
}

// SetDefaultCreateTimeout sets the resource's default Create timeout value.
func (w *WithTimeouts) SetDefaultCreateTimeout(timeout time.Duration) {
	w.defaultCreateTimeout = timeout
}

// SetDefaultReadTimeout sets the resource's default Read timeout value.
func (w *WithTimeouts) SetDefaultReadTimeout(timeout time.Duration) {
	w.defaultReadTimeout = timeout
}

// SetDefaultUpdateTimeout sets the resource's default Update timeout value.
func (w *WithTimeouts) SetDefaultUpdateTimeout(timeout time.Duration) {
	w.defaultUpdateTimeout = timeout
}

// SetDefaultDeleteTimeout sets the resource's default Delete timeout value.
func (w *WithTimeouts) SetDefaultDeleteTimeout(timeout time.Duration) {
	w.defaultDeleteTimeout = timeout
}

// CreateTimeout returns any configured Create timeout value or the default value.
func (w *WithTimeouts) CreateTimeout(ctx context.Context, timeouts timeouts.Value) time.Duration {
	timeout, diags := timeouts.Create(ctx, w.defaultCreateTimeout)

	if errors := diags.Errors(); len(errors) > 0 {
		tflog.Warn(ctx, "reading configured Create timeout", map[string]interface{}{
			"summary": errors[0].Summary(),
			"detail":  errors[0].Detail(),
		})

		return w.defaultCreateTimeout
	}

	return timeout
}

// ReadTimeout returns any configured Read timeout value or the default value.
func (w *WithTimeouts) ReadTimeout(ctx context.Context, timeouts timeouts.Value) time.Duration {
	timeout, diags := timeouts.Read(ctx, w.defaultReadTimeout)

	if errors := diags.Errors(); len(errors) > 0 {
		tflog.Warn(ctx, "reading configured Read timeout", map[string]interface{}{
			"summary": errors[0].Summary(),
			"detail":  errors[0].Detail(),
		})

		return w.defaultReadTimeout
	}

	return timeout
}

// UpdateTimeout returns any configured Update timeout value or the default value.
func (w *WithTimeouts) UpdateTimeout(ctx context.Context, timeouts timeouts.Value) time.Duration {
	timeout, diags := timeouts.Update(ctx, w.defaultUpdateTimeout)

	if errors := diags.Errors(); len(errors) > 0 {
		tflog.Warn(ctx, "reading configured Update timeout", map[string]interface{}{
			"summary": errors[0].Summary(),
			"detail":  errors[0].Detail(),
		})

		return w.defaultUpdateTimeout
	}

	return timeout
}

// DeleteTimeout returns any configured Delete timeout value or the default value.
func (w *WithTimeouts) DeleteTimeout(ctx context.Context, timeouts timeouts.Value) time.Duration {
	timeout, diags := timeouts.Delete(ctx, w.defaultDeleteTimeout)

	if errors := diags.Errors(); len(errors) > 0 {
		tflog.Warn(ctx, "reading configured Delete timeout", map[string]interface{}{
			"summary": errors[0].Summary(),
			"detail":  errors[0].Detail(),
		})

		return w.defaultDeleteTimeout
	}

	return timeout
}
