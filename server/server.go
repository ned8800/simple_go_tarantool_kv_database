package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"simple_go_tarantool_kv_database/config"
	"simple_go_tarantool_kv_database/delivery"

	"github.com/rs/zerolog/log"
	"github.com/tarantool/go-tarantool"
)

type Config struct {
	Address         string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type Server struct {
	Config        *config.Config
	TarantoolConn *tarantool.Connection
	httpServer    *http.Server
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down server")
	return s.httpServer.Shutdown(ctx)
}

func New(cfg *config.Config, tarantoolConn *tarantool.Connection) *Server {
	log.Info().Msg("Initializing server")

	s := &Server{
		Config:        cfg,
		TarantoolConn: tarantoolConn,
	}

	log.Info().Msg("Server initialized successfully")
	return s
}

func (s *Server) Run() error {

	mx := delivery.NewRouter()

	log.Info().Msg("Configuring routes")

	delivery.ApplyMiddlewares(mx)
	delivery.SetupRoutes(mx, s.TarantoolConn)

	log.Info().Msg("Routes configured successfully")

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Config.Server.Address, s.Config.Server.Port),
		ReadTimeout:  s.Config.Server.ReadTimeout,
		WriteTimeout: s.Config.Server.WriteTimeout,
		IdleTimeout:  s.Config.Server.IdleTimeout,
		Handler:      mx,
	}

	s.httpServer = srv

	log.Info().Msg("Running server")
	return s.httpServer.ListenAndServe()
}
