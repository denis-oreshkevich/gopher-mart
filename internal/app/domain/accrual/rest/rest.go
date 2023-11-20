package rest

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/accrual"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

type AccrualRepository struct {
	client  *resty.Client
	baseURL string
}

func NewAccrualRepository(client *resty.Client, baseURL string) *AccrualRepository {
	return &AccrualRepository{
		client:  client,
		baseURL: baseURL,
	}
}

var _ accrual.Repository = (*AccrualRepository)(nil)

func (a *AccrualRepository) FindByOrderNum(ctx context.Context, num string) (accrual.Accrual, error) {
	response, err := a.client.R().SetContext(ctx).Get(a.baseURL + `/api/orders/` + num)
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
