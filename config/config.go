package config

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
)

var (
	Conf Configuration

	FilePath = "config.json"
)

type Configuration struct {
	Port     string `json:"Port" default:"4444" environment:"PORT"`
	BasePath string `json:"BasePath" default:"/notifyme"`

	PosterURL string `json:"PosterURL" default:"https://joinposter.com"`
	Token     string `json:"Token"`
}

func Load() error {
	if err := readFile(&Conf); err != nil {
		return err
	}

	if err := readEnv(&Conf); err != nil {
		return err
	}

	fmt.Printf("%+v\n", Conf)
	return nil
}

func readFile(cfg *Configuration) error {
	fileContent, err := os.Open(FilePath)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(fileContent).Decode(&Conf); err != nil {
		return err
	}
	return nil
}

func readEnv(cfg *Configuration) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}
	return nil
}