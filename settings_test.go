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

func TestMustGetenv(t *testing.T) { //nolint:tparallel
	t.Run("existing variable", func(t *testing.T) {
		t.Setenv("TEST_VAR", "test_value")
		if got := mustGetenv("TEST_VAR"); got != "test_value" {
			t.Errorf("mustGetenv() = %v, want %v", got, "test_value")
		}
	})

	t.Run("empty variable", func(t *testing.T) {
		t.Setenv("TEST_VAR", "")
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty variable")
			}
		}()
		mustGetenv("TEST_VAR")
	})

	t.Run("non-existent variable", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for non-existent variable")
			}
		}()
		mustGetenv("DOES_NOT_EXIST")
	})
}

func TestParseRepositories(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    []Repository
		wantErr bool
	}{
		{
			name:  "valid repositories",
			input: "owner1/repo1,owner2/repo2",
			want: []Repository{
				{Owner: "owner1", Name: "repo1"},
				{Owner: "owner2", Name: "repo2"},
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "skip invalid entries",
			input: "owner1/repo1,invalid,owner2/repo2",
			want: []Repository{
				{Owner: "owner1", Name: "repo1"},
				{Owner: "owner2", Name: "repo2"},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		name := test.name
		input := test.input
		want := test.want
		wantErr := test.wantErr
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if wantErr {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic")
					}
				}()
			}
			got := parseRepositories(input)
			if !wantErr && !reflect.DeepEqual(got, want) {
				t.Errorf("parseRepositories() = %v, want %v", got, want)
			}
		})
	}
}

func TestParseParams(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "empty input",
			input: "",
			want:  map[string]string{},
		},
		{
			name:  "single key-value",
			input: "key=value",
			want:  map[string]string{"key": "value"},
		},
		{
			name:  "multiple key-values",
			input: "key1=value1,key2=value2",
			want:  map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:  "multiple values for same key",
			input: "key=value1,key=value2",
			want:  map[string]string{"key": "value1,value2"},
		},
		{
			name:  "ignore empty keys and values",
			input: "=value,,key=,=,key=value",
			want:  map[string]string{"key": "value"},
		},
		{
			name:  "trim spaces",
			input: " key1 = value1 , key2 = value2 ",
			want:  map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, test := range tests {
		name := test.name
		input := test.input
		want := test.want
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := parseParams(input)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("parseParams() = %v, want %v", got, want)
			}
		})
	}
}
