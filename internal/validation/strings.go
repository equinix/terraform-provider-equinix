package validation

import (
	"regexp"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	StringIsMetroCode      = validation.StringMatch(regexp.MustCompile("^[A-Z]{2}$"), "MetroCode must consist of two capital letters")
	StringIsEmailAddress   = validation.StringMatch(regexp.MustCompile("^[^ @]+@[^ @]+$"), "not valid email address")
	StringIsPortDefinition = validation.StringMatch(
		regexp.MustCompile("^(([0-9]+(,[0-9]+){0,9})|([0-9]+-[0-9]+)|(any))$"),
		"port definition has to be: up to 10 comma sepparated numbers (22,23), range (20-23) or word 'any'")
	StringIsSpeedBand      = validation.StringMatch(regexp.MustCompile("^[0-9]+(MB|GB)$"), "SpeedBand should consist of digit followed by MB or GB")
	StringIsCountryCode    = stringvalidator.RegexMatches(StringToRegex("(?i)^[a-z]{2}$"), "Address country must be a two letter code (ISO 3166-1 alpha-2)")
)


func StringToRegex(pattern string) (regExp *regexp.Regexp) {
	regExp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}
	return regExp
}
