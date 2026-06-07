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

## hash routing 

```
with hash routing
    RELIANCE tick 1  →  decoder 1  →  interpreter 4
    RELIANCE tick 2  →  decoder 2  →  interpreter 4
    
    same decoder always handles RELIANCE
    same interpreter always handles RELIANCE
    ordering guaranteed 
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
