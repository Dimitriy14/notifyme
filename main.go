package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dimitriy14/notifyme/apploader"
	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/services"
	"github.com/urfave/negroni"
)

func main() {
	if err := apploader.LoadApplicationServices(); err != nil {
		log.Fatal(err)
	}

	handler := services.NewRouter()

	middlewareManager := negroni.New()
	middlewareManager.Use(negroni.NewRecovery())

	middlewareManager.UseHandler(handler)

	fmt.Printf("Started listening on port: %s\n", config.Conf.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Conf.Port), middlewareManager); err != nil {
		log.Fatal(err)
	}
}
