package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"simple_go_tarantool_kv_database/config"
	"simple_go_tarantool_kv_database/server"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
)

const (
	ErrLoadConfig       = "Error loading config"
	ErrStartServer      = "Error starting server"
	ErrShutdown         = "Error shutting down"
	ErrTarantoolConnect = "Error connecting to tarantool"

	tarantoolHostEnv     = "TARANTOOL_HOST"
	tarantoolUsernameEnv = "TARANTOOL_USER"
	tarantoolPasswordEnv = "TARANTOOL_PASSWORD"
)

type ValueObject struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

type TarantoolDB struct {
	db *tarantool.Connection
}

func main() {
	log.Info().Msg("Starting server")

	// setup config
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(errors.Wrap(err, ErrLoadConfig)).Msg(errors.Wrap(err, ErrLoadConfig).Error())
	}

	// setup terantool connection
	opts := tarantool.Opts{User: viper.GetString(tarantoolUsernameEnv), Pass: viper.GetString(tarantoolPasswordEnv)}
	conn, err := tarantool.Connect(viper.GetString(tarantoolHostEnv), opts)
	if err != nil {
		log.Fatal().Err(errors.Wrap(err, ErrTarantoolConnect)).Msg(errors.Wrap(err, ErrTarantoolConnect).Error())
	}
	defer conn.Close()

	// setup http server
	srv := server.New(cfg, conn)
	log.Info().Msg("Starting server")

	go func() {
		if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(errors.Wrap(err, ErrStartServer)).Msg(errors.Wrap(err, ErrStartServer).Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	///////////////////////////////////////////////////////////////////////////////////
	var tuples []ValueObject
	err = conn.InsertTyped("json_kv_database", []interface{}{
		"Indisko",
		"Indian Food"},
		&tuples)
	if err != nil {
		log.Error().Msgf("error while inserting: %v", err)
	} else {
		log.Info().Msgf("inserted value: %v", tuples)
	}

	tuples = []ValueObject{}
	err = conn.SelectTyped("json_kv_database", "primary", 0, 10, tarantool.IterEq,
		[]interface{}{"Indisko"},
		&tuples)
	if err != nil {
		log.Error().Msgf("error while deleting: %v", err)
	} else {
		log.Info().Msgf("selected value: %v", tuples)
	}

	tuples = []ValueObject{}
	err = conn.DeleteTyped("json_kv_database", "primary", []interface{}{"Indisko"}, &tuples)
	if err != nil {
		log.Error().Msgf("error while deleting: %v", err)
	} else {
		log.Info().Msgf("deleted value: %v", tuples)
	}
	///////////////////////////////////////////////////

	<-stop
	log.Info().Msg("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(errors.Wrap(err, ErrShutdown)).Msg(errors.Wrap(err, ErrShutdown).Error())
	}
	log.Info().Msg("Server is shut down gracefully")
}
