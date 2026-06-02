package model

import "time"

// normalized model published to Kafka topic: price.updated
type MarketTick struct {
	InstrumentKey string    `json:"instrument_key"` 
	Symbol        string    `json:"symbol"`         
	Exchange      string    `json:"exchange"`       
	LTP           float64   `json:"ltp"`            
	ClosePrice    float64   `json:"close_price"`    
	Change        float64   `json:"change"`         // ltp - close
	ChangePct     float64   `json:"change_pct"`     // (change / close) * 100
	Volume        int64     `json:"volume"`        
	ATP           float64   `json:"atp"`           
	OI            float64   `json:"oi"`             
	TBQ           float64   `json:"tbq"`           
	TSQ           float64   `json:"tsq"`            
	Timestamp     time.Time `json:"timestamp"`      
	ReceivedAt    time.Time `json:"received_at"`    
	Sequence      int64     `json:"sequence"`      
}
