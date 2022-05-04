package metal

import (
	"fmt"
	"testing"
	"time"

	"github.com/packethost/packngo"
)

type mockHWService struct {
	GetFn  func(string, *packngo.GetOptions) (*packngo.HardwareReservation, *packngo.Response, error)
	ListFn func(string, *packngo.ListOptions) ([]packngo.HardwareReservation, *packngo.Response, error)
	MoveFn func(string, string) (*packngo.HardwareReservation, *packngo.Response, error)
}

func (m *mockHWService) Get(id string, opt *packngo.GetOptions) (*packngo.HardwareReservation, *packngo.Response, error) {
	return m.GetFn(id, opt)
}
func (m *mockHWService) List(project string, opt *packngo.ListOptions) ([]packngo.HardwareReservation, *packngo.Response, error) {
	return m.ListFn(project, opt)
}
func (m *mockHWService) Move(id string, dest string) (*packngo.HardwareReservation, *packngo.Response, error) {
	return m.MoveFn(id, dest)
}

var _ packngo.HardwareReservationService = (*mockHWService)(nil)

func Test_waitUntilReservationProvisionable(t *testing.T) {
	type args struct {
		reservationId string
		instanceId    string
		meta          *packngo.Client
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				reservationId: "reservationId",
				instanceId:    "instanceId",
				meta: &packngo.Client{
					HardwareReservations: &mockHWService{
						GetFn: func(_ string, _ *packngo.GetOptions) (*packngo.HardwareReservation, *packngo.Response, error) {
							return nil, nil, fmt.Errorf("boom")
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "provisionable",
			args: args{
				reservationId: "reservationId",
				instanceId:    "instanceId",
				meta: &packngo.Client{
					HardwareReservations: (func() *mockHWService {
						invoked := new(int)

						responses := map[int]struct {
							id            string
							provisionable bool
						}{
							0: {"instanceId", false}, // should retry
							1: {"", true},            // should return success
						}

						return &mockHWService{
							GetFn: func(_ string, opts *packngo.GetOptions) (*packngo.HardwareReservation, *packngo.Response, error) {
								response := responses[*invoked]
								*invoked++

								var device *packngo.Device
								if opts != nil && contains(opts.Includes, "device") {
									device = &packngo.Device{ID: response.id}
								}
								return &packngo.HardwareReservation{
									Device: device, Provisionable: response.provisionable,
								}, nil, nil
							},
						}
					})(),
				},
			},
			wantErr: false,
		},
		{
			name: "reprovisioned",
			args: args{
				reservationId: "reservationId",
				instanceId:    "instanceId",
				meta: &packngo.Client{
					HardwareReservations: (func() *mockHWService {
						responses := map[int]struct {
							id            string
							provisionable bool
						}{
							0: {"instanceId", false},      // should retry
							1: {"new instance id", false}, // should return success
						}
						invoked := new(int)

						return &mockHWService{
							GetFn: func(_ string, opts *packngo.GetOptions) (*packngo.HardwareReservation, *packngo.Response, error) {
								response := responses[*invoked]
								*invoked++

								var device *packngo.Device
								if opts != nil && contains(opts.Includes, "device") {
									device = &packngo.Device{ID: response.id}
								}
								return &packngo.HardwareReservation{
									Device: device, Provisionable: response.provisionable,
								}, nil, nil
							},
						}
					})(),
				},
			},
			wantErr: false,
		},
		{
			name: "foreverDeprovisioning",
			args: args{
				reservationId: "reservationId",
				instanceId:    "instanceId",
				meta: &packngo.Client{
					HardwareReservations: &mockHWService{
						GetFn: func(_ string, _ *packngo.GetOptions) (*packngo.HardwareReservation, *packngo.Response, error) {
							return &packngo.HardwareReservation{
								Device: nil, Provisionable: false,
							}, nil, nil
						},
					},
				},
			},
			wantErr: true,
		},
	}

	// delay and minTimeout * 2 should be less than timeout for each test.
	// timeout * number of tests that reach timeout must be less than 30s (default go test timeout).
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := waitUntilReservationProvisionable(tt.args.meta, tt.args.reservationId, tt.args.instanceId, 50*time.Millisecond, 1*time.Second, 50*time.Millisecond); (err != nil) != tt.wantErr {
				t.Errorf("waitUntilReservationProvisionable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
