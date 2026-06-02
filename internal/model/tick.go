package model

import "time"

// raw protobuf bytes from feed
type RawTick struct {
	Data      []byte
	ReceivedAt time.Time
}

type DecodeBytes struct {
	InstrumentKey string    // e.g. "NSE_EQ|INE020B01018"
	Symbol        string    // e.g. "RELIANCE"
	Exchnage      string    // e.g. "NSE_EQ"
	LTP           string    // last traded price
	LTQ           int64     // last traded quantity
	LTT           time.Time // last traded time
	ClosePrice    float64   // previous close
	ATP           float64   //Avg traded price
	Volume        float64   // /day volume
	OpenIntrest   float64
	TBQ           float64 // total buy quantity
	TSQ           float64 // total sell quantity
	Timestamp     time.Time
	Sequence      int64 // for duplicate detection
}
