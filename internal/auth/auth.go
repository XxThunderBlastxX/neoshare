package auth

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/XxThunderBlastxX/neoshare/internal/config"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

// New instantiates the *Authenticator.
func New(authConfig *config.AuthConfig) (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+authConfig.Domain+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     authConfig.ClientId,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  authConfig.CallbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

// RevokeToken TODO: Fix this revoke token function.
func (a *Authenticator) RevokeToken() string {
	var buff bytes.Buffer

	logoutURL := strings.Replace(a.Provider.Endpoint().AuthURL+"/", "/authorize", "/v2/logout", 1)
	buff.WriteString(logoutURL)

	params := url.Values{
		"client_id": {a.ClientID},
		"returnTo":  {a.RedirectURL},
	}

	buff.WriteString(params.Encode())
	return buff.String()
}
