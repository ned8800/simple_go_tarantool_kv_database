package repo

import (
	"context"
	errorconstants "simple_go_tarantool_kv_database/error_constants"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/tarantool/go-tarantool"
)

type ValueObject struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

// post
func InsertValueByKey(ctx context.Context, conn *tarantool.Connection, key, data string) (*ValueObject, error) {
	logger := log.Ctx(ctx)

	tuples := []ValueObject{}
	err := conn.InsertTyped("json_kv_database", []interface{}{
		key,
		data},
		&tuples)
	if err != nil {
		if strings.Contains(err.Error(), errorconstants.ErrDuplicateKey) {
			errMsg := errors.Wrap(errorconstants.ErrKeyAlreadyExists, err.Error())
			logger.Error().Err(err).Msgf("%s: %v", errorconstants.ErrKeyAlreadyExists.Error(), err)
			return nil, errMsg
		}
		logger.Error().Err(err).Msgf("%s: %v", errorconstants.ErrInsertValue.Error(), err)
		return nil, err
	}

	if len(tuples) != 0 {
		logger.Info().Msgf("inserted value: %v", tuples)
		return &tuples[0], nil
	}

	logger.Error().Err(errorconstants.ErrInsertValue).Msgf("%s: %s, %s", errorconstants.ErrInsertValue.Error(), key, data)
	return nil, errorconstants.ErrInsertValue
}

// get
func GetValueByKey(ctx context.Context, conn *tarantool.Connection, key string) (*ValueObject, error) {
	logger := log.Ctx(ctx)

	tuples := []ValueObject{}
	err := conn.SelectTyped("json_kv_database", "primary", 0, 1, tarantool.IterEq,
		[]interface{}{key},
		&tuples)
	if err != nil {
		logger.Error().Err(err).Msgf("error while getting value by id: %v", err)
		return nil, err
	}

	if len(tuples) != 0 {
		logger.Info().Msgf("selected value: %v", tuples)
		return &tuples[0], nil
	}

	logger.Error().Err(errorconstants.ErrNotFoundById).Msgf("%s: %s", errorconstants.ErrNotFoundById.Error(), key)
	return nil, errorconstants.ErrNotFoundById
}

// put
func UpdateValueByKey(ctx context.Context, conn *tarantool.Connection, key, data string) (*ValueObject, error) {
	logger := log.Ctx(ctx)

	tuples := []ValueObject{}

	err := conn.UpdateTyped("json_kv_database", "primary", []interface{}{key},
		[][]interface{}{
			{"=", "data", data},
		},
		&tuples)
	if err != nil {
		logger.Error().Err(err).Msgf("%s: %v", errorconstants.ErrUpdateValue.Error(), err)
		return nil, err
	}

	if len(tuples) != 0 {
		logger.Info().Msgf("updated value: %v", tuples)
		return &tuples[0], nil
	}

	logger.Error().Err(errorconstants.ErrUpdateValue).Msgf("%s: %s, %s", errorconstants.ErrUpdateValue.Error(), key, data)
	return nil, errorconstants.ErrUpdateValue
}

// delete
func DeleteValueByKey(ctx context.Context, conn *tarantool.Connection, key string) (*ValueObject, error) {
	logger := log.Ctx(ctx)

	tuples := []ValueObject{}
	err := conn.DeleteTyped("json_kv_database", "primary", []interface{}{key}, &tuples)
	if err != nil {
		logger.Error().Err(err).Msgf("%s: %v", errorconstants.ErrDeleteValue.Error(), err)
		return nil, err
	}

	if len(tuples) != 0 {
		logger.Info().Msgf("deleted value: %v", tuples)
		return &tuples[0], nil
	}

	logger.Error().Err(errorconstants.ErrDeleteValue).Msgf("%s: %s", errorconstants.ErrDeleteValue.Error(), key)
	return nil, errorconstants.ErrDeleteValue
}
