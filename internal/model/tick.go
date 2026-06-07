// model/tick.go

package model

import "time"

// raw protobuf bytes from feed
type RawTick struct {
	Data       []byte
	ReceivedAt time.Time
}

// after protobuf decode — matches Upstox V3 fullFeed structure
type DecodedTick struct {
	InstrumentKey string    // e.g. "NSE_FO|45450"
	LTP           float64   // ltp — last traded price
	LTT           int64     // ltt — last traded time (epoch ms)
	LTQ           int64     // ltq — last traded quantity
	CP            float64   // cp — close price (previous day)
	ATP           float64   // atp — average traded price
	VTT           int64     // vtt — volume traded today
	OI            float64   // oi — open interest
	TBQ           int64     // tbq — total buy quantity
	TSQ           int64     // tsq — total sell quantity
	ReceivedAt    time.Time // when we received the raw tick
}
