package mock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// ---- Upstox V3 exact payload structure ----

type FeedResponse struct {
	Type  string                    `json:"type"`
	Feeds map[string]InstrumentFeed `json:"feeds"`
}

type InstrumentFeed struct {
	FullFeed FullFeed `json:"fullFeed"`
}

type FullFeed struct {
	MarketFF MarketFF `json:"marketFF"`
}

type MarketFF struct {
	LTPC        LTPC        `json:"ltpc"`
	MarketLevel MarketLevel `json:"marketLevel"`
	MarketOHLC  MarketOHLC  `json:"marketOHLC"`
	ATP         float64     `json:"atp"`
	VTT         string      `json:"vtt"`
	OI          float64     `json:"oi"`
	TBQ         int64       `json:"tbq"`
	TSQ         int64       `json:"tsq"`
}

type LTPC struct {
	LTP float64 `json:"ltp"`
	LTT string  `json:"ltt"` // epoch ms as string
	LTQ string  `json:"ltq"` // quantity as string
	CP  float64 `json:"cp"`  // previous close
}

type MarketLevel struct {
	BidAskQuote []BidAskQuote `json:"bidAskQuote"`
}

type BidAskQuote struct {
	BidQ string  `json:"bidQ"`
	BidP float64 `json:"bidP"`
	AskQ string  `json:"askQ"`
	AskP float64 `json:"askP"`
}

type MarketOHLC struct {
	OHLC []OHLC `json:"ohlc"`
}

type OHLC struct {
	Interval string  `json:"interval"`
	Open     float64 `json:"open"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Close    float64 `json:"close"`
	Vol      string  `json:"vol"`
	TS       string  `json:"ts"` // epoch ms as string
}

// ---- Mock instruments with realistic base prices ----

type instrument struct {
	key        string
	basePrice  float64
	closePrice float64
}

var instruments = []instrument{
	{key: "NSE_EQ|INE002A01018", basePrice: 2500.00, closePrice: 2480.00},    // RELIANCE
	{key: "NSE_EQ|INE040A01034", basePrice: 1650.00, closePrice: 1640.00},    // HDFC BANK
	{key: "NSE_EQ|INE009A01021", basePrice: 820.00, closePrice: 815.00},      // INFOSYS
	{key: "NSE_EQ|INE467B01029", basePrice: 3800.00, closePrice: 3790.00},    // TCS
	{key: "NSE_INDEX|Nifty 50", basePrice: 22000.00, closePrice: 21950.00},   // NIFTY 50
	{key: "NSE_INDEX|Nifty Bank", basePrice: 47000.00, closePrice: 46900.00}, // BANK NIFTY
}

// Generator generates mock Upstox V3 payloads
type Generator struct {
	prices map[string]float64
	rng    *rand.Rand
}

func NewGenerator() *Generator {
	g := &Generator{
		prices: make(map[string]float64),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	// init prices from base
	for _, inst := range instruments {
		g.prices[inst.key] = inst.basePrice
	}
	return g
}

// Next generates one realistic tick payload — same structure as Upstox V3
func (g *Generator) Next() ([]byte, error) {
	feeds := make(map[string]InstrumentFeed)

	for _, inst := range instruments {
		// simulate realistic price movement
		// small random walk ±0.1% per tick
		change := g.prices[inst.key] * (g.rng.Float64()*0.002 - 0.001)
		g.prices[inst.key] = roundTo(g.prices[inst.key]+change, 2)
		ltp := g.prices[inst.key]

		// realistic volume
		ltq := int64(g.rng.Intn(500)+1) * 25
		vtt := int64(g.rng.Intn(5000000) + 100000)
		tbq := int64(g.rng.Intn(1000000) + 50000)
		tsq := int64(g.rng.Intn(1000000) + 50000)

		// bid/ask spread — realistic 0.05% spread
		spread := ltp * 0.0005
		bidAsk := []BidAskQuote{
			{BidQ: fmt.Sprintf("%d", ltq*3), BidP: roundTo(ltp-spread, 2), AskQ: fmt.Sprintf("%d", ltq*2), AskP: roundTo(ltp+spread, 2)},
			{BidQ: fmt.Sprintf("%d", ltq*2), BidP: roundTo(ltp-spread*2, 2), AskQ: fmt.Sprintf("%d", ltq*3), AskP: roundTo(ltp+spread*2, 2)},
			{BidQ: fmt.Sprintf("%d", ltq), BidP: roundTo(ltp-spread*3, 2), AskQ: fmt.Sprintf("%d", ltq), AskP: roundTo(ltp+spread*3, 2)},
			{BidQ: fmt.Sprintf("%d", ltq*4), BidP: roundTo(ltp-spread*4, 2), AskQ: fmt.Sprintf("%d", ltq*2), AskP: roundTo(ltp+spread*4, 2)},
			{BidQ: fmt.Sprintf("%d", ltq*2), BidP: roundTo(ltp-spread*5, 2), AskQ: fmt.Sprintf("%d", ltq*4), AskP: roundTo(ltp+spread*5, 2)},
		}

		now := time.Now()
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 9, 15, 0, 0, now.Location())

		feeds[inst.key] = InstrumentFeed{
			FullFeed: FullFeed{
				MarketFF: MarketFF{
					LTPC: LTPC{
						LTP: ltp,
						LTT: fmt.Sprintf("%d", now.UnixMilli()),
						LTQ: fmt.Sprintf("%d", ltq),
						CP:  inst.closePrice,
					},
					MarketLevel: MarketLevel{
						BidAskQuote: bidAsk,
					},
					MarketOHLC: MarketOHLC{
						OHLC: []OHLC{
							{
								Interval: "1d",
								Open:     inst.basePrice,
								High:     roundTo(inst.basePrice*1.02, 2),
								Low:      roundTo(inst.basePrice*0.98, 2),
								Close:    ltp,
								Vol:      fmt.Sprintf("%d", vtt),
								TS:       fmt.Sprintf("%d", dayStart.UnixMilli()),
							},
							{
								Interval: "I1",
								Open:     roundTo(ltp-change*2, 2),
								High:     roundTo(ltp+spread, 2),
								Low:      roundTo(ltp-spread, 2),
								Close:    ltp,
								Vol:      fmt.Sprintf("%d", ltq*10),
								TS:       fmt.Sprintf("%d", now.Add(-time.Minute).UnixMilli()),
							},
						},
					},
					ATP: roundTo((ltp+inst.closePrice)/2, 2),
					VTT: fmt.Sprintf("%d", vtt),
					OI:  float64(g.rng.Intn(500000) + 10000),
					TBQ: tbq,
					TSQ: tsq,
				},
			},
		}
	}

	response := FeedResponse{
		Type:  "live_feed",
		Feeds: feeds,
	}

	return json.Marshal(response)
}

func roundTo(val float64, places int) float64 {
	pow := 1.0
	for i := 0; i < places; i++ {
		pow *= 10
	}
	return float64(int(val*pow+0.5)) / pow
}
