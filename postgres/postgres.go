package postgres

import (
	"fmt"

	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/logger"
	"github.com/Dimitriy14/notifyme/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const dbInfo = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"

var Client PGClient

type PGClient struct {
	Session *gorm.DB
}

func Load() error {
	url := config.Conf.HerokuPg
	if url == "" {
		url = fmt.Sprintf(
			dbInfo,
			"localhost",
			"5432",
			"admin",
			"1488",
			"notifyme",
		)
	}

	db, err := gorm.Open("postgres", url)
	if err != nil {
		return err
	}

	Client = PGClient{Session: db}
	db.SetLogger(logger.NewGormLogger(logger.Log))
	db.LogMode(true)

	db.AutoMigrate(&models.ProductFiler{})
	return nil
}
