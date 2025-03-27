package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"simple_go_tarantool_kv_database/config"
	errorconstants "simple_go_tarantool_kv_database/error_constants"
	"simple_go_tarantool_kv_database/server"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
)

const (
	tarantoolHostEnv     = "TARANTOOL_HOST"
	tarantoolUsernameEnv = "TARANTOOL_USER"
	tarantoolPasswordEnv = "TARANTOOL_PASSWORD"
)

type TarantoolDB struct {
	db *tarantool.Connection
}

func main() {
	log.Info().Msg("Starting server")

	// setup config
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(errors.Wrap(err, errorconstants.ErrLoadConfig)).Msg(errors.Wrap(err, errorconstants.ErrLoadConfig).Error())
	}

	// setup terantool connection
	opts := tarantool.Opts{User: viper.GetString(tarantoolUsernameEnv), Pass: viper.GetString(tarantoolPasswordEnv)}
	conn, err := tarantool.Connect(viper.GetString(tarantoolHostEnv), opts)
	if err != nil {
		log.Fatal().Err(errors.Wrap(err, errorconstants.ErrTarantoolConnect)).Msg(errors.Wrap(err, errorconstants.ErrTarantoolConnect).Error())
	}
	defer conn.Close()

	// setup http server
	srv := server.New(cfg, conn)
	log.Info().Msg("Starting server")

	go func() {
		if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(errors.Wrap(err, errorconstants.ErrStartServer)).Msg(errors.Wrap(err, errorconstants.ErrStartServer).Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-stop
	log.Info().Msg("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(errors.Wrap(err, errorconstants.ErrShutdown)).Msg(errors.Wrap(err, errorconstants.ErrShutdown).Error())
	}
	log.Info().Msg("Server is shut down gracefully")
}
