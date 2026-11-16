package schema

import (
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

func TestChangeLogGoToTerraformMapsAllSchemaAttributes(t *testing.T) {
	changeLog := fabricv4.Changelog{}
	changeLog.SetCreatedByEmail("creator@example.com")
	changeLog.SetUpdatedByEmail("updater@example.com")
	changeLog.SetDeletedByEmail("deleter@example.com")

	set := ChangeLogGoToTerraform(&changeLog)
	if set.Len() != 1 {
		t.Fatalf("expected one change_log element, got %d", set.Len())
	}
	mapped := set.List()[0].(map[string]any)

	for key := range ChangeLogSch() {
		if _, ok := mapped[key]; !ok {
			t.Errorf("change_log attribute %q is not mapped", key)
		}
	}
	for key, want := range map[string]string{
		"created_by_email": "creator@example.com",
		"updated_by_email": "updater@example.com",
		"deleted_by_email": "deleter@example.com",
	} {
		if got := mapped[key]; got != want {
			t.Errorf("%s = %v, want %q", key, got, want)
		}
	}
}
