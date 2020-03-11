package auth

import (
	"context"
	"errors"
)

var ErrNotAuthorized = errors.New("not authorized")

type StaticAuthenticator struct {
	AllowedTokens []string
}

func (a *StaticAuthenticator) Authenticate(ctx context.Context, token string) error {
	for _, t := range a.AllowedTokens {
		if token == t {
			return nil
		}
	}
	return ErrNotAuthorized
}
