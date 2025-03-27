package tarantoolmanager

import (
	"errors"
	"net/http"
	"simple_go_tarantool_kv_database/jsonutil"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/tarantool/go-tarantool"
)

type TarantoolManager struct {
	TarantoolConn *tarantool.Connection
}

func NewTarantoolManager(tarantoolConn *tarantool.Connection) *TarantoolManager {
	tm := &TarantoolManager{
		TarantoolConn: tarantoolConn,
	}

	return tm
}

// post
func InsertValueHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка при разборе данных формы", http.StatusBadRequest)
		logger.Error().Msgf("Ошибка при разборе данных формы: %v", err)
		return
	}

}

// get
func GetValueHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	valID := mux.Vars(r)["id"]
	if valID == "" {
		errMsg := errors.New("bad_request")
		logger.Error().Err(errMsg).Msg("bad_request")
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errMsg.Error())

		return
	}

}

// put
func UpdateValueHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	valID := mux.Vars(r)["id"]
	if valID == "" {
		errMsg := errors.New("bad_request")
		logger.Error().Err(errMsg).Msg("bad_request")
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errMsg.Error())

		return
	}

}

// delete
func DeleteValueHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	valID := mux.Vars(r)["id"]
	if valID == "" {
		errMsg := errors.New("bad_request")
		logger.Error().Err(errMsg).Msg("bad_request")
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errMsg.Error())

		return
	}

}
