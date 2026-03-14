package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	actual, err := NewSettingsFromEnv()

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetenv(t *testing.T) {
	t.Run("existing variable", func(t *testing.T) {
		t.Setenv("TEST_VAR", "test_value")
		val, err := getenv("TEST_VAR")
		require.NoError(t, err)
		assert.Equal(t, "test_value", val)
	})

	t.Run("empty variable", func(t *testing.T) {
		t.Setenv("TEST_VAR", "")
		_, err := getenv("TEST_VAR")
		assert.Error(t, err)
	})

	t.Run("non-existent variable", func(t *testing.T) {
		t.Parallel()
		_, err := getenv("DOES_NOT_EXIST")
		assert.Error(t, err)
	})
}

func TestParseRepositories(t *testing.T) {
	t.Run("valid repositories", func(t *testing.T) {
		t.Parallel()
		given := "owner1/repo1,owner2/repo2"
		expected := []Repository{{Owner: "owner1", Name: "repo1"}, {Owner: "owner2", Name: "repo2"}}
		actual, err := parseRepositories(given)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		_, err := parseRepositories("")
		assert.Error(t, err)
	})

	t.Run("invalid format", func(t *testing.T) {
		t.Parallel()
		_, err := parseRepositories("invalid")
		assert.Error(t, err)
	})

	t.Run("skip invalid entries", func(t *testing.T) {
		t.Parallel()
		given := "owner1/repo1,invalid,owner2/repo2"
		expected := []Repository{{Owner: "owner1", Name: "repo1"}, {Owner: "owner2", Name: "repo2"}}
		actual, err := parseRepositories(given)
		require.NoError(t, err)
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
