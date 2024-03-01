package validation

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// uuidValidator validates that a string Attribute's value is valid UUID.
type uuidValidator struct{}

// Description describes the validation in plain text formatting.
func (validator uuidValidator) Description(_ context.Context) string {
	return "value must be valid UUID"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator uuidValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator uuidValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	configValue := request.ConfigValue

	if configValue.IsNull() || configValue.IsUnknown() {
		return
	}

	valueString := configValue.ValueString()

	if _, err := uuid.ParseUUID(valueString); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			validator.Description(ctx),
			valueString,
		))
		return
	}
}

// UUID returns a string validator which ensures that any configured
// attribute value:
//
//   - Is a string, which represents valid UUID.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func UUID() validator.String {
	return uuidValidator{}
}
