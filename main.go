package main

import (
	"fmt"
	"github.com/Dimitriy14/notifyme/apploader"
	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/services"
	"log"
	"net/http"
)

func main() {
	if err := apploader.LoadApplicationServices(); err != nil {
		log.Fatal(err)
	}

	handler := services.NewRouter()

	fmt.Printf("Started listening on port: %s\n", config.Conf.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Conf.Port), handler); err != nil {
		log.Fatal(err)
	}
}
