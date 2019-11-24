package shift

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
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
					result = append(result, f)
				}
			}
		}

		mail := models.Mail{
			SpotID:         spotID,
			SpotName:       shift.SpotName,
			SpotAddress:    shift.SpotAddress,
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

		m := fmt.Sprintf("Comment: %s", string(content))
		mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
		subject := "Subject: " + "info" + "!\n"
		msg := []byte(subject + mime + "\n" + m)

		logger.Log.Debugf("mail body: %s", string(content))
		a := smtp.PlainAuth("", "yankovskiydy98@gmail.com", config.Conf.GmailPassword, "smtp.gmail.com")
		err = smtp.SendMail("smtp.gmail.com:587", a, "yankovskiydy98@gmail.com", []string{"road2ps@gmail.com"}, []byte(msg))
		if err != nil {
			logger.Log.Errorf("Sending mail: err=%s", err)
			common.SendError(w, http.StatusInternalServerError, "sending mail err= %s\n", err)
			return
		}
	}

	common.RenderJSON(w, &cashShifts)
}
