package shift

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Dimitriy14/notifyme/config"

	"github.com/Dimitriy14/notifyme/repository"

	"github.com/Dimitriy14/notifyme/integration"
	"github.com/Dimitriy14/notifyme/logger"
	"github.com/Dimitriy14/notifyme/models"
	"github.com/Dimitriy14/notifyme/services/common"
)

type Closer interface {
	Close(w http.ResponseWriter, r *http.Request)
}

func NewShiftService(poster integration.Poster, repo repository.Repository) Closer {
	return &closerImpl{
		poster: poster,
		repo:   repo,
	}
}

type closerImpl struct {
	poster integration.Poster
	repo   repository.Repository
}

func (c *closerImpl) Close(w http.ResponseWriter, r *http.Request) {
	var (
		tx models.WebHookTransaction
	)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log.Errorf("Reading body: err=%s", err)
		common.SendError(w, http.StatusBadRequest, "Reading body err= %s\n", err)
	}
	defer common.CloseReqBody(r)

	logger.Log.Debugf("Received body: %s", string(body))
	if err = json.Unmarshal(body, &tx); err != nil {
		logger.Log.Errorf("Unmarshaling: err=%s", err)
		common.SendError(w, http.StatusBadRequest, "Unmarshal body err= %s\n", err)
		return
	}

	logger.Log.Debugf("data: %s", tx.Time.String())

	if !strings.Contains(tx.Data, "shift_close") {
		return
	}

	cashShifts, err := c.poster.GetCashShifts(tx.Time, tx.Time)
	if err != nil {
		logger.Log.Errorf("GetCashShiftByID: err=%s", err)
		common.SendError(w, http.StatusInternalServerError, "Unmarshal body err= %s\n", err)
		return
	}

	for spotID, shift := range cashShifts.GetMapShifts() {
		ps, err := c.poster.GetProducts(spotID, tx.Time, tx.Time)
		if err != nil {
			logger.Log.Errorf("GetProducts: err=%s", err)
			common.SendError(w, http.StatusInternalServerError, "GetProducts err= %s\n", err)
			return
		}

		filters, err := c.repo.GetFilters()
		if err != nil {
			logger.Log.Errorf("GetFilters: err=%s", err)
			common.SendError(w, http.StatusInternalServerError, "GetProducts err= %s\n", err)
			return
		}

		var result []models.ProductFiler

		for _, f := range filters {
			for _, p := range ps {
				if p.ProductID == f.ProductID {
					p.UserEmail = f.UserEmail
					result = append(result, p)
				}
			}
		}

		mail := models.Mail{
			Date:           tx.Time.Format(),
			SpotID:         spotID,
			SpotName:       shift.SpotName,
			AmountSellCash: shift.AmountSellCash,
			AmountSellCard: shift.AmountSellCard,
			Products:       result,
		}

		content, err := json.Marshal(&mail)
		if err != nil {
			logger.Log.Errorf("Result marshal: err=%s", err)
			common.SendError(w, http.StatusInternalServerError, "Result marshal err= %s\n", err)
			return
		}

		logger.Log.Debugf("Mail content: %s", string(content))
		response, err := http.Post(config.Conf.MailServiceURL, "application/json", bytes.NewReader(content))
		if err != nil {
			logger.Log.Errorf("Sending to mail service(url=%s): err=%s", config.Conf.MailServiceURL, err)
			common.SendError(w, http.StatusInternalServerError, "Result marshal err= %s\n", err)
			return
		}

		bodyResp, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Log.Errorf("Sending to mail service(url=%s): err=%s", config.Conf.MailServiceURL, err)
			common.SendError(w, http.StatusInternalServerError, "Result marshal err= %s\n", err)
			return
		}
		logger.Log.Debugf("Mail service %s response status: %d   %s", config.Conf.MailServiceURL, response.StatusCode, string(bodyResp))
	}

	common.RenderJSON(w, &cashShifts)
}
