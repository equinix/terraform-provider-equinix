package config

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAddModuleToNEUserAgent_doesNotMutateSharedClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"uuid":"test"}`))
	}))
	defer srv.Close()

	c := &Config{BaseURL: srv.URL, Token: "test-token"}
	if err := c.Load(context.Background()); err != nil {
		t.Fatalf("Load: %v", err)
	}
	shared := c.Ne.(*ne.RestClient)
	sharedUserAgent := shared.Header.Get("User-agent")

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}, map[string]interface{}{})

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := c.Ne
			c.AddModuleToNEUserAgent(&client, d)
			if client.(*ne.RestClient) == shared {
				t.Error("AddModuleToNEUserAgent returned the shared client")
			}
			if _, err := client.GetSSHPublicKey("test"); err != nil {
				t.Errorf("GetSSHPublicKey: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := shared.Header.Get("User-agent"); got != sharedUserAgent {
		t.Errorf("shared client User-agent changed from %q to %q", sharedUserAgent, got)
	}
}
