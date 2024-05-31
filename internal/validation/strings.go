package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	StringIsMetroCode      = validation.StringMatch(regexp.MustCompile("^[A-Z]{2}$"), "MetroCode must consist of two capital letters")
	StringIsEmailAddress   = validation.StringMatch(regexp.MustCompile("^[^ @]+@[^ @]+$"), "not valid email address")
	StringIsPortDefinition = validation.StringMatch(
		regexp.MustCompile("^(([0-9]+(,[0-9]+){0,9})|([0-9]+-[0-9]+)|(any))$"),
		"port definition has to be: up to 10 comma sepparated numbers (22,23), range (20-23) or word 'any'")
	StringIsSpeedBand   = validation.StringMatch(regexp.MustCompile("^[0-9]+(MB|GB)$"), "SpeedBand should consist of digit followed by MB or GB")
	StringIsCountryCode = stringvalidator.RegexMatches(regexp.MustCompile("(?i)^[a-z]{2}$"), "Address country must be a two letter code (ISO 3166-1 alpha-2)")
)

// StringInEnumSlice checks if a string is in a slice of enum
func StringInEnumSlice[T ~string](valid []T, ignoreCase bool) func(i interface{}, k string) (warnings []string, errors []error) {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		for _, item := range valid {
			str := string(item)
			if str == v || (ignoreCase && strings.EqualFold(v, str)) {
				return warnings, errors
			}
		}

		errors = append(errors, fmt.Errorf("expected %s to be one of %q, got %s", k, valid, v))
		return warnings, errors
	}
}
