package services

import (
	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/integration"
	"github.com/Dimitriy14/notifyme/services/shift"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true).PathPrefix(config.Conf.BasePath).Subrouter()
	poster := integration.NewPoster()
	shiftService := shift.NewShiftService(poster)

	router.HandleFunc("/close", shiftService.Close).Methods(http.MethodPost)
	return router
}
