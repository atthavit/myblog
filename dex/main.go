package main

import (
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func main() {
	provider, err := oidc.NewProvider(context.TODO(), "http://localhost:5556/dex")
	if err != nil {
		log.Fatal(err)
	}
	a := app{
		oauth2Config: &oauth2.Config{
			ClientID:     "example",
			ClientSecret: "secret",
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{"openid", "profile", "email", "federated:id"},
			RedirectURL:  "http://localhost:8000",
		},
		verifier: provider.Verifier(&oidc.Config{ClientID: "example"}),
	}
	http.HandleFunc("/login", a.handleLogin)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type app struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func (a *app) handleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.oauth2Config.AuthCodeURL("example-state"), http.StatusSeeOther)
}
