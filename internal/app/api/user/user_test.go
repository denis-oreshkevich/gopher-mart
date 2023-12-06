package user

import (
	"bytes"
	"context"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/auth"
	usr "github.com/denis-oreshkevich/gopher-mart/internal/app/domain/user"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository/mock"
	usvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/user"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UserSuite struct {
	suite.Suite
}

func (s *UserSuite) SetupTest() {
	err := logger.Initialize(zapcore.DebugLevel.String())
	s.Require().NoError(err, "logger.Initialize")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

type test struct {
	name    string
	svcFunc func(*minimock.Controller) *usvc.Service
	args    string
	assert  func(w *httptest.ResponseRecorder)
}

func (s *UserSuite) TestController_HandleLoginUser() {
	password, err := auth.EncryptPassword("password")
	s.Require().NoError(err, "auth.EncryptPassword")
	musr := usr.New("user", password)
	uID := uuid.NewString()
	musr.ID = uID
	tests := []test{
		{
			name: "simple test #1",
			svcFunc: func(mc *minimock.Controller) *usvc.Service {
				userRepoMock := mock.NewUserRepositoryMock(mc).FindUserByLoginMock.
					Expect(context.Background(), "user").
					Return(musr, nil)

				balanceRepoMock := mock.NewBalanceRepositoryMock(mc)
				transactMock := mock.NewTransactMock(mc)
				usrsvc := usvc.NewService(userRepoMock, balanceRepoMock, transactMock)
				return usrsvc
			},
			args: `{
    				"login": "user",
    				"password": "password"
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusOK, w.Code)
				authH := w.Header().Get("Authorization")
				s.Assert().True(authH != "", "Authorization header is missing")
			},
		},
		{
			name: "empty password test #2",
			svcFunc: func(mc *minimock.Controller) *usvc.Service {
				userRepoMock := mock.NewUserRepositoryMock(mc)
				balanceRepoMock := mock.NewBalanceRepositoryMock(mc)
				transactMock := mock.NewTransactMock(mc)
				usrsvc := usvc.NewService(userRepoMock, balanceRepoMock, transactMock)
				return usrsvc
			},
			args: `{
    				"login": "user",
    				"password": ""
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "whitespace login test #3",
			svcFunc: func(mc *minimock.Controller) *usvc.Service {
				userRepoMock := mock.NewUserRepositoryMock(mc)
				balanceRepoMock := mock.NewBalanceRepositoryMock(mc)
				transactMock := mock.NewTransactMock(mc)
				usrsvc := usvc.NewService(userRepoMock, balanceRepoMock, transactMock)
				return usrsvc
			},
			args: `{
    				"login": " ",
    				"password": "password"
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusBadRequest, w.Code)
			},
		},
	}
	for _, tt := range tests {
		mc := minimock.NewController(s.T())
		svc := tt.svcFunc(mc)

		cont := NewController(svc)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/user/login",
			bytes.NewBufferString(tt.args))
		cont.HandleLoginUser(w, req)
		tt.assert(w)
		mc.Finish()
	}
}

func (s *UserSuite) TestController_HandleRegisterUser() {
	tests := []test{
		{
			name: "simple test #1",
			svcFunc: func(mc *minimock.Controller) *usvc.Service {
				password, err := auth.EncryptPassword("password")
				s.Require().NoError(err, "auth.EncryptPassword")
				u := usr.New("user", password)
				uID := uuid.NewString()
				u.ID = uID
				userRepoMock := mock.NewUserRepositoryMock(mc).
					CreateUserMock.Return(u, nil).
					FindUserByLoginMock.
					Expect(context.Background(), "user").Return(u, nil)

				balanceRepoMock := mock.NewBalanceRepositoryMock(mc).
					CreateBalanceMock.Expect(context.Background(), uID).Return(nil)

				transactMock := mock.NewTransactMock(mc).InTransactionMock.
					Set(func(ctx context.Context, transact func(context.Context) error) (err error) {
						return transact(ctx)
					})
				usrsvc := usvc.NewService(userRepoMock, balanceRepoMock, transactMock)
				return usrsvc
			},
			args: `{
    				"login": "user",
    				"password": "password"
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusOK, w.Code)
				authH := w.Header().Get("Authorization")
				s.Assert().True(authH != "", "Authorization header is missing")
			},
		},
		{
			name: "empty password test #2",
			svcFunc: func(mc *minimock.Controller) *usvc.Service {
				userRepoMock := mock.NewUserRepositoryMock(mc)
				balanceRepoMock := mock.NewBalanceRepositoryMock(mc)
				transactMock := mock.NewTransactMock(mc)
				usrsvc := usvc.NewService(userRepoMock, balanceRepoMock, transactMock)
				return usrsvc
			},
			args: `{
    				"login": "user",
    				"password": ""
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "whitespace login test #3",
			svcFunc: func(mc *minimock.Controller) *usvc.Service {
				userRepoMock := mock.NewUserRepositoryMock(mc)
				balanceRepoMock := mock.NewBalanceRepositoryMock(mc)
				transactMock := mock.NewTransactMock(mc)
				usrsvc := usvc.NewService(userRepoMock, balanceRepoMock, transactMock)
				return usrsvc
			},
			args: `{
    				"login": " ",
    				"password": "password"
					}`,
			assert: func(w *httptest.ResponseRecorder) {
				s.Assert().Equal(http.StatusBadRequest, w.Code)
			},
		},
	}
	for _, tt := range tests {
		mc := minimock.NewController(s.T())

		svc := tt.svcFunc(mc)
		cont := NewController(svc)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/user/register",
			bytes.NewBufferString(tt.args))
		cont.HandleRegisterUser(w, req)
		tt.assert(w)
		mc.Finish()
	}
}
