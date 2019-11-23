package models

import "time"

type WebHookTransaction struct {
	ID   int       `json:"object_id"`
	Time time.Time `json:"time"`
}

type CashShift struct {
	ID             uint64    `json:"cash_shift_id"`
	SpotID         uint64    `json:"spot_id"`
	TimeStart      time.Time `json:"time_start"`
	TimeEnd        time.Time `json:"time_end"`
	AmountSellCash uint64    `json:"amount_sell_cash"`
	AmountSellCard uint64    `json:"amount_sell_card"`
}

type Shifts []CashShift
