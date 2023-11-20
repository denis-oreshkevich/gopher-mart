package accrual

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/accrual"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	balsvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/balance"
	osvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/order"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	repo     accrual.Repository
	orderSvc *osvc.Service
	balSvc   *balsvc.Service
}

func NewService(repo accrual.Repository,
	orderSvc *osvc.Service, balSvc *balsvc.Service) *Service {
	return &Service{
		repo:     repo,
		orderSvc: orderSvc,
		balSvc:   balSvc,
	}
}

func (s *Service) Process(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("interrupt process accruals")
			return
		case <-ticker.C:
			logger.Log.Info("start process accruals")
			s.process(ctx)
			logger.Log.Info("finish process accruals")
		}
	}
}

func (s *Service) process(ctx context.Context) {
	ords, err := s.orderSvc.StartProcessing(ctx, 3)
	if err != nil {
		logger.Log.Error("find new orders", zap.Error(err))
	}
	//TODO tx
	for _, o := range ords {
		acc, err := s.repo.FindByOrderNum(ctx, o.Number)
		st := order.ProcessedStatus
		if err != nil {
			logger.Log.Error("find accrual by order number", zap.Error(err))
			if errors.Is(err, accrual.ErrOrderNotRegistered) {
				st = order.InvalidStatus
			}
			if errors.Is(err, accrual.ErrTooManyRequests) {
				st = order.NewStatus
				//возвращаем статус в новый, чтобы попытаться еще раз
			}
			logger.Log.Debug(fmt.Sprintf("order num = %s set status = %s", o.Number, st))
		}
		upErr := s.orderSvc.UpdateStatusByID(ctx, o.ID, st)
		if upErr != nil {
			logger.Log.Error("update order status", zap.Error(err))
			continue
		}
		refErr := s.balSvc.RefillByUserID(ctx, acc.Accrual, o.UserID)
		if refErr != nil {
			logger.Log.Error("update order status", zap.Error(err))
			continue
		}
	}
}
