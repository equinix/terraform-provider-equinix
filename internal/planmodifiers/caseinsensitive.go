package planmodifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CaseInsensitiveString() planmodifier.String {
	return &caseInsensitivePlanModifier{}
}

type caseInsensitivePlanModifier struct{}

func (d *caseInsensitivePlanModifier) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	oldValue := request.StateValue.ValueString()
	newValue := request.PlanValue.ValueString()

	result := oldValue
	if !strings.EqualFold(strings.ToLower(newValue), strings.ToLower(oldValue)) {
		result = newValue
	}

	response.PlanValue = types.StringValue(result)
}

func (d *caseInsensitivePlanModifier) Description(ctx context.Context) string {
	return "For same string but different cases, does not trigger diffs in the plan"
}

func (d *caseInsensitivePlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}
