package filters

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Dimitriy14/notifyme/integration"
	"github.com/Dimitriy14/notifyme/logger"
	"github.com/Dimitriy14/notifyme/models"
	"github.com/Dimitriy14/notifyme/repository"
	"github.com/Dimitriy14/notifyme/services/common"
)

type Filter interface {
	GetFilter(w http.ResponseWriter, r *http.Request)
	AddFilter(w http.ResponseWriter, r *http.Request)
	DeleteFilter(w http.ResponseWriter, r *http.Request)
}

func NewFilterService(poster integration.Poster, repo repository.Repository) Filter {
	return &filterImpl{
		poster: poster,
		repo:   repo,
	}
}

type filterImpl struct {
	poster integration.Poster
	repo   repository.Repository
}

func (f *filterImpl) GetFilter(w http.ResponseWriter, r *http.Request) {
	filters, err := f.repo.GetFilters()
	if err != nil {
		logger.Log.Errorf("getting filters err=%s", err)
		common.SendInternalServerError(w, "getting filters", err)
		return
	}

	common.RenderJSON(w, &filters)
}

func (f *filterImpl) AddFilter(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log.Errorf("reading body err=%s", err)
		common.SendError(w, http.StatusBadRequest, "reading body", err)
		return
	}
	defer common.CloseReqBody(r)

	var filters []models.ProductFiler
	err = json.Unmarshal(body, &filters)
	if err != nil {
		logger.Log.Errorf("unmarshal body err=%s", err)
		common.SendError(w, http.StatusBadRequest, "unmarshal body", err)
		return
	}

	filters, err = f.repo.SaveFilter(filters)
	if err != nil {
		logger.Log.Errorf("saving filter err=%s", err)
		common.SendInternalServerError(w, "saving filter", err)
		return
	}

	common.RenderJSON(w, &filters)
}

func (f *filterImpl) DeleteFilter(w http.ResponseWriter, r *http.Request) {
}
