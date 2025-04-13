package main

import (
	"os"
	"strings"
)

const EnvPrefix = "PLUGIN_"

type Settings struct {
	Server       string            `json:"server"`       // Drone server url
	Token        string            `json:"-"`            // Drone access token
	Repositories []Repository      `json:"repositories"` // List of repositories (owner/name) to trigger
	Params       map[string]string `json:"params"`       // List of parameters (key=value) to pass to triggered builds
}

type Repository struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func NewSettingsFromEnv() Settings {
	return Settings{
		Server:       mustGetenv(EnvPrefix + "SERVER"),
		Token:        mustGetenv(EnvPrefix + "TOKEN"),
		Repositories: parseRepositories(mustGetenv(EnvPrefix + "REPOSITORIES")),
		Params:       parseParams(mustGetenv(EnvPrefix + "PARAMS")),
	}
}

func mustGetenv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		panic("missing required environment variable: " + key)
	}
	return value
}

func parseRepositories(str string) []Repository {
	repositories := make([]Repository, 0)
	for _, slug := range strings.Split(str, ",") {
		ownerAndName := strings.Split(slug, "/")
		if len(ownerAndName) != 2 {
			continue
		}
		owner := strings.TrimSpace(ownerAndName[0])
		if owner == "" {
			continue
		}
		name := strings.TrimSpace(ownerAndName[1])
		if name == "" {
			continue
		}
		repositories = append(repositories, Repository{Owner: owner, Name: name})
	}
	if len(repositories) == 0 {
		panic("no valid repositories specified in PLUGIN_REPOSITORIES")
	}
	return repositories
}

func parseParams(str string) map[string]string {
	lastKey := ""
	keyToValues := make(map[string][]string)
	for _, part := range strings.Split(str, ",") {
		keyAndValue := strings.Split(part, "=")
		if len(keyAndValue) == 2 {
			lastKey = strings.TrimSpace(keyAndValue[0])
			if lastKey == "" {
				continue
			}
			value := strings.TrimSpace(keyAndValue[1])
			if value != "" {
				keyToValues[lastKey] = append(keyToValues[lastKey], value)
			}
		} else if len(keyAndValue) == 1 && lastKey != "" {
			value := strings.TrimSpace(keyAndValue[0])
			if value != "" {
				keyToValues[lastKey] = append(keyToValues[lastKey], value)
			}
		}
	}
	params := make(map[string]string)
	for key, values := range keyToValues {
		if len(values) > 0 {
			params[key] = strings.Join(values, ",")
		}
	}
	return params
}
