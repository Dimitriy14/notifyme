package apploader

import (
	"fmt"

	"github.com/Dimitriy14/notifyme/postgres"

	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/logger"
	"github.com/pkg/errors"
)

// LoaderList is a collection of Load() functions
type LoaderList []struct {
	name string
	load func() error
}

var basicLoaders = LoaderList{
	{"config", config.Load}, //config should be loaded first
	{"logger", logger.Load},
}

var clientLoaders = LoaderList{
	{"postgres", postgres.Load},
}

func LoadApplicationServices() error {
	err := executeLoaders(basicLoaders)
	if err != nil {
		return err
	}

	err = executeLoaders(clientLoaders)
	if err != nil {
		return err
	}

	return nil
}

func executeLoaders(loaders LoaderList) error {
	for _, loader := range loaders {
		err := loader.load()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to execute %s.Load()", loader.name))
		}
	}
	return nil
}
