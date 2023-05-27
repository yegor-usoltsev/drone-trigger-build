package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestBuildCreate(t *testing.T) {
	token := "DroneExampleAccessToken"
	repoOwner := "example"
	repoName := "backend"
	paramKey := "key1"
	paramValue := "value1,value2"
	expected := &Build{ID: 42}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost &&
			r.URL.Path == fmt.Sprintf("/api/repos/%s/%s/builds", repoOwner, repoName) &&
			r.Header.Get("Authorization") == fmt.Sprintf("Bearer %s", token) &&
			r.URL.Query().Get(paramKey) == paramValue {
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(expected)
			return
		}
		w.WriteHeader(500)
	}))
	defer server.Close()
	droneClient := NewDroneClient(server.URL, token)
	actual, err := droneClient.BuildCreate(repoOwner, repoName, map[string]string{paramKey: paramValue})

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v; Actual: %+v", expected, actual)
	}
}
