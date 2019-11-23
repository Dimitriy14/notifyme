package shift

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Dimitriy14/notifyme/integration"
	"github.com/Dimitriy14/notifyme/logger"
	"github.com/Dimitriy14/notifyme/models"
	"github.com/Dimitriy14/notifyme/services/common"
)

type Closer interface {
	Close(w http.ResponseWriter, r *http.Request)
}

func NewShiftService(poster integration.Poster) Closer {
	return &closerImpl{
		poster: poster,
	}
}

type closerImpl struct {
	poster integration.Poster
}

func (c *closerImpl) Close(w http.ResponseWriter, r *http.Request) {
	var (
		shift models.ClosedShift
	)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log.Errorf("Reading body: err=%s", err)
		common.SendError(w, http.StatusBadRequest, "Reading body err= %s\n", err)
	}
	defer common.CloseReqBody(r)

	if err = json.Unmarshal(body, &shift); err != nil {
		logger.Log.Errorf("Unmarshaling: err=%s", err)
		common.SendError(w, http.StatusBadRequest, "Unmarshal body err= %s\n", err)
		return
	}

	cashShift, err := c.poster.GetCashShiftByID(shift.ID)
	if err != nil {
		logger.Log.Errorf("GetCashShiftByID: err=%s", err)
		common.SendError(w, http.StatusInternalServerError, "Unmarshal body err= %s\n", err)
		return
	}

	common.RenderJSON(w, &cashShift)
}
