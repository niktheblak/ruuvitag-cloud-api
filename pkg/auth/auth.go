package auth

import (
	"context"
)

type Authenticator interface {
	Authenticate(ctx context.Context, token string) error
}
