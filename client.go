package main

import (
	"net/http"

	"github.com/drone/drone-go/drone"
)

func NewDroneClient(server, token string) drone.Client { //nolint:ireturn
	return drone.NewClient(server, NewAuthorizedHTTPClient(token))
}

func NewAuthorizedHTTPClient(token string) *http.Client {
	return &http.Client{Transport: NewAuthorizedHTTPTransport(token)} //nolint:exhaustruct
}

type AuthorizedHTTPTransport struct {
	base  http.RoundTripper
	token string
}

func NewAuthorizedHTTPTransport(token string) *AuthorizedHTTPTransport {
	return &AuthorizedHTTPTransport{base: http.DefaultTransport, token: token}
}

func (t *AuthorizedHTTPTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(r) //nolint:wrapcheck
}
