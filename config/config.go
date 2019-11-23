package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
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

	LogLevel string `json:"LogLevel" default:"debug"`
	LogFile  string

	HerokuPg string `json:"HerokuPg" environment:"DATABASE_URL"`

	Postgres struct {
		Host     string `json:"Host"`
		Port     string `json:"Port"`
		DBName   string `json:"DBName" default:"notifyme"`
		User     string `json:"User" default:"admin"`
		Password string `json:"Password" default:"1488"`
	} `json:"Postgres"`
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
