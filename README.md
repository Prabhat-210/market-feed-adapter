# market-feed-adapter

This service has one job — **get price data from the stock exchange and send it to Kafka**.

No database. No business logic. Just data flowing through stages, each one refining it further.

---

## Flow

```
Exchange (Upstox)
      ↓
  Stage 0 — Connect        keep the connection alive
      ↓
  Stage 1 — Ingestor       spread the work across decoder workers
      ↓
  Stage 2 — Decoder        decode, normalize, split by symbol, hash route
      ↓
  Stage 3 — Interpreter    validate and enrich per symbol
      ↓
  Stage 4 — Egress         send to Kafka
      ↓
   Kafka (upstox.market.tick)
```

---

## The Stages

### Stage 0 — Connect
Connects to Upstox WebSocket feed and keeps it alive.
If the connection drops, it reconnects automatically with exponential backoff.
Outputs raw bytes — knows nothing about what the data means.

### Stage 1 — Ingestor
Receives raw bytes and spreads them across a pool of decoders using round-robin.
Does no processing itself — just distributes evenly for CPU parallelism.

```
tick 1  →  decoder 1
tick 2  →  decoder 2
tick 3  →  decoder 3
tick 4  →  decoder 1  (back to start)
```

### Stage 2 — Decoder
Does the heavy lifting:
- Decodes raw bytes (Upstox sends protobuf binary)
- Normalizes fields to standard format
- Splits one message into per-symbol ticks (one Upstox message contains all subscribed symbols)
- Hash routes each symbol to the same interpreter every time — guarantees ordering

```
one raw message (all symbols)
      ↓
decode
      ↓
RELIANCE tick  →  hash("NSE_EQ|RELIANCE") % interpreters  →  interpreter 4  (always)
HDFC tick      →  hash("NSE_EQ|HDFC")     % interpreters  →  interpreter 7  (always)
NIFTY tick     →  hash("NSE_INDEX|Nifty") % interpreters  →  interpreter 2  (always)
```

### Stage 3 — Interpreter
Each interpreter handles a fixed subset of symbols — no locking needed.
- Throws away duplicate or stale ticks
- Calculates price change and change percentage
- Prepares the final clean MarketTick ready for Kafka

### Stage 4 — Egress
Sends data to Kafka.
If Kafka is down — retries 3 times with backoff.
If still failing — saves to Dead Letter Queue so nothing is lost silently.

---

## Why hash routing matters

```
without hash routing
    RELIANCE tick 1  →  decoder 1  →  interpreter 4
    RELIANCE tick 2  →  decoder 2  →  interpreter 4
    
    two decoders feeding same interpreter
    tick 2 might arrive before tick 1
    ordering broken ❌

with hash routing
    RELIANCE tick 1  →  decoder 1  →  interpreter 4
    RELIANCE tick 2  →  decoder 2  →  interpreter 4
    
    same decoder always handles RELIANCE
    same interpreter always handles RELIANCE
    ordering guaranteed ✅
```

---

## What happens when things go wrong

| Problem | What we do |
|---|---|
| Exchange disconnects | Reconnect with exponential backoff |
| Same tick arrives twice | Sequence check, discard duplicate |
| Kafka is down | Retry 3 times, then route to DLQ |
| Decoder full | Drop tick, log warning |
| Interpreter full | Drop tick, log warning |

---

## Folder Structure

```
market-feed-adapter/
├── cmd/
│   ├── main.go                  # production entry point
│   └── mockserver/
│       └── main.go              # development mock server
│
├── internal/
│   ├── feed/
│   │   ├── connection.go        # struct + constructor
│   │   └── connect.go           # stage 0: connect, reconnect, read
│   │
│   ├── pipeline/
│   │   ├── pipeline.go          # wires all stages
│   │   ├── ingestor.go          # stage 1: round-robin to decoders
│   │   ├── decoder.go           # stage 2: decode + hash route
│   │   ├── interpreter.go       # stage 3: validate + enrich
│   │   └── egress.go            # stage 4: kafka publish + retry + dlq
│   │
│   ├── model/
│   │   ├── tick.go              # RawTick, DecodedTick
│   │   └── market_tick.go       # MarketTick (published to Kafka)
│   │
│   ├── publisher/
│   │   └── kafka.go             # kafka producer wrapper
│   │
│   ├── platform/
│   │   ├── config/
│   │   │   └── config.go
│   │   └── logger/
│   │       └── logger.go
│   │
│   └── bootstrap/
│       └── bootstrap.go
│
├── mock/
│   ├── generator.go             # generates realistic Upstox V3 ticks
│   └── server.go                # mock WebSocket server
│
├── .env
├── go.mod
└── go.sum
```

---

## Configuration

```env
# service
SERVICE_NAME=market-feed-adapter
ENVIRONMENT=dev
LOG_LEVEL=debug

# feed
AUTHORIZE_URL=https://api.upstox.com/v3/feed/market-data-feed/authorize
UPSTOX_ACCESS_TOKEN=your_token_here
INSTRUMENTS=NSE_INDEX|Nifty 50,NSE_INDEX|Nifty Bank
RECONNECT_MAX_MS=30000
BUFFER_SIZE=10000

# pipeline
DECODER_COUNT=10
INTERPRETER_COUNT=20

# mock
MOCK_ADDR=localhost:8765
MOCK_INTERVAL_MS=500

# kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=upstox.market.tick
KAFKA_DLQ_TOPIC=upstox.market.tick.dlq
```

---

## Running

```bash
# development — run mock server first
go run cmd/mockserver/main.go

# then run feed adapter (point AUTHORIZE_URL to mock in .env)
go run cmd/main.go

# production
go run cmd/main.go
```