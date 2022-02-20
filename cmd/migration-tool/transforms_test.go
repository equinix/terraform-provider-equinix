package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrationReplaceTemplateTokens_basic(t *testing.T) {
	const original = `
resource "metal_port" "bond0" {
	port_id = local.bond0_id
	layer2 = false
	bonded = true
	vlan_ids = [metal_vlan.test.id]
}
	
resource "metal_vlan" "test" {
	description = "test"
	metro = "sv"
	project = metal_project.test.id
}
`

	const expected = `
resource "equinix_metal_port" "bond0" {
	port_id = local.bond0_id
	layer2 = false
	bonded = true
	vlan_ids = [equinix_metal_vlan.test.id]
}
	
resource "equinix_metal_vlan" "test" {
	description = "test"
	metro = "sv"
	project = equinix_metal_project.test.id
}
`

	var actual strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(original))
	for scanner.Scan() {
		str := scanner.Text()
		str = replaceTemplateTokens(str)
		actual.WriteString(fmt.Sprintf("%s\n", str))
	}

	assert.Equal(t, expected, actual.String(), "Result matches expected result")
}

func TestMigrationReplaceTemplateTokens_multiMatchPerLine(t *testing.T) {
	const original = `
terraform {
	required_providers {
		packet = {
			source = "packethost/packet"
			version = "1.0.0"
		}
	}
}

# Configure Packet Provider.
provider "packet" {
	auth_token = var.auth_token
}

data "packet_project" "test" {
	name = var.packet_project_name
}

resource "packet_connection" "test" {
	name            = var.packet_connection_name
	organization_id = data.packet_project.test.organization_id
	project_id      = data.packet_project.test.project_id
	metro           = var.packet_connection_metro
	redundancy      = var.packet_connection_redundancy
	type            = "shared"
	description     = var.packet_connection_description
	tags            = var.packet_connection_tags
}

resource "packet_device" "test" {
	count = 3

	hostname         = "tf.coreos2"
	plan             = "c3.small.x86"
	metro            = "sv"
	operating_system = "ubuntu_20_04"
	billing_cycle    = "hourly"
	project_id       = local.project_id
}

data "packet_device_bgp_neighbors" "test" {
	device_id = packet_device.test[1].id
}

locals {
	address_family = "${lookup(data.packet_device_bgp_neighbors.test.bgp_neighbors[0], "address_family")}"
	ips = tolist([packet_device.test[*].access_public_ipv4])
	ip_address = var.packet_network_type == "public" ? metal_device.foodevice.access_public_ipv4 : metal_device.foodevice.access_private_ipv4
}
`

	const expected = `
terraform {
	required_providers {
		packet = {
			source = "packethost/packet"
			version = "1.0.0"
		}
	}
}

# Configure Packet Provider.
provider "packet" {
	auth_token = var.auth_token
}

data "equinix_metal_project" "test" {
	name = var.packet_project_name
}

resource "equinix_metal_connection" "test" {
	name            = var.packet_connection_name
	organization_id = data.equinix_metal_project.test.organization_id
	project_id      = data.equinix_metal_project.test.project_id
	metro           = var.packet_connection_metro
	redundancy      = var.packet_connection_redundancy
	type            = "shared"
	description     = var.packet_connection_description
	tags            = var.packet_connection_tags
}

resource "equinix_metal_device" "test" {
	count = 3

	hostname         = "tf.coreos2"
	plan             = "c3.small.x86"
	metro            = "sv"
	operating_system = "ubuntu_20_04"
	billing_cycle    = "hourly"
	project_id       = local.project_id
}

data "equinix_metal_device_bgp_neighbors" "test" {
	device_id = equinix_metal_device.test[1].id
}

locals {
	address_family = "${lookup(data.equinix_metal_device_bgp_neighbors.test.bgp_neighbors[0], "address_family")}"
	ips = tolist([equinix_metal_device.test[*].access_public_ipv4])
	ip_address = var.packet_network_type == "public" ? equinix_metal_device.foodevice.access_public_ipv4 : equinix_metal_device.foodevice.access_private_ipv4
}
`

	var actual strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(original))
	for scanner.Scan() {
		str := scanner.Text()
		str = replaceTemplateTokens(str)
		actual.WriteString(fmt.Sprintf("%s\n", str))
	}

	assert.Equal(t, expected, actual.String(), "Result matches expected result")
}

func TestMigrationBraceIndexing_basic(t *testing.T) {
	const str = `{{}}`
	expectStart := 0
	expectEnd := 3
	start, end := indexOpenCloseTokens('{', '}', str)

	if start != expectStart {
		t.Errorf("expected %d, got %d\n", expectStart, start)
	}

	if end != expectEnd {
		t.Errorf("expected %d, got %d\n", expectEnd, end)
	}
}

func TestMigrationBraceIndexingWith_strings(t *testing.T) {
	const str = ` " " { { } } `
	expectStart := 5
	expectEnd := 11
	start, end := indexOpenCloseTokens('{', '}', str)

	if start != expectStart {
		t.Errorf("expected %d, got %d\n", expectStart, start)
	}

	if end != expectEnd {
		t.Errorf("expected %d, got %d\n", expectEnd, end)
	}
}

func TestMigrationBraceIndexing_Provider(t *testing.T) {
	const str = `provider "metal" {
  auth_token = "${var.auth_token}"
}`

	expectStart := 17
	expectEnd := len(str) - 1
	start, end := indexOpenCloseTokens('{', '}', str)

	if start != expectStart {
		t.Errorf("expected %d, got %d\n", expectStart, start)
	}

	if end != expectEnd {
		t.Errorf("expected %d, got %d\n", expectEnd, end)
	}
}

func TestMigrationFindOpeningBrace(t *testing.T) {
	const str = `}{}`
	expect := 1
	start, _ := indexOpenCloseTokens('{', '}', str)

	if start != expect {
		t.Errorf("expected %d, got %d\n", expect, start)
	}
}

func TestMigrationFindClosingBrace(t *testing.T) {
	const str = `{}}`
	expect := 1
	_, end := indexOpenCloseTokens('{', '}', str)
	if end != expect {
		t.Errorf("expected %d, got %d\n", expect, end)
	}
}

func TestMigrationMissingOpeningBrace(t *testing.T) {
	const str = `}}`
	expect := -1
	start, _ := indexOpenCloseTokens('{', '}', str)

	if start != expect {
		t.Errorf("expected %d, got %d\n", expect, start)
	}
}

func TestMigrationMissingClosingBrace(t *testing.T) {
	const str = `{{}`
	expect := -1
	_, end := indexOpenCloseTokens('{', '}', str)

	if end != expect {
		t.Errorf("expected %d, got %d\n", expect, end)
	}
}

func TestMigrationReplaceProvider(t *testing.T) {
	// given
	context := `
provider "metal" {
	auth_token = var.auth_token
}`
	expected := `
provider "equinix" {
	auth_token = var.auth_token
}`
	// when
	result, _ := scanAndUpdateProvider(context)

	// then
	assert.Equal(t, expected, result, "Result matches expected result")
}

func TestMigrationReplaceRequiredProvider(t *testing.T) {
	// given
	context := `
terraform {
	required_providers {
		metal = {
			source = "equinix/metal"
			#commment
			version = "3.2.1"
		}
		foo = {
			source = "foo/fooprovider"
			version = "1.0.0"
		}
	}
}`
	expected := `
terraform {
	required_providers {
		equinix = {
			source = "equinix/equinix"
			#commment
			#version = "3.2.1"
		}
		foo = {
			source = "foo/fooprovider"
			version = "1.0.0"
		}
	}
}`
	// when
	result, _ := scanAndUpdateRequiredProvider(context)

	// then
	assert.Equal(t, expected, result, "Result matches expected result")
}
