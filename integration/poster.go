package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/Dimitriy14/notifyme/logger"

	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/models"
	"github.com/Dimitriy14/notifyme/services/common"
)

type Poster interface {
	GetCashShifts(date time.Time) ([]models.CashShift, error)
}

func NewPoster() Poster {
	return &posterImpl{}
}

type posterImpl struct {
}

type posterResponse struct {
	Response []models.CashShift `json:"response"`
}

func (p *posterImpl) GetCashShifts(date time.Time) ([]models.CashShift, error) {
	posterURL, err := url.Parse(config.Conf.PosterURL)
	if err != nil {
		return nil, err
	}
	val := url.Values{}
	val.Set("token", config.Conf.Token)
	val.Add("dateFrom", fmt.Sprintf("%s", date.String()))
	val.Add("dateTo", fmt.Sprintf("%s", date.String()))

	posterURL.Path = "/api/finance.getCashShift"
	posterURL.RawQuery = val.Encode()

	resp, err := http.Get(posterURL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer common.CloseRespBody(resp)

	logger.Log.Debugf("Received body from %s : %s", posterURL.String(), string(body))

	var posterResp posterResponse
	err = json.Unmarshal(body, &posterResp)
	if err != nil {
		return nil, err
	}

	return posterResp.Response, nil
}
