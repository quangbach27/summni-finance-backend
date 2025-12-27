package auth

import (
	"context"
	"errors"
	"fmt"
	"sumni-finance-backend/internal/config"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type keycloakClient struct {
	config   oauth2.Config
	verifier *oidc.IDTokenVerifier
}

func NewKeycloakClient() (*keycloakClient, error) {
	kcConfig := config.GetConfig().Keycloak()

	maxRetry := 10
	retryInterval := 10 * time.Second

	var provider *oidc.Provider
	var err error

	for i := 0; i < maxRetry; i++ {
		// Use a background context with a timeout for each specific attempt
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		provider, err = oidc.NewProvider(ctx, kcConfig.RealmURL())
		cancel()

		if err == nil {
			break
		}

		if i < maxRetry-1 {
			time.Sleep(retryInterval)
		}
	}

	if provider == nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("can not initizalize OIDC provider")
	}

	return &keycloakClient{
		config: oauth2.Config{
			ClientID:     kcConfig.ClientID(),
			ClientSecret: kcConfig.ClientSecret(),
			RedirectURL:  kcConfig.CallbackURL(),
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		verifier: provider.Verifier(&oidc.Config{
			ClientID: kcConfig.ClientID(),
		}),
	}, nil
}

func (k *keycloakClient) GetAuthorizationCodeURL(state string) string {
	return k.config.AuthCodeURL(state)
}

func (k *keycloakClient) Authenticate(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	freshToken, err := k.getFreshToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	rawIDToken, exist := freshToken.Extra("id_token").(string)
	if !exist {
		return nil, errors.New("missing id_token from token response")
	}

	_, err = k.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.New("token verification failed")
	}

	return freshToken, nil
}

func (k *keycloakClient) GetLogoutURL(ctx context.Context, token *oauth2.Token) (string, error) {
	freshToken, err := k.getFreshToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("token refresh failed: %w", err)
	}

	// 2. Extract ID Token for the 'id_token_hint'
	idTokenHint, exist := freshToken.Extra("id_token").(string)
	if !exist {
		return "", errors.New("missing id_token in freshToken")
	}

	kcConfig := config.GetConfig().Keycloak()
	logoutURL := fmt.Sprintf("%s/protocol/openid-connect/logout?id_token_hint=%s&client_id=%s",
		kcConfig.RealmURL(),
		idTokenHint,
		kcConfig.ClientID(),
	)

	return logoutURL, nil
}

func (k *keycloakClient) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return k.config.Exchange(ctx, code)
}

func (k *keycloakClient) getFreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	tokenSource := k.config.TokenSource(ctx, token)
	freshToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	return freshToken, nil
}
