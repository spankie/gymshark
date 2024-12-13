package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/spankie/gymshark/config"
)

func TestHealthHandler(t *testing.T) {
	conf := &config.Configuration{
		Port:       "8080",
		DbHost:     "localhost",
		DbPort:     5432,
		DbUsername: "spankie",
		DbPassword: "spankie",
		DbName:     "gymshark",
	}

	createDBAndHTTPServer(t, conf)

	// Make a request to the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", conf.Port))
	if err != nil {
		t.Errorf("failed to make request to server: %v", err)
	}

	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	var respMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&respMap)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}
	expected := "all systems are healthy"
	if respMap["message"] != expected {
		t.Errorf("expected db_health response %s, got %s", expected, respMap["message"])
	}
}
