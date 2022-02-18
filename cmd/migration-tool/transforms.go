package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// matches block headers, ex:
//   resource "metal_project" "fooproject" {
//   data "packet_vlan" "foovlan" {
var matchBlockHeader = regexp.MustCompile(`(resource|data)(\s+")(metal|packet)(.*?)`)

// matches resource interpolation strings (Terraform v0.11 and earlier), ex:
//   device_id = "${metal_device.foodevice.id}"
var matchResourceInterpolation = regexp.MustCompile(`(.*?)(\${\s*)(metal|packet)(_.*?)`)

// matches resource reference (Terraform v0.12+), ex:
//   device_id = metal_device.foodevice.id
var matchResourceReference = regexp.MustCompile(`(.*?)(=\s*)(metal|packet)(_.*?)`)

// matches resource reference in function, ex:
//   cidr_notation = join("/", [cidrhost(metal_reserved_ip_block.fooblock.cidr_notation, 0), "32"])
var matchResourceFunction = regexp.MustCompile(`(.*?)(\(\s*)(metal|packet)(_.*?)`)

// matches resource reference in conditional, ex:
//   ip_address = "${var.network_type == "public" ? metal_device.foodevice.access_public_ipv4 : metal_device.foodevice.access_private_ipv4}"
//   ip_address = var.network_type == "public" ? metal_device.foodevice.access_public_ipv4 : metal_device.foodevice.access_private_ipv4
var matchResourceConditional = regexp.MustCompile(`(.*?[:|\?])(\s*)(metal|packet)(_.*?)`)

// matches resource reference in for loop,ex:
//   toset([for network in metal_device.foodevice.network : network.family])
var matchResourceForLoop = regexp.MustCompile(`(.*?)(in\s*)(metal|packet)(_.*?)`)

// matches resource in expression,ex:
//   tolist([metal_device.foodevice[*].access_public_ipv4])
//   !metal_ip_attachment.fooattach.public
//   totalSpeed = metal_connection.fooconnA.speed + metal_connection.fooconnB.speed
var matchResourceExpression = regexp.MustCompile(`(.*?[\+|-|\*|\/|>|<|&|\|\||%|!|\[]\s*)(metal|packet)(_.*?)`)

// matches datasource references, ex:
//   address_family = "${lookup(data.packet_device_bgp_neighbors.test.bgp_neighbors[0], "address_family")}"
var matchDatasourceReference = regexp.MustCompile(`(.*?data)(\.)(metal|packet)(_.*?)`)

// replace specific string patterns in template files
func replaceTemplateTokens(str string) string {
	// resources
	str = matchBlockHeader.ReplaceAllString(str, `$1 "equinix_metal$4`)
	str = matchResourceInterpolation.ReplaceAllString(str, `$1${equinix_metal$4`)
	str = matchResourceReference.ReplaceAllString(str, `${1}= equinix_metal$4`)
	str = matchResourceFunction.ReplaceAllString(str, `$1(equinix_metal$4`)
	str = matchResourceConditional.ReplaceAllString(str, `$1 equinix_metal$4`)
	str = matchResourceForLoop.ReplaceAllString(str, `${1}in equinix_metal$4`)
	str = matchResourceExpression.ReplaceAllString(str, `${1}equinix_metal$3`)
	// datasources
	return matchDatasourceReference.ReplaceAllString(str, `$1.equinix_metal$4`)
}

// matches '"metal_' or '"packet_' prefixes in statefile
var matchStatePrefixes = regexp.MustCompile(`(.*")(metal|packet)(_.*)`)

// matches provider url in statefile
var matchStateProvider = regexp.MustCompile(`(.*?)(equinix\/metal|packethost\/packet)(\\".*?)`)

// replace metal|packet in statefile
func replaceStatefileTokens(str string) string {
	// provider
	str = matchStateProvider.ReplaceAllString(str, `${1}equinix/equinix$3`)
	// datasources
	str = matchDatasourceReference.ReplaceAllString(str, `$1.equinix_metal$4`)
	// metal and prefixes
	return matchStatePrefixes.ReplaceAllString(str, `${1}equinix_metal$3`)
}

// rewrite matching required provider to have equinix provider with no version
func updateRequiredProvider(content string) (string, error) {
	idx, _ := findToken("metal", content)
	if idx == -1 {
		idx, _ = findToken("packet", content)
	}
	if idx == -1 {
		return content, nil
	}

	subStr := content[idx:] // ignore everything before metal/packet provider

	// replace provider name
	subStr = strings.Replace(subStr, "metal", "equinix", 1)
	subStr = strings.Replace(subStr, "packet", "equinix", 1)

	blockStart, blockEnd := indexOpenCloseTokens('{', '}', subStr) // limit search to logical provider block
	if blockStart == -1 || blockEnd == -1 {
		return content, fmt.Errorf("required Provider metal/packet block start or end not detected")
	}

	blkContents := subStr[:blockEnd] // get just from provider name to the end of logical block
	// replace source
	blkContents = strings.Replace(blkContents, "equinix/metal", "equinix/equinix", 1)
	blkContents = strings.Replace(blkContents, "packethost/packet", "equinix/equinix", 1)

	// comment version
	blkContents = strings.Replace(blkContents, "version", "#version", 1)

	return content[:idx] + blkContents + subStr[blockEnd:], nil
}

// find all required_providers definitions and make required transforms
func scanAndUpdateRequiredProvider(content string) (string, error) {
	for start, i := 0, -1; ; {
		i, _ = findTokenAfter("required_providers", content, start)

		// "required_providers" block not present in file
		if i == -1 {
			return content, nil
		}

		start += i

		blockStart, blockEnd := indexOpenCloseTokens('{', '}', content[start:])

		if blockStart == -1 {
			return content, fmt.Errorf("required provider detected, block start not found")
		}

		if blockEnd == -1 {
			return content, fmt.Errorf("required provider detected, block end not found")
		}

		end := start + blockEnd + 1

		res, err := updateRequiredProvider(content[start:end])
		if err != nil {
			return content, fmt.Errorf("problem parsing terraform:required_providers block\n %s", err)
		}

		content = content[:start] + res + content[end:]

		start = end
	}
}

// rewrite matching provider block to have equinix provider
func updateProviderBlock(content string) (string, error) {
	idx, _ := findToken("metal", content)
	if idx == -1 {
		idx, _ = findToken("packet", content)
	}
	if idx == -1 {
		return content, nil
	}

	subStr := content[idx:] // ignore everything before metal/packet provider

	// replace provider name
	subStr = strings.Replace(subStr, "metal", "equinix", 1)
	subStr = strings.Replace(subStr, "packet", "equinix", 1)

	return content[:idx] + subStr, nil
}

// find all providers blocks and make required transforms
func scanAndUpdateProvider(content string) (string, error) {
	for start, i := 0, -1; ; {
		i, _ = findTokenAfter("provider", content, start)

		// "providers" block not present in file
		if i == -1 {
			return content, nil
		}

		start += i

		blockStart, blockEnd := indexOpenCloseTokens('{', '}', content[start:])

		if blockStart == -1 {
			return content, fmt.Errorf("provider detected, block start not found")
		}

		if blockEnd == -1 {
			return content, fmt.Errorf("provider detected, block end not found")
		}

		end := start + blockEnd + 1

		res, err := updateProviderBlock(content[start:end])
		if err != nil {
			return content, fmt.Errorf("problem parsing provider block\n %s", err)
		}

		content = content[:start] + res + content[end:]

		start = end
	}
}

// return the text extent of a token match in a string
func findToken(token string, content string) (start int, end int) {
	idx := strings.Index(content, token)
	return idx, idx + len(token)
}

// return the text extent of a token match in a string after a specified index
func findTokenAfter(token string, content string, begin int) (start int, end int) {
	newStr := content[begin:]
	idx := strings.Index(newStr, token)

	if idx == -1 {
		return -1, -1
	}

	return idx, idx + len(token)
}

// parse logical terraform blocks to find open and closing braces
func indexOpenCloseTokens(open rune, close rune, content string) (start int, end int) {
	ct := 0
	start = -1

	for idx := 0; idx < len(content); {
		rn, rnWidth := utf8.DecodeRuneInString(content[idx:])

		// keep track of opening brackets to account for nesting
		if rn == open {
			ct++
			if start < 0 { // start index still -1, record the first opening bracket
				start = idx
			}
		}

		// closing brackets decrement nest level
		if rn == close {
			ct--
			if ct == 0 { // bracket count back to 0, record the final closing bracket
				return start, idx
			}
		}

		idx += rnWidth
		nextRn, nextRnWidth := utf8.DecodeRuneInString(content[idx:])

		// match " and advance idx to closing "
		if rn == '"' {
			for idx < len(content)-1 {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				rn2, w2 := utf8.DecodeRuneInString(content[idx+w1:])

				if rn1 == '\\' && rn2 == '"' {
					idx += w1 + w2
					continue
				}

				idx += w1
				if rn1 == '"' {
					break
				}
			}
			continue
		}

		// match '#' and advance idx to line end
		if rn == '#' {
			for idx < len(content) {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				idx += w1

				if rn1 == '\n' {
					break
				}
			}
			continue
		}

		// match '//' and advance idx to line end
		if rn == '/' && nextRn == '/' {
			idx += nextRnWidth
			for idx < len(content) {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				if rn1 == '\n' {
					break
				}
				idx += w1
			}
			continue
		}

		// match '/*' and advance idx to closing '*/'
		if rn == '/' && nextRn == '*' {
			idx += nextRnWidth
			for idx < len(content)-1 {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				rn2, w2 := utf8.DecodeRuneInString(content[idx+w1:])
				idx += w1
				if rn1 == '*' && rn2 == '/' {
					idx += w2
					break
				}
			}
			continue
		}

		// match '${' and advance idx to closing '}'
		if rn == '$' && nextRn == '{' {
			idx += rnWidth + nextRnWidth
			for idx < len(content)-1 {
				rn1, w1 := utf8.DecodeRuneInString(content[idx:])
				idx += w1
				if rn1 == '}' {
					break
				}
			}
			continue
		}
	}

	return start, -1
}
