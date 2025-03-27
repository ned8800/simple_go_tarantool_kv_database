package delivery

import (
	"net/http"
	tarantoolmanager "simple_go_tarantool_kv_database/delivery/tarantool_manager"
	"simple_go_tarantool_kv_database/middleware"

	"github.com/gorilla/mux"
	"github.com/tarantool/go-tarantool"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	return router
}

func SetupRoutes(router *mux.Router, tarantoolConn *tarantool.Connection) {

	tm := tarantoolmanager.NewTarantoolManager(tarantoolConn)

	router.HandleFunc("/kv", tm.InsertValueHandler).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/kv/{id}", tm.UpdateValueHandler).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/kv/{id}", tm.GetValueHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/kv/{id}", tm.DeleteValueHandler).Methods(http.MethodDelete, http.MethodOptions)

}

func ApplyMiddlewares(router *mux.Router) {
	router.Use(middleware.AccessLogMiddleware)
	router.Use(middleware.WithLogger)
}
