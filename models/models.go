package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type WebHookTransaction struct {
	ID   int      `json:"object_id"`
	Time UnixTime `json:"time"`
}

type CashShift struct {
	ID             uint64   `json:"cash_shift_id"`
	SpotID         uint64   `json:"spot_id"`
	TimeStart      UnixTime `json:"time_start"`
	TimeEnd        UnixTime `json:"time_end"`
	AmountSellCash uint64   `json:"amount_sell_cash"`
	AmountSellCard uint64   `json:"amount_sell_card"`
}

type Shifts []CashShift

type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *UnixTime) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q/1000, 0)
	return
}

func (t UnixTime) String() string { return fmt.Sprintf("%d", time.Time(t).Unix()) }
