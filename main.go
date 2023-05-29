package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	settings := NewSettingsFromEnv()
	settingsJson, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsed settings: %s\n", settingsJson)
	droneClient := NewDroneClient(settings.Server, settings.Token)

	for _, repo := range settings.Repositories {
		build, err := droneClient.BuildCreate(repo.Owner, repo.Name, "", "", settings.Params)
		if err != nil {
			panic(fmt.Sprintf("Cannot create a new build for repository %s/%s", repo.Owner, repo.Name))
		}
		fmt.Printf("Successfully created a new build %d for repository %s/%s\n", build.ID, repo.Owner, repo.Name)
	}
}
