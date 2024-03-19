package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/server/api"
	"github.com/denis-oreshkevich/gophkeeper/internal/server/repo/postgres"
	"github.com/denis-oreshkevich/gophkeeper/internal/server/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/config"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Run() error {
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

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	pgRepo, err := postgres.NewRepository(ctx, conf.DataBaseURI)
	if err != nil {
		return fmt.Errorf("postgres.NewRepository: %w", err)
	}
	defer pgRepo.Close()

	serverService := service.NewServerService(pgRepo)
	controller := api.NewController(serverService)

	router, err := SetUpRouter(ctx, controller)
	if err != nil {
		return fmt.Errorf("SetUpRouter: %w", err)
	}

	srv := &http.Server{
		Addr:    conf.ServerAddress,
		Handler: router,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log.Error("HTTP server Shutdown", zap.Error(err))
		}
	}()

	err = os.Mkdir("certs", 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("os.Mkdir: %w", err)
		}
	}

	manager, errHTTPS := auth.NewCertManager("certs/cert.pem", "certs/key.pem")
	if errHTTPS != nil {
		return fmt.Errorf("auth.NewCertManager: %w", errHTTPS)
	}
	logger.Log.Info("server started")
	err = srv.ListenAndServeTLS(manager.CertPath, manager.KeyPath)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("srv.ListenAndServe(): %s\n", err)
	}

	logger.Log.Info("server started")

	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("Server closed")
		} else {
			return fmt.Errorf("router run %w", err)
		}
	}

	wg.Wait()
	return nil
}

func SetUpRouter(ctx context.Context, controller *api.Controller) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(api.Auth)
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", controller.HandleRegisterUser)
			r.Post("/login", controller.HandleLoginUser)
			r.Route("/client", func(r chi.Router) {
				r.Post("/", controller.HandlePostClient)
				r.Put("/", controller.HandlePutClient)
			})
			r.Route("/credentials", func(r chi.Router) {
				r.Get("/", controller.HandleGetUserCredentials)
				r.Get("/{id}", controller.HandleGetCredentialsByID)
				r.Post("/", controller.HandlePostCredentials)
				r.Delete("/{id}", controller.HandleDeleteCredentialsByID)
				r.Post("/sync", controller.HandlePostSyncCredentials)
			})
			r.Route("/cards", func(r chi.Router) {
				r.Get("/", controller.HandleGetUserCards)
				r.Get("/{id}", controller.HandleGetCardByID)
				r.Post("/", controller.HandlePostCard)
				r.Delete("/{id}", controller.HandleDeleteCardByID)
				r.Post("/sync", controller.HandlePostSyncCard)
			})
			r.Route("/texts", func(r chi.Router) {
				r.Get("/", controller.HandleGetUserTexts)
				r.Get("/{id}", controller.HandleGetTextByID)
				r.Post("/", controller.HandlePostText)
				r.Delete("/{id}", controller.HandleDeleteTextByID)
				r.Post("/sync", controller.HandlePostSyncText)
			})
			r.Route("/binaries", func(r chi.Router) {
				r.Get("/", controller.HandleGetUserBinaries)
				r.Get("/{id}", controller.HandleGetBinaryByID)
				r.Post("/", controller.HandlePostBinary)
				r.Delete("/{id}", controller.HandleDeleteBinaryByID)
				r.Post("/sync", controller.HandlePostSyncBinary)
			})

		})
	})
	return r, nil
}
