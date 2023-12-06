package balance

import (
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/service/balance"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type Controller struct {
	svc *balance.Service
}

func NewController(svc *balance.Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (c *Controller) HandleGetUserBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bal, err := c.svc.FindUserBalance(ctx)
	if err != nil {
		logger.Log.Error("svc.FindUserBalance", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := easyjson.Marshal(bal)
	if err != nil {
		logger.Log.Error("easyjson.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
