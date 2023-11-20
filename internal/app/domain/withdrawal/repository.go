package withdrawal

import (
	"context"
)

type Repository interface {
	Register(ctx context.Context, withdraw Withdrawal) error
	FindByUserID(ctx context.Context, userID string) ([]Withdrawal, error)
}
