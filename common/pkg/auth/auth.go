package auth

import (
	"context"
	"errors"
)

var ErrNotAuthorized = errors.New("not authorized")

type Authenticator interface {
	Authenticate(ctx context.Context, token string) error
}
