package streamsubscription

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clearWhenWriteOnlySet empties a stateful secret whose write-only twin is set; UseStateForUnknown would keep the old value.
func clearWhenWriteOnlySet(writeOnly path.Path) planmodifier.String {
	return clearWhenWriteOnlySetModifier{writeOnly: writeOnly}
}

type clearWhenWriteOnlySetModifier struct {
	writeOnly path.Path
}

func (m clearWhenWriteOnlySetModifier) Description(_ context.Context) string {
	return fmt.Sprintf("Plans an empty value when %s is configured.", m.writeOnly)
}

func (m clearWhenWriteOnlySetModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m clearWhenWriteOnlySetModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var writeOnly types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, m.writeOnly, &writeOnly)...)
	if resp.Diagnostics.HasError() || writeOnly.IsNull() {
		return
	}
	resp.PlanValue = types.StringValue("")
}
