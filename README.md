# market-feed-adapter

This service has one job — **get price data from the stock exchange and send it to Kafka**.

No database. No business logic.

---

## Flow

here each step is for filtering data(refining data)

```
Exchange (NSE/BSE)
      ↓
  Stage 0 — Connect        keep the connection alive
      ↓
  Stage 1 — Ingestor       spread the work across workers
      ↓
  Stage 2 — Collector      decode and clean the raw data
      ↓
  Stage 3 — Interpreter    validate and enrich the data
      ↓
  Stage 4 — Egress         send to Kafka
      ↓
   Kafka 
```

---

## The Stages

### Stage 0 — Connect
Connects to the exchange WebSocket and keeps it alive.
If the connection drops, it reconnects automatically.
Outputs raw bytes — knows nothing about what the data means.

### Stage 1 — Ingestor
Receives raw bytes and spreads them across a pool of collectors.
Like a manager handing tasks to a team — one task per worker, cycling through.
Does no processing itself — just distributes.

### Stage 2 — Collector
Does the heavy lifting:
- Decodes the raw bytes (exchange sends compressed/binary data)
- Converts to a readable format
- Groups multiple updates for the same stock together
- Sends each stock's data to the same interpreter every time (so order is never mixed up)

### Stage 3 — Interpreter
Checks and enriches the data:
- Throws away duplicate or old ticks
- Calculates price change, percentage change
- Prepares the final clean model ready for Kafka

### Stage 4 — Egress
Sends data to Kafka.
If Kafka is down — retries 3 times.
If still failing — saves to a Dead Letter Queue so nothing is lost silently.

---

## What happens when things go wrong

| Problem | What we do |
|---|---|
| Exchange disconnects | Reconnect automatically |
| Same tick arrives twice | Detect and discard the duplicate |
| Kafka is down | Retry 3 times, then save to DLQ |
| Data arrives too fast | Drop oldest data, keep newest |


## Configuration

```env
FEED_URL=wss://exchange.example.com/feed
COLLECTOR_COUNT=10
INTERPRETER_COUNT=20
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=price.updated
```

---