package sdkcompat

import "github.com/equinix/equinix-sdk-go/services/fabricv4"

const stateKey = "state"

type ConnectionState string

type RouteFilterState string

type RouteAggregationState string

const (
	CONNECTIONSTATE_ACTIVE         ConnectionState = "ACTIVE"
	CONNECTIONSTATE_CANCELLED      ConnectionState = "CANCELLED"
	CONNECTIONSTATE_DEPROVISIONED  ConnectionState = "DEPROVISIONED"
	CONNECTIONSTATE_DEPROVISIONING ConnectionState = "DEPROVISIONING"
	CONNECTIONSTATE_DRAFT          ConnectionState = "DRAFT"
	CONNECTIONSTATE_FAILED         ConnectionState = "FAILED"
	CONNECTIONSTATE_PENDING        ConnectionState = "PENDING"
	CONNECTIONSTATE_PROVISIONED    ConnectionState = "PROVISIONED"
	CONNECTIONSTATE_PROVISIONING   ConnectionState = "PROVISIONING"
	CONNECTIONSTATE_REPROVISIONING ConnectionState = "REPROVISIONING"

	ROUTEFILTERSTATE_PROVISIONING      RouteFilterState = "PROVISIONING"
	ROUTEFILTERSTATE_REPROVISIONING    RouteFilterState = "REPROVISIONING"
	ROUTEFILTERSTATE_DEPROVISIONING    RouteFilterState = "DEPROVISIONING"
	ROUTEFILTERSTATE_PROVISIONED       RouteFilterState = "PROVISIONED"
	ROUTEFILTERSTATE_DEPROVISIONED     RouteFilterState = "DEPROVISIONED"
	ROUTEFILTERSTATE_NOT_PROVISIONED   RouteFilterState = "NOT_PROVISIONED"
	ROUTEFILTERSTATE_NOT_DEPROVISIONED RouteFilterState = "NOT_DEPROVISIONED"

	ROUTEAGGREGATIONSTATE_PROVISIONING      RouteAggregationState = "PROVISIONING"
	ROUTEAGGREGATIONSTATE_REPROVISIONING    RouteAggregationState = "REPROVISIONING"
	ROUTEAGGREGATIONSTATE_DEPROVISIONING    RouteAggregationState = "DEPROVISIONING"
	ROUTEAGGREGATIONSTATE_PROVISIONED       RouteAggregationState = "PROVISIONED"
	ROUTEAGGREGATIONSTATE_DEPROVISIONED     RouteAggregationState = "DEPROVISIONED"
	ROUTEAGGREGATIONSTATE_NOT_PROVISIONED   RouteAggregationState = "NOT_PROVISIONED"
	ROUTEAGGREGATIONSTATE_NOT_DEPROVISIONED RouteAggregationState = "NOT_DEPROVISIONED"
)

func get(props map[string]interface{}) string {
	if s, ok := props[stateKey].(string); ok {
		return s
	}
	return ""
}

func GetConnectionState(c *fabricv4.Connection) ConnectionState {
	if c == nil {
		return ""
	}
	return ConnectionState(get(c.AdditionalProperties))
}

func SetConnectionState(c *fabricv4.Connection, s ConnectionState) {
	if c.AdditionalProperties == nil {
		c.AdditionalProperties = map[string]interface{}{}
	}
	c.AdditionalProperties[stateKey] = string(s)
}

func GetRouteFilterState(r *fabricv4.RouteFiltersData) RouteFilterState {
	if r == nil {
		return ""
	}
	return RouteFilterState(get(r.AdditionalProperties))
}

func GetRouteAggregationState(r *fabricv4.RouteAggregationsData) RouteAggregationState {
	if r == nil {
		return ""
	}
	return RouteAggregationState(get(r.AdditionalProperties))
}
