package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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
	http.HandleFunc("/", a.handleIndex)
	http.HandleFunc("/login", a.handleLogin)
	http.HandleFunc("/logout", a.handleLogout)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

var cookieName = "token"

type app struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

type Claims struct {
	Email           string `json:"email"`
	FederatedClaims struct {
		ConnectorID string `json:"connector_id"`
	} `json:"federated_claims"`
}

func (a *app) handleIndex(w http.ResponseWriter, r *http.Request) {
	if code := r.FormValue("code"); code != "" {
		ctx := oidc.ClientContext(r.Context(), http.DefaultClient)
		token, err := a.oauth2Config.Exchange(ctx, code)
		if err != nil {
			log.Println(err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			log.Println("no id_token in token response")
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    cookieName,
			Value:   rawIDToken,
			Expires: time.Now().Add(time.Minute),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if token, err := r.Cookie("token"); err == nil {
		if err := a.printUser(w, token.Value); err != nil {
			log.Println(err)
			a.handleLogout(w, r)
		}
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `not logged in, <a href="/login">log in</a>`)
}

func (a *app) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *app) handleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.oauth2Config.AuthCodeURL("example-state"), http.StatusSeeOther)
}

func (a *app) printUser(w http.ResponseWriter, token string) error {
	idToken, err := a.verifier.Verify(context.TODO(), token)
	if err != nil {
		return fmt.Errorf("Cannot verify token: %w", err)
	}

	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `logged in as %s from %s, <a href="/logout">logout</a>`, claims.Email, claims.FederatedClaims.ConnectorID)
	return nil
}
