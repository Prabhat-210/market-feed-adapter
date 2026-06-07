package model

import "time"

type MarketTick struct {
	InstrumentKey string    `json:"instrument_key"` // NSE_FO|45450
	LTP           float64   `json:"ltp"`
	CP            float64   `json:"cp"`              // previous close
	Change        float64   `json:"change"`          // ltp - cp
	ChangePct     float64   `json:"change_pct"`      // (change/cp)*100
	LTQ           int64     `json:"ltq"`
	ATP           float64   `json:"atp"`
	VTT           int64     `json:"vtt"`
	OI            float64   `json:"oi"`
	TBQ           int64     `json:"tbq"`
	TSQ           int64     `json:"tsq"`
	LTT           time.Time `json:"ltt"`             // last traded time
	ReceivedAt    time.Time `json:"received_at"`
}