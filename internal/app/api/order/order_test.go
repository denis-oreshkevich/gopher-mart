package order

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/auth"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository/mock"
	osvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/order"
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

type OrderSuite struct {
	suite.Suite
}

func (s *OrderSuite) SetupTest() {
	err := logger.Initialize(zapcore.DebugLevel.String())
	s.Require().NoError(err, "logger.Initialize")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

type test struct {
	name    string
	svcFunc func(*minimock.Controller) *osvc.Service
	args    string
	assert  func(w *httptest.ResponseRecorder)
}

func (s *OrderSuite) TestController_HandleGetUserOrders() {
	bctx := context.Background()
	uID := uuid.NewString()
	ctx := context.WithValue(bctx, auth.UserIDKey{}, uID)

	ords := make([]order.Order, 0)

	tests := []test{
		{
			name: "simple test #1",
			svcFunc: func(mc *minimock.Controller) *osvc.Service {
				now := time.Now()
				o1 := order.New(uuid.NewString(), "97329", uID, order.StatusNew, 0, now)
				o2 := order.New(uuid.NewString(), "20628", uID, order.StatusNew, 22.33, now)
				ords = append(ords, o1, o2)

				ordRepo := mock.NewOrderMock(mc).FindOrdersByUserIDMock.
					Expect(ctx, uID).Return(ords, nil)
				svc := osvc.NewService(ordRepo)
				return svc
			},
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusOK, w.Code)
				orders := ords
				exp, err := json.Marshal(orders)
				s.Require().NoError(err)
				bytes, err := io.ReadAll(w.Body)
				s.Require().NoError(err)
				s.Assert().JSONEq(string(exp), string(bytes))
			},
		},
	}
	for _, tt := range tests {
		mc := minimock.NewController(s.T())
		svc := tt.svcFunc(mc)

		cont := NewController(svc)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/user/orders",
			nil)
		req = req.WithContext(ctx)

		cont.HandleGetUserOrders(w, req)
		tt.assert(w)
		mc.Finish()
	}
}

func (s *OrderSuite) TestController_HandlePostOrder() {
	bctx := context.Background()
	uID := uuid.NewString()
	ctx := context.WithValue(bctx, auth.UserIDKey{}, uID)

	tests := []test{
		{
			name: "simple test #1",
			svcFunc: func(mc *minimock.Controller) *osvc.Service {
				ordRepo := mock.NewOrderMock(mc).CreateOrderMock.
					Expect(ctx, "20628", uID).
					Return(nil)

				svc := osvc.NewService(ordRepo)
				return svc
			},
			args: `20628`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusAccepted, w.Code)
			},
		},
		{
			name: "already exist test #2",
			svcFunc: func(mc *minimock.Controller) *osvc.Service {
				ordRepo := mock.NewOrderMock(mc).CreateOrderMock.
					Expect(ctx, "20628", uID).
					Return(order.ErrOrderAlreadyExist).FindOrderByNumMock.
					Expect(ctx, "20628").
					Return(order.Order{UserID: uID}, nil)

				svc := osvc.NewService(ordRepo)
				return svc
			},
			args: `20628`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusOK, w.Code)
			},
		},
		{
			name: "already exist another user test #3",
			svcFunc: func(mc *minimock.Controller) *osvc.Service {
				newUserID := uuid.NewString()
				ordRepo := mock.NewOrderMock(mc).CreateOrderMock.
					Expect(ctx, "20628", uID).
					Return(order.ErrOrderAlreadyExist).FindOrderByNumMock.
					Expect(ctx, "20628").
					Return(order.Order{UserID: newUserID}, nil)

				svc := osvc.NewService(ordRepo)
				return svc
			},
			args: `20628`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusConflict, w.Code)
			},
		},
	}
	for _, tt := range tests {
		mc := minimock.NewController(s.T())
		svc := tt.svcFunc(mc)

		cont := NewController(svc)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/user/orders",
			bytes.NewBufferString(tt.args))
		req = req.WithContext(ctx)

		cont.HandlePostOrder(w, req)
		tt.assert(w)
		mc.Finish()
	}
}
