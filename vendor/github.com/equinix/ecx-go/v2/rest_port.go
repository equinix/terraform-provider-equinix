package ecx

import (
	"fmt"
	"net/http"

	"github.com/equinix/ecx-go/v2/internal/api"
)

//GetUserPorts operation retrieves Equinix Fabric user ports
func (c RestClient) GetUserPorts() ([]Port, error) {
	path := "/ecx/v3/port/userport"
	respBody := []api.Port{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	mapped := make([]Port, len(respBody))
	for i := range respBody {
		mapped[i] = mapPortAPIToDomain(respBody[i])
	}
	return mapped, nil
}

func mapPortAPIToDomain(apiPort api.Port) Port {
	return Port{
		UUID:          apiPort.UUID,
		Name:          apiPort.Name,
		Region:        apiPort.Region,
		IBX:           apiPort.IBX,
		MetroCode:     apiPort.MetroCode,
		Priority:      apiPort.DevicePriority,
		Encapsulation: apiPort.Encapsulation,
		Buyout:        apiPort.Buyout,
		Bandwidth:     String(fmt.Sprintf("%d", Int64Value(apiPort.TotalBandwidth))),
		Status:        apiPort.ProvisionStatus,
	}
}
