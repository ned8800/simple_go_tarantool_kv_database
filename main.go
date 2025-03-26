package main

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/tarantool/go-tarantool"
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

	opts := tarantool.Opts{User: "storage_user", Pass: "passw0rd"}
	conn, err := tarantool.Connect("127.0.0.1:3301", opts)
	if err != nil {
		log.Fatal().Err(errors.Wrap(err, "error happened")).Msg(errors.Wrap(err, "error happened").Error())
	}
	defer conn.Close()

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

}
