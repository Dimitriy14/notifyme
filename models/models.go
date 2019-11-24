package models

import (
	"strconv"
	"strings"
	"time"
)

type Mail struct {
	Date           string         `json:"date"`
	SpotID         string         `json:"spot_id"`
	SpotName       string         `json:"spot_name"`
	AmountSellCash int            `json:"amount_sell_cash,string"`
	AmountSellCard int            `json:"amount_sell_card,string"`
	Products       []ProductFiler `json:"products"`
}

type WebHookTransaction struct {
	ID   int      `json:"object_id"`
	Time UnixTime `json:"time"`
	Data string   `json:"data"`
}

type CashShift struct {
	ID             string `json:"cash_shift_id"`
	SpotID         string `json:"spot_id"`
	SpotName       string `json:"spot_name"`
	SpotAddress    string `json:"spot_adress"`
	AmountSellCash int    `json:"amount_sell_cash,string"`
	AmountSellCard int    `json:"amount_sell_card,string"`
}

type Shifts []CashShift

func (s Shifts) GetMapShifts() map[string]CashShift {
	m := make(map[string]CashShift)
	for _, shift := range s {
		currentShift, ok := m[shift.SpotID]
		if !ok {
			m[shift.SpotID] = shift
			continue
		}

		currentShift.AmountSellCard += shift.AmountSellCard
		currentShift.AmountSellCash += shift.AmountSellCash

		m[shift.SpotID] = currentShift
	}

	return m
}

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
