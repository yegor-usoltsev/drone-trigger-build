package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, expected, actual)
}

func TestMustGetenv(t *testing.T) {
	t.Run("existing variable", func(t *testing.T) {
		t.Setenv("TEST_VAR", "test_value")
		assert.Equal(t, "test_value", mustGetenv("TEST_VAR"))
	})

	t.Run("empty variable", func(t *testing.T) {
		t.Setenv("TEST_VAR", "")
		assert.Panics(t, func() { mustGetenv("TEST_VAR") })
	})

	t.Run("non-existent variable", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() { mustGetenv("DOES_NOT_EXIST") })
	})
}

func TestParseRepositories(t *testing.T) {
	t.Run("valid repositories", func(t *testing.T) {
		t.Parallel()
		given := "owner1/repo1,owner2/repo2"
		expected := []Repository{{Owner: "owner1", Name: "repo1"}, {Owner: "owner2", Name: "repo2"}}
		actual := parseRepositories(given)
		assert.Equal(t, expected, actual)
	})

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() { parseRepositories("") })
	})

	t.Run("invalid format", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() { parseRepositories("invalid") })
	})

	t.Run("skip invalid entries", func(t *testing.T) {
		t.Parallel()
		given := "owner1/repo1,invalid,owner2/repo2"
		expected := []Repository{{Owner: "owner1", Name: "repo1"}, {Owner: "owner2", Name: "repo2"}}
		actual := parseRepositories(given)
		assert.Equal(t, expected, actual)
	})
}

func TestParseParams(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		actual := parseParams("")
		assert.Equal(t, map[string]string{}, actual)
	})

	t.Run("single key-value", func(t *testing.T) {
		t.Parallel()
		given := "key=value"
		expected := map[string]string{"key": "value"}
		actual := parseParams(given)
		assert.Equal(t, expected, actual)
	})

	t.Run("multiple key-values", func(t *testing.T) {
		t.Parallel()
		given := "key1=value1,key2=value2"
		expected := map[string]string{"key1": "value1", "key2": "value2"}
		actual := parseParams(given)
		assert.Equal(t, expected, actual)
	})

	t.Run("multiple values for same key", func(t *testing.T) {
		t.Parallel()
		given := "key=value1,key=value2"
		expected := map[string]string{"key": "value1,value2"}
		actual := parseParams(given)
		assert.Equal(t, expected, actual)
	})

	t.Run("ignore empty keys and values", func(t *testing.T) {
		t.Parallel()
		given := "=value,,key=,=,key=value"
		expected := map[string]string{"key": "value"}
		actual := parseParams(given)
		assert.Equal(t, expected, actual)
	})

	t.Run("trim spaces", func(t *testing.T) {
		t.Parallel()
		given := " key1 = value1 , key2 = value2 "
		expected := map[string]string{"key1": "value1", "key2": "value2"}
		actual := parseParams(given)
		assert.Equal(t, expected, actual)
	})
}
