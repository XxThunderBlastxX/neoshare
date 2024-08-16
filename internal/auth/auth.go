package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"

	"github.com/XxThunderBlastxX/neoshare/internal/config"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
	AuthConfig *config.AuthConfig
}

// New instantiates the *Authenticator.
func New(authConfig *config.AuthConfig) (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		authConfig.Domain,
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     authConfig.ClientId,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  authConfig.CallbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email"},
	}

	return &Authenticator{
		Provider:   provider,
		Config:     conf,
		AuthConfig: authConfig,
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

func (a *Authenticator) parseWellKnown() (map[string]interface{}, error) {
	wellKnownUrl := strings.TrimSuffix(a.AuthConfig.Domain, "/") + "/.well-known/openid-configuration"
	var data map[string]interface{}

	client := resty.New()
	res, err := client.R().Get(wellKnownUrl)
	if err != nil {
		return nil, err
	}
	defer res.RawBody().Close()

	json.Unmarshal(res.Body(), &data)
	return data, nil
}

func (a *Authenticator) JwksUri() (string, error) {
	data, err := a.parseWellKnown()
	if err != nil {
		return "", err
	}
	return data["jwks_uri"].(string), nil
}

func (a *Authenticator) LogoutURL() (string, error) {
	data, err := a.parseWellKnown()
	if err != nil {
		return "", err
	}

	return data["end_session_endpoint"].(string) + fmt.Sprintf("?client_id=%s", a.ClientID), nil
}

func (a *Authenticator) VerifyUserInfo(token string) (bool, error) {
	userInfoUri := a.AuthConfig.UserInfoURL

	client := resty.New()
	client.Header.Set("Authorization", "Bearer "+token)
	res, err := client.R().Get(userInfoUri)
	if err != nil {
		return false, err
	}

	if res.StatusCode() != 200 {
		return false, errors.New("failed to get user info")
	}

	return true, nil
}
