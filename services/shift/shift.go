package shift

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

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

	var products []models.ProductFiler
	for spotID := range cashShifts.GetMapShifts() {
		ps, err := c.poster.GetProducts(spotID, tx.Time, tx.Time)
		if err != nil {
			logger.Log.Errorf("GetProducts: err=%s", err)
			common.SendError(w, http.StatusInternalServerError, "GetProducts err= %s\n", err)
			return
		}

		products = append(products, ps...)
	}

	filters, err := c.repo.GetFilters()
	if err != nil {
		logger.Log.Errorf("GetFilters: err=%s", err)
		common.SendError(w, http.StatusInternalServerError, "GetProducts err= %s\n", err)
		return
	}

	var result []models.ProductFiler

	for _, f := range filters {
		for _, p := range products {
			if p.ProductID == f.ProductID {
				result = append(result, f)
			}
		}
	}

	common.RenderJSON(w, &result)
}
