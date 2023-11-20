package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/api"
	bapi "github.com/denis-oreshkevich/gopher-mart/internal/app/api/balance"
	oapi "github.com/denis-oreshkevich/gopher-mart/internal/app/api/order"
	uapi "github.com/denis-oreshkevich/gopher-mart/internal/app/api/user"
	wapi "github.com/denis-oreshkevich/gopher-mart/internal/app/api/withdrawal"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/config"
	accr "github.com/denis-oreshkevich/gopher-mart/internal/app/domain/accrual/rest"
	bp "github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance/postgres"
	op "github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order/postgres"
	up "github.com/denis-oreshkevich/gopher-mart/internal/app/domain/user/postgres"
	wp "github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal/postgres"
	ipg "github.com/denis-oreshkevich/gopher-mart/internal/app/infra/postgres"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	accsvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/accrual"
	bsvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/balance"
	osvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/order"
	usvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/user"
	wsvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/withdrawal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	err := logger.Initialize(zapcore.DebugLevel.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "logger initialize", err)
		os.Exit(1)
	}
	defer logger.Log.Sync()

	conf, err := config.Parse()
	if err != nil {
		logger.Log.Fatal("parse config", zap.Error(err))
	}

	if err := run(ctx, conf); err != nil {
		logger.Log.Fatal("main error", zap.Error(err))
	}
	logger.Log.Info("server exited properly")
}

func run(ctx context.Context, conf *config.Config) error {
	pg, err := ipg.New(ctx, conf.DataBaseURI())
	if err != nil {
		return fmt.Errorf("ipg.New: %w", err)
	}
	defer pg.Close()
	urepo := up.NewUserRepository(pg)
	brepo := bp.NewBalanceRepository(pg)
	orepo := op.NewOrderRepository(pg)
	wrepo := wp.NewWithdrawalRepository(pg)
	accrepo := accr.NewAccrualRepository(resty.New(), conf.AccrualSystemAddress())

	balSvc := bsvc.NewService(brepo)
	userSvc := usvc.NewService(urepo, balSvc)
	ordSvc := osvc.NewService(orepo)
	withSvc := wsvc.NewService(wrepo, ordSvc, balSvc)

	accSvc := accsvc.NewService(accrepo, ordSvc, balSvc)

	uAPI := uapi.NewController(userSvc)
	balAPI := bapi.NewController(balSvc)
	ordAPI := oapi.NewController(ordSvc)
	withAPI := wapi.NewController(withSvc)

	bctx := context.Background()
	ctx, cancel := context.WithCancel(bctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go accSvc.Process(ctx)

	r := setUpRouter(uAPI, balAPI, ordAPI, withAPI)

	srv := &http.Server{
		Addr:    conf.ServerAddress(),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Log.Info("server Started")

	select {
	case <-c:
		logger.Log.Info("cancelling context")
		cancel()
	case <-ctx.Done():
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	return nil
}

func setUpRouter(uAPI *uapi.Controller,
	balAPI *bapi.Controller,
	ordAPI *oapi.Controller,
	withAPI *wapi.Controller) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(api.Auth)
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", uAPI.HandleRegisterUser)
			r.Post("/login", uAPI.HandleLoginUser)
			r.Post("/orders", ordAPI.HandlePostOrder)
			r.Get("/orders", ordAPI.HandleGetUserOrders)
			r.Get("/withdrawals", withAPI.HandleGetUserWithdrawals)
			r.Get("/balance", balAPI.HandleGetUserBalance)
			r.Route("/balance", func(r chi.Router) {
				r.Post("/withdraw", withAPI.HandlePostWithdraw)
			})

		})
	})
	return r
}
