package withdrawal

import (
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	wsvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/withdrawal"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Controller struct {
	svc *wsvc.Service
}

func NewController(svc *wsvc.Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (a *Controller) HandlePostWithdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ww withdrawal.Withdrawal
	if err = easyjson.Unmarshal(body, &ww); err != nil {
		logger.Log.Error("easyjson.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = a.svc.Withdraw(ctx, ww)
	if err != nil {
		if errors.Is(err, wsvc.ErrOrderNotFound) {
			logger.Log.Debug(fmt.Sprintf("order not found by num = %s", ww.Order))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, wsvc.ErrInvalidSum) {
			logger.Log.Debug(fmt.Sprintf("sum is negative, sum = %f", ww.Sum))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, wsvc.ErrInsufficientFunds) {
			logger.Log.Debug(fmt.Sprintf("insufficient funds for sum = %f", ww.Sum))
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		logger.Log.Error("svc.Withdraw", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Controller) HandleGetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ws, err := a.svc.FindUserWithdrawals(ctx)
	if err != nil {
		logger.Log.Error("svc.FindUserWithdrawals", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(ws) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := easyjson.Marshal(ws)
	if err != nil {
		logger.Log.Error("easyjson.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
