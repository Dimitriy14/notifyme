package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Dimitriy14/notifyme/config"
	"github.com/Dimitriy14/notifyme/models"
	"github.com/Dimitriy14/notifyme/services/common"
)

type Poster interface {
	GetCashShiftByID(id int) (models.CashShift, error)
}

func NewPoster() Poster {
	return &posterImpl{}
}

type posterImpl struct {
}

type posterResponse struct {
	Response models.CashShift `json:"response"`
}

func (p *posterImpl) GetCashShiftByID(id int) (models.CashShift, error) {
	posterURL, err := url.Parse(config.Conf.PosterURL)
	if err != nil {
		return models.CashShift{}, err
	}
	val := url.Values{}
	val.Set("token", config.Conf.Token)
	val.Add("cash_shift_id", fmt.Sprintf("%d", id))

	posterURL.Path = "/api/finance.getCashShift"
	posterURL.RawQuery = val.Encode()

	resp, err := http.Get(posterURL.String())
	if err != nil {
		return models.CashShift{}, err
	}
	defer common.CloseRespBody(resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.CashShift{}, err
	}

	var posterResp posterResponse
	err = json.Unmarshal(body, &posterResp)
	if err != nil {
		return models.CashShift{}, err
	}

	return posterResp.Response, nil
}
