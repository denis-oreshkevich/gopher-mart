package accrual

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/accrual"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	accRepo   accrual.Repository
	orderRepo order.Repository
	balRepo   balance.Repository
	tr        repository.Transactor
}

func NewService(accRepo accrual.Repository,
	orderRepo order.Repository, balRepo balance.Repository, tr repository.Transactor) *Service {
	return &Service{
		accRepo:   accRepo,
		orderRepo: orderRepo,
		balRepo:   balRepo,
		tr:        tr,
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
	ords, err := s.orderRepo.StartOrderProcessing(ctx, 3)
	if err != nil {
		logger.Log.Error("find new orders", zap.Error(err))
		return
	}

	for _, o := range ords {
		log := logger.Log.With(zap.String("order", o.Number))
		err = s.tr.InTransaction(ctx, func(ctx context.Context) error {
			acc, err := s.accRepo.FindAccrualByOrderNum(ctx, o.Number)
			st := order.StatusProcessed
			log.Debug(fmt.Sprintf("processing status = %s, accrual = %f",
				acc.Status, acc.Accrual))
			if err != nil {
				log.Error("find accrual by order number", zap.Error(err))
				if errors.Is(err, accrual.ErrOrderNotRegistered) {
					st = order.StatusInvalid
				}
				if errors.Is(err, accrual.ErrTooManyRequests) {
					//возвращаем статус в новый, чтобы попытаться еще раз
					st = order.StatusNew
					log.Error("accRepo.FindAccrualByOrderNum", zap.Error(err))
				}
				log.Debug(fmt.Sprintf("set status = %s", st))
			}
			upErr := s.orderRepo.UpdateOrderStatusByID(ctx, o.ID, acc.Accrual, st)
			if upErr != nil {
				log.Error("update order status", zap.Error(err))
				return upErr
			}
			if st == order.StatusNew {
				log.Debug("returned from processing")
				return nil
			}
			refErr := s.balRepo.RefillBalanceByUserID(ctx, acc.Accrual, o.UserID)
			if refErr != nil {
				log.Error("refill balance", zap.Error(err))
				return refErr
			}
			return nil
		})
		if err != nil {
			log.Error("transaction", zap.Error(err))
		}
	}
}
