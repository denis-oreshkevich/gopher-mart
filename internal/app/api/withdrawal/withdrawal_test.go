package withdrawal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/auth"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository/mock"
	wsvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/withdrawal"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type WithdrawalSuite struct {
	suite.Suite
}

func (s *WithdrawalSuite) SetupTest() {
	err := logger.Initialize(zapcore.DebugLevel.String())
	s.Require().NoError(err, "logger.Initialize")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawalSuite))
}
func (s *WithdrawalSuite) TestController_HandleGetUserWithdrawals() {
	bctx := context.Background()
	uID := uuid.NewString()
	ctx := context.WithValue(bctx, auth.UserIDKey{}, uID)
	wl := make([]withdrawal.Withdrawal, 0)

	tests := []struct {
		name    string
		svcFunc func(mc *minimock.Controller) *wsvc.Service
		assert  func(w *httptest.ResponseRecorder)
	}{
		{
			name: "simple get test #1",
			svcFunc: func(mc *minimock.Controller) *wsvc.Service {
				now := time.Now()
				w1 := withdrawal.New("97329", 321.3, now)
				w2 := withdrawal.New("20628", 11.93, now)
				wl = append(wl, w1, w2)

				wRepo := mock.NewWithdrawalMock(mc).FindWithdrawalsByUserIDMock.
					Expect(ctx, uID).Return(wl, nil)

				balRepo := mock.NewBalanceRepositoryMock(mc)
				ordRepo := mock.NewOrderMock(mc)

				transactMock := mock.NewTransactMock(mc)

				svc := wsvc.NewService(wRepo, ordRepo, balRepo, transactMock)
				return svc
			},
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusOK, w.Code)
				wws := wl
				exp, err := json.Marshal(wws)
				s.Require().NoError(err)
				bytes, err := io.ReadAll(w.Body)
				s.Require().NoError(err)
				s.Assert().JSONEq(string(exp), string(bytes))
			},
		},
		{
			name: "no records test #2",
			svcFunc: func(mc *minimock.Controller) *wsvc.Service {

				wRepo := mock.NewWithdrawalMock(mc).FindWithdrawalsByUserIDMock.
					Expect(ctx, uID).Return(nil, nil)

				balRepo := mock.NewBalanceRepositoryMock(mc)
				ordRepo := mock.NewOrderMock(mc)

				transactMock := mock.NewTransactMock(mc)

				svc := wsvc.NewService(wRepo, ordRepo, balRepo, transactMock)
				return svc
			},
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusNoContent, w.Code)
			},
		},
	}
	for _, tt := range tests {
		mc := minimock.NewController(s.T())
		defer mc.Finish()

		svc := tt.svcFunc(mc)
		cont := NewController(svc)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
		req = req.WithContext(ctx)
		cont.HandleGetUserWithdrawals(w, req)
		tt.assert(w)
	}
}

func (s *WithdrawalSuite) TestController_HandlePostWithdraw() {

	bctx := context.Background()
	uID := uuid.NewString()
	ctx := context.WithValue(bctx, auth.UserIDKey{}, uID)

	tests := []struct {
		name    string
		svcFunc func(mc *minimock.Controller) *wsvc.Service
		args    string
		assert  func(w *httptest.ResponseRecorder)
	}{
		{
			name: "simple get test #1",
			svcFunc: func(mc *minimock.Controller) *wsvc.Service {
				sum := 31.8
				ordNum := "97329"
				w := withdrawal.Withdrawal{Order: ordNum, Sum: sum}
				wRepo := mock.NewWithdrawalMock(mc).RegisterWithdrawalMock.
					Expect(ctx, w).
					Return(nil)

				balRepo := mock.NewBalanceRepositoryMock(mc).WithdrawBalanceByUserIDMock.
					Expect(ctx, sum, uID).
					Return(nil)

				ordRepo := mock.NewOrderMock(mc).CreateOrderMock.
					Expect(ctx, ordNum, uID).
					Return(nil)

				transactMock := mock.NewTransactMock(mc).InTransactionMock.
					Set(func(ctx context.Context, transact func(context.Context) error) (err error) {
						return transact(ctx)
					})

				svc := wsvc.NewService(wRepo, ordRepo, balRepo, transactMock)
				return svc
			},
			args: `{
    				"order": "97329",
    				"sum": 31.8
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusOK, w.Code)
			},
		},
		{
			name: "duplicate order test #2",
			svcFunc: func(mc *minimock.Controller) *wsvc.Service {
				ordNum := "97329"
				wRepo := mock.NewWithdrawalMock(mc)

				balRepo := mock.NewBalanceRepositoryMock(mc)

				ordRepo := mock.NewOrderMock(mc).CreateOrderMock.
					Expect(ctx, ordNum, uID).
					Return(order.ErrOrderAlreadyExist)

				transactMock := mock.NewTransactMock(mc).InTransactionMock.
					Set(func(ctx context.Context, transact func(context.Context) error) (err error) {
						return transact(ctx)
					})

				svc := wsvc.NewService(wRepo, ordRepo, balRepo, transactMock)
				return svc
			},
			args: `{
    				"order": "97329",
    				"sum": 12.4
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "not enough balance test #3",
			svcFunc: func(mc *minimock.Controller) *wsvc.Service {
				sum := 10000.2
				ordNum := "97329"
				wRepo := mock.NewWithdrawalMock(mc)

				balRepo := mock.NewBalanceRepositoryMock(mc).WithdrawBalanceByUserIDMock.
					Expect(ctx, sum, uID).Return(balance.ErrCheckConstraint)

				ordRepo := mock.NewOrderMock(mc).CreateOrderMock.
					Expect(ctx, ordNum, uID).
					Return(nil)

				transactMock := mock.NewTransactMock(mc).InTransactionMock.
					Set(func(ctx context.Context, transact func(context.Context) error) (err error) {
						return transact(ctx)
					})

				svc := wsvc.NewService(wRepo, ordRepo, balRepo, transactMock)
				return svc
			},
			args: `{
    				"order": "97329",
    				"sum": 10000.2
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusPaymentRequired, w.Code)
			},
		},
	}
	for _, tt := range tests {
		mc := minimock.NewController(s.T())
		defer mc.Finish()

		svc := tt.svcFunc(mc)
		cont := NewController(svc)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/user/balance/withdraw",
			bytes.NewBufferString(tt.args))

		req = req.WithContext(ctx)
		cont.HandlePostWithdraw(w, req)
		tt.assert(w)
	}
}
