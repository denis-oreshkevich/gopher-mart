package rest

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/accrual"
	"github.com/mailru/easyjson"
	"net/http"
)

var _ accrual.Repository = (*Repository)(nil)

func (r *Repository) FindAccrualByOrderNum(ctx context.Context,
	num string) (accrual.Accrual, error) {
	response, err := r.client.R().
		SetContext(ctx).Get(r.conf.AccrualSystemAddress + `/api/orders/` + num)
	if err != nil {
		return accrual.Accrual{}, fmt.Errorf("client.R().Get: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusOK {
		if status == http.StatusTooManyRequests {
			return accrual.Accrual{}, accrual.ErrTooManyRequests
		}
		if status >= http.StatusNoContent {
			return accrual.Accrual{}, accrual.ErrOrderNotRegistered
		}
		return accrual.Accrual{}, fmt.Errorf("invalid response status = %d", status)
	}
	body := response.Body()
	var acc accrual.Accrual
	err = easyjson.Unmarshal(body, &acc)
	if err != nil {
		return accrual.Accrual{}, fmt.Errorf("easyjson.Unmarshal: %w", err)
	}
	return acc, nil
}
