package services

import (
	"net/http"

	"github.com/Dimitriy14/notifyme/postgres"
	"github.com/Dimitriy14/notifyme/repository"

	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/integration"
	"github.com/Dimitriy14/notifyme/services/filters"
	"github.com/Dimitriy14/notifyme/services/shift"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true).PathPrefix(config.Conf.BasePath).Subrouter()
	poster := integration.NewPoster()
	shiftService := shift.NewShiftService(poster)
	repo := repository.NewRepo(postgres.Client)
	filterService := filters.NewFilterService(poster, repo)

	router.HandleFunc("/close", shiftService.Close).Methods(http.MethodPost)
	router.HandleFunc("/filter", filterService.GetFilter).Methods(http.MethodGet)
	router.HandleFunc("/filter", filterService.AddFilter).Methods(http.MethodPost)
	router.HandleFunc("/filter", filterService.DeleteFilter).Methods(http.MethodDelete)
	return router
}
