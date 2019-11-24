package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Dimitriy14/notifyme/logger"

	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/models"
	"github.com/Dimitriy14/notifyme/services/common"
)

type Poster interface {
	GetCashShifts(dateFrom models.UnixTime, dateTo models.UnixTime) (models.Shifts, error)
	GetProducts(spotID string, dateFrom models.UnixTime, dateTo models.UnixTime) ([]models.ProductFiler, error)
}

func NewPoster() Poster {
	return &posterImpl{}
}

type posterImpl struct {
}

type posterResponse struct {
	Response []models.CashShift `json:"response"`
}

type posterProductResponse struct {
	Response []models.ProductFiler `json:"response"`
}

func (p *posterImpl) GetCashShifts(dateFrom models.UnixTime, dateTo models.UnixTime) (models.Shifts, error) {
	posterURL, err := url.Parse(config.Conf.PosterURL)
	if err != nil {
		return nil, err
	}
	val := url.Values{}
	val.Set("token", config.Conf.Token)
	val.Add("dateFrom", fmt.Sprintf("%s", dateFrom.Format()))
	val.Add("dateTo", fmt.Sprintf("%s", dateTo.Format()))

	posterURL.Path = "/api/finance.getCashShifts"
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

func (p *posterImpl) GetProducts(spotID string, dateFrom models.UnixTime, dateTo models.UnixTime) ([]models.ProductFiler, error) {
	posterURL, err := url.Parse(config.Conf.PosterURL)
	if err != nil {
		return nil, err
	}
	val := url.Values{}
	val.Set("token", config.Conf.Token)
	val.Add("dateFrom", fmt.Sprintf("%s", dateFrom.Format()))
	val.Add("dateTo", fmt.Sprintf("%s", dateTo.Format()))
	val.Add("spot_id", spotID)

	posterURL.Path = "/api/dash.getProductsSales"
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

	var posterResp posterProductResponse
	err = json.Unmarshal(body, &posterResp)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Received data from %s : %v", posterURL.String(), posterResp.Response)

	return posterResp.Response, nil
}
