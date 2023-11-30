package withdrawal

import (
	"context"
)

type Repository interface {
	RegisterWithdrawal(ctx context.Context, withdraw Withdrawal) error
	FindWithdrawalsByUserID(ctx context.Context, userID string) ([]Withdrawal, error)
}
