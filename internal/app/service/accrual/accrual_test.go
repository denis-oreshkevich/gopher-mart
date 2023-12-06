package accrual

import (
	"context"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/auth"
	accr "github.com/denis-oreshkevich/gopher-mart/internal/app/domain/accrual"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestService_process(t *testing.T) {
	err := logger.Initialize(zapcore.DebugLevel.String())
	require.NoError(t, err, "logger.Initialize")

	bctx := context.Background()
	uID := uuid.NewString()
	ctx := context.WithValue(bctx, auth.UserIDKey{}, uID)

	tests := []struct {
		name    string
		svcFunc func(mc *minimock.Controller) *Service
	}{
		{
			name: "simple test #1",
			svcFunc: func(mc *minimock.Controller) *Service {
				ordNum := "97329"
				sum := 85.7
				acc := accr.New(ordNum, accr.StatusProcessed, sum)
				aRepo := mock.NewAccrualMock(mc).FindAccrualByOrderNumMock.
					Expect(ctx, ordNum).Return(acc, nil)

				balRepo := mock.NewBalanceRepositoryMock(mc).RefillBalanceByUserIDMock.
					Expect(ctx, sum, uID).
					Return(nil)

				ordID := uuid.NewString()
				now := time.Now()
				ord := order.New(ordID, ordNum, uID, order.StatusNew, 0, now)
				ords := []order.Order{ord}
				ordRepo := mock.NewOrderMock(mc).StartOrderProcessingMock.
					Expect(ctx, 3).
					Return(ords, nil).
					UpdateOrderStatusByIDMock.
					Expect(ctx, ordID, sum, order.StatusProcessed).
					Return(nil)

				transactMock := mock.NewTransactMock(mc).InTransactionMock.
					Set(func(ctx context.Context, transact func(context.Context) error) (err error) {
						return transact(ctx)
					})

				svc := NewService(aRepo, ordRepo, balRepo, transactMock)
				return svc
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			svc := tt.svcFunc(mc)
			svc.process(ctx)
		})
	}
}
