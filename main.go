package main

import (
	"encoding/json"
	"log"
)

func main() {
	settings := NewSettingsFromEnv()
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		log.Fatalf("failed to marshal settings: %v", err)
	}
	log.Printf("parsed settings: %s", settingsJSON)

	droneClient := NewDroneClient(settings.Server, settings.Token)

	for _, repo := range settings.Repositories {
		build, err := droneClient.BuildCreate(repo.Owner, repo.Name, "", "", settings.Params)
		if err != nil {
			log.Fatalf("failed to create build for repository %s/%s: %v", repo.Owner, repo.Name, err)
		}
		log.Printf("created build %d for repository %s/%s", build.ID, repo.Owner, repo.Name)
	}
}
