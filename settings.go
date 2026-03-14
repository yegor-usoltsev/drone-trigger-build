package main

import (
	"fmt"
	"os"
	"strings"
)

const envPrefix = "PLUGIN_"

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

func NewSettingsFromEnv() (Settings, error) {
	server, err := getenv(envPrefix + "SERVER")
	if err != nil {
		return Settings{}, err
	}
	token, err := getenv(envPrefix + "TOKEN")
	if err != nil {
		return Settings{}, err
	}
	reposStr, err := getenv(envPrefix + "REPOSITORIES")
	if err != nil {
		return Settings{}, err
	}
	repos, err := parseRepositories(reposStr)
	if err != nil {
		return Settings{}, err
	}
	paramsStr, err := getenv(envPrefix + "PARAMS")
	if err != nil {
		return Settings{}, err
	}
	return Settings{
		Server:       server,
		Token:        token,
		Repositories: repos,
		Params:       parseParams(paramsStr),
	}, nil
}

func getenv(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return value, nil
}

func parseRepositories(str string) ([]Repository, error) {
	repositories := make([]Repository, 0)
	for slug := range strings.SplitSeq(str, ",") {
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
		return nil, fmt.Errorf("no valid repositories specified")
	}
	return repositories, nil
}

func parseParams(str string) map[string]string {
	lastKey := ""
	keyToValues := make(map[string][]string)
	for part := range strings.SplitSeq(str, ",") {
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
