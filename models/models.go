package models

import (
	"strconv"
	"strings"
	"time"
)

type WebHookTransaction struct {
	ID   int      `json:"object_id"`
	Time UnixTime `json:"time"`
	Data string   `json:"data"`
}

type CashShift struct {
	ID             string   `json:"cash_shift_id"`
	SpotID         string   `json:"spot_id"`
	TimeStart      UnixTime `json:"time_start"`
	TimeEnd        UnixTime `json:"time_end"`
	AmountSellCash string   `json:"amount_sell_cash"`
	AmountSellCard string   `json:"amount_sell_card"`
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
	*t = UnixTime(time.Unix(q, 0))
	return nil
}

func (t UnixTime) String() string {
	return time.Time(t).String()
}

func (t UnixTime) Format() string {
	return time.Time(t).Format("20060102")
}
