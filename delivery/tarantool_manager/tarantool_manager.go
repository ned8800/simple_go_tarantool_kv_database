package tarantoolmanager

import (
	"encoding/json"
	"net/http"

	errorconstants "simple_go_tarantool_kv_database/error_constants"
	"simple_go_tarantool_kv_database/jsonutil"
	"simple_go_tarantool_kv_database/repo"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/tarantool/go-tarantool"
)

type DataKeyValue struct {
	Key   string                 `json:"key"`
	Value map[string]interface{} `json:"value"`
}

type DataValue struct {
	Value map[string]interface{} `json:"value"`
}

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
func (tm *TarantoolManager) InsertValueHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	valueToInsert := DataKeyValue{}
	if err := jsonutil.ReadJSON(r.Context(), r, &valueToInsert); err != nil {
		msg := errorconstants.ErrBadPayload
		logger.Error().Err(errors.Wrap(err, errorconstants.ErrParseJSON)).Msg(errors.Wrap(err, errorconstants.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, msg)
		return
	}

	// Преобразовываем JSON в строку
	JSONString, err := json.Marshal(valueToInsert.Value)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	valObject, err := repo.InsertValueByKey(r.Context(), tm.TarantoolConn, valueToInsert.Key, string(JSONString))
	if err != nil {
		if errors.Is(err, errorconstants.ErrNotFoundById) {
			logger.Error().Err(err).Msg(err.Error())
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, errorconstants.ErrKeyAlreadyExists) {
			logger.Error().Err(err).Msg(errors.Wrap(errorconstants.ErrKeyAlreadyExists, err.Error()).Error())
			jsonutil.SendError(r.Context(), w, http.StatusConflict, errorconstants.ErrKeyAlreadyExists.Error())
			return
		}
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	var dataJSON map[string]interface{}
	err = json.Unmarshal([]byte(valObject.Data), &dataJSON)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	dKV := DataKeyValue{
		Key:   valObject.Id,
		Value: dataJSON,
	}

	if err := jsonutil.SendJSON(r.Context(), w, dKV); err != nil {
		logger.Error().Err(errors.Wrap(err, errorconstants.ErrSendJSON)).Msg(errors.Wrap(err, errorconstants.ErrSomethingWentWrong).Error())
		return
	}
}

// get
func (tm *TarantoolManager) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	valID := mux.Vars(r)["id"]
	if valID == "" {
		errMsg := errors.New("bad request")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errMsg.Error())
		return
	}

	valObject, err := repo.GetValueByKey(r.Context(), tm.TarantoolConn, valID)
	if err != nil {
		if errors.Is(err, errorconstants.ErrNotFoundById) {
			logger.Error().Err(err).Msg(err.Error())
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, err.Error())
			return
		}
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	var dataJSON map[string]interface{}
	err = json.Unmarshal([]byte(valObject.Data), &dataJSON)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	dKV := DataKeyValue{
		Key:   valObject.Id,
		Value: dataJSON,
	}

	if err := jsonutil.SendJSON(r.Context(), w, dKV); err != nil {
		logger.Error().Err(errors.Wrap(err, errorconstants.ErrSendJSON)).Msg(errors.Wrap(err, errorconstants.ErrSomethingWentWrong).Error())
		return
	}
}

// put
func (tm *TarantoolManager) UpdateValueHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	valID := mux.Vars(r)["id"]
	if valID == "" {
		errMsg := errors.New("bad_request")
		logger.Error().Err(errMsg).Msg("bad_request")
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errMsg.Error())
		return
	}

	valueToInsert := DataValue{}
	if err := jsonutil.ReadJSON(r.Context(), r, &valueToInsert); err != nil {
		msg := errorconstants.ErrBadPayload
		logger.Error().Err(errors.Wrap(err, errorconstants.ErrParseJSON)).Msg(errors.Wrap(err, errorconstants.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, msg)
		return
	}

	// Преобразовываем JSON в строку
	JSONString, err := json.Marshal(valueToInsert.Value)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	valObject, err := repo.UpdateValueByKey(r.Context(), tm.TarantoolConn, valID, string(JSONString))
	if err != nil {
		if errors.Is(err, errorconstants.ErrUpdateValue) {
			logger.Error().Err(err).Msg(err.Error())
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errorconstants.ErrNotFoundById.Error())
			return
		}
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	var dataJSON map[string]interface{}
	err = json.Unmarshal([]byte(valObject.Data), &dataJSON)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	dKV := DataKeyValue{
		Key:   valObject.Id,
		Value: dataJSON,
	}

	if err := jsonutil.SendJSON(r.Context(), w, dKV); err != nil {
		logger.Error().Err(errors.Wrap(err, errorconstants.ErrSendJSON)).Msg(errors.Wrap(err, errorconstants.ErrSomethingWentWrong).Error())
		return
	}

}

// delete
func (tm *TarantoolManager) DeleteValueHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	valID := mux.Vars(r)["id"]
	if valID == "" {
		errMsg := errors.New("bad_request")
		logger.Error().Err(errMsg).Msg("bad_request")
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errMsg.Error())
		return
	}

	valObject, err := repo.DeleteValueByKey(r.Context(), tm.TarantoolConn, valID)
	if err != nil {
		if errors.Is(err, errorconstants.ErrDeleteValue) {
			logger.Error().Err(err).Msg(err.Error())
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errorconstants.ErrNotFoundById.Error())
			return
		}
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	var dataJSON map[string]interface{}
	err = json.Unmarshal([]byte(valObject.Data), &dataJSON)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errorconstants.ErrSomethingWentWrong)
		return
	}

	dKV := DataKeyValue{
		Key:   valObject.Id,
		Value: dataJSON,
	}

	if err := jsonutil.SendJSON(r.Context(), w, dKV); err != nil {
		logger.Error().Err(errors.Wrap(err, errorconstants.ErrSendJSON)).Msg(errors.Wrap(err, errorconstants.ErrSomethingWentWrong).Error())
		return
	}

}
