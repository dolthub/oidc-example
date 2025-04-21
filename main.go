package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
)

func main() {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURL := os.Getenv("REDIRECT_URL")
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" || keycloakURL == "" {
		log.Fatal("Missing environment variables")
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, keycloakURL)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	fmt.Println("Server starting at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<a href="/login">Login with Keycloak</a>`)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	state := "random-state" // In production, generate a random string and validate
	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token in token response", http.StatusInternalServerError)
		return
	}

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var claims struct {
		Email    string `json:"email"`
		Username string `json:"preferred_username"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Hello %s (%s)", claims.Username, claims.Email)
}
