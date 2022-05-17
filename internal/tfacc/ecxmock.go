package tfacc

import "github.com/equinix/ecx-go/v2"

type mockECXClient struct {
	GetUserPortsFn func() ([]ecx.Port, error)

	GetL2OutgoingConnectionsFn     func(statuses []string) ([]ecx.L2Connection, error)
	GetL2ConnectionFn              func(uuid string) (*ecx.L2Connection, error)
	CreateL2ConnectionFn           func(conn ecx.L2Connection) (*string, error)
	CreateL2RedundantConnectionFn  func(priConn, secConn ecx.L2Connection) (*string, *string, error)
	NewL2ConnectionUpdateRequestFn func(uuid string) ecx.L2ConnectionUpdateRequest
	DeleteL2ConnectionFn           func(uuid string) error
	ConfirmL2ConnectionFn          func(uuid string, confirmConn ecx.L2ConnectionToConfirm) (*ecx.L2ConnectionConfirmation, error)

	GetL2SellerProfilesFn    func() ([]ecx.L2ServiceProfile, error)
	GetL2ServiceProfileFn    func(uuid string) (*ecx.L2ServiceProfile, error)
	CreateL2ServiceProfileFn func(sp ecx.L2ServiceProfile) (*string, error)
	UpdateL2ServiceProfileFn func(sp ecx.L2ServiceProfile) error
	DeleteL2ServiceProfileFn func(uuid string) error
}

func (m *mockECXClient) GetUserPorts() ([]ecx.Port, error) {
	return m.GetUserPortsFn()
}

func (m *mockECXClient) GetL2OutgoingConnections(statuses []string) ([]ecx.L2Connection, error) {
	return m.GetL2OutgoingConnectionsFn(statuses)
}

func (m *mockECXClient) GetL2Connection(uuid string) (*ecx.L2Connection, error) {
	return m.GetL2ConnectionFn(uuid)
}

func (m *mockECXClient) CreateL2Connection(conn ecx.L2Connection) (*string, error) {
	return m.CreateL2ConnectionFn(conn)
}

func (m *mockECXClient) CreateL2RedundantConnection(priConn, secConn ecx.L2Connection) (*string, *string, error) {
	return m.CreateL2RedundantConnectionFn(priConn, secConn)
}

func (m *mockECXClient) NewL2ConnectionUpdateRequest(uuid string) ecx.L2ConnectionUpdateRequest {
	return m.NewL2ConnectionUpdateRequestFn(uuid)
}

func (m *mockECXClient) DeleteL2Connection(uuid string) error {
	return m.DeleteL2ConnectionFn(uuid)
}

func (m *mockECXClient) ConfirmL2Connection(uuid string, confirmConn ecx.L2ConnectionToConfirm) (*ecx.L2ConnectionConfirmation, error) {
	return m.ConfirmL2ConnectionFn(uuid, confirmConn)
}

func (m *mockECXClient) GetL2SellerProfiles() ([]ecx.L2ServiceProfile, error) {
	return m.GetL2SellerProfilesFn()
}

func (m *mockECXClient) GetL2ServiceProfile(uuid string) (*ecx.L2ServiceProfile, error) {
	return m.GetL2ServiceProfileFn(uuid)
}

func (m *mockECXClient) CreateL2ServiceProfile(sp ecx.L2ServiceProfile) (*string, error) {
	return m.CreateL2ServiceProfileFn(sp)
}

func (m *mockECXClient) UpdateL2ServiceProfile(sp ecx.L2ServiceProfile) error {
	return m.UpdateL2ServiceProfileFn(sp)
}

func (m *mockECXClient) DeleteL2ServiceProfile(uuid string) error {
	return m.DeleteL2ServiceProfileFn(uuid)
}

var _ ecx.Client = (*mockECXClient)(nil)
