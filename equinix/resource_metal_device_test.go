package equinix

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccMetalDevice_readErrorHandling(t *testing.T) {
	type args struct {
		newResource bool
		handler     func(w http.ResponseWriter, r *http.Request)
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "forbiddenAfterProvision",
			args: args{
				newResource: false,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusForbidden)
				},
			},
			wantErr: false,
		},
		{
			name: "notFoundAfterProvision",
			args: args{
				newResource: false,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusNotFound)
				},
			},
			wantErr: false,
		},
		{
			name: "forbiddenWaitForActiveDeviceProvision",
			args: args{
				newResource: true,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusForbidden)
				},
			},
			wantErr: true,
		},
		{
			name: "notFoundProvision",
			args: args{
				newResource: true,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusNotFound)
				},
			},
			wantErr: true,
		},
		{
			name: "errorProvision",
			args: args{
				newResource: true,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusBadRequest)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			d := new(schema.ResourceData)
			if tt.args.newResource {
				d.MarkNewResource()
			} else {
				d.SetId(uuid.New().String())
			}

			mockAPI := httptest.NewServer(http.HandlerFunc(tt.args.handler))
			meta := &config.Config{
				BaseURL: mockAPI.URL,
				Token:   "fakeTokenForMock",
			}
			meta.Load(ctx)

			if err := resourceMetalDeviceRead(ctx, d, meta); (err != nil) != tt.wantErr {
				t.Errorf(" ResourceMetalDeviceRead() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockAPI.Close()
		})
	}
}
