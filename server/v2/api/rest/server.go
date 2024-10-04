package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/server/v2/appmanager"
)

const (
	ServerName = "rest-v2"
)

type Server[T transaction.Tx] struct {
	logger log.Logger
	router *mux.Router

	httpServer *http.Server
	config     *Config
	cfgOptions []CfgOption
}

func New[T transaction.Tx](cfgOptions ...CfgOption) *Server[T] {
	return &Server[T]{
		cfgOptions: cfgOptions,
	}
}

func (s *Server[T]) Name() string {
	return ServerName
}

func (s *Server[T]) Init(appI serverv2.AppI[T], cfg map[string]any, logger log.Logger) error {
	s.logger = logger.With(log.ModuleKey, s.Name())

	s.config = s.Config().(*Config)

	var appManager *appmanager.AppManager[T]
	appManager = appI.GetAppManager()

	s.router = mux.NewRouter()
	s.router.PathPrefix("/").Handler(NewDefaultHandler(appManager))

	return nil
}

func (s *Server[T]) Start(ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:    s.config.Address,
		Handler: s.router,
	}

	go func() {
		s.logger.Info("Starting HTTP server", "address", s.config.Address)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Failed to start HTTP server", "error", err)
		}
	}()

	return nil
}

func (s *Server[T]) Stop(ctx context.Context) error {
	if !s.config.Enable {
		return nil
	}

	s.logger.Info("Stopping HTTP server")

	return s.httpServer.Shutdown(ctx)
}

func (s *Server[T]) Config() any {
	if s.config == nil || s.config == (&Config{}) {
		cfg := DefaultConfig()

		for _, opt := range s.cfgOptions {
			opt(cfg)
		}

		return cfg
	}

	return s.config
}
