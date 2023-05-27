package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type DroneClient struct {
	server string
	token  string
}

func NewDroneClient(server string, token string) *DroneClient {
	return &DroneClient{server: server, token: token}
}

func (c *DroneClient) BuildCreate(owner string, name string, params map[string]string) (*Build, error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/repos/%s/%s/builds", c.server, owner, name), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("status: %s", response.Status)
	}

	var build Build
	err = json.NewDecoder(response.Body).Decode(&build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

type Build struct {
	ID int64 `json:"id"`
}
