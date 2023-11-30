package rest

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/accrual"
	"github.com/mailru/easyjson"
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
	if status != 200 {
		if status == 429 {
			return accrual.Accrual{}, accrual.ErrTooManyRequests
		}
		if status >= 204 {
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
