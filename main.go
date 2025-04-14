package main

import (
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	settings := NewSettingsFromEnv()
	slog.Info("parsed settings", "settings", settings)

	droneClient := NewDroneClient(settings.Server, settings.Token)

	for _, repo := range settings.Repositories {
		build, err := droneClient.BuildCreate(repo.Owner, repo.Name, "", "", settings.Params)
		if err != nil {
			slog.Error("failed to create build", "owner", repo.Owner, "repo", repo.Name, "error", err)
			continue
		}
		slog.Info("created build", "owner", repo.Owner, "repo", repo.Name, "build_id", build.ID)
	}
}
