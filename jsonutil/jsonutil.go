package jsonutil

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	ErrEncodeJSON      = "Error encoding JSON"
	ErrEncodeJSONShort = "encode_json_error"
	ErrCloseBody       = "Error closing body"
	ErrParseJSON       = "Error parsing JSON"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func SendError(ctx context.Context, w http.ResponseWriter, errCode int, errResp string) {
	logger := log.Ctx(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)

	errResponse := ErrorResponse{
		Error: errResp,
	}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		logger.Error().Err(errors.Wrap(err, ErrEncodeJSON)).Msg(errors.Wrap(err, ErrEncodeJSON).Error())
	}
}

func ReadJSON(ctx context.Context, r *http.Request, data interface{}) error {
	logger := log.Ctx(ctx)

	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Error().Err(errors.Wrap(err, ErrCloseBody)).Msg(errors.Wrap(err, ErrCloseBody).Error())
		}
	}()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return errors.Wrap(err, ErrParseJSON)
	}
	return nil
}

func SendJSON(ctx context.Context, w http.ResponseWriter, data interface{}) error {
	logger := log.Ctx(ctx)

	w.Header().Set("Content-Type", "application/json")

	code := http.StatusOK
	w.WriteHeader(code)

	if data == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error().Err(errors.Wrap(err, ErrEncodeJSON)).Msg(errors.Wrap(err, ErrEncodeJSON).Error())
		SendError(ctx, w, http.StatusInternalServerError, errors.Wrap(err, ErrEncodeJSONShort).Error())
		return errors.Wrap(err, ErrParseJSON)
	}
	return nil
}
