package order

import (
	"errors"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	osvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/order"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
)

const (
	orderNumberRegex = "^\\d{1,32}$"
)

type Controller struct {
	svc *osvc.Service
}

func NewController(svc *osvc.Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

var orderNumberMatcher = regexp.MustCompile(orderNumberRegex)

func (a *Controller) HandlePostOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bStr := string(body)
	if !isOrderNumberValid(bStr) {
		logger.Log.Debug("order number is not valid")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = a.svc.Create(ctx, bStr)
	if err != nil {
		if errors.Is(err, osvc.ErrOrderCreatedByAnotherUser) {
			logger.Log.Debug("order was already created by another user")
			w.WriteHeader(http.StatusConflict)
			return
		}
		if errors.Is(err, order.ErrOrderAlreadyExist) {
			logger.Log.Debug("order was already created by this user")
			w.WriteHeader(http.StatusOK)
			return
		}
		logger.Log.Error("mart.CreateOrder", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (a *Controller) HandleGetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ord, err := a.svc.FindUserOrders(ctx)
	if err != nil {
		logger.Log.Error("mart.FindUserOrders", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(ord) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := easyjson.Marshal(ord)
	if err != nil {
		logger.Log.Error("easyjson.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func isOrderNumberValid(num string) bool {
	return orderNumberMatcher.MatchString(num)
}
