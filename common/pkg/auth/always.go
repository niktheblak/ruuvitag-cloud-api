package auth

import (
	"context"
)

type AlwaysAllowAuthenticator struct {
}

func (a *AlwaysAllowAuthenticator) Authenticate(ctx context.Context, token string) error {
	return nil
}

func AlwaysAllow() Authenticator {
	return &AlwaysAllowAuthenticator{}
}
