package connection_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"strconv"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFabricCreatePort2SPConnection_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	connectionsTestData := testinghelpers.GetFabricEnvConnectionTestData(t)
	var publicSPName, portUUID string
	if len(ports) > 0 && len(connectionsTestData) > 0 {
		publicSPName = connectionsTestData["pfcr"]["publicSPName"]
		portUUID = ports["pfcr"]["dot1q"][0].GetUuid()
	}

	targetVlan, err := testinghelpers.RandomVlan(portUUID)

	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2SPConnectionConfig(publicSPName, "port2sp_PFCR", portUUID, "DC", targetVlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_connection.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "port2sp_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "1-323292"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.vlan_tag", strconv.Itoa(targetVlan)),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "SP"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.profile.0.type", "L2_PROFILE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.profile.0.name", publicSPName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.location.0.metro_code", "DC"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

func testAccFabricCreatePort2SPConnectionConfig(spName, name, portUUID, zSideMetro string, targetVlan int) string {
	return fmt.Sprintf(`

	data "equinix_fabric_service_profiles" "this" {
	  filter {
		property = "/name"
		operator = "="
		values   = ["%s"]
	  }
	}


	resource "equinix_fabric_connection" "test" {
		name = "%s"
		type = "EVPL_VC"
		notifications{
			type="ALL" 
			emails=["example@equinix.com"]
		} 
		bandwidth = 50
		geo_scope = "CONUS"
		redundancy {priority= "PRIMARY"}
		order {
			purchase_order_number= "1-323292"
		}
		a_side {
			access_point {
				type= "COLO"
				port {
					uuid= "%s"
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= %d
				}
			}
		}
		z_side {
			access_point {
				type= "SP"
				profile {
					type= "L2_PROFILE"
					uuid= data.equinix_fabric_service_profiles.this.data.0.uuid
				}
				location {
					metro_code= "%s"
				}
			}
		}
	}`, spName, name, portUUID, targetVlan, zSideMetro)
}

func TestAccFabricCreatePort2NonGenericSPConnection_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	connectionsTestData := testinghelpers.GetFabricEnvConnectionTestData(t)
	var nonGenericSPName, nonGenericSPAuthKey, portUUID string
	if len(ports) > 0 && len(connectionsTestData) > 0 {
		nonGenericSPName = connectionsTestData["pfcr"]["nonGenericSPName"]
		nonGenericSPAuthKey = connectionsTestData["pfcr"]["nonGenericSPAuthKey"]
		portUUID = ports["pfcr"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: port2NonGenericSPConfig,
				ConfigVariables: config.Variables{
					"sp_name":            config.StringVariable(nonGenericSPName),
					"authentication_key": config.StringVariable(nonGenericSPAuthKey),
					"port_uuid":          config.StringVariable(portUUID),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("name"), knownvalue.StringExact("port2nonG_sp_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("bandwidth"), knownvalue.Int32Exact(50)),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("type"), knownvalue.StringExact("EVPL_VC")),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("redundancy"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"priority": knownvalue.StringExact("PRIMARY"),
								"group":    knownvalue.NotNull(),
							}),
						}),
					),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("order"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"purchase_order_number": knownvalue.StringExact("1-323292"),
							}),
						}),
					),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("a_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("COLO"),
										"link_protocol": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type":     knownvalue.StringExact("DOT1Q"),
												"vlan_tag": knownvalue.Int32Exact(3769),
											}),
										}),
									}),
								}),
							}),
						})),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("z_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("SP"),
										"profile": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type": knownvalue.StringExact("L2_PROFILE"),
												"name": knownvalue.StringExact(nonGenericSPName),
											}),
										}),
									}),
								}),
							}),
						})),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

var port2NonGenericSPConfig = `
variable "sp_name" {
  type = string
}

variable "port_uuid" {
  type = string
}

variable "authentication_key" {
  type = string
  sensitive = true
}


data "equinix_fabric_service_profiles" "this" {
  filter {
    property = "/name"
    operator = "="
    values   = [var.sp_name]
  }
}

resource "equinix_fabric_connection" "test" {
  name = "port2nonG_sp_PFCR"
  type = "EVPL_VC"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com"]
  }
  bandwidth = 50
  geo_scope = "CONUS"
  redundancy { priority = "PRIMARY" }
  order {
    purchase_order_number = "1-323292"
  }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = var.port_uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = 3769
      }
    }
  }
  z_side {
    access_point {
      type               = "SP"
      authentication_key = var.authentication_key
      seller_region      = "us-east-1"
      profile {
        type = "L2_PROFILE"
        uuid = data.equinix_fabric_service_profiles.this.data.0.uuid
      }
      location {
        metro_code = "DC"
      }
    }
  }
}
`

func TestAccFabricCreatePort2PortConnection_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var aSidePortUUID, zSidePortUUID string
	if len(ports) > 0 {
		aSidePortUUID = ports["pfcr"]["dot1q"][0].GetUuid()
		zSidePortUUID = ports["pfcr"]["dot1q"][1].GetUuid()
	}

	asideVlan, err := testinghelpers.RandomVlan(aSidePortUUID)
	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	zsideVlan, err := testinghelpers.RandomVlan(zSidePortUUID)
	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: port2PortConnectionConfig,
				ConfigVariables: config.Variables{
					"aside_vlan":      config.IntegerVariable(asideVlan),
					"aside_port_uuid": config.StringVariable(aSidePortUUID),
					"zside_vlan":      config.IntegerVariable(zsideVlan),
					"zside_port_uuid": config.StringVariable(zSidePortUUID),
					"bandwidth":       config.IntegerVariable(50),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("name"), knownvalue.StringExact("port_test_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("bandwidth"), knownvalue.Int32Exact(50)),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("type"), knownvalue.StringExact("EVPL_VC")),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("order"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"purchase_order_number": knownvalue.StringExact("1-129105284100"),
							}),
						}),
					),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("a_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("COLO"),
										"link_protocol": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type":     knownvalue.StringExact("DOT1Q"),
												"vlan_tag": knownvalue.Int32Exact(int32(asideVlan)),
											}),
										}),
									}),
								}),
							}),
						})),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("z_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("COLO"),
										"link_protocol": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type":     knownvalue.StringExact("DOT1Q"),
												"vlan_tag": knownvalue.Int32Exact(int32(zsideVlan)),
											}),
										}),
									}),
								}),
							}),
						})),
				},
				ExpectNonEmptyPlan: true,
			},
			{
				Config: port2PortConnectionConfig,
				ConfigVariables: config.Variables{
					"aside_vlan":      config.IntegerVariable(asideVlan),
					"aside_port_uuid": config.StringVariable(aSidePortUUID),
					"zside_vlan":      config.IntegerVariable(zsideVlan),
					"zside_port_uuid": config.StringVariable(zSidePortUUID),
					"bandwidth":       config.IntegerVariable(100),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("name"), knownvalue.StringExact("port_test_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("bandwidth"), knownvalue.Int32Exact(100)),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("type"), knownvalue.StringExact("EVPL_VC")),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("order"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"purchase_order_number": knownvalue.StringExact("1-129105284100"),
							}),
						}),
					),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("a_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("COLO"),
										"link_protocol": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type":     knownvalue.StringExact("DOT1Q"),
												"vlan_tag": knownvalue.Int32Exact(int32(asideVlan)),
											}),
										}),
									}),
								}),
							}),
						})),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("z_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("COLO"),
										"link_protocol": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type":     knownvalue.StringExact("DOT1Q"),
												"vlan_tag": knownvalue.Int32Exact(int32(zsideVlan)),
											}),
										}),
									}),
								}),
							}),
						})),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

var port2PortConnectionConfig = `
variable "aside_vlan" {
  type = number
}

variable "aside_port_uuid" {
  type = string
}

variable "zside_vlan" {
  type = number
}

variable "zside_port_uuid" {
  type = string
}

variable "bandwidth" {
  type = number
}

resource "equinix_fabric_connection" "test" {
		type = "EVPL_VC"
		name = "port_test_PFCR"
		notifications{
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		order {
			purchase_order_number = "1-129105284100"
		}
		bandwidth = var.bandwidth
		a_side {
			access_point {
				type = "COLO"
				port {
				 uuid = var.aside_port_uuid
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= var.aside_vlan
				}
				location {
					metro_code = "SV"
				}
			}
		}
		z_side {
			access_point {
				type = "COLO"
				port{
				 uuid = var.zside_port_uuid
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= var.zside_vlan
				}
				location {
					metro_code= "SV"
				}
			}
		}
	}`

func TestAccFabricCreateCloudRouter2PortConnection_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var portUUID string
	if len(ports) > 0 {
		portUUID = ports["pfcr"]["dot1q"][1].GetUuid()
	}

	targetVlan, err := testinghelpers.RandomVlan(portUUID)
	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: cloudRouter2PortConnectionConfig,
				ConfigVariables: config.Variables{
					"vlan_tag":  config.IntegerVariable(targetVlan),
					"port_uuid": config.StringVariable(portUUID),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("name"), knownvalue.StringExact("fcr_test_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("bandwidth"), knownvalue.Int32Exact(50)),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("type"), knownvalue.StringExact("IP_VC")),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("redundancy"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"priority": knownvalue.StringExact("PRIMARY"),
								"group":    knownvalue.NotNull(),
							}),
						}),
					),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("order"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"purchase_order_number": knownvalue.StringExact("123485"),
							}),
						}),
					),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("project"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"),
							}),
						}),
					),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("a_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("CLOUD_ROUTER"),
										"router": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"uuid": knownvalue.NotNull(),
											}),
										}),
									}),
								}),
							}),
						})),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test", tfjsonpath.New("z_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("COLO"),
										"link_protocol": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type":     knownvalue.StringExact("DOT1Q"),
												"vlan_tag": knownvalue.Int32Exact(int32(targetVlan)),
											}),
										}),
									}),
								}),
							}),
						})),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

var cloudRouter2PortConnectionConfig = `
variable "vlan_tag" {
  type = number
}

variable "port_uuid" {
  type = string
}

resource "equinix_fabric_cloud_router" "this" {
  type = "XF_ROUTER"
  name = "Conn_Test_PFCR"
  location {
    metro_code = "SV"
  }
  order {
    purchase_order_number = "1-234567"
  }
  notifications {
    type = "ALL"
    emails = [
      "test@equinix.com",
      "test1@equinix.com"
    ]
  }
  project {
    project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
  }
  account {
    account_number = 201257
  }
  package {
    code = "STANDARD"
  }
}

resource "equinix_fabric_connection" "test" {
  type = "IP_VC"
  name = "fcr_test_PFCR"
  notifications {
    type   = "ALL"
    emails = ["test@equinix.com", "test1@equinix.com"]
  }
  order {
    purchase_order_number = "123485"
  }
  bandwidth = 50
  redundancy {
    priority = "PRIMARY"
  }
  a_side {
    access_point {
      type = "CLOUD_ROUTER"
      router {
        uuid = equinix_fabric_cloud_router.this.id
      }
    }
  }
  project {
    project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
  }
  z_side {
    access_point {
      type = "COLO"
      port {
        uuid = var.port_uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = var.vlan_tag
      }
      location {
        metro_code = "SV"
      }
    }
  }
}
	`

func TestAccFabricCreateVirtualDevice2NetworkConnection_PNFV(t *testing.T) {
	connectionTestData := testinghelpers.GetFabricEnvConnectionTestData(t)
	var virtualDevice string
	if len(connectionTestData) > 0 {
		virtualDevice = connectionTestData["pnfv"]["virtualDevice"]
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateVirtualDevice2NetworkConnectionConfig("vd2network_PNFV", virtualDevice),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "vd2network_PNFV"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPLAN_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "123485"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "VD"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.virtual_device.0.uuid", virtualDevice),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "NETWORK"),
					resource.TestCheckResourceAttrSet(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.network.0.uuid"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

func testAccFabricCreateVirtualDevice2NetworkConnectionConfig(name, virtualDeviceUUID string) string {
	return fmt.Sprintf(`
resource "equinix_fabric_network" "this" {
  type  = "EVPLAN"
  name  = "Tf_Network_PNFV"
  scope = "REGIONAL"
  notifications {
    type   = "ALL"
    emails = ["test@equinix.com", "test1@equinix.com"]
  }
  location {
    region = "AMER"
  }
  project {
    project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
  }
}

resource "equinix_fabric_connection" "test" {
  type = "EVPLAN_VC"
  name = "%s"
  notifications {
    type   = "ALL"
    emails = ["test@equinix.com", "test1@equinix.com"]
  }
  order {
    purchase_order_number = "123485"
  }
  bandwidth = 50
  redundancy {
    priority = "PRIMARY"
  }
  a_side {
    access_point {
      type = "VD"
      virtual_device {
        type = "EDGE"
        uuid = "%s"
      }
    }
  }
  z_side {
    access_point {
      type = "NETWORK"
      network {
        uuid = equinix_fabric_network.this.id
      }
    }
  }
}
`, name, virtualDeviceUUID)
}

func TestAccFabricCreatePort2EtreeNetworkConnection_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var portUUID string
	if len(ports) > 0 {
		portUUID = ports["pfcr"]["dot1q"][1].GetUuid()
	}

	targetVlan, err := testinghelpers.RandomVlan(portUUID)
	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: port2EtreeNetworkConfig,
				ConfigVariables: config.Variables{
					"vlan_tag":  config.IntegerVariable(targetVlan),
					"port_uuid": config.StringVariable(portUUID),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("name"), knownvalue.StringExact("port2etreenetwork_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("bandwidth"), knownvalue.Int32Exact(50)),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("type"), knownvalue.StringExact("EVPTREE_VC")),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("redundancy"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"priority": knownvalue.StringExact("PRIMARY"),
								"group":    knownvalue.NotNull(),
							}),
						}),
					),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("order"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"purchase_order_number": knownvalue.StringExact("123485"),
							}),
						}),
					),

					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("project"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"),
							}),
						}),
					),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("a_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("COLO"),
										"port": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"uuid": knownvalue.StringExact(portUUID),
											}),
										}),
									}),
								}),
							}),
						})),
					statecheck.ExpectKnownValue("equinix_fabric_connection.test_etree", tfjsonpath.New("z_side"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_point": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"type": knownvalue.StringExact("NETWORK"),
										"role": knownvalue.StringExact("LEAF"),
										"network": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"uuid": knownvalue.NotNull(),
											}),
										}),
									}),
								}),
							}),
						})),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

var port2EtreeNetworkConfig = `
variable "vlan_tag" {
  type = number
}

variable "port_uuid" {
  type = string
}

	resource "equinix_fabric_network" "this" {
		type = "EVPTREE"
		name = "Tf_EtreeNetwork_PFCR"
		scope = "REGIONAL"
		notifications {
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		location {
			region = "AMER"
		}
		project{
			project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
		}
	}

	resource "equinix_fabric_connection" "test_etree" {
		type = "EVPTREE_VC"
		name = "port2etreenetwork_PFCR"
		notifications{
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		order {
			purchase_order_number = "123485"
		}
		bandwidth = 50
		redundancy {
			priority= "PRIMARY"
		}
		a_side {
			access_point {
				type= "COLO"
				port {
					uuid = var.port_uuid
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= var.vlan_tag
				}
			}
		}
		z_side {
			access_point {
				type = "NETWORK"
				network {
					uuid = equinix_fabric_network.this.id
				}
                role = "LEAF"
			}
		}
	}`

func CheckConnectionDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection" {
			continue
		}

		err := connection.WaitUntilConnectionDeprovisioned(ctx, rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for connection deletion. ID: %s, Err: %s", rs.Primary.ID, err)
		}
	}
	return nil
}
