package equinix_test

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/equinix/ecx-go/v2"
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Legacy test helper functions
//
// These are deprecated functions that should not be used in new tests
// and should be removed from existing tests when the opportunity arises
//_______________________________________________________________________

// Deprecated: the logic here should be taken care of in your templates instead
func nprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		var strVal string
		switch val.(type) {
		case []string:
			r := regexp.MustCompile(`" "`)
			strVal = r.ReplaceAllString(fmt.Sprintf("%q", val), `", "`)
		default:
			strVal = fmt.Sprintf("%v", val)
		}
		format = strings.Replace(format, "%{"+key+"}", strVal, -1)
	}
	return format
}

// Deprecated: use stdlib maps.Copy instead
func copyMap(source map[string]interface{}) map[string]interface{} {
	target := make(map[string]interface{})
	for k, v := range source {
		target[k] = v
	}
	return target
}

// Deprecated: use stdlib slices.Equal instead
func slicesMatch(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	visited := make([]bool, len(s1))
	for i := 0; i < len(s1); i++ {
		found := false
		for j := 0; j < len(s2); j++ {
			if visited[j] {
				continue
			}
			if s1[i] == s2[j] {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Deprecated: use stdlib slices.EqualFunc instead
func slicesMatchCaseInsensitive(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	visited := make([]bool, len(s1))
	for i := 0; i < len(s1); i++ {
		found := false
		for j := 0; j < len(s2); j++ {
			if visited[j] {
				continue
			}
			if strings.EqualFold(s1[i], s2[j]) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type testAccConfig struct {
	ctx    map[string]interface{}
	config string
}

func newTestAccConfig(ctx map[string]interface{}) *testAccConfig {
	return &testAccConfig{
		ctx:    ctx,
		config: "",
	}
}

func (t *testAccConfig) build() string {
	return t.config
}

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
