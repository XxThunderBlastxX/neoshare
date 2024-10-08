package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
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
		ClientID:     authConfig.ClientID,
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

func (a *Authenticator) parseWellKnown() (map[string]any, error) {
	wellKnownURL := strings.TrimSuffix(a.AuthConfig.Domain, "/") + "/.well-known/openid-configuration"
	var data map[string]any

	client := resty.New()
	res, err := client.R().Get(wellKnownURL)
	if err != nil {
		return nil, err
	}
	defer res.RawBody().Close()

	_ = json.Unmarshal(res.Body(), &data)
	return data, nil
}

func (a *Authenticator) JwksURI() (string, error) {
	data, err := a.parseWellKnown()
	if err != nil {
		return "", err
	}
	return data["jwks_uri"].(string), nil
}

func (a *Authenticator) LogoutURL(token string) string {
	data, _ := a.parseWellKnown()

	var buff bytes.Buffer
	buff.WriteString(data["end_session_endpoint"].(string))
	v := url.Values{
		"client_id":                {a.ClientID},
		"token":                    {token},
		"post_logout_redirect_uri": {a.AuthConfig.LogoutCallbackURL},
	}
	buff.WriteString(fmt.Sprintf("?%s", v.Encode()))

	return buff.String()
}

func (a *Authenticator) VerifyUserInfo(token string) (bool, error) {
	userInfoURI := a.AuthConfig.UserInfoURL

	client := resty.New()
	client.Header.Set("Authorization", "Bearer "+token)
	res, err := client.R().Get(userInfoURI)
	if err != nil {
		return false, err
	}

	if res.StatusCode() != fiber.StatusOK {
		return false, errors.New("failed to get user info")
	}

	return true, nil
}
