package main

import (
	"context"

	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

func NewDroneClient(server, token string) drone.Client {
	config := new(oauth2.Config)
	httpClient := config.Client(context.Background(), &oauth2.Token{AccessToken: token})
	return drone.NewClient(server, httpClient)
}
