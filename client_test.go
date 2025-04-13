package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/drone/drone-go/drone"
)

func TestBuildCreate(t *testing.T) {
	t.Parallel()
	token := "DroneExampleAccessToken"
	repoOwner := "example"
	repoName := "backend"
	paramKey := "key1"
	paramValue := "value1,value2"
	expected := &drone.Build{ID: 42}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost &&
			r.URL.Path == fmt.Sprintf("/api/repos/%s/%s/builds", repoOwner, repoName) &&
			r.Header.Get("Authorization") == "Bearer "+token &&
			r.URL.Query().Get(paramKey) == paramValue {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(expected)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	droneClient := NewDroneClient(server.URL, token)
	actual, err := droneClient.BuildCreate(repoOwner, repoName, "", "", map[string]string{paramKey: paramValue})

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v; Actual: %+v", expected, actual)
	}
}
