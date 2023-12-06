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
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository/postgres"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository/rest"
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
	"syscall"
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
	pgRepo, err := postgres.NewRepository(ctx, conf.DataBaseURI)
	if err != nil {
		return fmt.Errorf("postgres.NewRepository: %w", err)
	}
	defer pgRepo.Close()
	restRepo := rest.NewRepository(resty.New(), conf)

	balSvc := bsvc.NewService(pgRepo)
	userSvc := usvc.NewService(pgRepo, pgRepo, pgRepo)
	ordSvc := osvc.NewService(pgRepo)
	withSvc := wsvc.NewService(pgRepo, pgRepo, pgRepo, pgRepo)

	accSvc := accsvc.NewService(restRepo, pgRepo, pgRepo, pgRepo)

	uAPI := uapi.NewController(userSvc)
	balAPI := bapi.NewController(balSvc)
	ordAPI := oapi.NewController(ordSvc)
	withAPI := wapi.NewController(withSvc)

	bctx := context.Background()
	ctx, cancel := context.WithCancel(bctx)
	c := make(chan os.Signal, 1)
	pc := make(chan struct{}, 1)
	go func() {
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		logger.Log.Info("cancelling context")
		cancel()
		close(c)

	}()

	go accSvc.Process(ctx, pc)

	r := setUpRouter(uAPI, balAPI, ordAPI, withAPI)

	srv := &http.Server{
		Addr:    conf.ServerAddress,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Log.Info("server started")

	<-pc
	logger.Log.Info("context cancelled")

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
			r.Route("/balance", func(r chi.Router) {
				r.Get("/", balAPI.HandleGetUserBalance)
				r.Post("/withdraw", withAPI.HandlePostWithdraw)
			})

		})
	})
	return r
}
