package main

import (
	"fmt"
	"net/http"

	"github.com/drone/drone-go/drone"
)

func NewDroneClient(server string, token string) drone.Client {
	return drone.NewClient(server, NewAuthorizedHttpClient(token))
}

func NewAuthorizedHttpClient(token string) *http.Client {
	return &http.Client{Transport: NewAuthorizedHttpTransport(token)}
}

type AuthorizedHttpTransport struct {
	token string
	base  http.RoundTripper
}

func NewAuthorizedHttpTransport(token string) *AuthorizedHttpTransport {
	return &AuthorizedHttpTransport{token: token, base: http.DefaultTransport}
}

func (t *AuthorizedHttpTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.base.RoundTrip(r)
}
