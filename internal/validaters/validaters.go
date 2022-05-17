package validaters

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func StringIsMetroCode() schema.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile("^[A-Z]{2}$"), "MetroCode must consist of two capital letters")
}

func StringIsEmailAddress() schema.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile("^[^ @]+@[^ @]+$"), "not valid email address")
}

func StringIsPortDefinition() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^(([0-9]+(,[0-9]+){0,9})|([0-9]+-[0-9]+)|(any))$"),
		"port definition has to be: up to 10 comma sepparated numbers (22,23), range (20-23) or word 'any'")
}

func StringIsSpeedBand() schema.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile("^[0-9]+(MB|GB)$"), "SpeedBand should consist of digit followed by MB or GB")
}
