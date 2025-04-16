package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCreate(t *testing.T) {
	t.Parallel()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	repoOwner := "acme"
	repoName := "payment-service"
	paramKey := "tags"
	paramValue := "v1,v1.0,v1.0.0"
	expected := &drone.Build{ID: 42}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, fmt.Sprintf("/api/repos/%s/%s/builds", repoOwner, repoName), r.URL.Path)
		assert.Equal(t, "Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal(t, paramValue, r.URL.Query().Get(paramKey))
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	droneClient := NewDroneClient(server.URL, token)
	actual, err := droneClient.BuildCreate(repoOwner, repoName, "", "", map[string]string{paramKey: paramValue})
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
