package delivery

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func insertHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка при разборе данных формы", http.StatusBadRequest)
		logger.Error().Msgf("Ошибка при разборе данных формы: %v", err)
		return
	}

}
