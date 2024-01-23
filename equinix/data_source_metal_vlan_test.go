package equinix

import (
	"reflect"
	"testing"

	"github.com/packethost/packngo"
)

func TestMetalVlan_matchingVlan(t *testing.T) {
	type args struct {
		vlans     []packngo.VirtualNetwork
		vxlan     int
		projectID string
		facility  string
		metro     string
	}
	tests := []struct {
		name    string
		args    args
		want    *packngo.VirtualNetwork
		wantErr bool
	}{
		{
			name: "MatchingVLAN",
			args: args{
				vlans:     []packngo.VirtualNetwork{{VXLAN: 123}},
				vxlan:     123,
				projectID: "",
				facility:  "",
				metro:     "",
			},
			want:    &packngo.VirtualNetwork{VXLAN: 123},
			wantErr: false,
		},
		{
			name: "MatchingFac",
			args: args{
				vlans:    []packngo.VirtualNetwork{{FacilityCode: "fac"}},
				facility: "fac",
			},
			want:    &packngo.VirtualNetwork{FacilityCode: "fac"},
			wantErr: false,
		},
		{
			name: "MatchingMet",
			args: args{
				vlans: []packngo.VirtualNetwork{{MetroCode: "met"}},
				metro: "met",
			},
			want:    &packngo.VirtualNetwork{MetroCode: "met"},
			wantErr: false,
		},
		{
			name: "SecondMatch",
			args: args{
				vlans: []packngo.VirtualNetwork{{FacilityCode: "fac"}, {MetroCode: "met"}},
				metro: "met",
			},
			want:    &packngo.VirtualNetwork{MetroCode: "met"},
			wantErr: false,
		},
		{
			name: "TwoMatches",
			args: args{
				vlans: []packngo.VirtualNetwork{{MetroCode: "met"}, {MetroCode: "met"}},
				metro: "met",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ComplexMatch",
			args: args{
				vlans: []packngo.VirtualNetwork{{VXLAN: 987, FacilityCode: "fac", MetroCode: "skip"}, {VXLAN: 123, FacilityCode: "fac", MetroCode: "met"}, {VXLAN: 456, FacilityCode: "fac", MetroCode: "nope"}},
				metro: "met",
			},
			want:    &packngo.VirtualNetwork{VXLAN: 123, FacilityCode: "fac", MetroCode: "met"},
			wantErr: false,
		},
		{
			name: "NoMatch",
			args: args{
				vlans:     nil,
				vxlan:     123,
				projectID: "pid",
				facility:  "fac",
				metro:     "met",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchingVlan(tt.args.vlans, tt.args.vxlan, tt.args.projectID, tt.args.facility, tt.args.metro)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchingVlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchingVlan() = %v, want %v", got, tt.want)
			}
		})
	}
}
