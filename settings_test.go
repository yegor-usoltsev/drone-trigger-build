package main

import (
	"reflect"
	"testing"
)

func TestNewSettingsFromEnv(t *testing.T) {
	expected := Settings{
		Server:       "https://drone.example.com",
		Token:        "DroneExampleAccessToken",
		Repositories: []Repository{{Owner: "example", Name: "backend"}, {Owner: "example", Name: "frontend"}},
		Params:       map[string]string{"key1": "value1", "key2": "value2,value3"},
	}

	t.Setenv("PLUGIN_SERVER", " https://drone.example.com ")
	t.Setenv("PLUGIN_TOKEN", " DroneExampleAccessToken ")
	t.Setenv("PLUGIN_REPOSITORIES", " example / backend , example / frontend , ignore / , / , ignore , , ")
	t.Setenv("PLUGIN_PARAMS", " key1 = value1 , key2 = value2 , ignore = , , = ignore , , key2 = , value3 ")
	actual := NewSettingsFromEnv()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v; Actual: %+v", expected, actual)
	}
}

func TestMustGetenv(t *testing.T) {
	defer func() { _ = recover() }()

	mustGetenv("DOES_NOT_EXIST") // should panic

	t.Errorf("did not panic")
}

func TestParseRepositories(t *testing.T) {
	defer func() { _ = recover() }()

	parseRepositories("not a valid repo slug") // should panic

	t.Errorf("did not panic")
}
