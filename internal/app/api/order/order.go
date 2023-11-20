package order

import (
	"errors"
	"fmt"
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
		logger.Log.Error("svc.FindUserOrders", zap.Error(err))
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
	match := orderNumberMatcher.MatchString(num)
	if !match {
		logger.Log.Debug(fmt.Sprintf("order number = %s is not match regex", num))
		return false
	}
	luhn := ValidLuhn(num)
	if !luhn {
		logger.Log.Debug(fmt.Sprintf("order number = %s is not valid by Luhn", num))
	}
	return luhn
}

func ValidLuhn(number string) bool {
	p := len(number) % 2
	sum := calculateLuhnSum(number, p)

	// If the total modulo 10 is not equal to 0, then the number is invalid.
	if sum%10 != 0 {
		return false
	}

	return true
}

const asciiZero = 48

func calculateLuhnSum(number string, parity int) int64 {
	var sum int64
	for i, d := range number {
		d = d - asciiZero
		// Double the value of every second digit.
		if i%2 == parity {
			d *= 2
			// If the result of this doubling operation is greater than 9.
			if d > 9 {
				// The same final result can be found by subtracting 9 from that result.
				d -= 9
			}
		}
		// Take the sum of all the digits.
		sum += int64(d)
	}
	return sum
}
