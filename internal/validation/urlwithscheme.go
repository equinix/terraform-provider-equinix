package validation

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// urlWithSchemeValidator validates that a string Attribute's value is valid URL with Scheme.
type urlWithSchemeValidator struct {
	schemes []string
}

// Description describes the validation in plain text formatting.
func (validator urlWithSchemeValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be valid URL with host and scheme (%s)", strings.Join(validator.schemes, ", "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator urlWithSchemeValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator urlWithSchemeValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	configValue := request.ConfigValue

	if configValue.IsNull() || configValue.IsUnknown() {
		return
	}

	valueString := configValue.ValueString()

	u, err := url.Parse(valueString)
	if err != nil || u.Host == "" || !slices.Contains(validator.schemes, u.Scheme) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			validator.Description(ctx),
			valueString,
		))
		return
	}
}

// URLWithScheme returns a string validator which ensures that any configured
// attribute value:
//
//   - Is a string, which represents a well-formed URL with host
//     and has a scheme that matches predefined schemes
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func URLWithScheme(schemes ...string) validator.String {
	return urlWithSchemeValidator{
		schemes: schemes,
	}
}
